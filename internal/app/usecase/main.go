package usecase

import "github.com/amberdance/url-shortener/internal/app/usecase/url"

type URLUseCases struct {
	GetByURL       url.GetByHashUseCase
	Create         url.CreateUseCase
	CreateBatch    url.BatchCreateURLUseCase
	GetAllByUserId url.GetURLsByUserIDUseCase
}
