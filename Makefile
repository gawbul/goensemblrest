.PHONY: all build test test-race test-live test-coverage lint format example clean

# Default target
all: lint test-race build

# Build packages and examples
build:
	go build -v ./...
	go build -v -o bin/example ./examples

# Run offline unit tests
test:
	go test -v ./...

# Run unit tests with race detector
test-race:
	go test -v -race -cover ./...

# Run live tests against Ensembl REST server
test-live:
	ENSEMBL_LIVE_TESTS=1 go test -v -tags=live -race -run TestLiveEnsemblAPI ./...

# Run tests and generate HTML coverage report
test-coverage:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

# Run linter
lint:
	go vet ./...

# Format code
format:
	go fmt ./...

# Run example
example:
	go run ./examples/main.go

# Clean build artifacts
clean:
	rm -rf bin/ coverage.out coverage.html
