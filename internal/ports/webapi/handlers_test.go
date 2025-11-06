package webapi

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/amberdance/url-shortener/internal/infrastructure/storage"
)

func setupTest() *Handler {
	st := storage.NewInMemoryStorage()
	return NewHandler(st, ":8080")
}

func TestHandlePost_Success(t *testing.T) {
	h := setupTest()
	body := bytes.NewBufferString("https://practicum.yandex.ru/")
	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	h.handlePost(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, res.StatusCode)
	}

	if res.Header.Get("Content-Type") != "text/plain" {
		t.Errorf("expected Content-Type text/plain")
	}

	respBody, _ := io.ReadAll(res.Body)
	respStr := string(respBody)
	if !strings.HasPrefix(respStr, "http://localhost:8080/") {
		t.Errorf("unexpected short URL: %s", respStr)
	}
}

func TestHandlePost_BadRequest(t *testing.T) {
	h := setupTest()

	tests := []struct {
		name        string
		body        io.Reader
		contentType string
	}{
		{"empty body", bytes.NewBuffer(nil), "text/plain"},
		{"wrong content-type", bytes.NewBufferString("https://example.com"), "application/json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", tt.body)
			req.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()

			h.handlePost(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != http.StatusBadRequest {
				t.Errorf("[%s] expected 400, got %d", tt.name, res.StatusCode)
			}
		})
	}
}

func TestHandleGet_Success(t *testing.T) {
	h := setupTest()
	_ = h.storage.Save("abc123", "https://practicum.yandex.ru/")

	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	w := httptest.NewRecorder()

	h.handleGet(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusTemporaryRedirect {
		t.Errorf("expected %d, got %d", http.StatusTemporaryRedirect, res.StatusCode)
	}

	if res.Header.Get("Location") != "https://practicum.yandex.ru/" {
		t.Errorf("unexpected redirect: %s", res.Header.Get("Location"))
	}
}

func TestHandleGet_BadRequest(t *testing.T) {
	h := setupTest()

	tests := []struct {
		name string
		path string
	}{
		{"empty path", "/"},
		{"unknown id", "/notfound"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			h.handleGet(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != http.StatusBadRequest {
				t.Errorf("[%s] expected 400, got %d", tt.name, res.StatusCode)
			}
		})
	}
}
