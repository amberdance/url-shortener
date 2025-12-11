package repository

import (
	"context"

	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/google/uuid"
)

type URLRepository interface {
	Create(ctx context.Context, url *model.URL) error
	CreateBatch(ctx context.Context, urls []*model.URL) error
	FindByHash(ctx context.Context, hash string) (*model.URL, error)
	FindByOriginalURL(ctx context.Context, originalURL string) (*model.URL, error)
	FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]*model.URL, error)
}
