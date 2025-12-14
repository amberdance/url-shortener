package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amberdance/url-shortener/internal/app/usecase"
	"github.com/amberdance/url-shortener/internal/app/usecase/url"
	"github.com/amberdance/url-shortener/internal/domain/errs"
	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/amberdance/url-shortener/internal/infrastructure/helpers"
	"github.com/amberdance/url-shortener/internal/mocks"
	"github.com/amberdance/url-shortener/internal/ports/webapi/dto"
	helpers2 "github.com/amberdance/url-shortener/internal/ports/webapi/helpers"
	"github.com/amberdance/url-shortener/internal/ports/webapi/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type URLShortenerHandlerTestSuite struct {
	suite.Suite
	ctx  context.Context
	ctrl *gomock.Controller

	host       string
	repository *mocks.MockURLRepository
	handler    *URLShortenerHandler
}

func (s *URLShortenerHandlerTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.ctrl = gomock.NewController(s.T())
	s.host = "http://127.0.0.1:9999/"
	s.repository = mocks.NewMockURLRepository(s.ctrl)
	useCases := usecase.URLUseCases{
		Create:      url.NewCreateUseCase(s.repository),
		CreateBatch: url.NewBatchCreateURLUseCase(s.repository),
		GetByURL:    url.NewGetByHashUseCase(s.repository),
	}
	s.handler = NewURLShortenerHandler(s.host, useCases, validator.New(), mocks.NewMockLogger(s.ctrl))
}

func (s *URLShortenerHandlerTestSuite) TearDownSuite() {
	s.ctrl.Finish()
}

func TestURLShortenerHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(URLShortenerHandlerTestSuite))
}

func (s *URLShortenerHandlerTestSuite) TestShortedPlainTextPost_Success() {
	var (
		originalURL  = "https://hard2code.ru"
		expectedHash = "hash"
	)

	s.repository.EXPECT().
		Create(gomock.Any(), gomock.AssignableToTypeOf(&model.URLEntry{})).
		DoAndReturn(func(_ context.Context, m *model.URLEntry) error {
			m.Hash = expectedHash
			return nil
		}).
		Times(1)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(originalURL))
	w := httptest.NewRecorder()
	s.handler.Routes().ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	s.Equal(http.StatusCreated, res.StatusCode)
	s.Equal("text/plain", res.Header.Get("Content-Type"))

	body, err := io.ReadAll(res.Body)
	s.NoError(err)

	expected := s.host + expectedHash
	s.Equal(expected, string(body))
}

func (s *URLShortenerHandlerTestSuite) Test_When_RequestPayloadInvalid_Then_ShortenPlainTextPost_ShouldReturn400HttpCode() {
	for _, tt := range invalidPayloadCases {
		s.T().Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", tt.body)
			w := httptest.NewRecorder()
			s.handler.Routes().ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			s.Equal(http.StatusBadRequest, res.StatusCode)
		})
	}
}

func (s *URLShortenerHandlerTestSuite) Test_When_RequestPayloadInvalid_Then_Shorten_ShouldReturn400HttpCode() {
	for _, tt := range invalidPayloadCases {
		s.T().Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(tt.json))
			w := httptest.NewRecorder()
			s.handler.Routes().ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			s.Equal(http.StatusBadRequest, res.StatusCode)
		})
	}
}

var invalidPayloadCases = []struct {
	name string
	body io.Reader
	json []byte
}{
	{name: "Nil payload", body: bytes.NewBuffer(nil), json: []byte(`{"url": null}`)},
	{name: "Only spaces", body: bytes.NewBufferString("   "), json: []byte(`{"url": "   "}`)},
	{name: "Empty payload", body: bytes.NewBufferString(""), json: []byte(`{"url": ""}`)},
}

func (s *URLShortenerHandlerTestSuite) TestShorten_Success() {
	var (
		originalURL = "https://hard2code.ru"
		entry       = &model.URLEntry{
			ID:          uuid.New(),
			Hash:        helpers.GenerateHash(),
			OriginalURL: originalURL,
		}
		reqDto = dto.ShortURLRequest{
			URL: originalURL,
		}
	)

	body, err := json.Marshal(reqDto)
	s.NoError(err)

	s.repository.EXPECT().
		Create(gomock.Any(), gomock.AssignableToTypeOf(&model.URLEntry{})).
		DoAndReturn(func(_ context.Context, m *model.URLEntry) error {
			m.ID = entry.ID
			m.Hash = entry.Hash
			m.OriginalURL = originalURL

			return nil
		}).
		Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	s.handler.Routes().ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	s.Equal(http.StatusCreated, res.StatusCode)
	s.Equal(middleware.ContentTypeJSONHeaderValue, res.Header.Get("Content-Type"))

	body, err = io.ReadAll(res.Body)
	s.NoError(err)

	var response dto.ShortURLResponse
	s.NoError(json.Unmarshal(body, &response))
	s.Equal(response.URL, s.host+entry.Hash)
}

func (s *URLShortenerHandlerTestSuite) Test_When_URLWithGivenHashDoesNotExists_Then_GetShouldReturn404HttpCode() {
	hash := helpers.GenerateHash()

	s.repository.EXPECT().
		FindByHash(gomock.Any(), hash).
		Return(nil, errs.ErrNotFound).
		Times(1)

	w := httptest.NewRecorder()
	s.handler.Routes().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/"+hash, nil))
	res := w.Result()
	defer res.Body.Close()

	s.Equal(http.StatusNotFound, res.StatusCode)
}

