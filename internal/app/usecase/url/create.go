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

func NewCreateUseCase(r repository.URLRepository) CreateUseCase {
	return CreateUseCase{repository: r}
}

func (uc CreateUseCase) Run(ctx context.Context, cmd command.CreateURLEntryCommand) (*model.URLEntry, error) {
	m, err := model.NewURLEntry(cmd.OriginalURL, helpers.GenerateHash(), cmd.CorrelationID, cmd.UserID)
	if err != nil {
		return nil, err
	}

	err = uc.repository.Create(ctx, m)
	if err != nil {
		var der errs.DuplicateEntryError
		if errors.As(err, &der) {
			existed, findErr := uc.repository.FindByOriginalURL(ctx, m.OriginalURL)
			if findErr != nil {
				return nil, findErr
			}
			if existed == nil {
				return nil, errs.ErrNotFound
			}
			return existed, der
		}
		return nil, err
	}

	return m, nil
}
