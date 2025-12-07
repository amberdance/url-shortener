package app

import (
	"github.com/amberdance/url-shortener/internal/app/usecase"
	"github.com/amberdance/url-shortener/internal/app/usecase/url"
	"github.com/go-playground/validator/v10"
)

type Container struct {
	RepositoryProvider RepositoryProvider
	Validator          *validator.Validate
	UseCases           struct {
		URL usecase.URLUseCases
	}
}

func buildContainer(r RepositoryProvider) *Container {
	return &Container{
		RepositoryProvider: r,
		Validator:          validator.New(),
		UseCases: struct {
			URL usecase.URLUseCases
		}{
			URL: usecase.URLUseCases{
				Create:   url.NewCreateURLUseCase(r.URLRepository()),
				GetByURL: url.NewGetByHashUseCase(r.URLRepository()),
			},
		},
	}
}
