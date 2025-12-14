package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const RequestIDCtxtKey requestIDCtxtKey = "request_id"

type requestIDCtxtKey string

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.NewString()
		context.WithValue(r.Context(), RequestIDCtxtKey, id)
		w.Header().Set("X-Request-Id", id)
		next.ServeHTTP(w, r)
	})
}

func GetRequestID(ctx context.Context) uuid.UUID {
	raw := ctx.Value(RequestIDCtxtKey).(string)
	if raw == "" {
		return uuid.New()
	}

	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.New()
	}

	return id
}
