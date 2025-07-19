# Infrastructure Layer

The Infrastructure Layer contains implementations of external dependencies and technical concerns. It implements the interfaces defined in the domain layer and provides concrete implementations for databases, caching, messaging, and other external services.

## Overview

The infrastructure layer is responsible for:

- **Database Implementation**: Concrete repository implementations
- **Caching**: Redis-based caching services
- **Messaging**: NATS message broker integration
- **External Services**: Third-party API integrations
- **Configuration**: Application configuration management
- **Logging**: Structured logging implementation

## Directory Structure

```
internal/infrastructure/
├── cache/
│   └── redis.go
├── database/
│   └── postgres/
│       ├── migrations/
│       │   ├── 000001_init_extensions.down.sql
│       │   ├── 000001_init_extensions.up.sql
│       │   ├── 000002_create_user_table.down.sql
│       │   └── 000002_create_user_table.up.sql
│       ├── queries/
│       │   └── users.sql
│       ├── sqlc/
│       │   ├── db.go
│       │   ├── models.go
│       │   └── users.sql.go
│       └── user_repository.go
└── message/
    └── nats.go
```

## Database Implementation

### PostgreSQL Integration

The service uses PostgreSQL as the primary database with the following components:

#### Database Configuration

**Connection Parameters:**
- Host, Port, Database Name
- User credentials and SSL settings
- Connection pool configuration
- Migration management

#### SQLC Code Generation

**Location:** `internal/infrastructure/database/postgres/sqlc/`

The service uses SQLC for type-safe SQL query generation:

- `db.go`: Database connection interface
- `models.go`: Generated Go structs matching database schema
- `users.sql.go`: Generated query functions

#### SQL Queries

**Location:** `internal/infrastructure/database/postgres/queries/users.sql`

Contains all SQL operations for user management:

```sql
-- Create user
INSERT INTO users (first_name, last_name, email, phone, password, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- Get user by ID
SELECT * FROM users WHERE id = $1;

-- Get user by email
SELECT * FROM users WHERE email = $1;

-- Update user
UPDATE users SET first_name = $2, last_name = $3, email = $4, phone = $5, updated_at = $6
WHERE id = $1;

-- Update password
UPDATE users SET password = $2, updated_at = $3 WHERE id = $1;

-- Get public profiles
SELECT id, first_name, last_name FROM users WHERE id = ANY($1::string[]);
```

#### User Repository Implementation

**Location:** `internal/infrastructure/database/postgres/user_repository.go`

Implements the domain repository interface:

```go
type PostgresUserRepository struct {
    db *sqlc.Queries
}

func NewPostgresUserRepository(db *sqlc.Queries) domain.UserRepository {
    return &PostgresUserRepository{
        db: db,
    }
}

func (r *PostgresUserRepository) Create(user *entity.User) error {
    // Implementation using SQLC generated queries
}

func (r *PostgresUserRepository) GetByID(id string) (*entity.User, error) {
    // Implementation using SQLC generated queries
}
```

### Database Migrations

**Location:** `internal/infrastructure/database/postgres/migrations/`

Migration files manage database schema evolution:

#### Extension Setup
- `000001_init_extensions.up.sql`: Initialize PostgreSQL extensions
- `000001_init_extensions.down.sql`: Remove extensions

#### User Table
- `000002_create_user_table.up.sql`: Create users table with indexes
- `000002_create_user_table.down.sql`: Drop users table

**User Table Schema:**
```sql
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  first_name VARCHAR(100) NOT NULL,
  last_name VARCHAR(100) NOT NULL,
  email VARCHAR(100) NOT NULL UNIQUE,
  phone VARCHAR(20) DEFAULT NULL,
  password VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
```

## Caching Implementation

### Redis Integration

**Location:** `internal/infrastructure/cache/redis.go`

Provides caching capabilities for frequently accessed data:

```go
type RedisCache struct {
    client redis.Client
}

func NewRedisCache(client redis.Client) *RedisCache {
    return &RedisCache{
        client: client,
    }
}
```

**Use Cases for Caching:**
- User profile data
- Authentication tokens
- Session management
- Frequently queried user information

## Messaging Implementation

### NATS Integration

**Location:** `internal/infrastructure/message/nats.go`

Provides message broker capabilities for inter-service communication:

```go
type NATSMessaging struct {
    conn *nats.Conn
}

func NewNATSMessaging(conn *nats.Conn) *NATSMessaging {
    return &NATSMessaging{
        conn: conn,
    }
}
```

**Message Types:**
- User creation events
- Profile update notifications
- Authentication events
- System health checks

## Configuration Management

### Configuration Structure

**Location:** `internal/config/config.go`

Centralized configuration management using Viper:

```go
type Config struct {
    Database DatabaseConfig `yaml:"database"`
    Redis    RedisConfig    `yaml:"redis"`
    Server   ServerConfig   `yaml:"server"`
    NATS     NATSConfig     `yaml:"nats"`
}

type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
    Name     string `yaml:"name"`
    SSLMode  string `yaml:"ssl_mode"`
}
```

### Environment Variables

Configuration supports environment variable binding:

```
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=password
DATABASE_DB_NAME=user
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

## Logging Implementation

### Structured Logging

**Location:** `internal/pkg/logger/logger.go`

Provides structured logging capabilities:

```go
type Logger interface {
    Info(msg string, fields ...interface{})
    Error(msg string, err error, fields ...interface{})
    Debug(msg string, fields ...interface{})
    Warn(msg string, fields ...interface{})
}
```

**Log Levels:**
- DEBUG: Detailed debugging information
- INFO: General application flow
- WARN: Warning conditions
- ERROR: Error conditions

## Utility Functions

### Identification Utilities

**Location:** `internal/pkg/utils/identify.go`

Provides UUID generation and validation:

```go
func NewUUID() string {
    return uuid.New().String()
}

func ValidateUUID(id string) error {
    _, err := uuid.Parse(id)
    return err
}
```

### Time Utilities

**Location:** `internal/pkg/utils/time.go`

Provides time handling utilities:

```go
func TimeNow() int64 {
    return time.Now().Unix()
}

func FormatTime(timestamp int64) string {
    return time.Unix(timestamp, 0).Format(time.RFC3339)
}
```

## Connection Management

### Database Connection Pool

Efficient database connection management:

```go
type ConnectionPool struct {
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
}
```

### Health Checks

Infrastructure components provide health check endpoints:

```go
func (r *PostgresUserRepository) HealthCheck() error {
    return r.db.Ping()
}

func (c *RedisCache) HealthCheck() error {
    return c.client.Ping().Err()
}
```

## Design Principles

### Interface Implementation

All infrastructure components implement domain interfaces, ensuring loose coupling and testability.

### Configuration Flexibility

Support for multiple configuration sources (YAML, environment variables, command-line flags).

### Error Handling

Comprehensive error handling with proper error wrapping and context preservation.

### Performance Optimization

- Connection pooling for database and cache
- Efficient query patterns
- Proper indexing strategies
- Connection reuse

### Security

- Secure credential management
- SQL injection prevention
- Connection encryption
- Authentication token handling