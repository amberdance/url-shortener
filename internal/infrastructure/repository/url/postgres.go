package url

import (
	"context"
	"errors"
	"fmt"

	"github.com/amberdance/url-shortener/internal/domain/errs"
	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresURLRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, m *model.URL) error {
	_, err := r.pool.Exec(ctx,
		"insert into urls (id, created_at, hash, original_url, correlation_id) values ($1, $2, $3, $4, $5)",
		m.ID,
		m.CreatedAt,
		m.Hash,
		m.OriginalURL,
		m.CorrelationID,
	)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return errs.DuplicateEntryError(pgErr.Message)
	}

	return err
}

func (r *PostgresRepository) CreateBatch(ctx context.Context, urls []*model.URL) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	batch := &pgx.Batch{}
	sql := `insert into urls (id, hash, original_url, created_at, correlation_id) values ($1, $2, $3, $4, $5)`

	for _, u := range urls {
		batch.Queue(sql, u.ID, u.Hash, u.OriginalURL, u.CreatedAt, u.CorrelationID)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	for range urls {
		_, err = br.Exec()
		if err != nil {
			return fmt.Errorf("batch insert failed: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}

func (r *PostgresRepository) FindByHash(ctx context.Context, hash string) (*model.URL, error) {
	row := r.pool.QueryRow(ctx,
		"select id, created_at, updated_at, hash, original_url, correlation_id from urls where hash = $1",
		hash,
	)

	var u model.URL
	err := row.Scan(
		&u.ID,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.Hash,
		&u.OriginalURL,
		&u.CorrelationID,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *PostgresRepository) FindByOriginalURL(ctx context.Context, original string) (*model.URL, error) {
	row := r.pool.QueryRow(ctx,
		`select id, hash, original_url, created_at, updated_at, correlation_id 
         from urls 
         where original_url = $1 
         limit 1`,
		original,
	)

	var m model.URL
	err := row.Scan(
		&m.ID,
		&m.Hash,
		&m.OriginalURL,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.CorrelationID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &m, nil
}
