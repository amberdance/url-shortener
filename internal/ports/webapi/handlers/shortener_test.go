package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/amberdance/url-shortener/internal/infrastructure/storage"
)

func setupTest() *URLShortenerHandler {
	st := storage.NewInMemoryStorage()
	return NewURLShortenerHandler(st, "http://localhost:8080/")
}

func TestPost_Success(t *testing.T) {
	h := setupTest()
	router := h.Routes()

	body := bytes.NewBufferString("https://hard2code.ru")
	req := httptest.NewRequest(http.MethodPost, "/", body)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", res.StatusCode)
	}

	respBody, _ := io.ReadAll(res.Body)
	if !strings.HasPrefix(string(respBody), "http://localhost:8080/") {
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

	_ = h.storage.Save("abc123", "https://hard2code.ru")

	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusTemporaryRedirect {
		t.Fatalf("expected 307, got %d", res.StatusCode)
	}

	if res.Header.Get("Location") != "https://yandex.ru" {
		t.Errorf("expected redirect to https://yandex.ru, got %s", res.Header.Get("Location"))
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

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", res.StatusCode)
	}
}
