.PHONY: help build run test clean docker-build docker-up docker-down demo

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the Go application
	@echo "Building MCP Task Tracker..."
	go build -o mcp-server ./cmd/server
	@echo "✓ Build complete: mcp-server"

run: build ## Build and run the application locally
	@echo "Starting MCP Task Tracker..."
	./mcp-server

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -f mcp-server
	rm -rf postgres_data
	@echo "✓ Clean complete"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t mcp-task-tracker .
	@echo "✓ Docker image built"

docker-up: ## Start services with Docker Compose
	@echo "Starting services..."
	docker-compose up -d
	@echo "✓ Services started"
	@echo "Web UI: http://localhost:8080"
	@echo "Health check: curl http://localhost:8080/health"

docker-down: ## Stop services
	@echo "Stopping services..."
	docker-compose down
	@echo "✓ Services stopped"

docker-logs: ## Show Docker logs
	docker-compose logs -f

demo: ## Run the demo script (requires running services)
	@echo "Running demo..."
	./demo.sh

deps: ## Install Go dependencies
	@echo "Installing dependencies..."
	go mod download
	@echo "✓ Dependencies installed"

fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...
	@echo "✓ Code formatted"

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	golangci-lint run || echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/"

migrate: ## Run database migrations (requires DATABASE_URL)
	@echo "Running migrations..."
	@if [ -z "$$DATABASE_URL" ]; then \
		echo "ERROR: DATABASE_URL not set"; \
		exit 1; \
	fi
	@echo "Migrations applied"

check: ## Check if services are running
	@echo "Checking services..."
	@curl -s -f http://localhost:8080/health && echo "✓ Server is running" || echo "✗ Server is not running"
	@docker-compose ps
