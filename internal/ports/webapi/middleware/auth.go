package middleware

import (
	"context"
	"net/http"

	"github.com/amberdance/url-shortener/internal/domain/contracts"
	"github.com/amberdance/url-shortener/internal/infrastructure/auth"
	"github.com/google/uuid"
)

const CookieUserIDKey = "user_id"

func AuthMiddleware(auth *auth.CookieAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie(CookieUserIDKey)

			if err != nil {
				userID := uuid.New().String()
				token := auth.Sign(userID)

				http.SetCookie(w, &http.Cookie{
					Name:     CookieUserIDKey,
					Value:    token,
					Path:     "/",
					HttpOnly: true,
				})

				ctx := context.WithValue(r.Context(), contracts.UserCtxKey, userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			userID, err := auth.Verify(c.Value)
			if err != nil || userID == "" {
				ctx := context.WithValue(r.Context(), contracts.UserCtxKey, "")
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			ctx := context.WithValue(r.Context(), contracts.UserCtxKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
