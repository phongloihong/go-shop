# Development Commands

## Docker Development Workflow (Primary)

### Service Management
```bash
make dev                # Start all services với hot reload
make dev-infra         # Start chỉ infrastructure (postgres, redis, nats)  
make dev-user          # Start chỉ user service
make stop              # Stop all services
make restart           # Restart all services
make restart-user      # Restart only user service
make clean             # Stop và remove containers/volumes
```

### Code Generation
```bash
make proto             # Generate protobuf files
make sqlc              # Generate SQLC code
make gen               # Generate both proto và sqlc
```

### Database Operations
```bash
make migrate-up        # Apply database migrations
make migrate-down      # Rollback migrations
make migrate-create NAME=migration_name  # Create new migration
make db-reset          # Reset database hoàn toàn
```

### Development Utilities
```bash
make shell-user        # Access user service container shell
make shell-db          # Access PostgreSQL shell
make shell-redis       # Access Redis CLI
make logs              # View all service logs
make logs-user         # View user service logs only
make health            # Check all service health
```

### Code Quality
```bash
make test              # Run tests
make test-coverage     # Run tests với coverage
make lint              # Run linter (nếu có)
make fmt               # Format Go code
make vet               # Run go vet
```

## Local Development (Alternative)
```bash
cd services/user-service
make gen-query         # Generate SQLC code locally
go build ./cmd/main.go # Build service
go run cmd/main.go     # Run service
go mod tidy            # Clean dependencies
```

## Environment Variables
- Docker: Configured trong docker-compose.yml
- Local: Configure trong shell hoặc .env file
- Database: postgresql://postgres:password@postgres:5432/user_db