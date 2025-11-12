package repository

import (
	"context"

	"github.com/amberdance/url-shortener/internal/domain/model"
)

type URLRepository interface {
	Create(ctx context.Context, url *model.URL) error
	FindByHash(ctx context.Context, hash string) (*model.URL, error)
}
