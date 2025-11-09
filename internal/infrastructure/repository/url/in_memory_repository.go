package url

import (
	"context"
	"errors"
	"sync"

	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/amberdance/url-shortener/internal/domain/repository"
	"github.com/google/uuid"
)

type inMemoryStorage struct {
	data map[uuid.UUID]*model.Url
	mu   sync.RWMutex
}

func newInMemoryStorage() *inMemoryStorage {
	return &inMemoryStorage{
		data: make(map[uuid.UUID]*model.Url),
	}
}

type inMemoryRepository struct {
	storage *inMemoryStorage
}

var _ repository.URLRepository = (*inMemoryRepository)(nil)

func NewInMemoryRepository() repository.URLRepository {
	return &inMemoryRepository{
		storage: newInMemoryStorage(),
	}
}

func (r *inMemoryRepository) Create(_ context.Context, m *model.Url) error {
	r.storage.mu.Lock()
	defer r.storage.mu.Unlock()

	if _, exists := r.storage.data[m.ID]; exists {
		return errors.New("url with this ID already exists")
	}

	r.storage.data[m.ID] = m
	return nil
}

func (r *inMemoryRepository) FindByHash(_ context.Context, url string) (*model.Url, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	for _, item := range r.storage.data {
		if item.Hash == url {
			return item, nil
		}
	}
	return nil, errors.New("url not found")
}
