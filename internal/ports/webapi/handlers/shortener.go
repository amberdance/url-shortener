package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/amberdance/url-shortener/internal/app/service"
	"github.com/go-chi/chi/v5"
)

type URLShortenerHandler struct {
	host    string
	service *service.URLShortenerService
}

func NewURLShortenerHandler(srv *service.URLShortenerService, host string) *URLShortenerHandler {
	return &URLShortenerHandler{
		service: srv,
		host:    strings.TrimRight(host, "/") + "/",
	}
}

func (h *URLShortenerHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.post)
	r.Get("/{id:[a-zA-Z0-9]+}", h.get)
	return r
}

func (h *URLShortenerHandler) post(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil || len(body) == 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateShortURL(r.Context(), string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL := h.host + id
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func (h *URLShortenerHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	original, err := h.service.ResolveURL(r.Context(), id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", original)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
