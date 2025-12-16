package handlers

import (
	"net/http"

	"github.com/amberdance/url-shortener/internal/domain/ports"
	"github.com/go-chi/chi/v5"
)

type PingHandler struct {
	pinger ports.Pinger
}

func NewPingHandler(s ports.Pinger) *PingHandler {
	return &PingHandler{pinger: s}
}

func (h *PingHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if err := h.pinger.Ping(r.Context()); err != nil {
			http.Error(w, "db not available", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	return r
}
