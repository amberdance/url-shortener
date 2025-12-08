package url

import (
	"context"
	"errors"

	"github.com/amberdance/url-shortener/internal/domain/errs"
	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/amberdance/url-shortener/internal/domain/repository"
	"github.com/amberdance/url-shortener/internal/infrastructure/storage"
)

type FileRepository struct {
	storage *storage.FileStorage
}

func NewFileURLRepository(path string) repository.URLRepository {
	return &FileRepository{
		storage: storage.NewFileStorage(path),
	}
}

func (r *FileRepository) Create(ctx context.Context, u *model.URL) error {
	if existing, _ := r.FindByOriginalURL(ctx, u.OriginalURL); existing != nil {
		return errs.DuplicateEntryError("url already exists")
	}

	return r.storage.Put(u)
}

func (r *FileRepository) CreateBatch(_ context.Context, urls []*model.URL) error {
	return r.storage.PutBatch(urls)
}

func (r *FileRepository) FindByHash(_ context.Context, hash string) (*model.URL, error) {
	u, ok := r.storage.GetByHash(hash)
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

func (r *FileRepository) FindByOriginalURL(_ context.Context, originalURL string) (*model.URL, error) {
	u, ok := r.storage.GetByOriginalURL(originalURL)
	if !ok {
		return nil, errs.NotFoundError("url not found")
	}
	return u, nil
}
