package config

import (
	"flag"
	"fmt"
	"sync"
)

// Config @TODO: env
type Config struct {
	Address  string
	BaseURL  string
	LogLevel string
}

var (
	cfg  *Config
	once sync.Once
)

func NewConfig() *Config {
	once.Do(func() {
		address := flag.String("a", "localhost:8080", "Адрес запуска HTTP-сервера (например, localhost:8080)")
		baseURL := flag.String("b", "http://localhost:8080", "Базовый адрес для коротких ссылок (например, http://localhost:8080)")

		flag.Parse()

		if *baseURL == "" {
			*baseURL = fmt.Sprintf("http://%s/", *address)
		}

		cfg = &Config{
			Address:  *address,
			BaseURL:  *baseURL,
			LogLevel: "info",
		}
	})

	return cfg
}
