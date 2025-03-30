.PHONY: build run test clean docker-build docker-run

# Binary name
BINARY_NAME=fetcher

# Build directory
BUILD_DIR=bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Main build target
build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/fetcher

# Run the application
run:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/fetcher
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./...

# Clean build files
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Build Docker image
docker-build:
	docker build -t codex-online-fetcher .

# Run Docker container
docker-run:
	docker run -p 8080:8080 --env-file .env codex-online-fetcher
