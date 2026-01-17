#!/bin/bash

# 🚀 GoFlow - Start Full Stack with Kong Gateway

set -e

echo "🚀 Starting GoFlow Platform with Kong Gateway..."
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running!"
    echo ""
    echo "Please start Docker Desktop and try again."
    echo ""
    exit 1
fi

echo "✅ Docker is running"
echo ""

# Navigate to project directory
cd "$(dirname "$0")/.."

echo "📦 Starting all services with Docker Compose..."
echo "   This may take 2-3 minutes on first run..."
echo ""

# Start all services
docker compose up -d --build

echo ""
echo "⏳ Waiting for services to be healthy..."
sleep 10

# Check if services are running
BACKEND_STATUS=$(docker compose ps backend --format json | grep -o '"State":"[^"]*"' | cut -d'"' -f4 || echo "not found")
KONG_STATUS=$(docker compose ps kong --format json | grep -o '"State":"[^"]*"' | cut -d'"' -f4 || echo "not found")
FRONTEND_STATUS=$(docker compose ps frontend --format json | grep -o '"State":"[^"]*"' | cut -d'"' -f4 || echo "not found")

echo ""
echo "📊 Service Status:"
echo "   Backend:  $BACKEND_STATUS"
echo "   Kong:     $KONG_STATUS"
echo "   Frontend: $FRONTEND_STATUS"
echo ""

# Wait a bit more for services to be fully ready
echo "⏳ Allowing services to initialize..."
sleep 10

echo ""
echo "✅ GoFlow Platform is starting!"
echo ""
echo "🌐 Access Points:"
echo "   📱 Frontend:        http://localhost:3000"
echo "   🔧 Backend API:     http://localhost:8080"
echo "   🌉 Kong Gateway:    http://localhost:8000"
echo "   ⚙️  Kong Admin:      http://localhost:8001"
echo "   📊 Kibana:          http://localhost:5601"
echo ""
echo "💡 Quick Start:"
echo "   1. Open http://localhost:3000 in your browser"
echo "   2. Click 'Skip Login - Dev Mode' for instant access"
echo "   3. Configure your API connections"
echo "   4. Build your first workflow!"
echo ""
echo "📋 Useful Commands:"
echo "   View logs:          docker compose logs -f"
echo "   Stop services:      docker compose down"
echo "   Restart a service:  docker compose restart <service>"
echo ""
echo "📚 Documentation:"
echo "   START_APP_AND_PROXY.md - Complete startup guide"
echo "   COMPONENT_RUNNING_GUIDE.md - Individual component guide"
echo ""

# Configure Kong (optional, on first run)
if [ "$1" == "--configure-kong" ]; then
    echo "🔧 Configuring Kong Gateway patterns..."
    echo ""
    sleep 5
    ./scripts/configure_kong_elk.sh
    echo ""
    echo "✅ Kong configuration complete!"
    echo ""
fi

echo "🎉 Ready to build integrations!"
echo ""
echo "Need help? Run: docker compose logs -f"
