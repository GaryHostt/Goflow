#!/bin/bash
# Fix go.sum by regenerating it locally

cd /Users/alex.macdonald/simple-ipass

echo "🔧 Fixing Go dependencies..."
echo ""

# Remove old go.sum
echo "Removing old go.sum..."
rm -f go.sum

# Regenerate go.sum with correct checksums
echo "Running go mod tidy..."
go mod tidy

echo ""
echo "✅ go.sum regenerated!"
echo ""
echo "Now rebuild Docker:"
echo "  docker compose down"
echo "  docker compose up -d --build"

