TEST_FLAGS=-v -coverprofile=coverage.out -covermode=atomic
TEST_PATH=./internal/...

address ?=
host ?=

test:
	go test $(TEST_PATH) $(TEST_FLAGS)

build:
	go build -o .bin/server cmd/shortener/main.go

run:
	go build -o .bin/server cmd/shortener/main.go
	.bin/server $(if $(address),-a $(address)) $(if $(host),-b $(host))
