# Database Queries

This document describes all available database queries in the User Service, including their purpose, parameters, and usage patterns.

## Overview

The User Service uses SQLC for type-safe SQL query generation. All queries are defined in SQL files and automatically converted to Go functions.

**Query Location:** `internal/infrastructure/database/postgres/queries/users.sql`

## Query Definitions

### 1. Insert User

**Purpose:** Create a new user record in the database.

**SQL Definition:**
```sql
-- name: InsertUser :one
INSERT INTO users (
  first_name,
  last_name,
  email,
  phone,
  password,
  created_at,
  updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;
```

**Parameters:**
1. `$1` - first_name (string)
2. `$2` - last_name (string)
3. `$3` - email (string)
4. `$4` - phone (string, nullable)
5. `$5` - password (string, bcrypt hash)
6. `$6` - created_at (timestamp)
7. `$7` - updated_at (timestamp)

**Returns:** Complete user record

**Usage:** User registration and account creation

**Generated Go Function:**
```go
func (q *Queries) InsertUser(ctx context.Context, arg InsertUserParams) (User, error)
```

### 2. Get User by ID

**Purpose:** Retrieve a user record by their unique ID.

**SQL Definition:**
```sql
-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;
```

**Parameters:**
1. `$1` - id (UUID)

**Returns:** Complete user record or error if not found

**Usage:** Profile retrieval, authentication verification

**Generated Go Function:**
```go
func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
```

### 3. Get User by Email

**Purpose:** Retrieve a user record by their email address.

**SQL Definition:**
```sql
-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;
```

**Parameters:**
1. `$1` - email (string)

**Returns:** Complete user record or error if not found

**Usage:** Login authentication, email uniqueness checks

**Generated Go Function:**
```go
func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error)
```

### 4. Update User

**Purpose:** Update user profile information (excluding password).

**SQL Definition:**
```sql
-- name: UpdateUser :execresult
UPDATE users
SET
  first_name = $2,
  last_name = $3,
  email = $4,
  phone = $5,
  updated_at = $6
WHERE id = $1;
```

**Parameters:**
1. `$1` - id (UUID)
2. `$2` - first_name (string)
3. `$3` - last_name (string)
4. `$4` - email (string)
5. `$5` - phone (string, nullable)
6. `$6` - updated_at (timestamp)

**Returns:** Execution result (rows affected)

**Usage:** Profile updates, user information changes

**Generated Go Function:**
```go
func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (sql.Result, error)
```

### 5. Update User Password

**Purpose:** Update user password separately from other profile data.

**SQL Definition:**
```sql
-- name: UpdateUserPassword :execresult
UPDATE users
SET
  password = $2,
  updated_at = $3
WHERE id = $1;
```

**Parameters:**
1. `$1` - id (UUID)
2. `$2` - password (string, bcrypt hash)
3. `$3` - updated_at (timestamp)

**Returns:** Execution result (rows affected)

**Usage:** Password changes, password resets

**Generated Go Function:**
```go
func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (sql.Result, error)
```

### 6. Get Public Profiles by IDs

**Purpose:** Retrieve public profile information for multiple users.

**SQL Definition:**
```sql
-- name: GetPublicProfileByIds :many
SELECT id, first_name, last_name FROM users
WHERE id = ANY(sqlc.arg(user_ids)::string[]);
```

**Parameters:**
1. `user_ids` - Array of user IDs ([]string)

**Returns:** Array of public profile records

**Usage:** Group member listings, friend lists, search results

**Generated Go Function:**
```go
func (q *Queries) GetPublicProfileByIds(ctx context.Context, userIds []string) ([]GetPublicProfileByIdsRow, error)
```

## Query Performance Analysis

### Index Usage

| Query | Primary Index | Performance | Notes |
|-------|---------------|-------------|-------|
| InsertUser | - | O(log n) | Uses primary key index for insertion |
| GetUserByID | PRIMARY (id) | O(log n) | Direct primary key lookup |
| GetUserByEmail | idx_users_email | O(log n) | Indexed email lookup |
| UpdateUser | PRIMARY (id) | O(log n) | Primary key-based update |
| UpdateUserPassword | PRIMARY (id) | O(log n) | Primary key-based update |
| GetPublicProfileByIds | PRIMARY (id) | O(k log n) | Multiple primary key lookups |

### Query Execution Plans

#### GetUserByEmail Example
```sql
EXPLAIN ANALYZE SELECT * FROM users WHERE email = 'user@example.com';

Index Scan using idx_users_email on users  (cost=0.15..8.17 rows=1 width=185)
  Index Cond: (email = 'user@example.com'::character varying)
  Actual time: 0.025..0.026 rows=1 loops=1
```

