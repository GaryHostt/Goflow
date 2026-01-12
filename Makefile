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

docker-up: ## Start all services with Docker Compose (PostgreSQL, Backend, Frontend, ELK)
	@echo "Starting iPaaS platform with Docker Compose..."
	docker-compose up -d
	@echo "✅ Platform is running!"
	@echo "Frontend: http://localhost:3000"
	@echo "Backend API: http://localhost:8080"
	@echo "Kibana (logs): http://localhost:5601"

docker-down: ## Stop all Docker services
	@echo "Stopping Docker services..."
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f

docker-build: ## Build Docker images
	@echo "Building Docker images..."
	docker-compose build

docker-clean: ## Remove Docker containers, volumes, and images
	@echo "Cleaning Docker resources..."
	docker-compose down -v
	docker system prune -f

test: ## Run unit tests (fast, uses MockStore)
	@echo "Running unit tests with MockStore..."
	go test ./internal/engine/... -v -count=1

test-integration: ## Run E2E integration tests
	@echo "Running end-to-end integration tests..."
	go test ./scripts/e2e_test.go -v

test-elk: ## Run E2E tests with ELK validation
	@echo "Running E2E tests with Elasticsearch validation..."
	ELASTICSEARCH_URL=http://localhost:9200 go test ./scripts/e2e_test.go -v

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test ./internal/... -v -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

test-bench: ## Run performance benchmarks
	@echo "Running benchmarks..."
	go test ./internal/engine/... -bench=. -benchmem

test-clean: ## Clean up test databases
	@echo "Cleaning test databases and coverage files..."
	rm -f ipaas_test.db test_*.db coverage.out coverage.html
	@echo "✅ Test artifacts cleaned!"

test-all: test test-integration test-coverage ## Run all tests with coverage

test-connectors: ## Run comprehensive connector validation tests
	@echo "Running connector validation tests..."
	go run scripts/connector_test.go

test-kong: ## Run Kong Gateway integration tests
	@echo "Running Kong Gateway integration tests..."
	@echo "⚠️  Make sure Kong is running: make docker-up"
	go run scripts/kong_test.go

test-full: test test-connectors test-kong ## Run all tests including connectors and Kong

configure-kong-elk: ## Configure Kong to ship logs to ELK
	@echo "Configuring Kong to send logs to Elasticsearch..."
	./scripts/configure_kong_elk.sh

.DEFAULT_GOAL := help

