package app

import "github.com/amberdance/url-shortener/internal/app/service"

type Container struct {
	Services struct {
		Shortener *service.URLShortenerService
	}
}

func buildContainer(a *App) *Container {
	container := &Container{
		Services: struct{ Shortener *service.URLShortenerService }{Shortener: service.NewURLShortenerService(a.storage)},
	}

	return container
}
