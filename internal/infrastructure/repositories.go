package app

import (
	"github.com/amberdance/url-shortener/internal/domain/repository"
	infr "github.com/amberdance/url-shortener/internal/infrastructure/repository/url"
)

type RepositoryProvider interface {
	URLRepository() repository.URLRepository
}

type repositories struct {
	urlRepo repository.URLRepository
}

func NewRepositories() RepositoryProvider {
	return &repositories{
		urlRepo: infr.NewInMemoryRepository(),
	}
}

func (r *repositories) URLRepository() repository.URLRepository {
	return r.urlRepo
}
