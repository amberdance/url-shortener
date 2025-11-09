package config

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Address  string `env:"HTTP_ADDR" env-default:"0.0.0.0:8080"`
	BaseURL  string `env:"HTTP_BASE_URL" env-default:"http://localhost"`
	LogLevel string `env:"LOG_LEVEL" env-default:"info"`
}

var (
	cfg  *Config
	once sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		cfg = &Config{}
		if err := godotenv.Load(); err != nil {
			log.Println("не найден файл .env, используются переменные окружения")
		}

		err := cleanenv.ReadEnv(cfg)
		if err != nil {
			log.Fatalln(err)
		}

		address := flag.String("a", "", "Адрес запуска HTTP-сервера (например, localhost:8080)")
		baseURL := flag.String("b", "", "Базовый адрес коротких ссылок (например, http://localhost:8080)")
		flag.Parse()

		if *address != "" {
			cfg.Address = *address
		}

		if *baseURL != "" {
			cfg.BaseURL = *baseURL
		}

		if !strings.HasPrefix(cfg.BaseURL, "http://") && !strings.HasPrefix(cfg.BaseURL, "https://") {
			cfg.BaseURL = fmt.Sprintf("http://%s", cfg.BaseURL)
		}

		cfg.BaseURL = strings.TrimRight(cfg.BaseURL, "/")
	})

	return cfg
}
