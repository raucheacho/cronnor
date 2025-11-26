.PHONY: build run clean migrate docker-build docker-run docker-stop help

# Variables
BINARY_NAME=cronnor
GO_FILES=$(shell find . -name '*.go')
DB_PATH=./data/cronnor.db
MIGRATION_PATH=./migrations/001_initial_schema.sql

help: ## Show this help message
	@echo "Cronnor - HTTP Cron Job Scheduler"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) ./cmd/server

run: build ## Run the application locally
	@echo "Starting Cronnor..."
	@./$(BINARY_NAME)

clean: ## Clean build artifacts and database
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf data/

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t cronnor:latest .

docker-run: docker-build ## Run with Docker Compose
	@echo "Starting Cronnor with Docker Compose..."
	@docker-compose up -d
	@echo "Cronnor is running at http://localhost:8080"

docker-stop: ## Stop Docker Compose
	@echo "Stopping Cronnor..."
	@docker-compose down

docker-logs: ## View Docker logs
	@docker-compose logs -f

dev: ## Run in development mode with auto-reload (requires air)
	@command -v air > /dev/null 2>&1 || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	@air

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

all: clean deps build ## Clean, download deps, and build
