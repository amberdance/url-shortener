package config

import (
	"flag"
	"fmt"
	"strings"
	"sync"
)

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
		baseURL := flag.String("b", "", "Базовый адрес для HTTP-сервера ссылок (например, http://localhost:8080)")
		flag.Parse()

		url := *baseURL
		if url == "" {
			url = fmt.Sprintf("http://%s", *address)
		}

		url = strings.TrimRight(url, "/")

		cfg = &Config{
			Address:  *address,
			BaseURL:  url,
			LogLevel: "info",
		}
	})
	return cfg
}
