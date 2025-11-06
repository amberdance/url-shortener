package storage

import (
	"context"
	"sync"
)

type InMemoryStorage struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]string),
	}
}

func (s *InMemoryStorage) Save(_ context.Context, id, originalURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[id] = originalURL
	return nil
}

func (s *InMemoryStorage) Get(_ context.Context, id string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.data[id]
	return url, ok
}
