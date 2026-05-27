.PHONY: build run clean test lint help

APP_NAME = gsystes-server
BUILD_DIR = bin
MAIN_PATH = cmd/server

help:
	@echo "Usage:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make test        - Run tests"
	@echo "  make lint        - Run linter"

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

run: build
	./$(BUILD_DIR)/$(APP_NAME) -config config/config.dev.yaml

clean:
	rm -rf $(BUILD_DIR)
	rm -rf logs/

test:
	go test ./... -v -cover

lint:
	golangci-lint run ./...

dev:
	go run $(MAIN_PATH)/main.go -config config/config.dev.yaml