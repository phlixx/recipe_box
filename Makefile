.PHONY: all build test vet fmt lint check clean run

# Default target
all: check build

# Build the binary
build:
	go build -o bin/recipe_box .

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run go vet
vet:
	go vet ./...

# Format code
fmt:
	go fmt ./...

# Run staticcheck if installed
lint:
	@which staticcheck > /dev/null && staticcheck ./... || echo "staticcheck not installed, skipping"

# Full verification pipeline (used by CI and agents)
check: fmt vet test
	@echo "✓ All checks passed"

# Clean build artifacts
clean:
	rm -rf bin/ coverage.out coverage.html

# Run the CLI
run:
	go run main.go $(ARGS)

# Quick smoke test
smoke:
	go run main.go --help
	go run main.go config --help
