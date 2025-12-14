address ?=
host ?=

up:
	docker compose -f docker-compose.local.yml -f docker-compose.local.override.yml up -d

down:
	docker compose -f docker-compose.local.yml -f docker-compose.local.override.yml down

log:
	docker compose logs -f

status:
	docker compose ps

test:
	gotestsum --format=testname -- -coverprofile=coverage.out -covermode=atomic ./internal/...

build:
	go build -o .bin/server cmd/shortener/main.go

run:
	go build -o .bin/server cmd/shortener/main.go
	.bin/server $(if $(address),-a $(address)) $(if $(host),-b $(host))

migrate:
	go build -o .bin/migrator cmd/migrator/main.go && .bin/migrator

generate:
	go generate ./...