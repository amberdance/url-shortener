package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/google/uuid"
)

type FileStorage struct {
	mu   sync.RWMutex
	data map[string]*model.URL
	path string
}

func NewFileStorage(path string) *FileStorage {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}

	s := &FileStorage{
		data: make(map[string]*model.URL),
		path: path,
	}

	if _, err := os.Stat(path); err == nil {
		if err := s.loadFromDisk(); err != nil {
			panic(err)
		}
	}

	return s
}

func (s *FileStorage) Ping(_ context.Context) error {
	return nil
}

func (s *FileStorage) Put(u *model.URL) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[u.Hash] = u
	return s.save()
}

func (s *FileStorage) PutBatch(urls []*model.URL) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, u := range urls {
		if _, exists := s.data[u.Hash]; exists {
			return fmt.Errorf("duplicate hash: %s", u.Hash)
		}
	}

	for _, u := range urls {
		s.data[u.Hash] = u
	}

	return s.save()
}

func (s *FileStorage) GetByHash(hash string) (*model.URL, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.data[hash]
	return u, ok
}

func (s *FileStorage) GetByOriginalURL(original string) (*model.URL, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, u := range s.data {
		if u.OriginalURL == original {
			return u, true
		}
	}
	return nil, false
}

func (s *FileStorage) loadFromDisk() error {
	file, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer file.Close()

	var records []model.URL
	if err := json.NewDecoder(file).Decode(&records); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, m := range records {
		mCopy := m
		s.data[m.Hash] = &mCopy
	}

	return nil
}

func (s *FileStorage) save() error {
	tmp := s.path + ".tmp"

	file, err := os.Create(tmp)
	if err != nil {
		return err
	}
	defer file.Close()

	records := make([]model.URL, 0, len(s.data))
	for _, u := range s.data {
		records = append(records, *u)
	}

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	if err := enc.Encode(records); err != nil {
		return err
	}

	return os.Rename(tmp, s.path)
}

func (s *FileStorage) GetByUserId(_ context.Context, userID uuid.UUID) ([]*model.URL, error) {
	var urls []*model.URL

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, item := range s.data {
		if item.UserID != nil && *item.UserID == userID {
			urls = append(urls, item)
		}
	}

	return urls, nil
}
