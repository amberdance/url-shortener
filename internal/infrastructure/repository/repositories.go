package repository

import (
	"github.com/amberdance/url-shortener/internal/config"
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

func NewRepositories(c *config.Config, s *storage.PostgresStorage) Provider {
	return &repositories{
		//urlRepo: url.NewPostgresRepository(s.Pool()),
		urlRepo: url.NewFileRepository(c.FileStoragePath),
	}
}
