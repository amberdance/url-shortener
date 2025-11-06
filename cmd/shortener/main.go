package main

import (
	"context"
	"log"

	"github.com/amberdance/url-shortener/internal/app"
	"github.com/amberdance/url-shortener/internal/ports/webapi"
)

func main() {
	a, err := app.GetApp()
	if err != nil {
		log.Fatalln(err)
	}

	defer a.Close()

	srv := webapi.NewServer(a)
	if err := srv.Run(context.Background()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
