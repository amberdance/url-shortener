package url

import (
	"context"
	"time"

	"github.com/amberdance/url-shortener/internal/app/command"
	"github.com/amberdance/url-shortener/internal/domain/contracts"
	"github.com/amberdance/url-shortener/internal/domain/repository"
)

type DeleteUserURLsBatchUseCase struct {
	repository repository.URLRepository
	logger     contracts.Logger
}

func NewDeleteUserURLsBatchUseCase(r repository.URLRepository, l contracts.Logger) DeleteUserURLsBatchUseCase {
	return DeleteUserURLsBatchUseCase{repository: r, logger: l}
}

func (uc DeleteUserURLsBatchUseCase) Run(ctx context.Context, cmd command.DeleteUserURLSCommand) error {
	return uc.repository.DeleteByUserIDAndHashes(ctx, cmd.UserID, cmd.Hashes)
}

func (uc DeleteUserURLsBatchUseCase) RunAsync(cmd command.DeleteUserURLSCommand) error {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err := uc.Run(ctx, cmd)
		if err != nil {
			uc.logger.Error("ошибка удаления записей асинхронно", err.Error())
		}
	}()

	return nil
}
