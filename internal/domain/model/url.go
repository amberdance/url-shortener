package model

import (
	"time"

	"github.com/google/uuid"
)

type Url struct {
	ID          uuid.UUID
	Hash        string
	OriginalURL string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

func NewURL(original string, hash string) *Url {
	return &Url{
		ID:          uuid.Must(uuid.NewV7()),
		OriginalURL: original,
		Hash:        hash,
		CreatedAt:   time.Now(),
	}
}
