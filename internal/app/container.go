package app

import (
	"github.com/amberdance/url-shortener/internal/app/usecase"
	"github.com/amberdance/url-shortener/internal/app/usecase/url"
	"github.com/amberdance/url-shortener/internal/config"
	"github.com/amberdance/url-shortener/internal/infrastructure/auth"
	"github.com/go-playground/validator/v10"
)

type Container struct {
	RepositoryProvider RepositoryProvider
	Validator          *validator.Validate
	Auth               *auth.CookieAuth
	UseCases           struct {
		URL usecase.URLUseCases
	}
}

func buildContainer(r RepositoryProvider, cfg *config.Config) *Container {
	rep := r.URLRepository()

	return &Container{
		RepositoryProvider: r,
		Validator:          validator.New(),
		Auth:               auth.NewCookieAuth(cfg.AuthSecret),
		UseCases: struct {
			URL usecase.URLUseCases
		}{
			URL: usecase.URLUseCases{
				Create:      url.NewCreateURLUseCase(rep),
				CreateBatch: url.NewBatchCreateURLUseCase(rep),
				GetByURL:    url.NewGetByHashUseCase(rep),
				GetByUserID: url.NewGetURLsByUserIDUseCase(rep),
			},
		},
	}
}
