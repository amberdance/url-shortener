package helpers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/amberdance/url-shortener/internal/domain/errs"
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func HandleError(w http.ResponseWriter, err error) {
	var code int
	var errorID string

	switch e := err.(type) {
	case errs.NotFoundError:
		code, errorID = http.StatusNotFound, e.ID()
	case errs.InvalidArgumentError:
		code, errorID = http.StatusUnprocessableEntity, e.ID()
	case errs.ValidationError:
		code, errorID = http.StatusBadRequest, e.ID()
	case errs.InternalError:
		code, errorID = http.StatusInternalServerError, e.ID()
		log.Println("Internal error:", err.Error())
	case errs.UnauthorizedError:
		code, errorID = http.StatusUnauthorized, e.ID()
	case errs.DuplicateEntryError:
		code, errorID = http.StatusConflict, e.ID()
	default:
		code, errorID = http.StatusInternalServerError, "internal_error"
		log.Println("Unexpected error:", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		ID:      errorID,
		Message: err.Error(),
	})
}

func Validate(w http.ResponseWriter, v *validator.Validate, dto any) error {
	err := v.Struct(dto)
	if err != nil {
		HandleError(w, errs.ValidationError(err.Error()))
		return err
	}
	return nil
}

func MustValidate(w http.ResponseWriter, v *validator.Validate, dto any) {
	err := Validate(w, v, dto)
	if err != nil {
		panic(err)
	}
}
