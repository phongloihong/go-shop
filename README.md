# Go Shop - Travel Planning Microservices

A **Clean Architecture** microservices application for collaborative travel planning with voting mechanisms, expense management, and itinerary creation.

## Quick Start with Docker

```bash
# Start all services
make dev

# View logs
make logs

# Access services
# User Service: http://localhost:8080
# Database: localhost:5432
# Redis: localhost:6379
```

## Table of Contents

- [Docker Development Setup](#docker-development-setup)
- [Services](#services)
- [Architecture](#architecture)
- [Development](#development)
- [Documentation](#documentation)

## Docker Development Setup

### Prerequisites

- Docker and Docker Compose
- Make (optional, for shortcuts)

### Available Commands

```bash
make dev          # Start all services
make dev-infra    # Start infrastructure only (postgres, redis, nats)
make dev-user     # Start user service only
make stop         # Stop all services
make logs         # View all logs
make shell        # Access user service container
make proto        # Generate protobuf files
make migrate-up   # Run database migrations
```

### Architecture Overview

**Infrastructure (Shared)**

- PostgreSQL: Single instance with multiple databases (user_db, product_db)
- Redis: Single instance with separate DB numbers per service
- NATS: Message broker for inter-service communication

**Services**

- **user-service** (Port 8080): User management, authentication, profiles
- **product-service** (Port 8081): Coming soon

## Services

### User Service

- **Status**: ✅ Active Development
- **Port**: 8080
- **Database**: user_db
- **Features**: User registration, authentication, profile management
- **Documentation**: [User Service Docs](services/user-service/docs/README.md)

### Product Service

- **Status**: 🔄 Planned
- **Port**: 8081
- **Database**: product_db

## Architecture

This project follows **Clean Architecture** principles with strict dependency rules:

1. **Domain Layer**: Business entities and value objects
2. **Application Layer**: Use cases and business logic orchestration
3. **Adapter Layer**: HTTP/gRPC handlers and repository interfaces
4. **Infrastructure Layer**: Database, cache, external service implementations

### Technology Stack

- **Language**: Go 1.24.2
- **Database**: PostgreSQL with pgx/v5
- **Cache**: Redis
- **API**: Echo framework + Connect-Go (planned)
- **Code Generation**: SQLC for type-safe SQL, Buf for protobuf
- **Migration**: golang-migrate
- **Messaging**: NATS (planned)

## Development

### Local Development (Docker)

```bash
# Start development environment
make dev

# Generate code
make proto  # Generate protobuf
make sqlc   # Generate SQLC code

# Database operations
make migrate-up     # Apply migrations
make shell-db      # Access database
make db-reset      # Reset database

# Development utilities
make shell         # Access service container
make logs-user     # View user service logs
make health        # Check service health
```

### Hot Reload

Code changes are automatically detected and the service restarts using Air. The container mounts the entire source code directory for immediate feedback.

### Proto Generation

```bash
# Access container and generate
make shell
cd external
buf generate
```

## Documentation

- [Docker Setup Guide](README-Docker.md)
- [System Architecture](docs/README.md)
- [User Service Documentation](services/user-service/docs/README.md)

## Project Structure

```
go-shop/
├── docker-compose.yml           # Main orchestration
├── Makefile                     # Development shortcuts
├── scripts/                     # Database and utility scripts
├── docs/                        # System-wide documentation
└── services/
    ├── user-service/           # User management service
    │   ├── docs/              # Service-specific documentation
    │   ├── docker/            # Service Dockerfile
    │   └── internal/          # Clean architecture layers
    └── service-product/       # Future product service
```

## Contributing

1. Follow Clean Architecture patterns
2. Use provided Makefile commands for development
3. Run tests and migrations before commits
4. Update documentation for new features

## Getting Help

- Check service logs: `make logs`
- Access service shell: `make shell`
- Check service health: `make health`
- View all commands: `make help`
