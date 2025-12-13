package url

import (
	"context"
	"errors"

	"github.com/amberdance/url-shortener/internal/app/command"
	"github.com/amberdance/url-shortener/internal/domain/errs"
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
	m, err := model.NewURL(cmd.OriginalURL, helpers.GenerateHash(), cmd.CorrelationID, cmd.UserID)
	if err != nil {
		return nil, err
	}

	err = uc.repository.Create(ctx, m)
	if err != nil {
		var dup errs.DuplicateEntryError
		if errors.As(err, &dup) {
			existed, findErr := uc.repository.FindByOriginalURL(ctx, m.OriginalURL)
			if findErr != nil {
				return nil, findErr
			}
			if existed == nil {
				return nil, errs.NotFoundError("URL not found")
			}
			return existed, dup
		}
		return nil, err
	}

	return m, nil
}
