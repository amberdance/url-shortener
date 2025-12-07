package url

import (
	"context"

	"github.com/amberdance/url-shortener/internal/app/command"
	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/amberdance/url-shortener/internal/domain/repository"
	"github.com/amberdance/url-shortener/internal/infrastructure/helpers"
)

type CreateUseCase struct {
	repository repository.URLRepository
}

func NewCreateURLUseCase(r repository.URLRepository) CreateUseCase {
	return CreateUseCase{repository: r}
}

func (uc CreateUseCase) Run(ctx context.Context, cmd command.CreateURLEntryCommand) (*model.URL, error) {
	// @TODO:
	//exists, err := uc.repository.FindByUrl(cmd.URLs)
	//if exists != nil {
	//	return exists
	//}

	m, err := model.NewURL(cmd.OriginalURL, helpers.GenerateHash(), cmd.CorrelationID)
	if err != nil {
		return nil, err
	}

	err = uc.repository.Create(ctx, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
