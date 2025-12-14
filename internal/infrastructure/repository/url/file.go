package url

import (
	"context"

	"github.com/amberdance/url-shortener/internal/domain/errs"
	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/amberdance/url-shortener/internal/domain/repository"
	"github.com/amberdance/url-shortener/internal/infrastructure/storage"
	"github.com/google/uuid"
)

type FileRepository struct {
	storage *storage.FileStorage
}

func NewFileURLRepository(s *storage.FileStorage) repository.URLRepository {
	return &FileRepository{
		storage: s,
	}
}

func (r *FileRepository) Create(ctx context.Context, u *model.URLEntry) error {
	if existing, _ := r.FindByOriginalURL(ctx, u.OriginalURL); existing != nil {
		return errs.ErrDuplicate
	}

	return r.storage.Put(u)
}

func (r *FileRepository) CreateBatch(_ context.Context, urls []*model.URLEntry) error {
	return r.storage.PutBatch(urls)
}

func (r *FileRepository) FindByHash(_ context.Context, hash string) (*model.URLEntry, error) {
	u, ok := r.storage.GetByHash(hash)
	if !ok {
		return nil, errs.ErrNotFound
	}

	return u, nil
}

func (r *FileRepository) FindByOriginalURL(_ context.Context, originalURL string) (*model.URLEntry, error) {
	u, ok := r.storage.GetByOriginalURL(originalURL)
	if !ok {
		return nil, errs.ErrNotFound
	}
	return u, nil
}

func (r *FileRepository) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]*model.URLEntry, error) {
	return r.storage.GetByUserID(ctx, userID)
}

func (r *FileRepository) DeleteByUserIDAndHashes(ctx context.Context, userID uuid.UUID, hashes []string) error {
	return r.storage.DeleteByUserIdBatch(ctx, userID, hashes)
}
