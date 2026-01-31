#!/bin/bash
set -e

echo "🧪 Running Scentora Backend Tests"
echo "=================================="
echo ""

# Ensure test database exists
echo "📦 Setting up test database..."
docker exec scentora-postgres psql -U admin -d postgres -c "SELECT 1 FROM pg_database WHERE datname='scentora_test'" | grep -q 1 || \
docker exec scentora-postgres psql -U admin -d postgres -c "CREATE DATABASE scentora_test;"

# Set test database URL
export TEST_DATABASE_URL="postgres://admin:password@localhost:5435/scentora_test?sslmode=disable"

echo "✅ Test database ready"
echo ""

# Run tests
echo "🏃 Running tests..."
echo ""

go test -p 1 -v -cover ./internal/config ./internal/repository "$@"

echo ""
echo "✅ All tests passed!"
