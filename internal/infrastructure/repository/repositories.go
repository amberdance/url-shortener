package repository

import (
	"github.com/amberdance/url-shortener/internal/domain/repository"
	"github.com/amberdance/url-shortener/internal/infrastructure/repository/url"
	"github.com/amberdance/url-shortener/internal/infrastructure/storage"
)

type Provider interface {
	URLRepository() repository.URLRepository
}

type repositories struct {
	urlRepo repository.URLRepository
}

func (r *repositories) URLRepository() repository.URLRepository {
	return r.urlRepo
}

func NewRepositories(s *storage.PostgresStorage) Provider {
	return &repositories{urlRepo: url.NewPostgresURLRepository(s.Pool())}
}

func NewFileRepositories(s *storage.FileStorage) Provider {
	return &repositories{urlRepo: url.NewFileURLRepository(s)}
}

func NewMemoryRepositories(s *storage.InMemoryStorage) Provider {
	return &repositories{urlRepo: url.NewInMemoryURLRepository(s)}
}
