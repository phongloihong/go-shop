# Environment Configuration

This document describes the environment setup and configuration requirements for the User Service.

## Environment Variables

The User Service requires the following environment variables to be configured:

### Database Configuration

```bash
DATABASE_HOST=localhost          # PostgreSQL host address
DATABASE_PORT=5432              # PostgreSQL port number
DATABASE_USER=postgres          # Database username
DATABASE_PASSWORD=password      # Database password
DATABASE_DB_NAME=user           # Database name
DATABASE_SSL_MODE=disable       # SSL mode (disable, require, verify-ca, verify-full)
```

### Redis Configuration

```bash
REDIS_HOST=localhost            # Redis host address
REDIS_PORT=6379                 # Redis port number
REDIS_PASSWORD=                 # Redis password (empty if no auth)
REDIS_DB=0                      # Redis database number
```

### Server Configuration

```bash
SERVER_HOST=localhost           # HTTP server host
SERVER_PORT=8080               # HTTP server port
GRPC_PORT=9090                 # gRPC server port (if implemented)
```

### Authentication Configuration

```bash
PASSWORD_SECRET=secret_pw       # Secret for password operations
ACCESS_SECRET=secret_ac_token   # JWT access token signing secret
REFRESH_SECRET=secret_rf_token  # JWT refresh token signing secret
```

### NATS Configuration (Optional)

```bash
NATS_URL=nats://localhost:4222  # NATS server URL
NATS_CLUSTER_ID=go-shop        # NATS cluster ID
NATS_CLIENT_ID=user-service    # NATS client ID
```

### Logging Configuration

```bash
LOG_LEVEL=info                  # Log level (debug, info, warn, error)
LOG_FORMAT=json                 # Log format (json, text)
```

## Configuration Files

### YAML Configuration

**Location:** `internal/config/config.yaml`

The service supports YAML-based configuration:

```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: password
  name: user
  ssl_mode: disable

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

auth:
  password_secret: secret_pw
  access_secret: secret_ac_token
  refresh_secret: secret_rf_token

server:
  host: localhost
  port: 8080

nats:
  url: nats://localhost:4222
  cluster_id: go-shop
  client_id: user-service

logging:
  level: info
  format: json
```

### Environment File (.env)

Create a `.env` file in the service root directory:

```bash
# Database
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=password
DATABASE_DB_NAME=user

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Authentication
PASSWORD_SECRET=secret_pw
ACCESS_SECRET=secret_ac_token
REFRESH_SECRET=secret_rf_token

# Server
SERVER_PORT=8080

# Logging
LOG_LEVEL=debug
```

## Configuration Priority

The configuration system loads settings in the following order (later sources override earlier ones):

1. Default values in code
2. YAML configuration file
3. Environment variables
4. Command-line flags

## Development Environment

For local development, use these settings:

```bash
export DATABASE_HOST=localhost
export DATABASE_PORT=5432
export DATABASE_USER=postgres
export DATABASE_PASSWORD=password
export DATABASE_DB_NAME=user_dev
export REDIS_HOST=localhost
export REDIS_PORT=6379
export PASSWORD_SECRET=secret_pw
export ACCESS_SECRET=secret_ac_token
export REFRESH_SECRET=secret_rf_token
export SERVER_PORT=8080
export LOG_LEVEL=debug
```

## Production Environment

For production deployment, ensure:

1. **Database Security:**
   - Use strong passwords
   - Enable SSL mode
   - Configure connection limits

2. **Redis Security:**
   - Use authentication
   - Configure appropriate ACLs
   - Enable SSL/TLS

3. **Authentication Security:**
   - Use strong, randomly generated secrets for JWT tokens
   - Rotate secrets regularly
   - Store secrets securely (e.g., AWS Secrets Manager, HashiCorp Vault)
   - Never commit secrets to version control

4. **Server Security:**
   - Use HTTPS (configure reverse proxy)
   - Set appropriate timeouts
   - Configure rate limiting

5. **Logging:**
   - Set log level to `info` or `warn`
   - Use structured JSON logging
   - Configure log rotation

## Docker Environment

When using Docker, environment variables can be set in `docker-compose.yaml`:

```yaml
version: '3.8'
services:
  user-service:
    build: .
    environment:
      - DATABASE_HOST=postgres
      - DATABASE_PORT=5432
      - DATABASE_USER=postgres
      - DATABASE_PASSWORD=password
      - DATABASE_DB_NAME=user
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - PASSWORD_SECRET=secret_pw
      - ACCESS_SECRET=secret_ac_token
      - REFRESH_SECRET=secret_rf_token
      - SERVER_PORT=8080
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis
      
  postgres:
    image: postgres:15
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=password
      - POSTGRES_DB=user
    ports:
      - "5432:5432"
      
  redis:
    image: redis:7
    ports:
      - "6379:6379"
```

## Configuration Validation

The service validates configuration on startup:

- Database connection testing
- Redis connection verification
- Required environment variables check
- Port availability validation

## Troubleshooting

### Common Configuration Issues

1. **Database Connection Failed:**
   - Verify PostgreSQL is running
   - Check host and port settings
   - Verify credentials
   - Test network connectivity

2. **Redis Connection Failed:**
   - Verify Redis is running
   - Check host and port settings
   - Verify authentication settings

3. **Port Already in Use:**
   - Check if another service is using the port
   - Use `netstat -tulpn | grep :8080` to find conflicting processes
   - Change SERVER_PORT to an available port

4. **Permission Denied:**
   - Verify database user has required permissions
   - Check file system permissions for log files
   - Verify service account permissions in production

### Configuration Testing

Test configuration with:

```bash
# Test database connection
psql -h $DATABASE_HOST -p $DATABASE_PORT -U $DATABASE_USER -d $DATABASE_DB_NAME

# Test Redis connection
redis-cli -h $REDIS_HOST -p $REDIS_PORT ping

# Validate service configuration
go run cmd/main.go --validate-config
```