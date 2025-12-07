package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func gzipBody(t *testing.T, data string) *bytes.Buffer {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write([]byte(data))
	assert.NoError(t, err)
	gz.Close()
	return &buf
}

func TestGzipDecompressMiddleware_ValidGzip(t *testing.T) {
	body := gzipBody(t, `{"test":"ok"}`)

	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	var received string
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		received = string(data)
	})

	w := httptest.NewRecorder()
	GzipDecompressMiddleware(h).ServeHTTP(w, req)

	assert.Equal(t, `{"test":"ok"}`, received)
}

func TestGzipDecompressMiddleware_InvalidGzip(t *testing.T) {
	req := httptest.NewRequest("POST", "/", bytes.NewBufferString("not gzip"))
	req.Header.Set("Content-Encoding", "gzip")

	w := httptest.NewRecorder()
	GzipDecompressMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestGzipDecompressMiddleware_NoGzip(t *testing.T) {
	req := httptest.NewRequest("POST", "/", bytes.NewBufferString("plain text"))

	var received string
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		received = string(data)
	})

	w := httptest.NewRecorder()
	GzipDecompressMiddleware(h).ServeHTTP(w, req)

	// Result не вызывается → ok
	assert.Equal(t, "plain text", received)
}

func TestGzipCompressMiddleware_CompressJSON(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"hello": "world"}`))
	})

	w := httptest.NewRecorder()
	GzipCompressMiddleware(h).ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"))

	gz, err := gzip.NewReader(res.Body)
	assert.NoError(t, err)
	defer gz.Close()

	data, _ := io.ReadAll(gz)
	assert.Equal(t, `{"hello": "world"}`, string(data))
}

func TestGzipCompressMiddleware_SkipUnsupportedTypes(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write([]byte("BIN"))
	})

	w := httptest.NewRecorder()
	GzipCompressMiddleware(h).ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.NotEqual(t, "gzip", res.Header.Get("Content-Encoding"))

	data, _ := io.ReadAll(res.Body)
	assert.Equal(t, "BIN", string(data))
}

func TestGzipCompressMiddleware_NoAcceptEncoding(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<h1>x</h1>"))
	})

	w := httptest.NewRecorder()
	GzipCompressMiddleware(h).ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Empty(t, res.Header.Get("Content-Encoding"))

	data, _ := io.ReadAll(res.Body)
	assert.Equal(t, "<h1>x</h1>", string(data))
}
