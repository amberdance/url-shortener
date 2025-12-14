package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/amberdance/url-shortener/internal/app/command"
	usecase "github.com/amberdance/url-shortener/internal/app/usecase/url"
	"github.com/amberdance/url-shortener/internal/domain/errs"
	"github.com/amberdance/url-shortener/internal/ports/webapi/dto"
	"github.com/amberdance/url-shortener/internal/ports/webapi/helpers"
	"github.com/amberdance/url-shortener/internal/ports/webapi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	baseURL                    string
	getURLsByUserIDUseCase     usecase.GetURLsByUserIDUseCase
	deleteUserURLsBatchUseCase usecase.DeleteUserURLsBatchUseCase
}

func NewUserHandler(u string, uc1 usecase.GetURLsByUserIDUseCase, uc2 usecase.DeleteUserURLsBatchUseCase) *UserHandler {
	return &UserHandler{baseURL: u, getURLsByUserIDUseCase: uc1, deleteUserURLsBatchUseCase: uc2}
}

func (h *UserHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Route("/urls", func(r chi.Router) {
		r.Get("/", h.getAll)
		r.Delete("/", h.deleteBatch)
	})

	return r
}

func (h *UserHandler) getAll(w http.ResponseWriter, r *http.Request) {
	userID := helpers.GetUserIDFromRequest(r)
	if userID == "" {
		helpers.HandleError(w, errs.ErrUnauthorized)
		return
	}

	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		helpers.HandleError(w, errs.ErrInvalidUserID)
		return
	}

	urls, err := h.getURLsByUserIDUseCase.Run(r.Context(), command.GetUrlsByUserIDCommand{UserID: parsedUUID})
	if err != nil {
		helpers.HandleError(w, errs.ErrNotFound)
		return
	}

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	result := make([]dto.UserURLsResponse, len(urls))
	for i, m := range urls {
		result[i] = dto.UserURLsResponse{
			ShortURL:    helpers.FormatFullURL(h.baseURL, m.Hash),
			OriginalURL: m.OriginalURL,
		}
	}

	w.Header().Set("Content-Type", middleware.ContentTypeJSONHeaderValue)
	json.NewEncoder(w).Encode(result)
}

func (h *UserHandler) deleteBatch(w http.ResponseWriter, r *http.Request) {
	userId := helpers.GetUserIDFromRequest(r)
	if userId == "" {
		helpers.HandleError(w, errs.ErrUnauthorized)
		return
	}

	parsedUUID, err := uuid.Parse(userId)
	if err != nil {
		helpers.HandleError(w, errs.ErrInvalidUserID)
		return
	}

	var hashes []string
	err = json.NewDecoder(r.Body).Decode(&hashes)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		helpers.HandleError(w, errs.ErrEmptyHashSet)
		return
	}

	_ = h.deleteUserURLsBatchUseCase.RunAsync(command.DeleteUserURLSCommand{
		UserID: parsedUUID,
		Hashes: hashes,
	})

	w.Header().Set("Content-Type", middleware.ContentTypeJSONHeaderValue)
	w.WriteHeader(http.StatusAccepted)
}
