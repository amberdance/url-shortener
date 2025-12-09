package model

import (
	"strings"
	"time"

	"github.com/amberdance/url-shortener/internal/domain/errs"
	"github.com/google/uuid"
)

type URL struct {
	ID            uuid.UUID
	Hash          string
	OriginalURL   string
	CorrelationID *string
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}

func NewURL(original string, hash string, correlationID *string) (*URL, error) {
	original = strings.TrimSpace(original)
	hash = strings.TrimSpace(hash)

	if original == "" {
		return nil, errs.ValidationError("empty url")
	}
	if hash == "" {
		return nil, errs.ValidationError("empty hash")
	}

	return &URL{
		ID:            uuid.Must(uuid.NewV7()),
		OriginalURL:   original,
		Hash:          hash,
		CorrelationID: correlationID,
		CreatedAt:     time.Now(),
	}, nil
}
