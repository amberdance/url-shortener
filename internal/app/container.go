package app

import (
	"github.com/amberdance/url-shortener/internal/app/usecase"
	"github.com/amberdance/url-shortener/internal/app/usecase/url"
	"github.com/amberdance/url-shortener/internal/config"
	"github.com/amberdance/url-shortener/internal/domain/contracts"
	"github.com/amberdance/url-shortener/internal/infrastructure/auth"
)

type Container struct {
	RepositoryProvider RepositoryProvider
	Auth               *auth.CookieAuth
	UseCases           struct {
		URL usecase.URLUseCases
	}
}

func buildContainer(r RepositoryProvider, cfg *config.Config, l contracts.Logger) *Container {
	rep := r.URLRepository()

	return &Container{
		RepositoryProvider: r,
		Auth:               auth.NewCookieAuth(cfg.AuthSecret),
		UseCases: struct {
			URL usecase.URLUseCases
		}{
			URL: usecase.URLUseCases{
				Create:                   url.NewCreateUseCase(rep),
				CreateBatch:              url.NewBatchCreateURLUseCase(rep),
				GetByURL:                 url.NewGetByHashUseCase(rep),
				GetByUserID:              url.NewGetURLsByUserIDUseCase(rep),
				DeleteByUserIDBatch:      url.NewDeleteUserURLsBatchUseCase(rep, l),
				DeleteByUserIDBatchAsync: url.NewDeleteUserURLsBatchAsyncUseCase(rep, l),
			},
		},
	}
}
