# Url shortener service

## Установка и запуск

```bash
# Запуск для локальной разработки 
cp .env.example .env
cp docker-compose.local.override.example.yml docker-compose.local.override.yml
make up

# Сборка сервера
make build

# Запуск тестов
make test

# Запуск сервера
make run