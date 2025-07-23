.PHONY: build clean test lint install run docker-build docker-run

# Variables
BINARY_NAME=ctrl-alt-play-agent
BINARY_PATH=./build/$(BINARY_NAME)
DOCKER_IMAGE=ctrl-alt-play-agent
MAIN_PATH=./cmd/agent

# Default target
all: clean lint test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p build
	go build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Binary built at $(BINARY_PATH)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf build/
	go clean

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Lint the code
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, running go vet instead"; \
		go vet ./...; \
	fi

# Install dependencies
install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	$(BINARY_PATH)

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p build
	GOOS=linux GOARCH=amd64 go build -o build/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=linux GOARCH=arm64 go build -o build/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 go build -o build/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 go build -o build/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "Cross-platform binaries built in build/ directory"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

# Run Docker container
docker-run: docker-build
	@echo "Running Docker container..."
	docker run --rm -it \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-e PANEL_URL=${PANEL_URL} \
		-e NODE_ID=${NODE_ID} \
		-e AGENT_SECRET=${AGENT_SECRET} \
		$(DOCKER_IMAGE)

# Development mode - run with file watching (requires 'air' tool)
dev:
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "Install 'air' for live reloading: go install github.com/cosmtrek/air@latest"; \
		make run; \
	fi

# Help
help:
	@echo "Available targets:"
	@echo "  build       - Build the binary"
	@echo "  clean       - Clean build artifacts"
	@echo "  test        - Run tests"
	@echo "  lint        - Run linter"
	@echo "  install     - Install dependencies"
	@echo "  run         - Build and run the application"
	@echo "  build-all   - Build for multiple platforms"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run  - Build and run Docker container"
	@echo "  dev         - Run in development mode with live reloading"
	@echo "  help        - Show this help"
