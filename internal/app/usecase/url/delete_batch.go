package url

import (
	"context"
	"time"

	"github.com/amberdance/url-shortener/internal/app/command"
	"github.com/amberdance/url-shortener/internal/domain/contracts"
	"github.com/amberdance/url-shortener/internal/domain/repository"
	"github.com/google/uuid"
)

type DeleteUserURLsBatchUseCase struct {
	repository repository.URLRepository
	logger     contracts.Logger
	input      chan command.DeleteUserURLSCommand
}

func NewDeleteUserURLsBatchUseCase(repo repository.URLRepository, l contracts.Logger) DeleteUserURLsBatchUseCase {
	uc := DeleteUserURLsBatchUseCase{
		repository: repo,
		logger:     l,
		input:      make(chan command.DeleteUserURLSCommand, 100),
	}

	go uc.worker()

	return uc
}

func (uc DeleteUserURLsBatchUseCase) Run(ctx context.Context, cmd command.DeleteUserURLSCommand) error {
	return uc.repository.DeleteByUserIDAndHashes(ctx, cmd.UserID, cmd.Hashes)
}

func (uc DeleteUserURLsBatchUseCase) RunAsync(cmd command.DeleteUserURLSCommand) {
	uc.input <- cmd
}

func (uc DeleteUserURLsBatchUseCase) worker() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	batch := make(map[uuid.UUID][]string)

	for {
		select {
		case cmd := <-uc.input:
			batch[cmd.UserID] = append(batch[cmd.UserID], cmd.Hashes...)

		case <-ticker.C:
			uc.flush(batch)
			batch = make(map[uuid.UUID][]string)
		}
	}
}

func (uc DeleteUserURLsBatchUseCase) flush(batch map[uuid.UUID][]string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for userID, hashes := range batch {
		err := uc.repository.DeleteByUserIDAndHashes(ctx, userID, hashes)
		if err != nil {
			uc.logger.Error("не удалось удалить ссылки", "user_id", userID)
		}
	}
}
