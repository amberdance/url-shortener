package url

import (
	"context"

	"github.com/amberdance/url-shortener/internal/app/command"
	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/amberdance/url-shortener/internal/domain/repository"
)

type GetURLsByUserIDUseCase struct {
	repository repository.URLRepository
}

func NewGetURLsByUserIDUseCase(r repository.URLRepository) GetURLsByUserIDUseCase {
	return GetURLsByUserIDUseCase{repository: r}
}

func (uc GetURLsByUserIDUseCase) Run(ctx context.Context, cmd command.GetUrlsByUserIDCommand) ([]*model.URL, error) {
	return uc.repository.FindAllByUserID(ctx, cmd.UserID)
}
