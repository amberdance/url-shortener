package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/amberdance/url-shortener/internal/app"
	"github.com/amberdance/url-shortener/internal/ports/webapi/helpers"
	"github.com/go-chi/chi/v5"
)

type URLShortenerHandler struct {
	host    string
	storage app.Storage
}

func NewURLShortenerHandler(st app.Storage, host string) *URLShortenerHandler {
	return &URLShortenerHandler{
		storage: st,
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

	originalURL := strings.TrimSpace(string(body))
	if originalURL == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	id := helpers.GenerateShortID()
	if err := h.storage.Save(r.Context(), id, originalURL); err != nil {
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	shortURL := h.host + id
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func (h *URLShortenerHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	originalURL, ok := h.storage.Get(r.Context(), id)
	if !ok {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
