package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/amberdance/url-shortener/internal/app/command"
	"github.com/amberdance/url-shortener/internal/app/usecase"
	"github.com/amberdance/url-shortener/internal/domain/errs"
	"github.com/amberdance/url-shortener/internal/domain/shared"
	"github.com/amberdance/url-shortener/internal/ports/webapi/dto"
	"github.com/amberdance/url-shortener/internal/ports/webapi/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

type URLShortenerHandler struct {
	baseURL   string
	usecases  usecase.URLUseCases
	validator *validator.Validate
	logger    shared.Logger
}

func NewURLShortenerHandler(host string, uc usecase.URLUseCases, v *validator.Validate, l shared.Logger) *URLShortenerHandler {
	return &URLShortenerHandler{host, uc, v, l}
}

func (h *URLShortenerHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Post("/", h.post)
	r.Get("/{hash:[a-zA-Z0-9]+}", h.get)
	r.Post("/api/shorten", h.shortenJSON)
	return r
}

func (h *URLShortenerHandler) post(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil || len(body) == 0 {
		helpers.HandleError(w, errs.ValidationError("Не передан URL"))
		return
	}

	original := strings.TrimSpace(string(body))
	if original == "" {
		helpers.HandleError(w, errs.ValidationError("Не передан URL"))
		return
	}

	model, err := h.usecases.Create.Run(r.Context(), command.CreateURLEntryCommand{
		OriginalURL: original,
	})
	if err != nil {
		h.logger.Error(err.Error())
		helpers.HandleError(w, errs.ValidationError("Не удалось сформировать ссылку"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(h.baseURL + model.Hash))
}

func (h *URLShortenerHandler) get(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	if hash == "" {
		helpers.HandleError(w, errs.ValidationError("Не передана ссылка"))
		return
	}

	model, err := h.usecases.GetByURL.Run(r.Context(), command.GetURLByHashCommand{
		Hash: hash,
	})
	if err != nil {
		helpers.HandleError(w, errs.NotFoundError("Не найден ресурс"))
		return
	}

	w.Header().Set("Location", model.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *URLShortenerHandler) shortenJSON(w http.ResponseWriter, r *http.Request) {
	var req dto.ShortURLRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	helpers.MustValidate(w, h.validator, req)

	model, err := h.usecases.Create.Run(r.Context(), command.CreateURLEntryCommand{
		OriginalURL: req.URL,
	})
	if err != nil {
		helpers.HandleError(w, errs.ValidationError("Не удалось сформировать ссылку"))
		return
	}

	shortURL := h.baseURL + model.Hash

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(dto.ShortURLResponse{Result: shortURL})
}
