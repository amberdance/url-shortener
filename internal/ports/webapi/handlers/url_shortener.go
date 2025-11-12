package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/amberdance/url-shortener/internal/app/command"
	"github.com/amberdance/url-shortener/internal/app/usecase"
	"github.com/go-chi/chi/v5"
)

type URLShortenerHandler struct {
	baseURL  string
	usecases usecase.URLUseCases
}

func NewURLShortenerHandler(host string, uc usecase.URLUseCases) *URLShortenerHandler {
	return &URLShortenerHandler{
		baseURL:  host,
		usecases: uc,
	}
}

func (h *URLShortenerHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.post)
	r.Get("/{hash:[a-zA-Z0-9]+}", h.get)
	return r
}

func (h *URLShortenerHandler) post(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil || len(body) == 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	original := strings.TrimSpace(string(body))
	if original == "" {
		http.Error(w, "Empty URL", http.StatusBadRequest)
		return
	}

	model, err := h.usecases.Create.Run(r.Context(), command.CreateURLEntryCommand{
		OriginalURL: original,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(h.baseURL + model.Hash))
}

func (h *URLShortenerHandler) get(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	if hash == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	model, err := h.usecases.GetByURL.Run(r.Context(), command.GetURLByHashCommand{
		Hash: hash,
	})
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", model.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
