package middleware

import "net/http"

const ContentTypeJSONHeaderValue = "application/json; charset=utf-8"

func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", ContentTypeJSONHeaderValue)
		next.ServeHTTP(w, r)
	})
}
