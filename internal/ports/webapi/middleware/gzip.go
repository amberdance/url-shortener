package middleware

import (
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"strings"
)

func GzipDecompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") != "gzip" {
			next.ServeHTTP(w, r)
			return
		}

		gr, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "invalid gzip body", http.StatusBadRequest)
			return
		}

		r.Body = gr
		next.ServeHTTP(w, r)
		gr.Close()
	})
}

func GzipCompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)

		ct := rec.Header().Get("Content-Type")
		if !(strings.Contains(ct, "application/json") || strings.Contains(ct, "text/html")) {
			for k, v := range rec.Header() {
				w.Header()[k] = v
			}
			w.WriteHeader(rec.Code)
			w.Write(rec.Body.Bytes())
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length")

		for k, v := range rec.Header() {
			w.Header()[k] = v
		}

		w.WriteHeader(rec.Code)

		gz := gzip.NewWriter(w)
		defer gz.Close()

		_, _ = gz.Write(rec.Body.Bytes())
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	gz    *gzip.Writer
	wrote bool
}

func (w *gzipResponseWriter) WriteHeader(statusCode int) {
	if !w.wrote {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length")
		w.wrote = true
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if !w.wrote {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length")
		w.wrote = true
	}
	return w.gz.Write(b)
}
