package app

import "github.com/amberdance/url-shortener/internal/domain/repository"

type RepositoryProvider interface {
	URLRepository() repository.URLRepository
}
