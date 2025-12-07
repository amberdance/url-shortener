TEST_FLAGS=-v -coverprofile=coverage.out -covermode=atomic
TEST_PATH=./internal/...

address ?=
host ?=

up:
	docker compose -f docker-compose.local.yml -f docker-compose.local.override.yml up -d

down:
	docker compose down

log:
	docker compose logs -f

status:
	docker compose ps

test:
	go test $(TEST_PATH) $(TEST_FLAGS)

build:
	go build -o .bin/server cmd/shortener/main.go

run:
	go build -o .bin/server cmd/shortener/main.go
	.bin/server $(if $(address),-a $(address)) $(if $(host),-b $(host))

migrate:
	go build -o .bin/migrator cmd/migrator/main.go && .bin/migrator