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

func (r *Postgres) Create(ctx context.Context, m *model.URL) error {
	_, err := r.pool.Exec(ctx,
		"insert into urls (id, created_at, hash, original_url) values ($1, $2, $3, $4)",
		m.ID,
		m.CreatedAt,
		m.Hash,
		m.OriginalURL,
	)
	return err
}

func (r *Postgres) FindByHash(ctx context.Context, hash string) (*model.URL, error) {
	row := r.pool.QueryRow(ctx,
		"select id, created_at, updated_at, hash, original_url from urls where hash = $1",
		hash,
	)

	var u model.URL
	err := row.Scan(
		&u.ID,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.Hash,
		&u.OriginalURL,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
