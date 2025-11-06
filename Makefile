TEST_FLAGS=-v -coverprofile=coverage.out -covermode=atomic
TEST_PATH=./internal/...

test:
	go test $(TEST_PATH) $(TEST_FLAGS)
	@echo ""
	@echo "===================="
	@echo "Code coverage report"
	@echo "===================="
	go tool cover -func=coverage.out | tail -n 10

run:
	go build -o .bin/server cmd/shortener/main.go && .bin/server