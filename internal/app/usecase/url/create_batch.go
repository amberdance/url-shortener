package url

import (
	"context"

	"github.com/amberdance/url-shortener/internal/app/command"
	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/amberdance/url-shortener/internal/domain/repository"
	"github.com/amberdance/url-shortener/internal/infrastructure/helpers"
)

type BatchCreateURLUseCase struct {
	repo repository.URLRepository
}

func NewBatchCreateURLUseCase(r repository.URLRepository) BatchCreateURLUseCase {
	return BatchCreateURLUseCase{repo: r}
}

func (uc *BatchCreateURLUseCase) Run(ctx context.Context, cmd command.CreateBatchURLEntryCommand) ([]*model.URL, error) {
	var urls []*model.URL
	for _, e := range cmd.Entries {
		m, err := model.NewURL(e.OriginalURL, helpers.GenerateHash(), e.CorrelationID, e.UserID)
		if err != nil {
			return nil, err
		}

		urls = append(urls, m)
	}

	if err := uc.repo.CreateBatch(ctx, urls); err != nil {
		return nil, err
	}

	return urls, nil
}
