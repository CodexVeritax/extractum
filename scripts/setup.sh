#!/bin/bash


echo "Creating extractum project structure...."

mkdir -p cmd/fetcher
mkdir -p config
mkdir -p internal/api/models
mkdir -p internal/parser
mkdir -p internal/analyzer
mkdir -p internal/service
mkdir -p pkg/model
mkdir -p examples

# Create basic files to ensure the directory structure is preserved in Git
touch cmd/fetcher/.gitkeep
touch config/.gitkeep
touch internal/api/.gitkeep
touch internal/api/models/.gitkeep
touch internal/parser/.gitkeep
touch internal/analyzer/.gitkeep
touch internal/service/.gitkeep
touch pkg/model/.gitkeep
touch examples/.gitkeep


go mod init github.com/CodexVeritax/extractum

# Install dependencies
echo "Installing dependencies..."
go get github.com/gorilla/mux
go get github.com/sirupsen/logrus

# Create .env.example
echo "Creating .env.example..."
cat > .env.example << 'EOF'
# Required: GitHub Personal Access Token
# Create one at https://github.com/settings/tokens
# Requires 'repo' scope for private repos, 'public_repo' for public repos
GITHUB_TOKEN=your_github_token_here

# Optional: GitHub API URL (change for GitHub Enterprise)
# GITHUB_API_URL=https://api.github.com

# Optional: User agent for API requests
# USER_AGENT=CodexOnlineFetcher/1.0

# Optional: Timeout for API requests (in seconds)
# API_TIMEOUT=30

# Optional: Maximum concurrent API calls
# MAX_CONCURRENT_CALLS=5

# Optional: HTTP server port
# PORT=8080

# Optional: Logging configuration
# LOG_LEVEL=info  # debug, info, warn, error
# LOG_FORMATTER=json  # json or text
EOF

# Create .gitignore
echo "Creating .gitignore..."
cat > .gitignore << 'EOF'
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib
bin/

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories
vendor/

# Go workspace file
go.work

# Environment variables
.env

# IDE files
.idea/
.vscode/
*.swp
*.swo

# OS-specific files
.DS_Store
Thumbs.db
EOF


# Create Makefile
echo "Creating Makefile..."
cat > Makefile << 'EOF'
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
EOF

echo "Project structure initialized successfully!"
echo "Remember to create a .env file from .env.example with your GitHub token"
