package main

import (
	"context"
	"log"

	"github.com/amberdance/url-shortener/internal/infrastructure/storage"
	"github.com/amberdance/url-shortener/internal/ports/webapi"
)

func main() {
	st := storage.NewInMemoryStorage()
	srv := webapi.NewServer(":8080", st)

	if err := srv.Run(context.Background()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
