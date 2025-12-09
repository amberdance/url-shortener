package storage

import (
	"sync"

	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/google/uuid"
)

type InMemoryStorage struct {
	Data map[uuid.UUID]*model.URL
	Mu   sync.RWMutex
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		Data: make(map[uuid.UUID]*model.URL),
	}
}
