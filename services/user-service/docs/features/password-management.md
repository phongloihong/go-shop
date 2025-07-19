# Password Management

This document describes the password management features in the User Service, including secure password handling, validation, and updates.

## Overview

Password management ensures secure handling of user credentials through industry-standard practices including bcrypt hashing, validation rules, and secure update mechanisms.

## Security Features

### Password Hashing

- **Algorithm**: bcrypt with default cost (currently 10)
- **Salt**: Automatically generated per password
- **Storage**: Only hashed passwords are stored in database
- **Comparison**: Secure hash comparison for authentication

### Password Validation

- **Minimum Length**: 8 characters required
- **Format**: Plain text validation before hashing
- **Error Handling**: Clear validation error messages

## Implementation

### Domain Layer

#### Password Value Object

**Location:** `internal/domain/valueObject/password.go`

```go
type Password string

func NewPassword(password string) Password {
    return Password(password)
}

func (p Password) String() string {
    return string(p)
}
```

#### Validation Logic

**Location:** `internal/domain/valueObject/password.go:19`

```go
func (p Password) Validate() error {
    if len(p) < 8 {
        return errors.New("password must be at least 8 characters long")
    }
    return nil
}
```

#### Password Hashing

**Location:** `internal/domain/valueObject/password.go:27`

```go
func (p Password) Hash() (string, error) {
    bytes, error := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
    return string(bytes), error
}
```

#### Password Verification

**Location:** `internal/domain/valueObject/password.go:32`

```go
func (p Password) CompareHash(hash Password) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(p))
}
```

### Application Layer

#### Password Update Use Case

**Expected Location:** `internal/usecase/user_usecase.go`

[TODO: Verify with team] - The current use case implementation appears incomplete.

**Expected Implementation:**

```go
func (uc *UserUseCase) UpdatePassword(ctx context.Context, userID string, req UpdatePasswordRequest) error {
    // Validate user exists
    user, err := uc.userRepo.GetByID(userID)
    if err != nil {
        return ErrUserNotFound
    }

    // Create and validate new password
    newPassword := valueobject.NewPassword(req.Password)
    if err := newPassword.Validate(); err != nil {
        return fmt.Errorf("password validation failed: %w", err)
    }

    // Hash the new password
    hashedPassword, err := newPassword.Hash()
    if err != nil {
        return fmt.Errorf("password hashing failed: %w", err)
    }

    // Update password in repository
    if err := uc.userRepo.UpdatePassword(userID, hashedPassword); err != nil {
        return fmt.Errorf("password update failed: %w", err)
    }

    return nil
}
```

#### Password Verification Use Case

```go
func (uc *UserUseCase) VerifyPassword(ctx context.Context, email, password string) (*entity.User, error) {
    // Get user by email
    user, err := uc.userRepo.GetByEmail(email)
    if err != nil {
        return nil, ErrUserNotFound
    }

    // Create password value object
    inputPassword := valueobject.NewPassword(password)
    
    // Compare with stored hash
    if err := inputPassword.CompareHash(user.Password); err != nil {
        return nil, ErrInvalidCredentials
    }

    return user, nil
}
```

### Infrastructure Layer

#### Database Query

**Location:** `internal/infrastructure/database/postgres/queries/users.sql:24`

```sql
-- name: UpdateUserPassword :execresult
UPDATE users
SET
  password = $2,
  updated_at = $3
WHERE id = $1;
```

#### Repository Implementation

**Expected Location:** `internal/infrastructure/database/postgres/user_repository.go`

```go
func (r *PostgresUserRepository) UpdatePassword(id string, hashedPassword string) error {
    _, err := r.db.UpdateUserPassword(context.Background(), sqlc.UpdateUserPasswordParams{
        ID:        uuid.MustParse(id),
        Password:  hashedPassword,
        UpdatedAt: time.Now(),
    })
    
    return err
}
```

### Delivery Layer

#### HTTP Handler

**Expected Implementation:**

```go
func (h *UserHandler) UpdatePassword(c echo.Context) error {
    userID := c.Param("id")
    
    var req UpdatePasswordRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{
            Error: "Invalid request format",
            Code:  "INVALID_REQUEST",
        })
    }
    
    if err := req.Validate(); err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{
            Error: err.Error(),
            Code:  "VALIDATION_ERROR",
        })
    }
    
    if err := h.userUseCase.UpdatePassword(c.Request().Context(), userID, req); err != nil {
        return h.handleError(c, err)
    }
    
    return c.JSON(http.StatusOK, map[string]string{
        "message": "Password updated successfully",
    })
}
```