#### GetPublicProfileByIds Example
```sql
EXPLAIN ANALYZE SELECT id, first_name, last_name FROM users 
WHERE id = ANY(ARRAY['uuid1', 'uuid2']::string[]);

Index Scan using users_pkey on users  (cost=0.15..16.34 rows=2 width=185)
  Index Cond: (id = ANY('{uuid1,uuid2}'::uuid[]))
  Actual time: 0.050..0.052 rows=2 loops=1
```

## Query Patterns and Best Practices

### Parameter Binding

All queries use parameterized statements to prevent SQL injection:

```go
// Safe - parameterized query
user, err := q.GetUserByEmail(ctx, "user@example.com")

// Unsafe - string concatenation (NOT used)
// query := "SELECT * FROM users WHERE email = '" + email + "'"
```

### Error Handling

Standard error patterns for each query type:

```go
user, err := q.GetUserByID(ctx, userID)
if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
        return nil, ErrUserNotFound
    }
    return nil, fmt.Errorf("database error: %w", err)
}
```

### Transaction Usage

For operations requiring multiple queries:

```go
tx, err := db.Begin()
if err != nil {
    return err
}
defer tx.Rollback()

qtx := q.WithTx(tx)
// Perform multiple operations with qtx

return tx.Commit()
```

## Query Optimization

### Batch Operations

For bulk operations, use array parameters:

```sql
-- Efficient bulk lookup
SELECT id, first_name, last_name FROM users 
WHERE id = ANY($1::uuid[]);

-- Inefficient multiple queries
-- SELECT id, first_name, last_name FROM users WHERE id = $1;
-- SELECT id, first_name, last_name FROM users WHERE id = $2;
-- ...
```

### Result Set Limiting

For large result sets, implement pagination:

```sql
-- Future enhancement: paginated queries
SELECT * FROM users 
ORDER BY created_at DESC 
LIMIT $1 OFFSET $2;
```

## SQLC Configuration

### Configuration File

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

### Generated Code Location

**Output Directory:** `internal/infrastructure/database/postgres/sqlc/`

Generated files:
- `db.go` - Database interface and connection handling
- `models.go` - Go structs matching database schema
- `users.sql.go` - Query functions and parameter structs

## Query Testing

### Unit Testing

Example test for query validation:

```go
func TestGetUserByEmail(t *testing.T) {
    db := setupTestDB(t)
    queries := sqlc.New(db)
    
    // Create test user
    user := createTestUser(t, queries)
    
    // Test query
    retrieved, err := queries.GetUserByEmail(context.Background(), user.Email)
    require.NoError(t, err)
    assert.Equal(t, user.ID, retrieved.ID)
}
```

### Integration Testing

Example integration test:

```go
func TestUserCreateAndRetrieve(t *testing.T) {
    db := setupIntegrationDB(t)
    queries := sqlc.New(db)
    
    // Test full create and retrieve cycle
    params := sqlc.InsertUserParams{
        FirstName: "John",
        LastName:  "Doe",
        Email:     "john@example.com",
        // ... other fields
    }
    
    created, err := queries.InsertUser(context.Background(), params)
    require.NoError(t, err)
    
    retrieved, err := queries.GetUserByID(context.Background(), created.ID)
    require.NoError(t, err)
    assert.Equal(t, created.Email, retrieved.Email)
}
```

## Monitoring and Observability

### Query Metrics

Key metrics to monitor:
- Query execution time
- Query frequency
- Error rates
- Connection pool usage

### Slow Query Logging

PostgreSQL configuration for monitoring:

```sql
-- Enable slow query logging
ALTER SYSTEM SET log_min_duration_statement = 1000;
ALTER SYSTEM SET log_statement = 'all';
SELECT pg_reload_conf();
```

### Performance Analysis

Regular performance analysis queries:

```sql
-- Check index usage
SELECT schemaname, tablename, indexname, idx_scan, idx_tup_read, idx_tup_fetch
FROM pg_stat_user_indexes
WHERE tablename = 'users';

-- Check table statistics
SELECT * FROM pg_stat_user_tables WHERE relname = 'users';
```

## Future Query Enhancements

### Planned Additions

1. **Search Queries**: Full-text search on user names
2. **Pagination**: Limit/offset queries for large result sets
3. **Sorting**: Ordered queries for user listings
4. **Filtering**: Complex WHERE clauses for admin queries
5. **Aggregations**: User count and statistics queries

### Performance Improvements

1. **Prepared Statements**: Pre-compiled queries for high-frequency operations
2. **Query Caching**: Result caching for read-heavy operations
3. **Connection Pooling**: Optimized connection management
4. **Read Replicas**: Separate read queries from write operations

## Security Considerations

### SQL Injection Prevention

- All queries use parameter binding
- No dynamic SQL construction
- Input validation at application layer

### Data Access Control

- Queries access only necessary columns
- Public profile queries exclude sensitive data
- Password queries handle hashed data only

### Audit Trail

Future enhancement for query auditing:
- Log all data modification queries
- Track user access patterns
- Monitor for suspicious query patterns