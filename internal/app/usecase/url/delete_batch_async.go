package url

import (
	"context"
	"sync"
	"time"

	"github.com/amberdance/url-shortener/internal/app/command"
	"github.com/amberdance/url-shortener/internal/domain/contracts"
	"github.com/amberdance/url-shortener/internal/domain/repository"
)

type DeleteUserURLsBatchAsyncUseCase struct {
	repository repository.URLRepository
	logger     contracts.Logger
	ch         chan command.DeleteUserURLSCommand
	stopCh     chan struct{}
	wg         *sync.WaitGroup
	once       *sync.Once
}

func NewDeleteUserURLsBatchAsyncUseCase(repo repository.URLRepository, l contracts.Logger) DeleteUserURLsBatchAsyncUseCase {
	uc := DeleteUserURLsBatchAsyncUseCase{
		repository: repo,
		logger:     l,
		ch:         make(chan command.DeleteUserURLSCommand, 100),
		stopCh:     make(chan struct{}),
		wg:         &sync.WaitGroup{},
		once:       &sync.Once{},
	}

	uc.wg.Add(1)
	go uc.processLoop()

	return uc
}

func (uc DeleteUserURLsBatchAsyncUseCase) Run(ctx context.Context, cmd command.DeleteUserURLSCommand) error {
	select {
	case uc.ch <- cmd:
		uc.logger.Debug("команда на удаление отправлена в очередь", "user_id", cmd.UserID)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-uc.stopCh:
		return context.Canceled
	}
}

func (uc DeleteUserURLsBatchAsyncUseCase) processLoop() {
	defer uc.wg.Done()

	for {
		select {
		case cmd := <-uc.ch:
			uc.processDelete(cmd)
		case <-uc.stopCh:
			uc.logger.Debug("обработчик команд на удаление остановлен")
			return
		}
	}
}

func (uc DeleteUserURLsBatchAsyncUseCase) processDelete(cmd command.DeleteUserURLSCommand) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := uc.repository.DeleteByUserIDAndHashes(ctx, cmd.UserID, cmd.Hashes)
	if err != nil {
		uc.logger.Error("не удалось пакетно удалить ссылки", "user_id", cmd.UserID, "error", err)
	} else {
		uc.logger.Debug("пакетное удаление завершено успешно", "user_id", cmd.UserID)
	}
}

func (uc DeleteUserURLsBatchAsyncUseCase) Shutdown() {
	uc.once.Do(func() {
		close(uc.stopCh)
		uc.wg.Wait()
		close(uc.ch)
	})
}
