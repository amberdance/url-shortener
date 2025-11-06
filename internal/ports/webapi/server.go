package webapi

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amberdance/url-shortener/internal/app"
	"github.com/amberdance/url-shortener/internal/ports/webapi/handlers"
	webmw "github.com/amberdance/url-shortener/internal/ports/webapi/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(a *app.App) *Server {
	router := buildRoutes(a)
	handler := cors.AllowAll().Handler(router)
	httpSrv := &http.Server{
		Addr:    a.Config().Address,
		Handler: handler,
	}

	return &Server{httpServer: httpSrv}
}

func (s *Server) Run(ctx context.Context) error {
	log.Printf("Server is running on %s\n", s.httpServer.Addr)

	idleConnsClosed := make(chan struct{})

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		select {
		case <-quit:
			log.Println("Shutdown signal received")
		case <-ctx.Done():
			log.Println("Context cancelled, shutting down")
		}

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve: %w", err)
	}

	<-idleConnsClosed
	log.Println("Server stopped gracefully")
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Println("Stopping server manually")
	return s.httpServer.Shutdown(ctx)
}

func buildRoutes(a *app.App) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Mount("/health", handlers.NewHealthcheckHandler().Routes())
	router.Group(func(r chi.Router) {
		r.Use(webmw.TextPlainHeaderMiddleware)
		r.Mount("/", handlers.NewURLShortenerHandler(a.Storage(), a.Config().BaseURL).Routes())
	})

	//router.NotFound(func(w http.ResponseWriter, r *http.Request) {
	//	http.Error(w, "404 page not found", http.StatusNotFound)
	//})

	return router
}
