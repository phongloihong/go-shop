# Authentication API

This document describes authentication-related endpoints and mechanisms.

## Overview

The User Service provides comprehensive JWT-based authentication with access/refresh token mechanism and secure password management capabilities.

## JWT Authentication Service

The service implements JWT authentication with the following components:
- **JWT Service**: `internal/infrastructure/auth/jwt_service.go`
- **Auth Domain Interface**: `internal/domain/service/auth_service.go`
- **Token Management**: Access tokens (30min) and refresh tokens (7 days)

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
func (p Password) CompareHash(passwordString string) error {
    err := bcrypt.CompareHashAndPassword([]byte(p), []byte(passwordString))
    if err == nil {
        return nil
    }
    if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
        return fmt.Errorf("password does not match: %w", err)
    }
    return err
}
```

## Authentication Flow

### Login Endpoint

**Endpoint**: `POST /user.v1.UserService/Login`  
**Implementation**: `internal/delivery/connect/user_service.go:47`

#### Request
```protobuf
message LoginRequest {
  string email = 1 [(buf.validate.field).string.email = true];
  string password = 2;
}
```

#### Response  
```protobuf
message LoginResponse {
  string access_token = 1;
  string refresh_token = 2;
  int64 expires_in = 3; // in seconds
}
```

#### Authentication Process
1. Validate email format and password
2. Retrieve user by email from database
3. Compare provided password with stored hash
4. Generate JWT access token (30 minutes expiry)
5. Generate JWT refresh token (7 days expiry)
6. Return token pair with expiration info

### Planned Authentication Endpoints

- ✅ `POST /user.v1.UserService/Login` - User login with email/password
- ❌ `POST /auth/logout` - User logout (planned)
- ❌ `POST /auth/refresh` - Token refresh (planned)
- ❌ `GET /auth/me` - Get current user information (planned)

## Security Considerations

- All passwords are hashed using bcrypt
- Email validation uses Go's `net/mail` package
- UUIDs are used for user identification
- Middleware authentication is available for protecting routes

## Error Handling

Authentication errors are handled through domain-specific error mapping:

**Implementation**: `internal/domain/domain_errors/errors.go`

- `InvalidCredentialsError`: Invalid email/password combination
- `UserNotFoundError`: User does not exist
- `InternalError`: System errors during authentication
- `ValidationError`: Invalid input format

Errors are mapped to appropriate Connect-RPC status codes via `domain_error.MapError()`.

## Implementation Status

- ✅ Password hashing and validation
- ✅ JWT token generation and validation  
- ✅ User validation (email, password, phone)
- ✅ Login endpoint with JWT response
- ✅ Domain error handling and mapping
- ✅ Access/refresh token mechanism
- ❌ Logout endpoint
- ❌ Token refresh endpoint
- ❌ Authentication middleware for protected routes
- ❌ Session management