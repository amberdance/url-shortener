package app

import (
	"fmt"
	"sync"

	"github.com/amberdance/url-shortener/internal/config"
	"github.com/amberdance/url-shortener/internal/domain/shared"
	"github.com/amberdance/url-shortener/internal/infrastructure/logging"
	infr "github.com/amberdance/url-shortener/internal/infrastructure/storage"
)

type App struct {
	config  *config.Config
	storage Storage
	logger  shared.Logger
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

func (a *App) Config() *config.Config {
	return a.config
}

func (a *App) Storage() Storage {
	return a.storage
}

func (a *App) Logger() shared.Logger { return a.logger }

func (a *App) init() error {
	a.config = config.NewConfig()
	a.logger = logging.NewLogger()

	// @TODO: не забыть скрыть за интерфейсом
	a.storage = infr.NewInMemoryStorage()

	return nil
}
