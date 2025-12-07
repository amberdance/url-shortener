package app

import (
	"fmt"
	"sync"

	"github.com/amberdance/url-shortener/internal/config"
	"github.com/amberdance/url-shortener/internal/domain/shared"
	"github.com/amberdance/url-shortener/internal/infrastructure/logging"
	"github.com/amberdance/url-shortener/internal/infrastructure/repository"
	"github.com/amberdance/url-shortener/internal/infrastructure/storage"
)

type App struct {
	config    *config.Config
	container *Container
	logger    shared.Logger
	storage   *storage.PostgresStorage
}

var (
	instance *App
	once     sync.Once
)

func GetApp() (*App, error) {
	var initErr error

	once.Do(func() {
		app := &App{}
		if err := app.init(); err != nil {
			initErr = fmt.Errorf("failed to initialize app: %w", err)
			return
		}
		instance = app
	})

	if initErr != nil {
		return nil, initErr
	}
	return instance, nil
}

func (a *App) Config() *config.Config { return a.config }

func (a *App) Container() *Container { return a.container }

func (a *App) Logger() shared.Logger { return a.logger }

func (a *App) Storage() *storage.PostgresStorage { return a.storage }

func (a *App) Close() {
	if a.logger != nil {
		a.logger.Close()
	}
	if a.storage != nil {
		a.storage.Close()
	}
}

func (a *App) init() error {
	a.config = config.GetConfig()
	a.logger = logging.NewLogger()

	p, err := a.resolveRepositoryProvider()
	if err != nil {
		return err
	}

	a.container = buildContainer(p)
	return nil
}

func (a *App) resolveRepositoryProvider() (repository.Provider, error) {
	if a.config.DatabaseDSN != "" {
		st, err := storage.NewPostgresStorage(a.config.DatabaseDSN)
		if err != nil {
			return nil, fmt.Errorf("database connection error: %w", err)
		}
		a.storage = st
		return repository.NewRepositories(st), nil
	}

	if a.config.FileStoragePath != "" {
		return repository.NewFileRepositories(a.config.FileStoragePath), nil
	}

	return repository.NewMemoryRepositories(), nil
}