func (s *URLShortenerHandlerTestSuite) setupMockFor409HttCode(originalURL string, hash string) {
	s.repository.EXPECT().
		Create(gomock.Any(), gomock.AssignableToTypeOf(&model.URLEntry{})).
		Return(errs.ErrDuplicate).
		Times(1)

	s.repository.EXPECT().
		FindByOriginalURL(gomock.Any(), originalURL).
		DoAndReturn(func(_ context.Context, originalURL string) (*model.URLEntry, error) {
			return &model.URLEntry{
					OriginalURL: originalURL,
					Hash:        hash,
				},
				nil
		}).
		Times(1)
}

func (s *URLShortenerHandlerTestSuite) Test_When_URLExists_Then_ShortenShouldReturn409HttpCodeAndURL() {
	var (
		originalURL = "https://hard2code.ru"
		hash        = helpers.GenerateHash()
		req         = dto.ShortURLRequest{
			URL: originalURL,
		}
	)

	s.setupMockFor409HttCode(originalURL, hash)

	body, err := json.Marshal(req)
	s.NoError(err)

	w := httptest.NewRecorder()
	s.handler.Routes().ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(body)))
	res := w.Result()
	defer res.Body.Close()

	s.Equal(http.StatusConflict, res.StatusCode)

	var resp dto.ShortURLResponse
	body, err = io.ReadAll(res.Body)
	err = json.Unmarshal(body, &resp)

	s.NoError(err)
	s.Equal(s.host+hash, resp.URL)
}

func (s *URLShortenerHandlerTestSuite) Test_When_URLExists_Then_ShortenPlainTextShouldReturn409HttpCodeAndURL() {
	var (
		originalURL = "https://hard2code.ru"
		hash        = helpers.GenerateHash()
	)

	s.setupMockFor409HttCode(originalURL, hash)

	w := httptest.NewRecorder()
	s.handler.Routes().ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(originalURL)))
	res := w.Result()
	defer res.Body.Close()

	s.Equal(http.StatusConflict, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	s.NoError(err)
	s.Equal(s.host+hash, string(body))
}

func (s *URLShortenerHandlerTestSuite) Test_ShortenBatch_Success() {
	var (
		dtos = []dto.BatchShortenURLRequest{
			{
				URL:           "https://google.com",
				CorrelationID: uuid.New().String(),
			},
			{
				URL:           "https://hard2code.ru",
				CorrelationID: uuid.New().String(),
			},
		}
		expectedURLs = []model.URLEntry{
			{
				ID:            uuid.New(),
				OriginalURL:   dtos[0].URL,
				CorrelationID: &dtos[0].CorrelationID,
				Hash:          helpers.GenerateHash(),
			},
			{
				ID:            uuid.New(),
				OriginalURL:   dtos[1].URL,
				CorrelationID: &dtos[1].CorrelationID,
				Hash:          helpers.GenerateHash(),
			},
		}
	)

	body, err := json.Marshal(dtos)
	s.NoError(err)

	s.repository.EXPECT().
		CreateBatch(gomock.Any(), gomock.AssignableToTypeOf([]*model.URLEntry{})).
		DoAndReturn(func(_ context.Context, m []*model.URLEntry) error {
			m[0].Hash = expectedURLs[0].Hash
			m[1].Hash = expectedURLs[1].Hash

			return nil
		}).
		Times(1)

	w := httptest.NewRecorder()
	s.handler.Routes().ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBuffer(body)))
	res := w.Result()
	defer res.Body.Close()

	s.Equal(http.StatusCreated, res.StatusCode)

	body, err = io.ReadAll(res.Body)
	s.NoError(err)

	var resp []dto.BatchShortenURLResponse
	err = json.Unmarshal(body, &resp)
	s.NoError(err)

	for i, r := range resp {
		s.Equal(s.host+expectedURLs[i].Hash, r.URL)
		s.Equal(*expectedURLs[i].CorrelationID, r.CorrelationID)
	}
}

func (s *URLShortenerHandlerTestSuite) Test_When_URLOrCorrelationIDExists_Then_ShortenBatchShouldReturn400HttpCode() {
	dtos := []dto.BatchShortenURLRequest{
		{
			URL:           "https://google.com",
			CorrelationID: uuid.New().String(),
		},
		{
			URL:           "https://hard2code.ru",
			CorrelationID: uuid.New().String(),
		},
	}

	body, err := json.Marshal(dtos)
	s.NoError(err)

	s.repository.EXPECT().
		CreateBatch(gomock.Any(), gomock.AssignableToTypeOf([]*model.URLEntry{})).
		Return(errs.ErrDuplicate).
		Times(1)

	w := httptest.NewRecorder()
	s.handler.Routes().ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBuffer(body)))
	res := w.Result()
	defer res.Body.Close()

	s.Equal(http.StatusBadRequest, res.StatusCode)

	body, err = io.ReadAll(res.Body)
	s.NoError(err)

	var resp helpers2.ErrorResponse
	err = json.Unmarshal(body, &resp)
	s.Equal(resp.ID, errs.ErrIncorrectURL.ID())
	s.Equal(resp.Message, errs.ErrIncorrectURL.Error())
}
