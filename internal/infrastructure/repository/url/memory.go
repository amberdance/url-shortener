package url

import (
	"context"
	"errors"
	"fmt"

	"github.com/amberdance/url-shortener/internal/domain/errs"
	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/amberdance/url-shortener/internal/domain/repository"
	"github.com/amberdance/url-shortener/internal/infrastructure/storage"
)

type inMemoryRepository struct {
	storage *storage.InMemoryStorage
}

var _ repository.URLRepository = (*inMemoryRepository)(nil)

func NewInMemoryURLRepository() repository.URLRepository {
	return &inMemoryRepository{
		storage: storage.NewInMemoryStorage(),
	}
}

func (r *inMemoryRepository) Create(ctx context.Context, m *model.URL) error {
	existing, _ := r.FindByOriginalURL(ctx, m.OriginalURL)
	if existing != nil {
		return errs.DuplicateEntryError("url already exists")
	}

	r.storage.Mu.Lock()
	defer r.storage.Mu.Unlock()

	r.storage.Data[m.ID] = m
	return nil
}

func (r *inMemoryRepository) CreateBatch(_ context.Context, urls []*model.URL) error {
	r.storage.Mu.Lock()
	defer r.storage.Mu.Unlock()
	for _, u := range urls {
		if _, ok := r.storage.Data[u.ID]; ok {
			return fmt.Errorf("duplicate hash: %s", u.Hash)
		}
	}
	for _, u := range urls {
		r.storage.Data[u.ID] = u
	}
	return nil
}

func (r *inMemoryRepository) FindByHash(_ context.Context, url string) (*model.URL, error) {
	r.storage.Mu.RLock()
	defer r.storage.Mu.RUnlock()

	for _, item := range r.storage.Data {
		if item.Hash == url {
			return item, nil
		}
	}
	return nil, errors.New("url not found")
}

func (r *inMemoryRepository) FindByOriginalURL(_ context.Context, originalURL string) (*model.URL, error) {
	r.storage.Mu.RLock()
	defer r.storage.Mu.RUnlock()

	for _, m := range r.storage.Data {
		if m.OriginalURL == originalURL {
			return m, nil
		}
	}

	return nil, nil
}
