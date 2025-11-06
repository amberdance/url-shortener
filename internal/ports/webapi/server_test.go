package webapi

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestServer() *Server {
	s := NewServer(":8080")
	return s
}

func TestHandlePost_Success(t *testing.T) {
	s := newTestServer()

	body := bytes.NewBufferString("https://practicum.yandex.ru/")
	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	s.handlePost(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")
	if contentType != "text/plain" {
		t.Errorf("expected Content-Type text/plain, got %s", contentType)
	}

	respBody, _ := io.ReadAll(res.Body)
	respStr := string(respBody)
	if !strings.HasPrefix(respStr, "http://localhost:8080/") {
		t.Errorf("unexpected short URL: %s", respStr)
	}
}

func TestHandlePost_BadRequest(t *testing.T) {
	s := newTestServer()

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

			s.handlePost(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != http.StatusBadRequest {
				t.Errorf("[%s] expected status 400, got %d", tt.name, res.StatusCode)
			}
		})
	}
}

// --- TEST: handleGet ---

func TestHandleGet_Success(t *testing.T) {
	s := newTestServer()
	s.storage["abc123"] = "https://practicum.yandex.ru/"

	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	w := httptest.NewRecorder()

	s.handleGet(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusTemporaryRedirect {
		t.Errorf("expected status %d, got %d", http.StatusTemporaryRedirect, res.StatusCode)
	}

	location := res.Header.Get("Location")
	if location != "https://practicum.yandex.ru/" {
		t.Errorf("expected redirect to original URL, got %s", location)
	}
}

func TestHandleGet_BadRequest(t *testing.T) {
	s := newTestServer()

	tests := []struct {
		name string
		path string
	}{
		{"empty path", "/"},
		{"not found id", "/unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			s.handleGet(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != http.StatusBadRequest {
				t.Errorf("[%s] expected status 400, got %d", tt.name, res.StatusCode)
			}
		})
	}
}
