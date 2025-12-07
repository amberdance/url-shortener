package url

import (
	"context"

	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *Postgres {
	return &Postgres{pool: pool}
}

func (r *Postgres) Create(ctx context.Context, url *model.URL) error {
	panic("implement me")
}

func (r *Postgres) FindByHash(ctx context.Context, hash string) (*model.URL, error) {
	panic("implement me")
}
