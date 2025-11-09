package repository

import (
	"context"

	"github.com/amberdance/url-shortener/internal/domain/model"
)

type URLRepository interface {
	Create(ctx context.Context, url *model.Url) error
	FindByHash(ctx context.Context, hash string) (*model.Url, error)
}
