# GoEdu Control Testing Platform - Makefile
# This file provides convenient development commands

.PHONY: help build test clean run dev setup lint format deps

# Default target
help: ## Show this help message
	@echo "GoEdu Control Testing Platform - Development Commands"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# Build targets
build: ## Build the server binary
	@echo "🔨 Building server..."
	@go build -o bin/server ./cmd/server
	@echo "✅ Server built successfully"

build-all: ## Build all binaries (server, migrate, seed)
	@echo "🔨 Building all binaries..."
	@go build -o bin/server ./cmd/server
	@go build -o bin/migrate ./cmd/migrate
	@go build -o bin/seed ./cmd/seed
	@echo "✅ All binaries built successfully"

# Development targets
dev: ## Run the server in development mode with hot reload
	@echo "🚀 Starting development server..."
	@go run ./cmd/server

run: build ## Build and run the server
	@echo "🚀 Starting server..."
	@./bin/server

# Setup and initialization
setup: ## Set up development environment
	@echo "📦 Setting up development environment..."
	@./scripts/setup-dev.sh

deps: ## Download and tidy dependencies
	@echo "📦 Downloading dependencies..."
	@go mod download
	@go mod tidy

# Database operations
migrate: ## Run database migrations
	@echo "🗄️ Running database migrations..."
	@go run ./cmd/migrate

seed: ## Seed database with development data
	@echo "🌱 Seeding database..."
	@go run ./cmd/seed

db-setup: migrate seed ## Run migrations and seed data

# Testing and quality
test: ## Run tests
	@echo "🧪 Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "🧪 Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

lint: ## Run linter
	@echo "🔍 Running linter..."
	@golangci-lint run

format: ## Format code
	@echo "💅 Formatting code..."
	@go fmt ./...

# Docker targets
docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	@docker build -t goedu/backend:latest .

docker-run: ## Run Docker container
	@echo "🐳 Running Docker container..."
	@docker run -p 8080:8080 --env-file .env goedu/backend:latest

# Utility targets
clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf tmp/
	@rm -f coverage.out coverage.html
	@echo "✅ Clean completed"

health: ## Check application health
	@echo "🏥 Checking application health..."
	@curl -s http://localhost:8080/health | jq '.' || echo "❌ Health check failed"

logs: ## Show application logs (if running in Docker)
	@echo "📋 Showing application logs..."
	@docker logs goedu-backend

# Environment targets
env-copy: ## Copy environment template
	@if [ ! -f .env ]; then \
		echo "📄 Copying environment template..."; \
		cp .env.template .env; \
		echo "⚠️  Please edit .env file with your configuration"; \
	else \
		echo "ℹ️  .env file already exists"; \
	fi

# Release targets
version: ## Show current version
	@echo "📋 Current version:"
	@grep 'version:' config.yaml | awk '{print $$2}' | tr -d '"'

# Development workflow
start: deps env-copy build-all db-setup run ## Complete development setup and start

# Production build
build-prod: ## Build for production
	@echo "🏭 Building for production..."
	@CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o bin/server ./cmd/server
	@echo "✅ Production build completed"

# Help with common tasks
install-tools: ## Install development tools
	@echo "🛠️  Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Development tools installed"

# Variables
BINARY_NAME=server
VERSION := $(shell grep 'version:' config.yaml | awk '{print $$2}' | tr -d '"')
BUILD_TIME := $(shell date +%Y%m%d_%H%M%S)
GIT_COMMIT := $(shell git rev-parse --short HEAD)

# Advanced build with version info
build-version: ## Build with version information
	@echo "🔨 Building with version information..."
	@go build -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)" -o bin/$(BINARY_NAME) ./cmd/server
	@echo "✅ Build completed with version $(VERSION)"