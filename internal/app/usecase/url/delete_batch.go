package url

import (
	"context"

	"github.com/amberdance/url-shortener/internal/app/command"
	"github.com/amberdance/url-shortener/internal/domain/contracts"
	"github.com/amberdance/url-shortener/internal/domain/repository"
)

type DeleteUserURLsBatchUseCase struct {
	repository repository.URLRepository
	logger     contracts.Logger
}

func NewDeleteUserURLsBatchUseCase(repo repository.URLRepository, l contracts.Logger) DeleteUserURLsBatchUseCase {
	return DeleteUserURLsBatchUseCase{
		repository: repo,
		logger:     l,
	}
}

func (uc DeleteUserURLsBatchUseCase) Run(ctx context.Context, cmd command.DeleteUserURLSCommand) error {
	return uc.repository.DeleteByUserIDAndHashes(ctx, cmd.UserID, cmd.Hashes)
}
