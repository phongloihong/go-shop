# Docker Development Setup

## Overview

This guide covers the complete Docker-based development environment for Go Shop. The setup uses Docker Compose to orchestrate multiple services with shared infrastructure.

## Prerequisites

- Docker Desktop (latest version)
- Docker Compose V2
- Make (optional, for shortcuts)
- Git

## Architecture

### Service Layout
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   user-service  │    │ product-service │    │  order-service  │
│     :8080       │    │     :8081       │    │     :8082       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
         ┌───────────────────────┼───────────────────────┐
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   PostgreSQL    │    │      Redis      │    │      NATS       │
│     :5432       │    │     :6379       │    │     :4222       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Database Strategy
- **Single PostgreSQL Instance**: Multiple databases (user_db, product_db, order_db)
- **Single Redis Instance**: Separate DB numbers per service (0, 1, 2)
- **Shared NATS**: Single message broker for all services

## Quick Start

### 1. Initial Setup
```bash
# Clone repository
git clone <repository-url>
cd go-shop

# Start all services
make dev

# Verify everything is running
make health
```

### 2. Service Access
- **User Service**: http://localhost:8080
- **PostgreSQL**: localhost:5432 (postgres/password)
- **Redis**: localhost:6379
- **NATS**: localhost:4222, Management: localhost:8222

## Development Commands

### Service Management
```bash
# Start all services
make dev

# Start infrastructure only
make dev-infra

# Start specific service
make dev-user

# Stop all services
make stop

# Restart services
make restart
make restart-user    # Restart specific service

# Clean everything (remove containers/volumes)
make clean
```

### Development Utilities
```bash
# View logs
make logs           # All services
make logs-user      # User service only
make logs-db        # Database only

# Access containers
make shell          # User service shell
make shell-db       # PostgreSQL shell
make shell-redis    # Redis CLI

# Health checks
make health         # Check all services
make ps             # Show running containers
```

### Code Generation
```bash
# Generate protobuf files
make proto

# Generate SQLC code
make sqlc

# Generate all code
make gen
```

### Database Operations
```bash
# Apply migrations
make migrate-up

# Rollback migrations  
make migrate-down

# Create new migration
make migrate-create NAME=add_user_table

# Reset database completely
make db-reset
```

## Development Features

### Hot Reload
- **Technology**: Air (Go hot reload tool)
- **Behavior**: Automatic rebuild and restart on code changes
- **Configuration**: `.air.toml` in each service directory
- **Exclusions**: Test files, vendor directory, generated code

### Volume Mounting
```yaml
# Source code is mounted for hot reload
volumes:
  - ./services/user-service:/app
  - /app/tmp              # Exclude build artifacts
  - /app/vendor           # Exclude vendor directory
```

### Code Generation Workflow
1. **Write SQL queries** in `queries/` directory
2. **Run** `make sqlc` to generate type-safe Go code  
3. **Write proto definitions** in `external/proto/`
4. **Run** `make proto` to generate Go and Connect code
5. **Code automatically reloads** in development container

## Environment Configuration

### Default Values
```bash
# Database (shared PostgreSQL)
DATABASE_HOST=postgres
DATABASE_PORT=5432
DATABASE_USER=postgres  
DATABASE_PASSWORD=password

# Service-specific databases
USER_SERVICE_DB=user_db
PRODUCT_SERVICE_DB=product_db
ORDER_SERVICE_DB=order_db

# Redis (shared instance)
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=""

# Service-specific Redis DBs
USER_SERVICE_REDIS_DB=0
PRODUCT_SERVICE_REDIS_DB=1
ORDER_SERVICE_REDIS_DB=2

# NATS (shared message broker)
NATS_URL=nats://nats:4222
```

### Customization
Create `.env` file in root directory to override defaults:
```bash
# .env
DATABASE_PASSWORD=custom_password
REDIS_PASSWORD=custom_redis_password
```

## Troubleshooting

### Common Issues

#### Container Won't Start
```bash
# Check logs for specific service
make logs-user

# Rebuild from scratch
make clean-build

# Check Docker resources
docker system df
docker system prune -f
```

#### Database Connection Issues
```bash
# Check PostgreSQL health
make shell-db
# Inside container: \l to list databases

# Reset database if corrupted
make db-reset
```

#### Hot Reload Not Working
```bash
# Check Air configuration
make shell
cat .air.toml

# Restart service to reload config
make restart-user
```

#### Build Failures
```bash
# Check if vendor directory conflicts
rm -rf services/user-service/vendor
make restart-user

# Check Go module issues
make shell
go mod tidy
```

#### Port Conflicts
```bash
# Check what's using ports
lsof -i :8080
lsof -i :5432

# Stop conflicting services
sudo lsof -ti:8080 | xargs sudo kill -9
```

### Performance Optimization

#### Volume Performance (macOS)
```yaml
# Use delegated mounting for better performance
volumes:
  - ./services/user-service:/app:delegated
```

#### Build Cache
```bash
# Prune build cache if too large
docker builder prune

# Check build cache usage
docker system df
```

## Advanced Configuration

### Adding New Service
1. **Create service directory**: `services/new-service/`
2. **Add Dockerfile**: Follow user-service pattern
3. **Update docker-compose.yml**: Add new service definition
4. **Update Makefile**: Add service-specific commands
5. **Create database**: Add to `POSTGRES_MULTIPLE_DATABASES`
6. **Assign Redis DB**: Use next available DB number

### Custom Development Tools
```bash
# Add tools to Dockerfile
RUN go install github.com/custom/tool@latest

# Or install in running container
make shell
go install github.com/custom/tool@latest
```

### Network Configuration
```yaml
# All services use shared network
networks:
  go-shop-network:
    driver: bridge
    name: go-shop-network
```

Services can communicate using service names as hostnames (e.g., `postgres`, `redis`, `user-service`).

## Best Practices

### Development Workflow
1. **Start with infrastructure**: `make dev-infra`
2. **Add services incrementally**: `make dev-user`
3. **Use health checks**: `make health` before development
4. **Monitor logs actively**: `make logs-user` in separate terminal
5. **Clean up regularly**: `make clean` to free resources

### Code Generation
1. **Generate before coding**: Run `make gen` after schema changes
2. **Commit generated code**: Include in version control
3. **Validate generation**: Check generated files compile correctly

### Database Management  
1. **Use migrations**: Never modify schema directly
2. **Test migrations**: Both up and down directions
3. **Backup before reset**: `make db-reset` loses all data
4. **Use transactions**: Wrap multiple operations

### Performance
1. **Monitor resource usage**: `make top`
2. **Clean build cache**: `docker builder prune` periodically
3. **Limit log retention**: Configure in docker-compose.yml
4. **Use .dockerignore**: Exclude unnecessary files from context