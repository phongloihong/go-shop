# Database Schema

This document describes the database schema design for the User Service.

## Overview

The User Service uses PostgreSQL as the primary database with a single `users` table that stores all user-related information.

## Database Configuration

- **Database Engine**: PostgreSQL 12+
- **Driver**: pgx/v5
- **Migration Tool**: golang-migrate
- **Query Generation**: SQLC for type-safe queries

## Schema Structure

### Users Table

**Location:** `internal/infrastructure/database/postgres/migrations/000002_create_user_table.up.sql`

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
```

### Column Specifications

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | Unique user identifier |
| `first_name` | VARCHAR(100) | NOT NULL | User's first name |
| `last_name` | VARCHAR(100) | NOT NULL | User's last name |
| `email` | VARCHAR(100) | NOT NULL, UNIQUE | User's email address |
| `phone` | VARCHAR(20) | DEFAULT NULL | User's phone number (optional) |
| `password` | VARCHAR(255) | NOT NULL | bcrypt hashed password |
| `created_at` | TIMESTAMP | DEFAULT NOW() | Record creation timestamp |
| `updated_at` | TIMESTAMP | DEFAULT NOW() | Record modification timestamp |

### Indexes

**Location:** `internal/infrastructure/database/postgres/migrations/000002_create_user_table.up.sql:12`

```sql
CREATE INDEX idx_users_email ON users(email);
```

| Index | Columns | Type | Purpose |
|-------|---------|------|---------|
| `PRIMARY` | `id` | B-tree | Primary key lookup |
| `idx_users_email` | `email` | B-tree | Email-based user lookup |
| `users_email_key` | `email` | Unique | Email uniqueness constraint |

## Data Types and Constraints

### UUID Primary Key

- **Type**: UUID (Universally Unique Identifier)
- **Generation**: PostgreSQL `gen_random_uuid()` function
- **Format**: Standard UUID v4 format (e.g., `123e4567-e89b-12d3-a456-426614174000`)
- **Benefits**: Globally unique, non-sequential, secure

### String Fields

- **first_name/last_name**: VARCHAR(100) for user names
- **email**: VARCHAR(100) with unique constraint
- **phone**: VARCHAR(20) for international phone number formats
- **password**: VARCHAR(255) for bcrypt hashes (up to 60 chars needed)

### Timestamp Fields

- **created_at**: Set once during record creation
- **updated_at**: Updated on every record modification
- **Format**: PostgreSQL TIMESTAMP type
- **Default**: Current timestamp via `NOW()` function

## Database Extensions

**Location:** `internal/infrastructure/database/postgres/migrations/000001_init_extensions.up.sql`

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

Required extensions:
- **uuid-ossp**: Provides UUID generation functions

## Constraints and Validation

### Business Rules Enforced at Database Level

1. **Email Uniqueness**: UNIQUE constraint prevents duplicate emails
2. **Required Fields**: NOT NULL constraints on essential fields
3. **Data Integrity**: Primary key ensures row uniqueness

### Application-Level Validation

The following validations are handled in the application layer:

- Email format validation (RFC compliance)
- Password strength requirements (minimum 8 characters)
- Phone number format validation
- Name length and character validation

## Query Patterns

### Primary Access Patterns

Based on the SQL queries in `internal/infrastructure/database/postgres/queries/users.sql`:

1. **User Creation**: Insert new user with all required fields
2. **Lookup by ID**: Fast UUID-based primary key lookup
3. **Lookup by Email**: Indexed email lookup for authentication
4. **Profile Updates**: Update user information (excluding password)
5. **Password Updates**: Separate password update operation
6. **Bulk Public Profiles**: Retrieve multiple users' public information

### Query Performance

- **Primary Key Lookups**: O(log n) via B-tree index
- **Email Lookups**: O(log n) via B-tree index on email
- **Bulk Operations**: Efficient using ANY() array operations

## Normalization

The schema follows **Third Normal Form (3NF)**:

- Each column is dependent on the primary key
- No transitive dependencies exist
- Minimal data redundancy

### Design Decisions

1. **Single Table**: User data is cohesive and accessed together
2. **No Foreign Keys**: User service is independent (microservice design)
3. **Denormalized Timestamps**: Created/updated timestamps in same table for simplicity

## Scalability Considerations

### Indexing Strategy

- **Selective Indexes**: Only on frequently queried columns
- **Email Index**: Supports login and uniqueness checks
- **No Phone Index**: Phone lookups not expected to be frequent

### Partitioning Strategy

For future scaling:
- **Horizontal Partitioning**: By user ID hash or geographic region
- **Vertical Partitioning**: Separate frequently vs rarely accessed columns

### Connection Pooling

- **Max Connections**: Configurable pool size
- **Connection Lifetime**: Automatic connection recycling
- **Idle Timeout**: Close idle connections

## Security Considerations

### Data Protection

1. **Password Storage**: Only bcrypt hashes stored
2. **No Sensitive Data**: Minimal personal information stored
3. **UUID IDs**: Non-sequential IDs prevent enumeration attacks

### Access Control

1. **Database User**: Dedicated application user with minimal privileges
2. **Connection Security**: SSL/TLS encryption for production
3. **Query Parameterization**: Prevents SQL injection attacks

## Backup and Recovery

### Backup Strategy

- **Full Backups**: Daily full database backups
- **Incremental Backups**: Hourly transaction log backups
- **Point-in-Time Recovery**: WAL archiving for precise recovery

### Data Retention

- **User Data**: Retained indefinitely unless user requests deletion
- **Audit Logs**: Database change logs for compliance
- **Backup Retention**: 30-day backup retention policy

## Migration Management

### Migration Files

All migrations located in `internal/infrastructure/database/postgres/migrations/`:

1. `000001_init_extensions.*` - Initialize required PostgreSQL extensions
2. `000002_create_user_table.*` - Create users table and indexes

### Migration Process

- **Sequential Numbering**: Migrations applied in order
- **Up/Down Scripts**: Both creation and rollback scripts
- **Version Tracking**: golang-migrate tracks applied migrations

### Schema Evolution

Future schema changes will follow these patterns:
- **Additive Changes**: New columns with default values
- **Non-Breaking Changes**: Index additions, constraint relaxation
- **Breaking Changes**: Require application coordination

## Performance Monitoring

### Key Metrics

- **Query Performance**: Slow query logging and analysis
- **Index Usage**: Monitor index hit ratios
- **Connection Pool**: Monitor pool utilization
- **Lock Contention**: Track blocking queries

### Optimization Strategies

- **Query Analysis**: Regular EXPLAIN ANALYZE review
- **Index Maintenance**: VACUUM and REINDEX operations
- **Statistics Updates**: ANALYZE command for query planner

## Future Enhancements

### Schema Additions

Planned schema enhancements:
- **Profile Pictures**: BLOB or external URL storage
- **User Preferences**: JSON column for flexible preferences
- **Account Status**: Active/inactive/suspended status
- **Email Verification**: Boolean flag for verified emails

### Performance Improvements

- **Read Replicas**: For read-heavy workloads
- **Caching Layer**: Redis for frequently accessed data
- **Database Sharding**: For horizontal scaling

## Compliance and Auditing

### Data Privacy

- **GDPR Compliance**: User data deletion capabilities
- **Data Minimization**: Only necessary data stored
- **Audit Trail**: Change tracking for compliance

### Security Auditing

- **Access Logging**: Database access monitoring
- **Change Tracking**: Record modification auditing
- **Backup Verification**: Regular backup integrity checks