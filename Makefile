# Go Shop Development Makefile

.PHONY: help dev dev-infra dev-user stop clean build logs shell proto migrate test lint

# Default target
help: ## Show this help message
	@echo "Go Shop Development Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development environment
dev: ## Start all services for development
	docker-compose up -d

dev-infra: ## Start only infrastructure (postgres, redis, nats)
	docker-compose up -d postgres redis nats

dev-user: ## Start only user service
	docker-compose up -d user-service

# Service management
stop: ## Stop all services
	docker-compose down

restart: ## Restart all services
	docker-compose restart

restart-user: ## Restart user service only
	docker-compose restart user-service

# Cleanup
clean: ## Stop and remove all containers, networks, and volumes
	docker-compose down -v --remove-orphans
	docker system prune -f

clean-build: ## Rebuild all services from scratch
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

# Logs
logs: ## Show logs for all services
	docker-compose logs -f

logs-user: ## Show logs for user service
	docker-compose logs -f user-service

logs-db: ## Show logs for database
	docker-compose logs -f postgres

# Development utilities
shell-user: ## Open shell in user service container
	docker-compose exec user-service /bin/sh

shell-db: ## Open psql shell in database
	docker-compose exec postgres psql -U postgres -d user_db

shell-redis: ## Open redis-cli shell
	docker-compose exec redis redis-cli

# Code generation
proto: ## Generate protobuf files
	docker-compose exec user-service sh -c "cd external && buf generate"

sqlc: ## Generate SQLC code
	docker-compose exec user-service make gen-query

gen: proto sqlc ## Generate all code (proto + sqlc)

# Dependency management
tidy: ## Clean up go modules (current: user-service only)
	docker-compose exec user-service go mod tidy

tidy-user: ## Clean up user service modules
	docker-compose exec user-service go mod tidy

download: ## Download go modules for user service
	docker-compose exec user-service go mod download

# Database operations
migrate-up: ## Run database migrations up
	docker-compose exec user-service migrate -path ./internal/infrastructure/database/postgres/migrations -database "postgresql://postgres:password@postgres:5432/user_db?sslmode=disable" up

migrate-down: ## Run database migrations down
	docker-compose exec user-service migrate -path ./internal/infrastructure/database/postgres/migrations -database "postgresql://postgres:password@postgres:5432/user_db?sslmode=disable" down

migrate-create: ## Create new migration file (usage: make migrate-create NAME=create_users_table)
	docker-compose exec user-service migrate create -ext sql -dir ./internal/infrastructure/database/postgres/migrations $(NAME)

# Testing
test: ## Run tests in user service
	docker-compose exec user-service go test ./...

test-coverage: ## Run tests with coverage
	docker-compose exec user-service go test -cover ./...

# Code quality
lint: ## Run linter (if available)
	docker-compose exec user-service sh -c "command -v golangci-lint >/dev/null 2>&1 && golangci-lint run || echo 'golangci-lint not installed'"

fmt: ## Format Go code
	docker-compose exec user-service go fmt ./...

vet: ## Run go vet
	docker-compose exec user-service go vet ./...

# Build
build: ## Build user service binary
	docker-compose exec user-service go build -o bin/user-service ./cmd/main.go

# Health checks
health: ## Check health of all services
	@echo "Checking service health..."
	@echo "PostgreSQL:"
	@docker-compose exec postgres pg_isready -U postgres || echo "PostgreSQL not ready"
	@echo "Redis:"
	@docker-compose exec redis redis-cli ping || echo "Redis not ready"
	@echo "User Service:"
	@curl -sf http://localhost:8080/health >/dev/null && echo "User service healthy" || echo "User service not ready"

# Database utilities
db-reset: ## Reset database (drop and recreate)
	docker-compose exec postgres psql -U postgres -c "DROP DATABASE IF EXISTS user_db;"
	docker-compose exec postgres psql -U postgres -c "CREATE DATABASE user_db;"
	$(MAKE) migrate-up

db-seed: ## Seed database with test data (implement as needed)
	@echo "Database seeding not implemented yet"

# Monitoring
ps: ## Show running containers
	docker-compose ps

top: ## Show container resource usage
	docker stats --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}\t{{.BlockIO}}"
