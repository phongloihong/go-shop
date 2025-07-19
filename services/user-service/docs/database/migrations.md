# Database Migrations

This document describes the database migration system used in the User Service for managing schema changes and database evolution.

## Overview

The User Service uses golang-migrate for database schema management, providing versioned, sequential migration capabilities with both up and down migration support.

## Migration System

### Migration Tool

- **Tool**: golang-migrate/migrate
- **Database**: PostgreSQL
- **Location**: `internal/infrastructure/database/postgres/migrations/`
- **Naming**: Sequential numbering with descriptive names

### Directory Structure

```
internal/infrastructure/database/postgres/migrations/
├── 000001_init_extensions.down.sql
├── 000001_init_extensions.up.sql
├── 000002_create_user_table.down.sql
└── 000002_create_user_table.up.sql
```

## Current Migrations

### Migration 000001: Initialize Extensions

**Purpose:** Set up required PostgreSQL extensions for the user service.

#### Up Migration (`000001_init_extensions.up.sql`)

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

**Actions:**
- Installs uuid-ossp extension for UUID generation
- Uses IF NOT EXISTS to prevent errors on re-run

#### Down Migration (`000001_init_extensions.down.sql`)

```sql
DROP EXTENSION IF EXISTS "uuid-ossp";
```

**Actions:**
- Removes uuid-ossp extension
- Uses IF EXISTS to prevent errors if extension not installed

### Migration 000002: Create User Table

**Purpose:** Create the main users table with all required columns and indexes.

#### Up Migration (`000002_create_user_table.up.sql`)

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

**Actions:**
1. Creates users table with all required columns
2. Sets up primary key with UUID auto-generation
3. Applies NOT NULL constraints on required fields
4. Creates unique constraint on email
5. Adds performance index on email column

#### Down Migration (`000002_create_user_table.down.sql`)

```sql
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
```

**Actions:**
1. Removes email index
2. Drops users table
3. Uses IF EXISTS to prevent errors

## Migration Management

### Running Migrations

#### Using Migration Script

**Location:** `scripts/migrate.sh`

```bash
# Run all pending migrations
./scripts/migrate.sh

# Check migration status  
./scripts/migrate.sh status

# Rollback specific number of migrations
./scripts/migrate.sh down 1
```

#### Using golang-migrate Directly

```bash
# Set database URL
export DATABASE_URL="postgres://user:password@localhost:5432/dbname?sslmode=disable"

# Run all pending migrations
migrate -path internal/infrastructure/database/postgres/migrations \
        -database $DATABASE_URL up

# Check current version
migrate -path internal/infrastructure/database/postgres/migrations \
        -database $DATABASE_URL version

# Rollback last migration
migrate -path internal/infrastructure/database/postgres/migrations \
        -database $DATABASE_URL down 1

# Rollback to specific version
migrate -path internal/infrastructure/database/postgres/migrations \
        -database $DATABASE_URL goto 1
```

### Creating New Migrations

```bash
# Create new migration files
migrate create -ext sql \
               -dir internal/infrastructure/database/postgres/migrations \
               -seq add_user_indexes

# This creates:
# 000003_add_user_indexes.up.sql
# 000003_add_user_indexes.down.sql
```

### Migration Naming Convention

Pattern: `{version}_{description}.{direction}.sql`

- **version**: 6-digit sequential number (000001, 000002, etc.)
- **description**: Snake_case description of the change
- **direction**: `up` for forward migration, `down` for rollback

Examples:
- `000003_add_user_status_column.up.sql`
- `000004_create_user_sessions_table.up.sql`
- `000005_add_email_verification_index.up.sql`

## Migration Best Practices

### Schema Changes

#### Safe Changes (Non-Breaking)

1. **Adding New Columns:**
```sql
-- Up migration
ALTER TABLE users ADD COLUMN status VARCHAR(20) DEFAULT 'active';

-- Down migration  
ALTER TABLE users DROP COLUMN IF EXISTS status;
```

2. **Adding Indexes:**
```sql
-- Up migration
CREATE INDEX CONCURRENTLY idx_users_status ON users(status);

-- Down migration
DROP INDEX IF EXISTS idx_users_status;
```

3. **Adding Constraints (with defaults):**
```sql
-- Up migration
ALTER TABLE users ADD CONSTRAINT check_email_format 
  CHECK (email ~ '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

-- Down migration
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_email_format;
```

#### Potentially Breaking Changes

1. **Dropping Columns (with caution):**
```sql
-- Up migration (after ensuring application compatibility)
ALTER TABLE users DROP COLUMN IF EXISTS old_column;

-- Down migration (may lose data)
ALTER TABLE users ADD COLUMN old_column VARCHAR(100);
```

2. **Changing Column Types:**
```sql
-- Up migration
ALTER TABLE users ALTER COLUMN phone TYPE VARCHAR(30);

-- Down migration
ALTER TABLE users ALTER COLUMN phone TYPE VARCHAR(20);
```

### Data Migrations

For migrations involving data transformation:

```sql
-- Up migration with data transformation
UPDATE users SET status = 'active' WHERE status IS NULL;
ALTER TABLE users ALTER COLUMN status SET NOT NULL;

-- Down migration
ALTER TABLE users ALTER COLUMN status DROP NOT NULL;
```

## Migration Validation

### Pre-Migration Checks

Before running migrations:

