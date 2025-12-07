package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/amberdance/url-shortener/internal/app/command"
	"github.com/amberdance/url-shortener/internal/app/usecase"
	"github.com/amberdance/url-shortener/internal/app/usecase/url"
	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/amberdance/url-shortener/internal/domain/repository"
	"github.com/amberdance/url-shortener/internal/domain/shared"
	infr "github.com/amberdance/url-shortener/internal/infrastructure/repository/url"
	"github.com/amberdance/url-shortener/internal/ports/webapi/dto"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

const testHost string = "http://127.0.0.1:9999/"

type MockLogger struct{}

func (m MockLogger) Debug(_ string, _ ...any) {}
func (m MockLogger) Info(_ string, _ ...any)  {}
func (m MockLogger) Error(_ string, _ ...any) {}
func (m MockLogger) Close() error             { return nil }

var repo repository.URLRepository

func setupTest() *URLShortenerHandler {
	var log shared.Logger = MockLogger{}

	repo = infr.NewInMemoryURLRepository()
	useCases := usecase.URLUseCases{
		Create:      url.NewCreateURLUseCase(repo),
		CreateBatch: url.NewBatchCreateURLUseCase(repo),
		GetByURL:    url.NewGetByHashUseCase(repo),
	}
	return NewURLShortenerHandler(testHost, useCases, validator.New(), log)
}

func TestPost_Success(t *testing.T) {
	h := setupTest()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://hard2code.ru"))
	w := httptest.NewRecorder()

	h.Routes().ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", res.StatusCode)
	}

	respBody, _ := io.ReadAll(res.Body)
	if !strings.HasPrefix(string(respBody), testHost) {
		t.Errorf("unexpected response: %s", respBody)
	}
}

func TestPost_BadRequest(t *testing.T) {
	h := setupTest()
	router := h.Routes()

	tests := []struct {
		name string
		body io.Reader
	}{
		{"empty body", bytes.NewBuffer(nil)},
		{"spaces only", bytes.NewBufferString("   ")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", tt.body)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != http.StatusBadRequest {
				t.Errorf("[%s] expected 400, got %d", tt.name, res.StatusCode)
			}
		})
	}
}

func TestGet_Success(t *testing.T) {
	h := setupTest()
	router := h.Routes()

	ctx := context.Background()
	entry, err := h.usecases.Create.Run(ctx, command.CreateURLEntryCommand{OriginalURL: "https://hard2code.ru"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/"+entry.Hash, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusTemporaryRedirect {
		t.Fatalf("expected 307, got %d", res.StatusCode)
	}

	if res.Header.Get("Location") != "https://hard2code.ru" {
		t.Errorf("expected redirect to https://hard2code.ru, got %s", res.Header.Get("Location"))
	}
}

func TestGet_NotFound(t *testing.T) {
	h := setupTest()
	router := h.Routes()

	req := httptest.NewRequest(http.MethodGet, "/notfound", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", res.StatusCode)
	}
}

func TestShorten_Success(t *testing.T) {
	h := setupTest()
	router := h.Routes()

	body := `{"url":"https://hard2code.ru"}`
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var resp dto.ShortURLResponse
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
}

func TestShortenJSON_BadRequest(t *testing.T) {
	h := setupTest()
	router := h.Routes()

	tests := []struct {
		name string
		body string
	}{
		{"missing field", `{"u":"wrong"}`},
		{"null field", `{"url":null}`},
		{"empty string", `{"url":""}`},
		{"spaces", `{"url":"   "}`},
		{"invalid json", `{invalid`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		})
	}
}

func TestShortenBatch_Success(t *testing.T) {
	h := setupTest()
	router := h.Routes()

	body := `[
			{
				"original_url": "https://google.com",
				"correlation_id": "11111111-1111-1111-1111-111111111111"
			},
			{
				"original_url": "https://hard2code.ru",
				"correlation_id": "22222222-2222-2222-2222-222222222222"
			}
		]`

	req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var resp []dto.BatchShortenURLResponse
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)

	assert.Len(t, resp, 2)

	var dtos []dto.BatchShortenURLRequest
	json.NewDecoder(res.Body).Decode(&dtos)

	for i := range dtos {
		assert.Equal(t, dtos[i], resp[0].CorrelationID)
		assert.NotEmpty(t, resp[0].URL)
	}
}

func TestShortenBatch_422Error(t *testing.T) {
	h := setupTest()
	router := h.Routes()

	body := `[
			
		]`

	req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
}
func TestShorten_409Error(t *testing.T) {
	h := setupTest()
	router := h.Routes()

	existing := &model.URL{
		OriginalURL: "https://hard2code.ru",
		Hash:        "hash",
	}
	err := repo.Create(context.Background(), existing)
	assert.NoError(t, err)

	body := `{"url":"https://hard2code.ru"}`
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusConflict, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var resp dto.ShortURLResponse
	err = json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)

	assert.Equal(t, h.baseURL+existing.Hash, resp.URL)
}
