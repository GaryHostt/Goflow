#!/bin/bash

echo "🚀 Starting Simple iPaaS Platform..."
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ from https://go.dev/dl/"
    exit 1
fi

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 18+ from https://nodejs.org/"
    exit 1
fi

echo "✅ Go version: $(go version)"
echo "✅ Node version: $(node --version)"
echo ""

# Install Go dependencies
echo "📦 Installing Go dependencies..."
go mod download
if [ $? -ne 0 ]; then
    echo "❌ Failed to download Go dependencies"
    exit 1
fi

# Build backend
echo "🔨 Building backend..."
go build -o bin/api cmd/api/main.go
if [ $? -ne 0 ]; then
    echo "❌ Failed to build backend"
    exit 1
fi

# Install frontend dependencies if not already installed
if [ ! -d "frontend/node_modules" ]; then
    echo "📦 Installing frontend dependencies..."
    cd frontend
    npm install
    if [ $? -ne 0 ]; then
        echo "❌ Failed to install frontend dependencies"
        exit 1
    fi
    cd ..
fi

# Start backend in background
echo "🚀 Starting backend on http://localhost:8080..."
./bin/api &
BACKEND_PID=$!

# Wait for backend to start
sleep 2

# Start frontend
echo "🚀 Starting frontend on http://localhost:3000..."
cd frontend
npm run dev &
FRONTEND_PID=$!

echo ""
echo "✅ iPaaS Platform is running!"
echo ""
echo "📍 Backend:  http://localhost:8080"
echo "📍 Frontend: http://localhost:3000"
echo ""
echo "Press Ctrl+C to stop all services"
echo ""

# Trap Ctrl+C to kill both processes
trap "kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit" INT

# Wait for processes
wait

