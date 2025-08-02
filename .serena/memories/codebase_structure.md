# Codebase Structure

## Root Level Structure
```
go-shop/
├── services/          # Microservices
├── cicd/             # CI/CD configurations  
├── docs/             # Documentation
├── scripts/          # Utility scripts
├── Makefile          # Development commands
├── docker-compose.yml # Development environment
└── README.md         # Project documentation
```

## User Service Structure (Clean Architecture)
```
services/user-service/
├── cmd/                    # Application entry points
│   └── main.go            # Main application
├── internal/
│   ├── domain/            # Business entities, value objects, domain services
│   │   ├── entity/        # Domain entities (User, etc.)
│   │   ├── valueObject/   # Value objects (Email, Password, Phone, etc.)
│   │   └── repository/    # Repository interfaces
│   ├── delivery/          # HTTP và gRPC handlers
│   │   ├── http/         # HTTP handlers và routing
│   │   └── connect/      # Connect-go handlers
│   ├── infrastructure/    # Database, cache, external services
│   │   └── database/postgres/ # PostgreSQL implementation
│   ├── usecase/          # Application business rules
│   ├── config/           # Configuration management
│   └── pkg/              # Shared utilities
├── external/             # Generated protobuf code
│   └── proto/user/v1/    # User service proto definitions
├── scripts/              # Database migration scripts  
└── sqlc.yaml            # SQL code generation config
```

## Key Architectural Patterns
- **Repository Pattern**: Domain interfaces, infrastructure implementations
- **Use Case Pattern**: Single file per business operation
- **Value Objects**: Immutable với built-in validation
- **Clean Architecture Flow**: HTTP → Handler → Use Case → Entity → Repository

## Important Files
- Main: `cmd/main.go`
- Config: `internal/config/config.go`
- SQL queries: `internal/infrastructure/database/postgres/queries/users.sql`
- Migrations: `internal/infrastructure/database/postgres/migrations/`
- Value objects: `internal/domain/valueObject/`