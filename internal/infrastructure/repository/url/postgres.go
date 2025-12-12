package url

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/amberdance/url-shortener/internal/domain/errs"
	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbFields = []string{
	"id",
	"created_at",
	"updated_at",
	"hash",
	"original_url",
	"correlation_id",
}

func getFormattedSelectFields() string {
	return strings.Join(dbFields, ", ")
}

type Mapper interface {
	Scan(dest ...any) error
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresURLRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, m *model.URL) error {
	_, err := r.pool.Exec(ctx,
		"insert into urls (id, created_at, hash, original_url, correlation_id, user_id) values ($1, $2, $3, $4, $5, $6)",
		m.ID,
		m.CreatedAt,
		m.Hash,
		m.OriginalURL,
		m.CorrelationID,
		m.UserID,
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
	sql := "insert into urls (id, created_at, hash, original_url, correlation_id, user_id) values ($1, $2, $3, $4, $5, $6)"

	for _, u := range urls {
		batch.Queue(sql, u.ID, u.CreatedAt, u.Hash, u.OriginalURL, u.CorrelationID, u.UserID)
	}

	br := tx.SendBatch(ctx, batch)

	for i := 0; i < len(urls); i++ {
		_, execErr := br.Exec()
		if execErr != nil {
			br.Close()
			return fmt.Errorf("batch insert failed: %w", execErr)
		}
	}

	if closeErr := br.Close(); closeErr != nil {
		return fmt.Errorf("batch close failed: %w", closeErr)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}

func (r *PostgresRepository) FindByHash(ctx context.Context, hash string) (*model.URL, error) {
	return r.mapToModel(r.pool.QueryRow(ctx, "select "+getFormattedSelectFields()+" from urls where hash = $1", hash))
}

func (r *PostgresRepository) FindByOriginalURL(ctx context.Context, original string) (*model.URL, error) {
	return r.mapToModel(r.pool.QueryRow(
		ctx,
		"select "+getFormattedSelectFields()+" from urls  where original_url = $1 limit 1",
		original,
	))
}

func (r *PostgresRepository) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]*model.URL, error) {
	rows, err := r.pool.Query(
		ctx,
		"select "+getFormattedSelectFields()+" from urls where user_id=$1",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.URL
	for rows.Next() {
		m, err := r.mapToModel(rows)
		if err != nil {
			return nil, err
		}

		result = append(result, m)
	}

	return result, nil
}

func (r *PostgresRepository) mapToModel(mapper Mapper) (*model.URL, error) {
	var m model.URL
	err := mapper.Scan(
		&m.ID,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.Hash,
		&m.OriginalURL,
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
