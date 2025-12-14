package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amberdance/url-shortener/internal/app"
	"github.com/amberdance/url-shortener/internal/app/usecase/url"
	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/amberdance/url-shortener/internal/infrastructure/helpers"
	"github.com/amberdance/url-shortener/internal/mocks"
	"github.com/amberdance/url-shortener/internal/ports/webapi/dto"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

const (
	userUrlsEndpoint = "/urls"
)

type UserTestSuite struct {
	suite.Suite
	ctrl *gomock.Controller
	ctx  context.Context

	repository *mocks.MockURLRepository
	host       string
	handler    *UserHandler
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (s *UserTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.ctrl = gomock.NewController(s.T())
	s.repository = mocks.NewMockURLRepository(s.ctrl)
	s.host = "http://localhost:8080/"
	s.handler = NewUserHandler(s.host, url.NewGetURLsByUserIDUseCase(s.repository), url.NewDeleteUserURLsBatchUseCase(s.repository, mocks.NewMockLogger(s.ctrl)))
}

func (s *UserTestSuite) TearDownSuite() {
	s.ctrl.Finish()
}

func (s *UserTestSuite) Test_GetURLS_When_UserHasUrls_Then_URLsReturned() {
	var (
		id, ctx  = generateUUIDWithContext(s.ctx)
		urls     = s.seedUrls(&id, 10)
		recorder = httptest.NewRecorder()
	)

	s.repository.EXPECT().
		FindAllByUserID(gomock.Any(), id).
		Return(urls, nil).
		Times(1)

	req := httptest.NewRequest(http.MethodGet, userUrlsEndpoint, nil).WithContext(ctx)
	s.handler.Routes().ServeHTTP(recorder, req)
	res := recorder.Result()
	defer res.Body.Close()

	response, _ := io.ReadAll(res.Body)

	var dtos []dto.UserURLsResponse
	err := json.Unmarshal(response, &dtos)
	s.NoError(err)

	s.Equal(http.StatusOK, res.StatusCode)
	s.Equal(len(urls), len(dtos))

	for i := range urls {
		s.Equal(urls[i].OriginalURL, dtos[i].OriginalURL)
	}
}

func (s *UserTestSuite) Test_GetURLs_When_UserDoesNotHasUrls_Then_204Returned() {
	var (
		id, ctx  = generateUUIDWithContext(s.ctx)
		recorder = httptest.NewRecorder()
	)

	s.repository.EXPECT().
		FindAllByUserID(gomock.Any(), id).
		Return(nil, nil).
		Times(1)

	req := httptest.NewRequest(http.MethodGet, userUrlsEndpoint, nil).WithContext(ctx)
	s.handler.Routes().ServeHTTP(recorder, req)
	res := recorder.Result()
	defer res.Body.Close()

	s.Equal(http.StatusNoContent, res.StatusCode)
}

func (s *UserTestSuite) Test_GetURLS_When_HeaderNotPresent_Then_401Returned() {
	var (
		w   = httptest.NewRecorder()
		req = httptest.NewRequest(http.MethodGet, userUrlsEndpoint, nil)
	)

	s.handler.Routes().ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()

	s.Equal(http.StatusUnauthorized, res.StatusCode)
}

func (s *UserTestSuite) Test_GetURLS_When_InvalidSignatureProvided_Then_401Returned() {
	var (
		badCookie = &http.Cookie{
			Name:  "user_id",
			Value: "broken.token.value",
			Path:  "/",
		}
		recorder = httptest.NewRecorder()
	)

	req := httptest.NewRequest(http.MethodGet, userUrlsEndpoint, nil)
	req.AddCookie(badCookie)

	s.handler.Routes().ServeHTTP(recorder, req)
	res := recorder.Result()
	defer res.Body.Close()

	s.Equal(http.StatusUnauthorized, res.StatusCode)
}

func (s *UserTestSuite) Test_DeleteBatch_When_ValidRequest_Then_202Returned() {
	var (
		id, ctx  = generateUUIDWithContext(s.ctx)
		recorder = httptest.NewRecorder()
		hashes   = []string{"hash1", "hash2", "hash3", "hash4"}
		body, _  = json.Marshal(hashes)
	)

	s.repository.EXPECT().
		DeleteByUserIDAndHashes(gomock.Any(), id, hashes).
		AnyTimes()

	req := httptest.NewRequest(http.MethodDelete, userUrlsEndpoint, bytes.NewBuffer(body)).WithContext(ctx)
	s.handler.Routes().ServeHTTP(recorder, req)
	res := recorder.Result()
	defer res.Body.Close()

	s.Equal(http.StatusAccepted, res.StatusCode)
}

func (s *UserTestSuite) Test_DeleteBatch_When_InvalidBody_Then_422Returned() {
	var (
		_, ctx   = generateUUIDWithContext(s.ctx)
		recorder = httptest.NewRecorder()
		body     = bytes.NewBufferString("invalid json")
	)

	req := httptest.NewRequest(http.MethodDelete, userUrlsEndpoint, body).WithContext(ctx)
	s.handler.Routes().ServeHTTP(recorder, req)
	s.Equal(http.StatusUnprocessableEntity, recorder.Code)
}

func (s *UserTestSuite) Test_DeleteBatch_When_Unauthorized_Then_401Returned() {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, userUrlsEndpoint, nil)

	s.handler.Routes().ServeHTTP(recorder, req)
	s.Equal(http.StatusUnauthorized, recorder.Code)
}

func generateUUIDWithContext(ctx context.Context) (uuid.UUID, context.Context) {
	id := uuid.New()
	c := context.WithValue(ctx, app.UserCtxKey, id.String())

	return id, c
}

func (s *UserTestSuite) seedUrls(userID *uuid.UUID, n int) []*model.URLEntry {
	urls := make([]*model.URLEntry, 0, n)

	for i := 0; i < n; i++ {
		m, _ := model.NewURLEntry(fmt.Sprintf("https://original-%d.ru", i), helpers.GenerateHash(), nil, userID)
		urls = append(urls, m)
	}

	return urls
}
