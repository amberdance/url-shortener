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
	address := flag.String("a", "localhost:8080", "Адрес запуска HTTP-сервера (localhost:8888)")
	baseURL := flag.String("b", "", "Базовый адрес результирующего сокращённого URL (http://localhost:8080)")

	flag.Parse()

	if *baseURL == "" {
		*baseURL = fmt.Sprintf("http://%s/", *address)
	}

	return &Config{
		Address: *address,
		BaseURL: *baseURL,
	}
}
