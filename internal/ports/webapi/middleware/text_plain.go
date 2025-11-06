package middleware

import "net/http"

func TextPlainHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if r.Header.Get("Content-Type") != "text/plain" {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		}

		next.ServeHTTP(w, r)
	})
}
