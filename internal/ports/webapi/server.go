package webapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amberdance/url-shortener/internal/app"
	"github.com/amberdance/url-shortener/internal/domain/contracts"
	"github.com/amberdance/url-shortener/internal/ports/webapi/handlers"
	mdw "github.com/amberdance/url-shortener/internal/ports/webapi/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

type Server struct {
	httpServer *http.Server
	logger     contracts.Logger
}

func NewServer(a *app.App) *Server {
	router := buildRoutes(a)
	handler := cors.AllowAll().Handler(router)
	httpSrv := &http.Server{
		Addr:    a.Config().Address,
		Handler: handler,
	}

	return &Server{httpServer: httpSrv, logger: a.Logger()}
}

func (s *Server) Run(ctx context.Context) error {
	l := s.logger
	l.Info(fmt.Sprintf("Server is running on %s", s.httpServer.Addr))

	idleConnsClosed := make(chan struct{})

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		select {
		case <-quit:
			l.Info("Shutdown signal received")
		case <-ctx.Done():
			l.Info("Context cancelled, shutting down")
		}

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			l.Error("HTTP server shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	<-idleConnsClosed
	l.Info("Server stopped gracefully")
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping server manually")
	return s.httpServer.Shutdown(ctx)
}

func buildRoutes(a *app.App) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Mount("/health", handlers.NewHealthcheckHandler().Routes())
	router.Mount("/ping", handlers.NewPingHandler(a.Pinger()).Routes())

	cont := a.Container()
	baseURL := a.Config().BaseURL

	router.Group(func(r chi.Router) {
		r.Use(mdw.JSONMiddleware)
		r.Use(mdw.GzipDecompressMiddleware)
		r.Use(mdw.GzipCompressMiddleware)
		r.Use(mdw.AuthMiddleware(a.Container().Auth))

		r.Mount("/", handlers.NewURLShortenerHandler(
			baseURL,
			cont.UseCases.URL,
			cont.Validator,
			a.Logger()).
			Routes(),
		)

		r.Mount("/", handlers.NewUserHandler(baseURL, cont.UseCases.URL.GetAllByUserId).Routes())
	})

	return router
}
