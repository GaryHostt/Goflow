#!/bin/bash
# GoFlow Test Suite Runner
# Run this script to execute all tests

set -e  # Exit on error

echo "🚀 GoFlow Test Suite Runner"
echo "============================"
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Step 1: Check if Docker is running
echo "📋 Step 1: Checking Docker..."
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker is not running. Please start Docker Desktop.${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Docker is running${NC}"
echo ""

# Step 2: Start the platform if not already running
echo "📋 Step 2: Starting platform services..."
if docker compose ps | grep -q "Up"; then
    echo -e "${YELLOW}⚠️  Services already running${NC}"
else
    echo "Starting services with docker compose up -d..."
    docker compose up -d
    echo ""
    echo "⏳ Waiting 60 seconds for services to be healthy..."
    sleep 60
fi
echo -e "${GREEN}✅ Platform is running${NC}"
echo ""

# Step 3: Check service health
echo "📋 Step 3: Checking service health..."
echo ""
docker compose ps
echo ""

# Step 4: Configure Kong ELK integration
echo "📋 Step 4: Configuring Kong ELK integration..."
if [ -f "./scripts/configure_kong_elk.sh" ]; then
    chmod +x ./scripts/configure_kong_elk.sh
    ./scripts/configure_kong_elk.sh
    echo -e "${GREEN}✅ Kong ELK integration configured${NC}"
else
    echo -e "${YELLOW}⚠️  Kong ELK script not found, skipping...${NC}"
fi
echo ""

# Step 5: Run connector tests
echo "📋 Step 5: Running Connector Tests..."
echo "======================================"
echo ""
cd "$(dirname "$0")/.."
go run scripts/validate_connectors.go
echo ""
echo -e "${GREEN}✅ Connector tests complete${NC}"
echo ""

# Step 6: Run Kong Gateway tests
echo "📋 Step 6: Running Kong Gateway Tests..."
echo "========================================="
echo ""
go run scripts/validate_kong.go
echo ""
echo -e "${GREEN}✅ Kong Gateway tests complete${NC}"
echo ""

# Step 7: Summary
echo "🎉 Test Suite Complete!"
echo "======================="
echo ""
echo "📊 View results in Kibana:"
echo "   http://localhost:5601"
echo ""
echo "🔍 Create data views:"
echo "   - connector-tests-*"
echo "   - kong-logs-*"
echo ""
echo "✅ All tests completed successfully!"

