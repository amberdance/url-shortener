package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/amberdance/url-shortener/internal/app/storage"
	"github.com/amberdance/url-shortener/internal/ports/webapi/helpers"
	"github.com/go-chi/chi/v5"
)

type UrlShortenerHandler struct {
	host    string
	storage storage.Storage
}

func NewUrlShortenerHandler(st storage.Storage, host string) *UrlShortenerHandler {
	return &UrlShortenerHandler{
		storage: st,
		host:    strings.TrimRight(host, "/") + "/",
	}
}

func (h *UrlShortenerHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.post)
	r.Get("/{id:[a-zA-Z0-9]+}", h.get)
	return r
}

func (h *UrlShortenerHandler) post(w http.ResponseWriter, r *http.Request) {
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
	if err := h.storage.Save(id, originalURL); err != nil {
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	shortURL := h.host + id
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(shortURL))
}

func (h *UrlShortenerHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	originalURL, ok := h.storage.Get(id)
	if !ok {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
