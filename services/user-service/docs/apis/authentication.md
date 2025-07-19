# Authentication API

This document describes authentication-related endpoints and mechanisms.

## Overview

The User Service includes authentication middleware and password management capabilities for securing user operations.

## Authentication Middleware

The service includes JWT-based authentication middleware located at:
- `internal/delivery/http/middleware/auth.go`

## Password Management

### Password Security

The service implements secure password handling with the following features:

- **Hashing**: Passwords are hashed using bcrypt with default cost
- **Validation**: Minimum 8 character requirement
- **Storage**: Passwords are never stored in plain text

### Password Operations

#### Hash Password

Passwords are automatically hashed when users are created or passwords are updated.

**Implementation Reference:** `internal/domain/valueObject/password.go:27`

```go
func (p Password) Hash() (string, error) {
    bytes, error := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
    return string(bytes), error
}
```

#### Verify Password

Password verification compares plain text passwords with hashed versions.

**Implementation Reference:** `internal/domain/valueObject/password.go:32`

```go
func (p Password) CompareHash(hash Password) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(p))
}
```

## Authentication Flow

[TODO: Verify with team] - Complete authentication flow endpoints (login, logout, token refresh) are not yet implemented in the current codebase.

### Planned Authentication Endpoints

The following endpoints are expected to be implemented:

- `POST /auth/login` - User login with email/password
- `POST /auth/logout` - User logout
- `POST /auth/refresh` - Token refresh
- `GET /auth/me` - Get current user information

## Security Considerations

- All passwords are hashed using bcrypt
- Email validation uses Go's `net/mail` package
- UUIDs are used for user identification
- Middleware authentication is available for protecting routes

## Error Handling

Authentication errors follow standard HTTP status codes:

- `401 Unauthorized`: Invalid credentials or missing authentication
- `403 Forbidden`: Valid authentication but insufficient permissions
- `422 Unprocessable Entity`: Invalid input format

## Implementation Status

- ✅ Password hashing and validation
- ✅ Authentication middleware structure
- ✅ User validation (email, password, phone)
- ❌ Login/logout endpoints
- ❌ JWT token generation and validation
- ❌ Session management