package handlers

import (
	"net/http"

	"github.com/amberdance/url-shortener/internal/infrastructure/storage"
	"github.com/go-chi/chi/v5"
)

type PingHandler struct {
	storage *storage.PostgresStorage
}

func NewPingHandler(s *storage.PostgresStorage) *PingHandler {
	return &PingHandler{storage: s}
}

func (h *PingHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if err := h.storage.Ping(r.Context()); err != nil {
			http.Error(w, "db not available", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	return r
}
