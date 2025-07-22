# Go Shop System Documentation

## System Overview

Go Shop is a microservices-based travel planning application built with Clean Architecture principles. The system enables collaborative trip planning with voting mechanisms, expense tracking, and itinerary management.

## Architecture

### Service Map

- **user-service** (Port 8080): User management, authentication, profile handling
- **product-service** (Port 8081): Planned - Product and service catalog management
- **order-service** (Planned): Order processing and booking management

### Shared Infrastructure

- **PostgreSQL**: Centralized database with separate schemas per service
- **Redis**: Shared cache layer with DB partitioning
- **NATS**: Message broker for inter-service communication

## Quick Navigation

### Setup & Development
- [Docker Development Setup](../README-Docker.md)
- [Local Development Guide](setup/local-development.md)
- [Environment Configuration](setup/environment.md)

### Architecture Documentation
- [System Architecture](architecture/system-design.md)
- [Service Communication](architecture/service-communication.md)
- [Database Design](architecture/database-architecture.md)

### Service Documentation
- [User Service](../services/user-service/docs/README.md)
- [Service Overview](services-overview.md)

### Deployment
- [Docker Deployment](deployment/docker.md)
- [Production Setup](deployment/production.md)

## Development Workflow

### Getting Started
1. **Start Development Environment**
   ```bash
   make dev
   ```

2. **Verify Services**
   ```bash
   make health
   ```

3. **Access Services**
   - User Service: http://localhost:8080
   - Database: localhost:5432
   - Redis: localhost:6379

### Common Tasks
- **Code Generation**: `make gen`
- **Database Migration**: `make migrate-up`
- **View Logs**: `make logs`
- **Access Shell**: `make shell`

### Development Tools
- **Hot Reload**: Automatic code recompilation using Air
- **Proto Generation**: buf for protocol buffer compilation
- **Type-Safe SQL**: SQLC for database query generation
- **Migration Management**: golang-migrate for schema changes

## Technology Stack

### Core Technologies
- **Language**: Go 1.24.2
- **Database**: PostgreSQL 16 with pgx/v5 driver
- **Cache**: Redis 7
- **API Framework**: Echo v4 + Connect-Go (planned)
- **Message Broker**: NATS JetStream

### Development Tools
- **Containerization**: Docker & Docker Compose
- **Code Generation**: SQLC, Buf
- **Hot Reload**: Air
- **Migration**: golang-migrate
- **Configuration**: Viper

### Architecture Patterns
- **Clean Architecture**: Strict dependency inversion
- **Domain-Driven Design**: Rich domain models
- **Repository Pattern**: Interface-based data access
- **Use Case Pattern**: Single-responsibility business logic

## Project Structure

```
go-shop/
├── docs/                        # System documentation
├── services/                    # Microservices
│   ├── user-service/           # User management
│   └── service-product/        # Product catalog (planned)
├── scripts/                    # Utility scripts
├── docker-compose.yml          # Development orchestration
├── Makefile                    # Development commands
└── README.md                   # Quick start guide
```

## Contributing

### Code Standards
1. Follow Clean Architecture principles
2. Use provided Makefile commands
3. Write comprehensive tests
4. Update documentation for new features
5. Include proper error handling

### Development Process
1. Create feature branch from main
2. Implement changes with tests
3. Run `make test` and `make lint`
4. Apply database migrations if needed
5. Update relevant documentation
6. Submit pull request

## Troubleshooting

### Common Issues
- **Container Issues**: `make clean && make dev`
- **Database Problems**: `make db-reset`
- **Hot Reload Not Working**: Check `.air.toml` configuration
- **Build Failures**: `make clean-build`

### Getting Help
- Check service logs: `make logs`
- Access service shell: `make shell`
- Health check: `make health`
- Review documentation: Browse `/docs` structure