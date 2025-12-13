package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/amberdance/url-shortener/internal/app/command"
	"github.com/amberdance/url-shortener/internal/app/usecase"
	"github.com/amberdance/url-shortener/internal/domain/contracts"
	"github.com/amberdance/url-shortener/internal/domain/errs"
	"github.com/amberdance/url-shortener/internal/ports/webapi/dto"
	"github.com/amberdance/url-shortener/internal/ports/webapi/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

const (
	writeRequestTimeout = 30 * time.Second
	readRequestTimeout  = 10 * time.Second
)

type URLShortenerHandler struct {
	baseURL   string
	usecases  usecase.URLUseCases
	validator *validator.Validate
	logger    contracts.Logger
}

func NewURLShortenerHandler(host string, uc usecase.URLUseCases, v *validator.Validate, l contracts.Logger) *URLShortenerHandler {
	return &URLShortenerHandler{host, uc, v, l}
}

func (h *URLShortenerHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.deprecatedPost)
	r.Get("/{hash:[a-zA-Z0-9]+}", h.get)
	r.Post("/api/shorten", h.shorten)
	r.Post("/api/shorten/batch", h.shortenBatch)
	return r
}

func (h *URLShortenerHandler) shorten(w http.ResponseWriter, r *http.Request) {
	var req dto.ShortURLRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	helpers.MustValidate(w, h.validator, req)

	ctx, cancel := context.WithTimeout(r.Context(), writeRequestTimeout)
	defer cancel()

	userID, _ := uuid.Parse(helpers.GetUserIDFromRequest(r))

	m, err := h.usecases.Create.Run(ctx, command.CreateURLEntryCommand{
		OriginalURL:   req.URL,
		CorrelationID: req.CorrelationID,
		UserID:        &userID,
	})

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}

		var conflictErr errs.DuplicateEntryError
		if errors.As(err, &conflictErr) {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(dto.ShortURLResponse{URL: helpers.FormatFullURL(h.baseURL, m.Hash)})
			return
		}

		helpers.HandleError(w, errs.ValidationError("Не удалось сформировать ссылку"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.ShortURLResponse{URL: helpers.FormatFullURL(h.baseURL, m.Hash)})
}

func (h *URLShortenerHandler) shortenBatch(w http.ResponseWriter, r *http.Request) {
	var reqDto, err = h.validateBatchRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		helpers.HandleError(w, err)
		return
	}

	cmd := command.CreateBatchURLEntryCommand{
		Entries: make([]command.CreateURLEntryCommand, 0, len(reqDto)),
	}

	userID, _ := uuid.Parse(helpers.GetUserIDFromRequest(r))
	for _, d := range reqDto {
		cmd.Entries = append(cmd.Entries, command.CreateURLEntryCommand{
			OriginalURL:   d.URL,
			CorrelationID: &d.CorrelationID,
			UserID:        &userID,
		})
	}

	ctx, cancel := context.WithTimeout(r.Context(), writeRequestTimeout)
	urls, err := h.usecases.CreateBatch.Run(ctx, cmd)
	defer cancel()

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}

		helpers.HandleError(w, errs.InvalidArgumentError("Не удалось создать записи"))
		return
	}

	res := make([]dto.BatchShortenURLResponse, 0, len(reqDto))
	for _, m := range urls {
		res = append(res, dto.BatchShortenURLResponse{
			CorrelationID: *m.CorrelationID,
			URL:           helpers.FormatFullURL(h.baseURL, m.Hash),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(res)
}

func (h *URLShortenerHandler) validateBatchRequest(r *http.Request) ([]dto.BatchShortenURLRequest, error) {
	var reqItems []dto.BatchShortenURLRequest
	err := json.NewDecoder(r.Body).Decode(&reqItems)
	if err != nil {
		return nil, errs.ValidationError(err.Error())
	}

	if len(reqItems) == 0 {
		return nil, errs.ValidationError("Не передано ни одного url")
	}

	return reqItems, nil
}

func (h *URLShortenerHandler) get(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	if hash == "" {
		helpers.HandleError(w, errs.ValidationError("Не передана ссылка"))
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
		helpers.HandleError(w, errs.NotFoundError("Не найден ресурс"))
		return
	}

	w.Header().Set("Location", m.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *URLShortenerHandler) deprecatedPost(w http.ResponseWriter, r *http.Request) {
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

	var (
		requestID = uuid.New().String()
		userID, _ = uuid.Parse(helpers.GetUserIDFromRequest(r))
	)

	m, err := h.usecases.Create.Run(r.Context(), command.CreateURLEntryCommand{
		OriginalURL:   original,
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

		h.logger.Error(err.Error())
		helpers.HandleError(w, errs.ValidationError("Не удалось сформировать ссылку"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(helpers.FormatFullURL(h.baseURL, m.Hash)))
}
