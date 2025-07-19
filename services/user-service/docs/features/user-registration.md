# User Registration

This document describes the user registration feature implementation in the User Service.

## Overview

User registration allows new users to create accounts in the travel planning application. The feature implements secure user creation with comprehensive validation and password hashing.

## Business Requirements

### Core Functionality

- New user account creation
- Email uniqueness validation
- Secure password handling
- User data validation
- Automatic timestamp management

### Validation Rules

- **First Name**: Required, 1-100 characters
- **Last Name**: Required, 1-100 characters  
- **Email**: Required, valid email format, unique in system
- **Phone**: Optional, valid phone number format
- **Password**: Required, minimum 8 characters

## Implementation

### Domain Layer

#### User Entity

**Location:** `internal/domain/entity/user.go:19`

The User entity provides a factory method for creating new users:

```go
func NewUser(firstName, lastName, email, phone, password, createdAt, updatedAt string) (*User, error) {
    passwordVO := valueobject.NewPassword(password)
    emailVO := valueobject.NewEmail(email)
    phoneVO := valueobject.NewPhone(phone)
    nowVO := valueobject.NewTime(utils.TimeNow())

    user := &User{
        ID:        utils.NewUUID(),
        FirstName: firstName,
        LastName:  lastName,
        Email:     emailVO,
        Phone:     phoneVO,
        Password:  passwordVO,
        CreatedAt: nowVO,
        UpdatedAt: nowVO,
    }

    if err := user.Validate(); err != nil {
        return nil, err
    }

    return user, nil
}
```

#### Value Object Validation

**Email Validation** (`internal/domain/valueObject/email.go:15`):
```go
func (e Email) Validate() error {
    _, err := mail.ParseAddress(string(e))
    return err
}
```

**Password Validation** (`internal/domain/valueObject/password.go:19`):
```go
func (p Password) Validate() error {
    if len(p) < 8 {
        return errors.New("password must be at least 8 characters long")
    }
    return nil
}
```

### Application Layer

#### Use Case Implementation

**Expected Location:** `internal/usecase/user_usecase.go`

[TODO: Verify with team] - The current use case implementation appears incomplete.

**Expected Implementation:**

```go
func (uc *UserUseCase) CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
    // Check if user already exists
    existingUser, err := uc.userRepo.GetByEmail(req.Email)
    if err == nil && existingUser != nil {
        return nil, ErrEmailAlreadyExists
    }

    // Create new user entity
    user, err := entity.NewUser(
        req.FirstName,
        req.LastName, 
        req.Email,
        req.Phone,
        req.Password,
        "", // createdAt will be auto-generated
        "", // updatedAt will be auto-generated
    )
    if err != nil {
        return nil, fmt.Errorf("user creation failed: %w", err)
    }

    // Hash password
    hashedPassword, err := user.Password.Hash()
    if err != nil {
        return nil, fmt.Errorf("password hashing failed: %w", err)
    }
    user.Password = valueobject.NewPassword(hashedPassword)

    // Save to database
    if err := uc.userRepo.Create(user); err != nil {
        return nil, fmt.Errorf("user persistence failed: %w", err)
    }

    // Return response (without password)
    return &UserResponse{
        ID:        user.ID,
        FirstName: user.FirstName,
        LastName:  user.LastName,
        Email:     user.Email.String(),
        Phone:     user.Phone.String(),
        CreatedAt: user.CreatedAt.Time(),
        UpdatedAt: user.UpdatedAt.Time(),
    }, nil
}
```

### Infrastructure Layer

#### Database Persistence

**SQL Query** (`internal/infrastructure/database/postgres/queries/users.sql:1`):
```sql
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

#### Repository Implementation

**Location:** `internal/infrastructure/database/postgres/user_repository.go`

Expected implementation:
```go
func (r *PostgresUserRepository) Create(user *entity.User) error {
    hashedPassword, err := user.Password.Hash()
    if err != nil {
        return err
    }

    _, err = r.db.InsertUser(context.Background(), sqlc.InsertUserParams{
        FirstName: user.FirstName,
        LastName:  user.LastName,
        Email:     user.Email.String(),
        Phone:     sql.NullString{String: user.Phone.String(), Valid: user.Phone.String() != ""},
        Password:  hashedPassword,
        CreatedAt: user.CreatedAt.Time(),
        UpdatedAt: user.UpdatedAt.Time(),
    })
    
    return err
}
```

### Delivery Layer

#### HTTP Handler

**Expected Location:** `internal/delivery/http/handler/user_handler.go`

Expected implementation:
```go
func (h *UserHandler) CreateUser(c echo.Context) error {
    var req CreateUserRequest
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

    user, err := h.userUseCase.CreateUser(c.Request().Context(), req)
    if err != nil {
        return h.handleError(c, err)
    }

    return c.JSON(http.StatusCreated, user)
}
```

## API Interface

### Request Format

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone": "+1234567890",
  "password": "securepassword123"
}
```

### Response Format

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone": "+1234567890",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Error Responses

```json
{
  "error": "Email already exists",
  "code": "EMAIL_ALREADY_EXISTS"
}
```

```json
{
  "error": "Password must be at least 8 characters long",
  "code": "VALIDATION_ERROR"
}
```

## Security Features

### Password Security

- **Hashing**: Passwords are hashed using bcrypt with default cost
- **Storage**: Plain text passwords are never stored
- **Response**: Passwords are never included in API responses

### Email Validation

- **Format Validation**: Uses Go's `net/mail` package for RFC compliance
- **Uniqueness**: Database constraint ensures email uniqueness
- **Case Sensitivity**: Email handling is case-sensitive

### Data Validation

- **Input Sanitization**: All inputs are validated before processing
- **SQL Injection Prevention**: Uses parameterized queries
- **XSS Prevention**: Proper data encoding in responses

## Business Rules

### User ID Generation

- UUID v4 format using `github.com/google/uuid`
- Automatic generation via `utils.NewUUID()`
- Immutable once created

### Timestamp Management

- `created_at`: Set once during user creation
- `updated_at`: Set during creation and updates
- Uses Unix timestamp format internally
- ISO 8601 format in API responses

### Email Uniqueness

- Enforced at database level with unique constraint
- Checked at application level before insertion
- Case-sensitive matching

## Testing Strategy

### Unit Tests

- Domain entity validation
- Value object validation
- Password hashing verification
- Email format validation

### Integration Tests

- End-to-end user creation flow
- Database persistence verification
- Duplicate email handling
- Error response validation

### Performance Tests

- Bulk user creation
- Password hashing performance
- Database insertion performance

## Error Handling

### Validation Errors

- Invalid email format
- Password too short
- Required field missing
- Invalid phone format

### Business Logic Errors

- Email already exists
- Database connection failure
- Password hashing failure

### System Errors

- Database timeout
- Memory allocation failure
- Network connectivity issues

## Future Enhancements

### Planned Features

- Email verification workflow
- Phone number verification
- Social login integration
- Account activation process
- Rate limiting for registration attempts

### Security Improvements

- Password strength requirements
- CAPTCHA integration
- IP-based rate limiting
- Suspicious activity detection