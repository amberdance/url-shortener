package config

import (
	"flag"
	"fmt"
)

type Config struct {
	Address string
	BaseURL string
}

func NewConfig() *Config {
	address := flag.String("a", "localhost:8080", "Адрес запуска HTTP-сервера (:8080)")
	baseURL := flag.String("b", "http://localhost:8080", "Полный адрес сервера (http://localhost:8080)")

	flag.Parse()

	if *baseURL == "" {
		*baseURL = fmt.Sprintf("http://%s/", *address)
	}

	return &Config{
		Address: *address,
		BaseURL: *baseURL,
	}
}
