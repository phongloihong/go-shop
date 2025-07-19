# Development Setup

This document describes how to set up a local development environment for the User Service.

## Prerequisites

### Required Software

- **Go 1.24.2 or higher**
- **PostgreSQL 12 or higher**
- **Redis 6 or higher**
- **Git**

### Development Tools

```bash
# Install SQLC for type-safe SQL generation
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Install golang-migrate for database migrations
brew install golang-migrate
# OR
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Install Air for hot reloading (optional)
go install github.com/cosmtrek/air@latest
```

## Project Setup

### 1. Clone Repository

```bash
git clone https://github.com/phongloihong/go-shop.git
cd go-shop/services/user-service
```

### 2. Install Dependencies

```bash
# Download Go modules
go mod download

# Verify dependencies
go mod tidy
```

### 3. Environment Configuration

Create a `.env` file in the user-service directory:

```bash
# Database
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=password
DATABASE_DB_NAME=user_dev

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Server
SERVER_PORT=8080

# Logging
LOG_LEVEL=debug
```

### 4. Database Setup

```bash
# Start PostgreSQL (if not running)
brew services start postgresql
# OR
sudo systemctl start postgresql

# Create development database
createdb user_dev

# Run database migrations
./scripts/migrate.sh
```

### 5. Redis Setup

```bash
# Start Redis (if not running)
brew services start redis
# OR
sudo systemctl start redis

# Test Redis connection
redis-cli ping
```

### 6. Generate SQL Code

```bash
# Generate type-safe SQL code
make gen-query

# Verify generation
ls internal/infrastructure/database/postgres/sqlc/
```

## Development Workflow

### Running the Service

```bash
# Run the service
go run cmd/main.go

# With hot reloading (if Air is installed)
air

# Build and run
go build -o bin/user-service cmd/main.go
./bin/user-service
```

### Code Generation

```bash
# Generate SQL code after query changes
make gen-query

# Verify SQLC configuration
sqlc vet

# Compile SQL queries without generation
sqlc compile
```

### Database Operations

```bash
# Run new migrations
./scripts/migrate.sh

# Check migration status
migrate -path internal/infrastructure/database/postgres/migrations -database "postgres://postgres:password@localhost:5432/user_dev?sslmode=disable" version

# Rollback last migration
migrate -path internal/infrastructure/database/postgres/migrations -database "postgres://postgres:password@localhost:5432/user_dev?sslmode=disable" down 1
```

## Docker Development

### Using Docker Compose

Create `docker-compose.dev.yaml`:

```yaml
version: '3.8'
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: user_dev
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  postgres_data:
  redis_data:
```

```bash
# Start services
/73/󱂷 /docs
docker-compose -f docker-compose.dev.yaml up -d

# Stop services
docker-compose -f docker-compose.dev.yaml down
```

## IDE Configuration

### VS Code

Recommended extensions:
- Go (by Google)
- PostgreSQL (by Chris Kolkman)  
- Redis (by Dunn)
- YAML (by Red Hat)

### Settings

Create `.vscode/settings.json`:

```json
{
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "package",
  "go.formatTool": "goimports",
  "go.useLanguageServer": true,
  "go.testFlags": ["-v"],
  "files.eol": "\n"
}
```

### Launch Configuration

Create `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch User Service",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/main.go",
      "cwd": "${workspaceFolder}",
      "envFile": "${workspaceFolder}/.env"
    }
  ]
}
```

## Testing Setup

### Unit Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific test package
go test ./internal/domain/entity/...
```

### Integration Tests

```bash
# Run integration tests (when implemented)
go test -tags=integration ./...
```

/73/󱂷 /docs
### Test Database

Create separate test database:

```bash
# Create test database
createdb user_test

# Set test environment
export DATABASE_DB_NAME=user_test

# Run migrations for test database
./scripts/migrate.sh
```

## Code Quality

### Linting

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run

# Run with auto-fix
golangci-lint run --fix
```

### Formatting

```bash
# Format code
go fmt ./...

# Import organization
goimports -w .
```

### Security Scanning

```bash
# Install gosec
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Run security scan
gosec ./...
```

## API Testing

### Using curl

```bash
# Test user creation
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe", 
    "email": "john@example.com",
    "phone": "+1234567890",
    "password": "password123"
  }'

# Test user retrieval
curl http://localhost:8080/api/v1/users/{user-id}
```

### Using HTTPie

```bash
# Install HTTPie
pip install httpie

# Test user creation
http POST localhost:8080/api/v1/users \
  first_name=John \
  last_name=Doe \
  email=john@example.com \
  phone=+1234567890 \
  password=password123
```

## Debugging

### Debugging with Delve

```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug the application
dlv debug cmd/main.go

# Debug tests
dlv test ./internal/domain/entity/
```

### Log Analysis

```bash
# Follow logs with structured output
tail -f logs/user-service.log | jq '.'

# Filter error logs
tail -f logs/user-service.log | jq 'select(.level == "error")'
```

## Performance Monitoring

### Profiling

```bash
# CPU profiling
go tool pprof http://localhost:8080/debug/pprof/profile

# Memory profiling
go tool pprof http://localhost:8080/debug/pprof/heap

# Goroutine profiling
go tool pprof http://localhost:8080/debug/pprof/goroutine
```

## Common Development Tasks

### Adding New SQL Queries

1. Add query to `internal/infrastructure/database/postgres/queries/users.sql`
2. Run `make gen-query` to generate Go code
3. Implement repository method using generated functions
4. Update use case to use new repository method

### Adding New Endpoints

1. Define request/response DTOs in use case layer
2. Implement use case method
3. Add handler method in delivery layer
4. Register route in router
5. Add tests and documentation

### Database Schema Changes

1. Create new migration files
2. Run migration with `./scripts/migrate.sh`
3. Update SQL queries if needed
4. Run `make gen-query` to regenerate code
5. Update repository implementations

## Troubleshooting

### Common Issues

1. **Port already in use:**
   ```bash
   lsof -ti:8080 | xargs kill -9
   ```

2. **Database connection issues:**
   ```bash
   psql -h localhost -U postgres -d user_dev
   ```

3. **Redis connection issues:**
   ```bash
   redis-cli ping
   ```

4. **SQLC generation fails:**
   ```bash
   sqlc vet
   sqlc compile
   ```

5. **Go module issues:**
   ```bash
   go clean -modcache
   go mod download
   ```