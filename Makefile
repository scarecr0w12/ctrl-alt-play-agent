.PHONY: build clean test lint install run docker-build docker-run release help dev

# Variables
BINARY_NAME=ctrl-alt-play-agent
VERSION=$(shell cat VERSION)
BINARY_PATH=./bin/$(BINARY_NAME)
DOCKER_IMAGE=ctrl-alt-play-agent
MAIN_PATH=./cmd/agent
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Default target
all: clean lint test build

# Help target
help:
	@echo "Available targets:"
	@echo "  build          - Build the binary"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  lint           - Run linter"
	@echo "  install        - Install dependencies"
	@echo "  run            - Build and run the application"
	@echo "  dev            - Run in development mode"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  release        - Build releases for multiple platforms"
	@echo "  clean-all      - Clean everything including Docker images"

# Build the binary
build:
	@echo "Building $(BINARY_NAME) v$(VERSION)..."
	@mkdir -p bin
	go build $(LDFLAGS) -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Binary built at $(BINARY_PATH)"

# Development build and run
dev:
	@echo "Running in development mode..."
	go run $(MAIN_PATH)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/ build/ dist/
	go clean

# Clean everything including Docker images
clean-all: clean
	@echo "Cleaning Docker images..."
	@docker rmi $(DOCKER_IMAGE):latest $(DOCKER_IMAGE):$(VERSION) 2>/dev/null || true

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p build
	go test -coverprofile=build/coverage.out ./...
	go tool cover -html=build/coverage.out -o build/coverage.html
	@echo "Coverage report generated at build/coverage.html"

# Lint the code
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, running go vet instead"; \
		go vet ./...; \
		echo "Install golangci-lint for better linting: https://golangci-lint.run/usage/install/"; \
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

# Build for multiple platforms (release builds)
release:
	@echo "Building release binaries v$(VERSION)..."
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-$(VERSION)-linux-amd64 $(MAIN_PATH)
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-$(VERSION)-linux-arm64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-$(VERSION)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-$(VERSION)-darwin-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-$(VERSION)-windows-amd64.exe $(MAIN_PATH)
	@echo "Release binaries built in dist/ directory"
	@ls -la dist/

# Build Docker image
docker-build:
	@echo "Building Docker image $(DOCKER_IMAGE):$(VERSION)..."
	docker build -t $(DOCKER_IMAGE):$(VERSION) -t $(DOCKER_IMAGE):latest .
	@echo "Docker image built: $(DOCKER_IMAGE):$(VERSION)"

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run --rm -it \
		-p 8081:8081 \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-e PANEL_URL=ws://localhost:8080 \
		-e NODE_ID=docker-agent \
		-e AGENT_SECRET=docker-secret \
		$(DOCKER_IMAGE):latest

# Push Docker image
docker-push: docker-build
	@echo "Pushing Docker image..."
	docker push $(DOCKER_IMAGE):$(VERSION)
	docker push $(DOCKER_IMAGE):latest

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Security scan
security:
	@echo "Running security scan..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not found. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Generate documentation
docs:
	@echo "Generating documentation..."
	@if command -v godoc >/dev/null 2>&1; then \
		echo "Starting godoc server on http://localhost:6060"; \
		godoc -http=:6060; \
	else \
		echo "godoc not found. Install with: go install golang.org/x/tools/cmd/godoc@latest"; \
	fi

# Initialize project (for new developers)
init: install
	@echo "Initializing project..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2; \
	fi
	@echo "Project initialized successfully!"

# Check for updates
check-updates:
	@echo "Checking for dependency updates..."
	go list -u -m all

# Integration test
test-integration:
	@echo "Running integration tests..."
	@if [ -f "tests/integration_test.go" ]; then \
		go test -tags=integration ./tests/...; \
	else \
		echo "No integration tests found"; \
	fi
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
