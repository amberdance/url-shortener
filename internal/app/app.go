package app

import (
	"fmt"
	"sync"

	"github.com/amberdance/url-shortener/internal/config"
	infr "github.com/amberdance/url-shortener/internal/infrastructure/storage"
)

type App struct {
	config  *config.Config
	storage Storage
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

func (a *App) init() error {
	a.config = config.NewConfig()

	// @TODO: не забыть скрыть за интерфейсом
	a.storage = infr.NewInMemoryStorage()

	return nil
}
