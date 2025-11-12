package url

import (
	"context"

	"github.com/amberdance/url-shortener/internal/app/command"
	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/amberdance/url-shortener/internal/domain/repository"
)

type GetByHashUseCase struct {
	repository repository.URLRepository
}

func NewGetByHashUseCase(r repository.URLRepository) GetByHashUseCase {
	return GetByHashUseCase{repository: r}
}

func (uc GetByHashUseCase) Run(ctx context.Context, cmd command.GetURLByHashCommand) (*model.URL, error) {
	return uc.repository.FindByHash(ctx, cmd.Hash)
}
