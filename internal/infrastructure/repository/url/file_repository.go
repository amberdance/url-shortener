package url

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/amberdance/url-shortener/internal/domain/repository"
)

type FileRepository struct {
	mu   sync.RWMutex
	data map[string]*model.URL
	path string
}

func NewFileRepository(path string) repository.URLRepository {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}

	repo := &FileRepository{
		data: make(map[string]*model.URL),
		path: path,
	}

	if _, err := os.Stat(path); err == nil {
		if err := repo.load(); err != nil {
			panic(err)
		}
	}

	return repo
}

func (r *FileRepository) Create(_ context.Context, u *model.URL) error {
	r.mu.Lock()
	if _, exists := r.data[u.Hash]; exists {
		r.mu.Unlock()
		return errors.New("duplicate hash")
	}
	r.data[u.Hash] = u
	r.mu.Unlock()

	return r.save()
}

func (r *FileRepository) FindByHash(_ context.Context, hash string) (*model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.data[hash]
	if !ok {
		return nil, errors.New("not found")
	}

	return u, nil
}

func (r *FileRepository) load() error {
	file, err := os.Open(r.path)
	if err != nil {
		return err
	}
	defer file.Close()

	var records []model.URL
	if err := json.NewDecoder(file).Decode(&records); err != nil {
		return err
	}

	r.loadEntries(records)
	return nil
}

func (r *FileRepository) loadEntries(records []model.URL) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, m := range records {
		r.data[m.Hash] = &model.URL{
			ID:          m.ID,
			Hash:        m.Hash,
			OriginalURL: m.OriginalURL,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
		}
	}
}

func (r *FileRepository) recordEntries() []model.URL {
	r.mu.RLock()
	defer r.mu.RUnlock()

	records := make([]model.URL, 0, len(r.data))
	for _, m := range r.data {
		records = append(records, model.URL{
			ID:          m.ID,
			Hash:        m.Hash,
			OriginalURL: m.OriginalURL,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
		})
	}
	return records
}

func (r *FileRepository) save() error {
	records := r.recordEntries()

	tmp := r.path + ".tmp"
	file, err := os.Create(tmp)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	if err := enc.Encode(records); err != nil {
		return err
	}

	return os.Rename(tmp, r.path)
}
