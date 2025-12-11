package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/amberdance/url-shortener/internal/app/command"
	usecase "github.com/amberdance/url-shortener/internal/app/usecase/url"
	"github.com/amberdance/url-shortener/internal/domain/contracts"
	"github.com/amberdance/url-shortener/internal/domain/errs"
	"github.com/amberdance/url-shortener/internal/ports/webapi/dto"
	"github.com/amberdance/url-shortener/internal/ports/webapi/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	baseURL                string
	getURLsByUserIDUseCase usecase.GetURLsByUserIDUseCase
}

func NewUserHandler(u string, uc usecase.GetURLsByUserIDUseCase) *UserHandler {
	return &UserHandler{baseURL: u, getURLsByUserIDUseCase: uc}
}

func (h *UserHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/api/user/urls", h.getAll)
	return r
}

func (h *UserHandler) getAll(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		helpers.HandleError(w, errs.UnauthorizedError("Unauthorized"))
		return
	}

	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		helpers.HandleError(w, errs.InvalidArgumentError("Incorrect user ID"))
		return
	}

	urls, err := h.getURLsByUserIDUseCase.Run(r.Context(), command.GetUrlsByUserIDCommand{UserID: parsedUUID})
	if err != nil {
		helpers.HandleError(w, errs.NotFoundError("Urls not found"))
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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(result)
}

func getUserID(r *http.Request) string {
	v := r.Context().Value(contracts.UserCtxKey)
	if v == nil {
		return ""
	}
	return v.(string)
}
