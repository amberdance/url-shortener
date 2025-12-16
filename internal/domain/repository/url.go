package repository

import (
	"context"

	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/google/uuid"
)

//go:generate mockgen -source=$GOFILE -destination=../../mocks/url_repository_mock.go -package=mocks
type URLRepository interface {
	Create(ctx context.Context, url *model.URLEntry) error
	CreateBatch(ctx context.Context, urls []*model.URLEntry) error
	FindByHash(ctx context.Context, hash string) (*model.URLEntry, error)
	FindByOriginalURL(ctx context.Context, originalURL string) (*model.URLEntry, error)
	FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]*model.URLEntry, error)
}
