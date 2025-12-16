package url

import (
	"context"

	"github.com/amberdance/url-shortener/internal/app/command"
	"github.com/amberdance/url-shortener/internal/domain/ports"
	"github.com/amberdance/url-shortener/internal/domain/repository"
)

type DeleteUserURLsBatchUseCase struct {
	repository repository.URLRepository
	logger     ports.Logger
}

func NewDeleteUserURLsBatchUseCase(repo repository.URLRepository, l ports.Logger) DeleteUserURLsBatchUseCase {
	return DeleteUserURLsBatchUseCase{
		repository: repo,
		logger:     l,
	}
}

func (uc DeleteUserURLsBatchUseCase) Run(ctx context.Context, cmd command.DeleteUserURLSCommand) error {
	return uc.repository.DeleteByUserIDAndHashes(ctx, cmd.UserID, cmd.Hashes)
}
