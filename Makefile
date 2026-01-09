.PHONY: help install build run dev test clean

help: ## Show this help message
	@echo "iPaaS Platform - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## Install all dependencies
	@echo "Installing Go dependencies..."
	go mod download
	@echo "Installing frontend dependencies..."
	cd frontend && npm install
	@echo "✅ All dependencies installed!"

build: ## Build the backend binary
	@echo "Building backend..."
	go build -o bin/api cmd/api/main.go
	@echo "✅ Backend built successfully: bin/api"

run: build ## Build and run the backend
	@echo "Starting backend on http://localhost:8080..."
	./bin/api

dev: ## Run backend in development mode with auto-reload
	@echo "Starting backend in dev mode..."
	go run cmd/api/main.go

frontend: ## Start the frontend development server
	@echo "Starting frontend on http://localhost:3000..."
	cd frontend && npm run dev

test-data: ## Generate test data
	@echo "Generating test data..."
	go run scripts/generate_test_data.go
	@echo "✅ Test data created!"
	@echo "Login with: demo@ipaas.com / password123"

clean: ## Clean build artifacts and database
	@echo "Cleaning..."
	rm -f bin/api
	rm -f ipaas.db
	rm -rf frontend/.next
	rm -rf frontend/node_modules
	@echo "✅ Cleaned!"

reset: clean install ## Clean everything and reinstall

start: ## Start both backend and frontend (requires two terminals)
	@echo "Use './start.sh' to start both services in one terminal"
	@echo "Or run 'make dev' in one terminal and 'make frontend' in another"

docker-build: ## Build Docker image (TODO: create Dockerfile)
	@echo "Docker support coming soon!"

.DEFAULT_GOAL := help

