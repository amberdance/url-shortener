package model

import (
	"strings"
	"time"

	"github.com/amberdance/url-shortener/internal/domain/errs"
	"github.com/google/uuid"
)

type URLEntry struct {
	ID            uuid.UUID
	Hash          string
	OriginalURL   string
	CorrelationID *string
	UserID        *uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     *time.Time
	DeletedAt     *time.Time
}

func NewURLEntry(original string, hash string, correlationID *string, userID *uuid.UUID) (*URLEntry, error) {
	original = strings.TrimSpace(original)
	hash = strings.TrimSpace(hash)

	if original == "" {
		return nil, errs.ValidationError("empty url")
	}
	if hash == "" {
		return nil, errs.ValidationError("empty hash")
	}

	return &URLEntry{
		ID:            uuid.Must(uuid.NewV7()),
		OriginalURL:   original,
		Hash:          hash,
		CorrelationID: correlationID,
		UserID:        userID,
		CreatedAt:     time.Now(),
	}, nil
}