```bash
# 1. Backup database
pg_dump -h localhost -U postgres -d user_dev > backup_pre_migration.sql

# 2. Validate migration syntax
migrate -path internal/infrastructure/database/postgres/migrations \
        -database $DATABASE_URL validate

# 3. Check current version
migrate -path internal/infrastructure/database/postgres/migrations \
        -database $DATABASE_URL version
```

### Post-Migration Validation

After running migrations:

```bash
# 1. Verify schema state
psql -h localhost -U postgres -d user_dev -c "\d users"

# 2. Check data integrity
psql -h localhost -U postgres -d user_dev -c "SELECT COUNT(*) FROM users;"

# 3. Verify indexes
psql -h localhost -U postgres -d user_dev -c "\di"
```

## Environment-Specific Migrations

### Development Environment

```bash
# Development database setup
createdb user_dev
export DATABASE_URL="postgres://postgres:password@localhost:5432/user_dev?sslmode=disable"
./scripts/migrate.sh
```

### Testing Environment

```bash
# Test database setup  
createdb user_test
export DATABASE_URL="postgres://postgres:password@localhost:5432/user_test?sslmode=disable"
./scripts/migrate.sh
```

### Production Environment

```bash
# Production migration (with extra caution)
# 1. Create backup
pg_dump -h prod-host -U prod-user -d user_prod > backup_$(date +%Y%m%d_%H%M%S).sql

# 2. Run migrations with monitoring
migrate -path internal/infrastructure/database/postgres/migrations \
        -database $PRODUCTION_DATABASE_URL up

# 3. Verify application functionality
```

## Error Handling and Recovery

### Common Migration Errors

1. **Migration Failed Mid-Execution:**
```bash
# Check current state
migrate -path internal/infrastructure/database/postgres/migrations \
        -database $DATABASE_URL version

# Force version if needed (use with caution)
migrate -path internal/infrastructure/database/postgres/migrations \
        -database $DATABASE_URL force 2
```

2. **Dirty Database State:**
```bash
# If migration left database in dirty state
# First, manually fix the issue, then:
migrate -path internal/infrastructure/database/postgres/migrations \
        -database $DATABASE_URL force 2

# Then retry migration
migrate -path internal/infrastructure/database/postgres/migrations \
        -database $DATABASE_URL up
```

3. **Rollback Required:**
```bash
# Rollback to previous version
migrate -path internal/infrastructure/database/postgres/migrations \
        -database $DATABASE_URL down 1

# Or rollback to specific version
migrate -path internal/infrastructure/database/postgres/migrations \
        -database $DATABASE_URL goto 1
```

### Recovery Procedures

1. **Backup Restoration:**
```bash
# Drop current database
dropdb user_dev

# Recreate database
createdb user_dev

# Restore from backup
psql -h localhost -U postgres -d user_dev < backup_pre_migration.sql
```

2. **Manual Schema Fixes:**
```sql
-- If migration partially completed
-- Manually apply remaining changes
-- Then force version
```

## Migration Monitoring

### Migration Status Tracking

```bash
# Current migration version
migrate -path internal/infrastructure/database/postgres/migrations \
        -database $DATABASE_URL version

# Migration history (using schema_migrations table)
psql -h localhost -U postgres -d user_dev \
     -c "SELECT * FROM schema_migrations ORDER BY version;"
```

### Performance Monitoring

For large migrations:

```sql
-- Monitor long-running migrations
SELECT pid, state, query, query_start 
FROM pg_stat_activity 
WHERE state = 'active' AND query LIKE '%ALTER TABLE%';

-- Check table locks
SELECT * FROM pg_locks WHERE relation = 'users'::regclass;
```

## Advanced Migration Techniques

### Concurrent Index Creation

```sql
-- For large tables, use concurrent index creation
CREATE INDEX CONCURRENTLY idx_users_created_at ON users(created_at);
```

### Zero-Downtime Migrations

For production environments:

1. **Add nullable column:**
```sql
ALTER TABLE users ADD COLUMN new_field VARCHAR(100);
```

2. **Deploy application with dual write:**
```go
// Application writes to both old and new fields
```

3. **Backfill data:**
```sql
UPDATE users SET new_field = old_field WHERE new_field IS NULL;
```

4. **Make column NOT NULL:**
```sql
ALTER TABLE users ALTER COLUMN new_field SET NOT NULL;
```

5. **Deploy application using new field:**
```go
// Application uses new field only
```

6. **Remove old column:**
```sql
ALTER TABLE users DROP COLUMN old_field;
```

## Future Migration Planning

### Anticipated Changes

1. **User Status Column**: Active, inactive, suspended states
2. **Email Verification**: Boolean flag for verified emails
3. **Profile Pictures**: URL or blob storage reference
4. **User Preferences**: JSON column for flexible settings
5. **Audit Columns**: Created_by, updated_by tracking

### Migration Strategy

- **Backward Compatibility**: Maintain API compatibility during migrations
- **Feature Flags**: Use feature flags for gradual rollouts
- **Monitoring**: Comprehensive monitoring during migrations
- **Rollback Plans**: Always have rollback procedures ready

## Security Considerations

### Migration Security

- **Access Control**: Limit migration execution to authorized personnel
- **Backup Requirements**: Mandatory backups before production migrations
- **Audit Trail**: Log all migration activities
- **Testing**: Thorough testing in staging environments

### Data Protection

- **Sensitive Data**: Handle sensitive data carefully during migrations
- **Encryption**: Maintain encryption during data transformations
- **Privacy Compliance**: Ensure GDPR/CCPA compliance during schema changes