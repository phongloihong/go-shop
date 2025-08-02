# Docker Development Setup

## Docker Architecture

### Shared Infrastructure

- **PostgreSQL**: Single instance with multiple databases (user_db, product_db, order_db)
- **Redis**: Single instance with multiple db numbers (0: user-service, 1: product-service)
- **NATS**: Message broker for inter-service communication

### Services

- **user-service**: Port 8080 with hot reload
- **product-service**: Port 8081 (commented out, uncomment when needed)

## Commands

### Start development environment

```bash
# From root project
docker-compose up -d

# Infrastructure only
docker-compose up -d postgres redis nats

# User service only
docker-compose up -d user-service
```

### Exec into container to generate proto

```bash
# Enter user-service container
docker-compose exec user-service /bin/sh

# Generate proto inside container
cd external
buf generate

# Generate SQLC
make gen-query

# Database migration
migrate -path ./internal/infrastructure/database/postgres/migrations -database "postgresql://postgres:password@postgres:5432/user_db?sslmode=disable" up
```

### Development workflow

```bash
# Watch logs
docker-compose logs -f user-service

# Restart service
docker-compose restart user-service

# Rebuild and restart
docker-compose up -d --build user-service
```

### Database operations

```bash
# Connect to PostgreSQL
docker-compose exec postgres psql -U postgres -d user_db

# List databases
docker-compose exec postgres psql -U postgres -c "\l"

# Redis CLI
docker-compose exec redis redis-cli
```

## File Structure

```
go-shop/
├── docker-compose.yml           # Main compose file at root
├── scripts/
│   └── create-multiple-databases.sh  # Script to create multiple DBs
└── services/
    └── user-service/
        └── docker/
            └── Dockerfile       # Development Dockerfile with hot reload
```

## Environment Variables

### User Service

- `DATABASE_HOST=postgres`
- `DATABASE_DB_NAME=user_db`
- `REDIS_DB=0`

### Product Service (future)

- `DATABASE_HOST=postgres`
- `DATABASE_DB_NAME=product_db`
- `REDIS_DB=1`

## Development Features

### Hot Reload

- Uses `air` for Go hot reload
- Code changes automatically rebuild and restart
- Excludes test files and proto files

### Proto Generation

- All tools pre-installed in container: `buf`, `protoc-gen-go`, `protoc-gen-connect-go`
- Run `buf generate` inside container to generate proto files

### Database Migration

- `golang-migrate` pre-installed in container
- Migration files mounted from host for development

## Troubleshooting

### Container won't start

```bash
docker-compose logs user-service
docker-compose down && docker-compose up -d
```

### Database connection issues

```bash
# Check postgres health
docker-compose exec postgres pg_isready -U postgres

# Check if databases exist
docker-compose exec postgres psql -U postgres -c "\l"
```

### Hot reload not working

```bash
# Check air config and volumes
docker-compose exec user-service cat .air.toml
```