## API Interface

### Update Password

**Endpoint:** `PUT /users/{id}/password`

**Request:**
```json
{
  "password": "newSecurePassword123"
}
```

**Response:**
```json
{
  "message": "Password updated successfully"
}
```

### Error Responses

#### Validation Error
```json
{
  "error": "Password must be at least 8 characters long",
  "code": "VALIDATION_ERROR"
}
```

#### User Not Found
```json
{
  "error": "User not found",
  "code": "USER_NOT_FOUND"
}
```

#### Server Error
```json
{
  "error": "Password update failed",
  "code": "INTERNAL_ERROR"
}
```

## Security Best Practices

### bcrypt Configuration

- **Cost Factor**: Uses bcrypt default cost (10)
- **Salt**: Automatically generated per password
- **Randomization**: Each hash is unique even for identical passwords

### Password Storage

- **Plain Text**: Never stored in database
- **JSON Responses**: Password field excluded with `json:"-"` tag
- **Logs**: Passwords excluded from all logging

### Password Transmission

- **HTTPS**: All password operations require encrypted connections
- **Request Validation**: Input sanitization and validation
- **Memory Safety**: Secure string handling in Go

## Validation Rules

### Password Requirements

1. **Minimum Length**: 8 characters
2. **Character Set**: No restrictions (allows special characters)
3. **Maximum Length**: Limited by bcrypt (72 bytes)

### Future Enhancements

Planned password requirements:
- At least one uppercase letter
- At least one lowercase letter
- At least one number
- At least one special character
- Maximum length limit
- Common password blacklist

## Authentication Flow

### Login Process

1. User provides email and password
2. System retrieves user by email
3. Password is compared with stored hash
4. Success/failure response returned

### Password Reset Flow

[TODO: Verify with team] - Password reset functionality not yet implemented.

Expected flow:
1. User requests password reset via email
2. System generates secure reset token
3. Token sent via email with expiration
4. User provides new password with token
5. System validates token and updates password

## Performance Considerations

### Hashing Performance

- **bcrypt Cost**: Balanced between security and performance
- **Parallel Processing**: Multiple password operations can run concurrently
- **Memory Usage**: bcrypt is memory-hard function

### Database Performance

- **Indexed Lookups**: Email field indexed for fast user retrieval
- **Update Efficiency**: Only password and timestamp updated
- **Connection Pooling**: Efficient database connection management

## Error Handling

### Common Errors

1. **Validation Failures**:
   - Password too short
   - Invalid character encoding
   - Empty password

2. **System Errors**:
   - Hashing failures
   - Database connection issues
   - User not found

3. **Authentication Errors**:
   - Invalid credentials
   - Hash comparison failures

### Error Recovery

- Graceful error handling with appropriate HTTP status codes
- Detailed error messages for validation issues
- Generic error messages for security-sensitive operations

## Monitoring and Logging

### Security Events

- Failed password attempts
- Successful password changes
- Authentication events
- Suspicious activity patterns

### Performance Metrics

- Password hashing duration
- Database update performance
- Authentication success/failure rates

### Audit Trail

- Password change timestamps
- User IP addresses for security events
- Failed authentication attempts

## Testing Strategy

### Unit Tests

- Password validation logic
- Hashing and verification functions
- Error handling scenarios
- Edge cases and boundary conditions

### Integration Tests

- End-to-end password update flow
- Authentication with various passwords
- Database persistence verification
- Error response validation

### Security Tests

- Timing attack resistance
- Hash collision testing
- Password strength validation
- Brute force protection

## Compliance Considerations

### Data Protection

- **GDPR**: Password data handling compliance
- **CCPA**: User data privacy requirements
- **SOC 2**: Security controls and monitoring

### Industry Standards

- **OWASP**: Password storage guidelines
- **NIST**: Digital identity guidelines
- **ISO 27001**: Information security management

## Migration Considerations

### Password Hash Upgrades

Future bcrypt cost increases:
1. Detect old cost factors during authentication
2. Re-hash passwords with new cost on successful login
3. Gradual migration without user disruption

### Algorithm Migration

If migration to different hashing algorithm needed:
1. Support multiple hash formats during transition
2. Re-hash on password changes or successful logins
3. Maintain backward compatibility during migration period