#!/bin/bash

# Scentora Go Backend - Development Script

set -e

BACKEND_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$BACKEND_DIR"

case "$1" in
  start)
    echo "🚀 Starting Scentora Go Backend..."
    
    # Check if .env exists
    if [ ! -f .env ]; then
      echo "⚠️  No .env file found. Creating from .env.example..."
      cp .env.example .env
      echo "✅ Please edit .env with your configuration"
    fi
    
    # Build and run
    echo "🔨 Building..."
    go build -o scentora-backend cmd/server/main.go
    
    echo "▶️  Starting server..."
    ./scentora-backend
    ;;
    
  build)
    echo "🔨 Building Scentora Go Backend..."
    go build -o scentora-backend cmd/server/main.go
    echo "✅ Build complete: ./scentora-backend"
    ;;
    
  dev)
    echo "🔧 Starting in development mode..."
    
    # Check if .env exists
    if [ ! -f .env ]; then
      echo "⚠️  No .env file found. Creating from .env.example..."
      cp .env.example .env
    fi
    
    # Use air for hot reload if available, otherwise go run
    if command -v air &> /dev/null; then
      air
    else
      echo "💡 Tip: Install 'air' for hot reload: go install github.com/cosmtrek/air@latest"
      go run cmd/server/main.go
    fi
    ;;
    
  test)
    echo "🧪 Running tests..."
    go test ./... -v
    ;;
    
  clean)
    echo "🧹 Cleaning build artifacts..."
    rm -f scentora-backend
    echo "✅ Clean complete"
    ;;
    
  db:up)
    echo "🐘 Starting PostgreSQL..."
    cd .. && docker compose up -d postgres
    echo "✅ PostgreSQL started on port 5435"
    ;;
    
  db:down)
    echo "🛑 Stopping PostgreSQL..."
    cd .. && docker compose stop postgres
    echo "✅ PostgreSQL stopped"
    ;;
    
  *)
    echo "Scentora Go Backend - Development Script"
    echo ""
    echo "Usage: $0 {start|build|dev|test|clean|db:up|db:down}"
    echo ""
    echo "Commands:"
    echo "  start    - Build and run the server"
    echo "  build    - Build the binary"
    echo "  dev      - Run in development mode (hot reload if 'air' installed)"
    echo "  test     - Run tests"
    echo "  clean    - Remove build artifacts"
    echo "  db:up    - Start PostgreSQL container"
    echo "  db:down  - Stop PostgreSQL container"
    echo ""
    exit 1
    ;;
esac
