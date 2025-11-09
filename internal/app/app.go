package app

import (
	"fmt"
	"sync"

	"github.com/amberdance/url-shortener/internal/config"
	"github.com/amberdance/url-shortener/internal/domain/shared"
	infr "github.com/amberdance/url-shortener/internal/infrastructure"
	"github.com/amberdance/url-shortener/internal/infrastructure/logging"
)

type App struct {
	config    *config.Config
	container *Container
	logger    shared.Logger
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

func (a *App) Close() {
	if a.logger != nil {
		a.logger.Close()
	}
}

func (a *App) init() error {
	a.config = config.GetConfig()
	a.logger = logging.NewLogger()
	a.container = buildContainer(infr.NewRepositories())

	return nil
}
