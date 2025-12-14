package handlers

import (
	"net/http"

	"github.com/amberdance/url-shortener/internal/ports/webapi/middleware"
	"github.com/go-chi/chi/v5"
)

type HealthcheckHandler struct {
}

func NewHealthcheckHandler() *HealthcheckHandler {
	return &HealthcheckHandler{}
}

func (h *HealthcheckHandler) Routes() chi.Router {
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", middleware.ContentTypeJSONHeaderValue)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	return router
}
