#!/bin/bash

# Run Frontend Locally Script
# This bypasses Docker issues and runs the frontend directly

set -e

echo "════════════════════════════════════════════════"
echo "🚀 Running GoFlow Frontend Locally"
echo "════════════════════════════════════════════════"
echo ""

cd /Users/alex.macdonald/simple-ipass/frontend

# Step 1: Create .env.local
echo "📋 Step 1/3: Creating environment configuration..."
cat > .env.local << 'EOF'
NEXT_PUBLIC_API_URL=http://localhost:8080/api
EOF
echo "✅ Environment configured"
echo ""

# Step 2: Check if node_modules exists
if [ ! -d "node_modules" ]; then
    echo "📋 Step 2/3: Installing dependencies (first time, takes 1-2 minutes)..."
    npm install
else
    echo "📋 Step 2/3: Dependencies already installed ✅"
fi
echo ""

# Step 3: Start dev server
echo "📋 Step 3/3: Starting development server..."
echo ""
echo "════════════════════════════════════════════════"
echo "✅ Frontend will start at: http://localhost:3000"
echo "════════════════════════════════════════════════"
echo ""
echo "📝 Press Ctrl+C to stop the server"
echo ""

npm run dev

