# Simple Makefile for Go API

APP_NAME=pz3
PORT?=8080

# Run the application
run:
	PORT=$(PORT) go run ./cmd/server/

# Build binary
build:
	go build -o bin/$(APP_NAME) ./cmd/server/

# Run tests
test:
	go test -v ./...

# Help
help:
	@echo "Available commands:"
	@echo "  make run    - Run application (PORT=3000 make run for custom port)"
	@echo "  make build  - Build binary"
	@echo "  make test   - Run tests"

.PHONY: run build test help