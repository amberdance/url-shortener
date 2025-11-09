package handlers

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/amberdance/url-shortener/internal/app/command"
	"github.com/amberdance/url-shortener/internal/app/usecase"
	"github.com/amberdance/url-shortener/internal/app/usecase/url"
	infr "github.com/amberdance/url-shortener/internal/infrastructure/repository/url"
)

const testHost string = "http://127.0.0.1:9999"

func setupTest() *URLShortenerHandler {
	repo := infr.NewInMemoryRepository()
	useCases := usecase.URLUseCases{
		Create:   url.NewCreateUrlUseCase(repo),
		GetByUrl: url.NewGetByHashUseCase(repo),
	}
	return NewURLShortenerHandler(
		testHost,
		useCases,
	)
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
