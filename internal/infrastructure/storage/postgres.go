package storage

import (
	"context"
	"time"

	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	dsn  string
	pool *pgxpool.Pool
}

func (s *PostgresStorage) Pool() *pgxpool.Pool {
	return s.pool
}

func (s *PostgresStorage) Ping(ctx context.Context) error {
	return s.pool.Ping(ctx)
}

func (s *PostgresStorage) Close() {
	s.pool.Close()
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &PostgresStorage{dsn: dsn, pool: pool}, nil
}
