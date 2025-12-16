package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/amberdance/url-shortener/internal/app/command"
	"github.com/amberdance/url-shortener/internal/app/usecase"
	"github.com/amberdance/url-shortener/internal/domain/errs"
	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/amberdance/url-shortener/internal/ports/webapi/dto"
	"github.com/amberdance/url-shortener/internal/ports/webapi/helpers"
	"github.com/amberdance/url-shortener/internal/ports/webapi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const (
	writeRequestTimeout = 30 * time.Second
	readRequestTimeout  = 10 * time.Second
)

type URLShortenerHandler struct {
	baseURL  string
	usecases usecase.URLUseCases
}

func NewURLShortenerHandler(baseURL string, uc usecase.URLUseCases) *URLShortenerHandler {
	return &URLShortenerHandler{
		baseURL:  baseURL,
		usecases: uc,
	}
}

func (h *URLShortenerHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/{hash:[a-zA-Z0-9]+}", h.get)
	r.Post("/", h.plainTextShorten)
	r.Route("/api/shorten", func(r chi.Router) {
		r.Post("/", h.shorten)
		r.Post("/batch", h.shortenBatch)
	})

	return r
}

func (h *URLShortenerHandler) shorten(w http.ResponseWriter, r *http.Request) {
	var req dto.ShortURLRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	validatedURL, err := helpers.ValidateURL(req.URL)
	if err != nil {
		helpers.HandleError(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), writeRequestTimeout)
	defer cancel()

	userID, _ := uuid.Parse(helpers.GetUserIDFromRequest(r))
	m, err := h.usecases.Create.Run(ctx, command.CreateURLEntryCommand{
		OriginalURL:   validatedURL,
		CorrelationID: req.CorrelationID,
		UserID:        &userID,
	})

	w.Header().Set("Content-Type", middleware.ContentTypeJSONHeaderValue)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}

		var conflictErr errs.DuplicateEntryError
		if errors.As(err, &conflictErr) {
			w.WriteHeader(http.StatusConflict)
			h.writeShortenDto(w, m)
			return
		}

		helpers.HandleError(w, errs.ErrIncorrectURL)
		return
	}

	w.WriteHeader(http.StatusCreated)
	h.writeShortenDto(w, m)
}

func (h *URLShortenerHandler) writeShortenDto(w http.ResponseWriter, m *model.URLEntry) {
	json.NewEncoder(w).Encode(dto.ShortURLResponse{URL: helpers.FormatFullURL(h.baseURL, m.Hash)})
}

func (h *URLShortenerHandler) shortenBatch(w http.ResponseWriter, r *http.Request) {
	var reqDto, err = h.validateBatchRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		helpers.HandleError(w, err)
		return
	}

	userID, _ := uuid.Parse(helpers.GetUserIDFromRequest(r))
	cmd := command.CreateBatchURLEntryCommand{
		Commands: make([]command.CreateURLEntryCommand, 0, len(reqDto)),
	}

	for _, d := range reqDto {
		cmd.Commands = append(cmd.Commands, command.CreateURLEntryCommand{
			OriginalURL:   d.URL,
			CorrelationID: &d.CorrelationID,
			UserID:        &userID,
		})
	}

	ctx, cancel := context.WithTimeout(r.Context(), writeRequestTimeout)
	defer cancel()

	urls, err := h.usecases.CreateBatch.Run(ctx, cmd)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}

		helpers.HandleError(w, errs.ErrIncorrectURL)
		return
	}

	res := make([]dto.BatchShortenURLResponse, 0, len(reqDto))
	for _, m := range urls {
		res = append(res, dto.BatchShortenURLResponse{
			CorrelationID: *m.CorrelationID,
			URL:           helpers.FormatFullURL(h.baseURL, m.Hash),
		})
	}

	w.Header().Set("Content-Type", middleware.ContentTypeJSONHeaderValue)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *URLShortenerHandler) validateBatchRequest(r *http.Request) ([]dto.BatchShortenURLRequest, error) {
	var reqItems []dto.BatchShortenURLRequest
	err := json.NewDecoder(r.Body).Decode(&reqItems)
	if err != nil {
		return nil, errs.ValidationError(err.Error())
	}

	if len(reqItems) == 0 {
		return nil, errs.ErrEmptyURLSet
	}

	return reqItems, nil
}

func (h *URLShortenerHandler) get(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	if hash == "" {
		helpers.HandleError(w, errs.ErrIncorrectURL)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), readRequestTimeout)
	defer cancel()

	m, err := h.usecases.GetByURL.Run(ctx, command.GetURLByHashCommand{
		Hash: hash,
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}
		helpers.HandleError(w, errs.ErrNotFound)
		return
	}

	if m.DeletedAt != nil {
		w.WriteHeader(http.StatusGone)
		return
	}

	w.Header().Set("Location", m.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *URLShortenerHandler) plainTextShorten(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil || len(body) == 0 {
		helpers.HandleError(w, errs.ErrIncorrectURL)
		return
	}

	validatedURL, err := helpers.ValidateURL(string(body))
	if err != nil {
		helpers.HandleError(w, err)
		return
	}

	var (
		requestID = uuid.New().String()
		userID, _ = uuid.Parse(helpers.GetUserIDFromRequest(r))
	)

	m, err := h.usecases.Create.Run(r.Context(), command.CreateURLEntryCommand{
		OriginalURL:   validatedURL,
		CorrelationID: &requestID,
		UserID:        &userID,
	})

	w.Header().Set("Content-Type", "text/plain")

	if err != nil {
		var conflictErr errs.DuplicateEntryError
		if errors.As(err, &conflictErr) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(helpers.FormatFullURL(h.baseURL, m.Hash)))
			return
		}

		helpers.HandleError(w, errs.ErrIncorrectURL)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(helpers.FormatFullURL(h.baseURL, m.Hash)))
}
