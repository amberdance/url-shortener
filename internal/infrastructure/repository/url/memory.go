package url

import (
	"context"
	"fmt"

	"github.com/amberdance/url-shortener/internal/domain/errs"
	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/amberdance/url-shortener/internal/domain/repository"
	"github.com/amberdance/url-shortener/internal/infrastructure/storage"
	"github.com/google/uuid"
)

type InMemoryRepository struct {
	storage *storage.InMemoryStorage
}

var _ repository.URLRepository = (*InMemoryRepository)(nil)

func NewInMemoryURLRepository(s *storage.InMemoryStorage) repository.URLRepository {
	return &InMemoryRepository{
		storage: s,
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, m *model.URLEntry) error {
	existing, _ := r.FindByOriginalURL(ctx, m.OriginalURL)
	if existing != nil {
		return errs.DuplicateEntryError("url already exists")
	}

	r.storage.Mu.Lock()
	defer r.storage.Mu.Unlock()

	r.storage.Data[m.ID] = m

	return nil
}

func (r *InMemoryRepository) CreateBatch(_ context.Context, urls []*model.URLEntry) error {
	r.storage.Mu.Lock()
	defer r.storage.Mu.Unlock()

	for _, u := range r.storage.Data {
		if u.OriginalURL == u.OriginalURL || u.Hash == u.Hash {
			return errs.DuplicateEntryError(fmt.Sprintf("url %s already exists", u.OriginalURL))
		}
	}

	for _, u := range urls {
		r.storage.Data[u.ID] = u
	}

	return nil
}

func (r *InMemoryRepository) FindByHash(_ context.Context, url string) (*model.URLEntry, error) {
	r.storage.Mu.RLock()
	defer r.storage.Mu.RUnlock()

	for _, item := range r.storage.Data {
		if item.Hash == url {
			return item, nil
		}
	}

	return nil, errs.NotFoundError("url not found")
}

func (r *InMemoryRepository) FindByOriginalURL(_ context.Context, originalURL string) (*model.URLEntry, error) {
	r.storage.Mu.RLock()
	defer r.storage.Mu.RUnlock()

	for _, m := range r.storage.Data {
		if m.OriginalURL == originalURL {
			return m, nil
		}
	}

	return nil, nil
}

func (r *InMemoryRepository) FindAllByUserID(_ context.Context, userID uuid.UUID) ([]*model.URLEntry, error) {
	var urls []*model.URLEntry

	r.storage.Mu.RLock()
	defer r.storage.Mu.RUnlock()

	for _, item := range r.storage.Data {
		if item.UserID != nil && *item.UserID == userID {
			urls = append(urls, item)
		}
	}

	return urls, nil
}
