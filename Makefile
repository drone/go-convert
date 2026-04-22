.PHONY: build run clean test docker-build docker-run docker-stop help

# Variables
SERVICE_NAME := go-convert-service
DOCKER_IMAGE := $(SERVICE_NAME):latest
PORT ?= 8090
LOG_LEVEL ?= debug
MAX_BATCH_SIZE ?= 100
MAX_YAML_BYTES ?= 1048576

# Default target
.DEFAULT_GOAL := help

## help: Display this help message
help:
	@echo "Available targets:"
	@echo ""
	@grep -E '^##' $(MAKEFILE_LIST) | sed 's/## /  /'
	@echo ""

## build: Build the service binary
build:
	@echo "Building $(SERVICE_NAME)..."
	@go build -o $(SERVICE_NAME) ./cmd/server
	@echo "Build complete: $(SERVICE_NAME)"

## run: Build and run the service locally
run: build
	@echo "Starting $(SERVICE_NAME) on port $(PORT)..."
	@PORT=$(PORT) \
	 LOG_LEVEL=$(LOG_LEVEL) \
	 MAX_BATCH_SIZE=$(MAX_BATCH_SIZE) \
	 MAX_YAML_BYTES=$(MAX_YAML_BYTES) \
	 ./$(SERVICE_NAME)

## run-debug: Run the service with debug logging
run-debug:
	@$(MAKE) run LOG_LEVEL=debug

## run-info: Run the service with info logging
run-info:
	@$(MAKE) run LOG_LEVEL=info

## test: Run all tests
test:
	@echo "Running tests..."
	@go test -v ./service/... ./cmd/...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./service/... ./cmd/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## clean: Remove build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f $(SERVICE_NAME)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image: $(DOCKER_IMAGE)..."
	@docker build -f Dockerfile.service -t $(DOCKER_IMAGE) .
	@echo "Docker image built: $(DOCKER_IMAGE)"

## docker-run: Run service in Docker container
docker-run: docker-build
	@echo "Running Docker container on port $(PORT)..."
	@docker run -d \
		--name $(SERVICE_NAME) \
		-p $(PORT):8090 \
		-e LOG_LEVEL=$(LOG_LEVEL) \
		-e MAX_BATCH_SIZE=$(MAX_BATCH_SIZE) \
		-e MAX_YAML_BYTES=$(MAX_YAML_BYTES) \
		$(DOCKER_IMAGE)
	@echo "Container started. Check logs with: make docker-logs"

## docker-run-foreground: Run service in Docker container (foreground)
docker-run-foreground: docker-build
	@echo "Running Docker container on port $(PORT)..."
	@docker run --rm \
		--name $(SERVICE_NAME) \
		-p $(PORT):8090 \
		-e LOG_LEVEL=$(LOG_LEVEL) \
		-e MAX_BATCH_SIZE=$(MAX_BATCH_SIZE) \
		-e MAX_YAML_BYTES=$(MAX_YAML_BYTES) \
		$(DOCKER_IMAGE)

## docker-stop: Stop and remove Docker container
docker-stop:
	@echo "Stopping Docker container..."
	@docker stop $(SERVICE_NAME) || true
	@docker rm $(SERVICE_NAME) || true
	@echo "Container stopped and removed"

## docker-logs: Show Docker container logs
docker-logs:
	@docker logs -f $(SERVICE_NAME)

## docker-shell: Open shell in running container
docker-shell:
	@docker exec -it $(SERVICE_NAME) /bin/sh

## health-check: Check if service is running
health-check:
	@echo "Checking service health..."
	@curl -s http://localhost:$(PORT)/healthz || echo "Service not responding"

## example-request: Send example batch request
example-request:
	@echo "Sending example batch request..."
	@curl -X POST http://localhost:$(PORT)/api/v1/convert/batch \
		-H "Content-Type: application/json" \
		-d @test_batch_request.json | jq

## dev: Run in development mode with live reload (requires air)
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "air not installed. Install with: go install github.com/air-verse/air@latest"; \
		echo "Falling back to normal run..."; \
		$(MAKE) run-debug; \
	fi

## install-tools: Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/air-verse/air@latest
	@echo "Tools installed"

## mod-download: Download Go module dependencies
mod-download:
	@echo "Downloading dependencies..."
	@go mod download
	@echo "Dependencies downloaded"

## mod-tidy: Tidy Go module dependencies
mod-tidy:
	@echo "Tidying dependencies..."
	@go mod tidy
	@echo "Dependencies tidied"

## fmt: Format Go code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted"

## lint: Run linters (requires golangci-lint)
lint:
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install from: https://golangci-lint.run/"; \
	fi

## all: Build, test, and create Docker image
all: clean test build docker-build
	@echo "All tasks complete"
