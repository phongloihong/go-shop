# Database Setup

This document describes how to set up and manage the PostgreSQL database for the User Service.

## Prerequisites

- PostgreSQL 12 or higher
- golang-migrate tool

### Installing golang-migrate

```bash
# Using Homebrew (macOS)
brew install golang-migrate

# Using Go install
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Using curl (Linux)
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
```

## Database Creation

### 1. Create Database

Connect to PostgreSQL and create the database:

```sql
-- Connect as postgres superuser
psql -U postgres

-- Create database
CREATE DATABASE user_dev;

-- Create user (optional)
CREATE USER user_service WITH ENCRYPTED PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE user_dev TO user_service;
```

### 2. Enable Required Extensions

The service requires PostgreSQL extensions for UUID generation:

```sql
-- Connect to the user database
\c user_dev

-- Enable uuid-ossp extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Verify extension
SELECT * FROM pg_extension WHERE extname = 'uuid-ossp';
```

## Migration Management

### Migration Files

**Location:** `internal/infrastructure/database/postgres/migrations/`

Current migrations:
- `000001_init_extensions.up.sql` - Initialize PostgreSQL extensions
- `000002_create_user_table.up.sql` - Create users table

### Running Migrations

#### Using the Migration Script

**Location:** `scripts/migrate.sh`

```bash
# Run all pending migrations
./scripts/migrate.sh

# Check migration status
./scripts/migrate.sh status
```

#### Using golang-migrate Directly

```bash
# Set database URL
export DATABASE_URL="postgres://user_service:password@localhost:5432/user_dev?sslmode=disable"

# Run migrations
migrate -path internal/infrastructure/database/postgres/migrations -database $DATABASE_URL up

# Check current version
migrate -path internal/infrastructure/database/postgres/migrations -database $DATABASE_URL version

# Rollback last migration
migrate -path internal/infrastructure/database/postgres/migrations -database $DATABASE_URL down 1
```

### Creating New Migrations

```bash
# Create a new migration
migrate create -ext sql -dir internal/infrastructure/database/postgres/migrations -seq add_user_index

# This creates two files:
# 000003_add_user_index.up.sql
# 000003_add_user_index.down.sql
```

## Database Schema

### Users Table

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

-- Indexes
CREATE INDEX idx_users_email ON users(email);
```

### Column Details

- **id**: Primary key using UUID with automatic generation
- **first_name**: User's first name (required, max 100 chars)
- **last_name**: User's last name (required, max 100 chars)  
- **email**: Unique email address (required, max 100 chars)
- **phone**: Optional phone number (max 20 chars)
- **password**: bcrypt hashed password (required, max 255 chars)
- **created_at**: Record creation timestamp
- **updated_at**: Record modification timestamp

### Indexes

- **idx_users_email**: Optimizes email-based lookups
- **Primary key (id)**: Automatic B-tree index for UUID lookups

## SQLC Integration

The service uses SQLC for type-safe SQL query generation.

### SQLC Configuration

**Location:** `sqlc.yaml`

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "internal/infrastructure/database/postgres/queries"
    schema: "internal/infrastructure/database/postgres/migrations"
    gen:
      go:
        package: "sqlc"
        out: "internal/infrastructure/database/postgres/sqlc"
```

### Generating Query Code

```bash
# Generate type-safe Go code from SQL queries
make gen-query

# Or run sqlc directly
sqlc generate
```

### SQL Queries

**Location:** `internal/infrastructure/database/postgres/queries/users.sql`

Available queries:
- `InsertUser`: Create new user
- `GetUserByID`: Retrieve user by ID
- `GetUserByEmail`: Retrieve user by email
- `UpdateUser`: Update user profile
- `UpdateUserPassword`: Update user password
- `GetPublicProfileByIds`: Get public profiles for multiple users

## Database Connection

### Connection Configuration

The service uses pgx/v5 driver for PostgreSQL connections:

```go
type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
    Name     string `yaml:"name"`
    SSLMode  string `yaml:"ssl_mode"`
}
```

### Connection Pool Settings

Recommended production settings:

```go
config.MaxConns = 25
config.MinConns = 5
config.MaxConnLifetime = time.Hour
config.MaxConnIdleTime = time.Minute * 30
```

## Backup and Recovery

### Database Backup

```bash
# Create database backup
pg_dump -h localhost -U user_service -d user_dev > user_backup.sql

# Create compressed backup
pg_dump -h localhost -U user_service -d user_dev | gzip > user_backup.sql.gz
```

### Database Restore

```bash
# Restore from backup
psql -h localhost -U user_service -d user_dev < user_backup.sql

# Restore from compressed backup
gunzip -c user_backup.sql.gz | psql -h localhost -U user_service -d user_dev
```

## Performance Tuning

### Index Optimization

Current indexes support:
- Email-based user lookups (login)
- UUID-based primary key access
- Unique constraint on email

### Query Performance

Monitor query performance using:

```sql
-- Enable query logging
ALTER SYSTEM SET log_statement = 'all';
ALTER SYSTEM SET log_min_duration_statement = 1000;

-- Reload configuration
SELECT pg_reload_conf();

-- Analyze query performance
EXPLAIN ANALYZE SELECT * FROM users WHERE email = 'user@example.com';
```

## Security Considerations

### Database Security

1. **User Permissions:**
   - Create dedicated database user for the service
   - Grant minimal required permissions
   - Avoid using superuser accounts

2. **Connection Security:**
   - Use SSL connections in production
   - Configure proper firewall rules
   - Use connection pooling

3. **Data Security:**
   - Passwords are bcrypt hashed
   - Sensitive data is properly handled
   - Regular security updates

### Production Checklist

- [ ] SSL/TLS enabled
- [ ] Strong passwords configured
- [ ] Regular backups scheduled
- [ ] Monitoring and alerting set up
- [ ] Connection limits configured
- [ ] Database user permissions reviewed
- [ ] Security patches applied

## Troubleshooting

### Common Issues

1. **Migration Failed:**
   ```bash
   # Check migration status
   migrate -path internal/infrastructure/database/postgres/migrations -database $DATABASE_URL version
   
   # Force version (if needed)
   migrate -path internal/infrastructure/database/postgres/migrations -database $DATABASE_URL force 2
   ```

2. **Connection Refused:**
   - Verify PostgreSQL is running
   - Check port and host configuration
   - Verify firewall settings

3. **Permission Denied:**
   - Check database user permissions
   - Verify database exists
   - Check authentication method

4. **SQLC Generation Failed:**
   ```bash
   # Verify SQLC configuration
   sqlc vet
   
   # Check SQL syntax
   sqlc compile
   ```