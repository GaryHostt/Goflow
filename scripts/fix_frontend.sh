#!/bin/bash

# GoFlow Frontend Fix Script
# Fixes "Failed to fetch" error by rebuilding frontend with correct API URL

set -e  # Exit on error

cd /Users/alex.macdonald/simple-ipass

echo "════════════════════════════════════════════════"
echo "🔧 GoFlow Frontend Fix Script"
echo "════════════════════════════════════════════════"
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker Desktop and try again."
    exit 1
fi

echo "✅ Docker is running"
echo ""

# Step 1: Stop frontend
echo "📋 Step 1/5: Stopping frontend container..."
docker compose stop frontend
echo "✅ Frontend stopped"
echo ""

# Step 2: Remove frontend container
echo "📋 Step 2/5: Removing old frontend container..."
docker compose rm -f frontend
echo "✅ Container removed"
echo ""

# Step 3: Rebuild frontend (no cache)
echo "📋 Step 3/5: Rebuilding frontend (this takes 2-3 minutes)..."
echo "   Building with API URL: http://localhost:8080/api"
docker compose build --no-cache --pull frontend
echo "✅ Frontend rebuilt"
echo ""

# Step 4: Start frontend
echo "📋 Step 4/5: Starting frontend..."
docker compose up -d frontend
echo "✅ Frontend started"
echo ""

# Step 5: Wait and verify
echo "📋 Step 5/5: Waiting for frontend to be ready (30 seconds)..."
for i in {30..1}; do
    echo -ne "   ⏳ $i seconds remaining...\r"
    sleep 1
done
echo ""
echo ""

# Check status
echo "📊 Container Status:"
docker compose ps frontend
echo ""

# Check if frontend is accessible
echo "🔍 Testing frontend accessibility..."
if curl -s -o /dev/null -w "%{http_code}" http://localhost:3000 | grep -q "200"; then
    echo "✅ Frontend is accessible!"
else
    echo "⚠️  Frontend may still be starting up..."
fi

echo ""
echo "════════════════════════════════════════════════"
echo "✅ Fix Complete!"
echo "════════════════════════════════════════════════"
echo ""
echo "📝 Next Steps:"
echo ""
echo "1. Clear your browser cache:"
echo "   - Chrome/Edge: Cmd+Shift+Delete"
echo "   - Safari: Cmd+Option+E"
echo "   - Or open in Incognito/Private window"
echo ""
echo "2. Open your browser:"
echo "   http://localhost:3000"
echo ""
echo "3. Try to register:"
echo "   - Email: newuser@example.com"
echo "   - Password: password123"
echo ""
echo "If you still see 'Failed to fetch':"
echo "  - Check browser console (F12) for errors"
echo "  - Run: docker compose logs frontend"
echo "  - See: COMPREHENSIVE_FIX.md for more solutions"
echo ""
echo "════════════════════════════════════════════════"

