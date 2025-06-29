#!/bin/bash

# GoEdu Control Testing Platform - Development Setup Script
# This script sets up the development environment for the backend application

set -e

echo "🚀 Setting up GoEdu development environment..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "📋 Go version: $GO_VERSION"

# Navigate to backend directory
cd "$(dirname "$0")"

# Download dependencies
echo "📦 Downloading Go dependencies..."
go mod download

# Create necessary directories
echo "📁 Creating necessary directories..."
mkdir -p bin tmp logs

# Copy environment template if .env doesn't exist
if [ ! -f .env ]; then
    echo "📄 Creating .env file from template..."
    cp .env.template .env
    echo "⚠️  Please edit .env file with your configuration before running the application"
fi

# Build the application
echo "🔨 Building the application..."
go build -o bin/server ./cmd/server

# Run tests if available
if [ -f "go.mod" ] && go list ./... | grep -q test; then
    echo "🧪 Running tests..."
    go test ./...
else
    echo "ℹ️  No tests found"
fi

echo "✅ Development environment setup complete!"
echo ""
echo "Next steps:"
echo "1. Review and update the .env file with your configuration"
echo "2. Start external dependencies (MongoDB, Redis) if needed"
echo "3. Run the application with: ./bin/server"
echo "4. Or use 'go run ./cmd/server' for development"
echo ""
echo "Health check: http://localhost:8080/health"
echo "Ready check: http://localhost:8080/ready"
echo ""
echo "Happy coding! 🎉"