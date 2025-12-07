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
	Address         string `env:"SERVER_ADDRESS" env-default:"0.0.0.0:8080"`
	BaseURL         string `env:"BASE_URL" env-default:""`
	LogLevel        string `env:"LOG_LEVEL" env-default:"info"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" env-default:"./db/db.json"`
	//DatabaseDSN     string `env:"DATABASE_DSN"`
	DatabaseDSN *string
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
		if err := cleanenv.ReadEnv(cfg); err != nil {
			log.Fatalln(err)
		}

		address := flag.String("a", "", "Адрес запуска HTTP-сервера (например, localhost:8080)")
		baseURL := flag.String("b", "", "Базовый адрес коротких ссылок (например, http://localhost:8080)")
		dbFilePath := flag.String("f", "", "Путь к файловому хранилищу")
		dsn := flag.String("d", "", "PostgreSQL DSN")
		flag.Parse()

		if *address != "" {
			cfg.Address = *address
		}

		if *baseURL != "" {
			cfg.BaseURL = *baseURL
		}

		if cfg.BaseURL == "" {
			cfg.BaseURL = fmt.Sprintf("http://%s", cfg.Address)
		}

		if !strings.HasPrefix(cfg.BaseURL, "http://") && !strings.HasPrefix(cfg.BaseURL, "https://") {
			cfg.BaseURL = "http://" + cfg.BaseURL
		}

		cfg.BaseURL = strings.TrimRight(cfg.BaseURL, "/") + "/"

		if *dbFilePath != "" {
			cfg.FileStoragePath = *dbFilePath
		}

		if *dsn != "" {
			cfg.DatabaseDSN = dsn
		}
	})

	return cfg
}
