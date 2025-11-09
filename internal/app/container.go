package app

import (
	"github.com/amberdance/url-shortener/internal/app/usecase"
	"github.com/amberdance/url-shortener/internal/app/usecase/url"
)

type Container struct {
	RepositoryProvider RepositoryProvider
	UseCases           struct {
		URL usecase.URLUseCases
	}
}

func buildContainer(r RepositoryProvider) *Container {
	return &Container{
		RepositoryProvider: r,
		UseCases: struct {
			URL usecase.URLUseCases
		}{
			URL: usecase.URLUseCases{
				Create:   url.NewCreateUrlUseCase(r.URLRepository()),
				GetByUrl: url.NewGetByHashUseCase(r.URLRepository()),
			},
		},
	}
}
