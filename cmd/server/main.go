package main

import (
	"log"

	"github.com/amberdance/url-shortener/internal/ports/webapi"
)

func main() {
	server := webapi.NewServer(":8080")

	if err := server.Run(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
