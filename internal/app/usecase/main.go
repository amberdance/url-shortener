package usecase

import "github.com/amberdance/url-shortener/internal/app/usecase/url"

type URLUseCases struct {
	GetByUrl url.GetByHashUseCase
	Create   url.CreateUseCase
}
