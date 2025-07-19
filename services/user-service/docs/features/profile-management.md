# Profile Management

This document describes the user profile management features in the User Service.

## Overview

Profile management allows users to view and update their personal information. The feature supports updating user details while maintaining data integrity and validation.

## Features

### Profile Retrieval

Users can retrieve their profile information using:
- User ID lookup
- Email address lookup
- Public profile information for multiple users

### Profile Updates

Users can update the following information:
- First name
- Last name
- Email address
- Phone number

### Data Validation

All profile updates are validated against business rules and constraints.

## Implementation

### Domain Layer

#### User Entity Structure

**Location:** `internal/domain/entity/user.go:8`

```go
type User struct {
    ID        string               `json:"id"`
    FirstName string               `json:"first_name"`
    LastName  string               `json:"last_name"`
    Email     valueobject.Email    `json:"email"`
    Phone     valueobject.Phone    `json:"phone"`
    Password  valueobject.Password `json:"-"`
    CreatedAt valueobject.DateTime `json:"created_at"`
    UpdatedAt valueobject.DateTime `json:"updated_at"`
}
```

#### Database Reconstruction

**Location:** `internal/domain/entity/user.go:43`

Factory method for reconstructing users from database:

```go
func UserFromDatabase(id, firstName, lastName, email, phone, password string, createdAt, updatedAt int64) *User {
    passwordVO := valueobject.NewPassword(password)
    emailVO := valueobject.NewEmail(email)
    phoneVO := valueobject.NewPhone(phone)
    createdAtVO := valueobject.NewTime(createdAt)
    updatedAtVO := valueobject.NewTime(updatedAt)

    user := &User{
        ID:        id,
        FirstName: firstName,
        LastName:  lastName,
        Email:     emailVO,
        Phone:     phoneVO,
        Password:  passwordVO,
        CreatedAt: createdAtVO,
        UpdatedAt: updatedAtVO,
    }

    return user
}
```

### Application Layer

#### Use Case Operations

**Expected Location:** `internal/usecase/user_usecase.go`

[TODO: Verify with team] - The current use case implementation appears incomplete.

**Expected Profile Operations:**

##### Get User Profile

```go
func (uc *UserUseCase) GetUserByID(ctx context.Context, userID string) (*UserResponse, error) {
    // Validate UUID format
    if err := utils.ValidateUUID(userID); err != nil {
        return nil, ErrInvalidUserID
    }

    // Retrieve user from repository
    user, err := uc.userRepo.GetByID(userID)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrUserNotFound
        }
        return nil, fmt.Errorf("failed to retrieve user: %w", err)
    }

    // Convert to response DTO
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

##### Update User Profile

```go
func (uc *UserUseCase) UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) (*UserResponse, error) {
    // Retrieve existing user
    existingUser, err := uc.userRepo.GetByID(userID)
    if err != nil {
        return nil, ErrUserNotFound
    }

    // Validate new email uniqueness (if changed)
    if req.Email != existingUser.Email.String() {
        emailUser, err := uc.userRepo.GetByEmail(req.Email)
        if err == nil && emailUser != nil {
            return nil, ErrEmailAlreadyExists
        }
    }

    // Create updated user entity
    updatedUser := &entity.User{
        ID:        existingUser.ID,
        FirstName: req.FirstName,
        LastName:  req.LastName,
        Email:     valueobject.NewEmail(req.Email),
        Phone:     valueobject.NewPhone(req.Phone),
        Password:  existingUser.Password, // Keep existing password
        CreatedAt: existingUser.CreatedAt,
        UpdatedAt: valueobject.NewTime(utils.TimeNow()),
    }

    // Validate updated user
    if err := updatedUser.Validate(); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // Update in repository
    if err := uc.userRepo.Update(updatedUser); err != nil {
        return nil, fmt.Errorf("update failed: %w", err)
    }

    // Return updated user
    return &UserResponse{
        ID:        updatedUser.ID,
        FirstName: updatedUser.FirstName,
        LastName:  updatedUser.LastName,
        Email:     updatedUser.Email.String(),
        Phone:     updatedUser.Phone.String(),
        CreatedAt: updatedUser.CreatedAt.Time(),
        UpdatedAt: updatedUser.UpdatedAt.Time(),
    }, nil
}
```

### Infrastructure Layer

#### Database Queries

**Location:** `internal/infrastructure/database/postgres/queries/users.sql`

##### Get User by ID

```sql
-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;
```

##### Get User by Email

```sql
-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;
```

##### Update User Profile

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

##### Get Public Profiles

```sql
-- name: GetPublicProfileByIds :many
SELECT id, first_name, last_name FROM users
WHERE id = ANY(sqlc.arg(user_ids)::string[]);
```

#### Repository Implementation

**Expected Location:** `internal/infrastructure/database/postgres/user_repository.go`

```go
func (r *PostgresUserRepository) GetByID(id string) (*entity.User, error) {
    user, err := r.db.GetUserByID(context.Background(), id)
    if err != nil {
        return nil, err
    }

    return entity.UserFromDatabase(
        user.ID.String(),
        user.FirstName,
        user.LastName,
        user.Email,
        user.Phone.String,
        user.Password,
        user.CreatedAt.Unix(),
        user.UpdatedAt.Unix(),
    ), nil
}

func (r *PostgresUserRepository) Update(user *entity.User) error {
    _, err := r.db.UpdateUser(context.Background(), sqlc.UpdateUserParams{
        ID:        uuid.MustParse(user.ID),
        FirstName: user.FirstName,
        LastName:  user.LastName,
        Email:     user.Email.String(),
        Phone:     sql.NullString{String: user.Phone.String(), Valid: true},
        UpdatedAt: user.UpdatedAt.Time(),
    })
    
    return err
}
```

### Delivery Layer

#### HTTP Endpoints

**Expected Implementation:**

##### Get User Profile

```go
func (h *UserHandler) GetUserByID(c echo.Context) error {
    userID := c.Param("id")
    
    user, err := h.userUseCase.GetUserByID(c.Request().Context(), userID)
    if err != nil {
        return h.handleError(c, err)
    }
    
    return c.JSON(http.StatusOK, user)
}
```

##### Update User Profile

```go
func (h *UserHandler) UpdateUser(c echo.Context) error {
    userID := c.Param("id")
    
    var req UpdateUserRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{
            Error: "Invalid request format",
            Code:  "INVALID_REQUEST",
        })
    }
    
    user, err := h.userUseCase.UpdateUser(c.Request().Context(), userID, req)
    if err != nil {
        return h.handleError(c, err)
    }
    
    return c.JSON(http.StatusOK, user)
}
```

## API Interface

### Get User by ID

**Endpoint:** `GET /users/{id}`

**Response:**
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

### Get User by Email

**Endpoint:** `GET /users/email/{email}`

**Response:** Same as Get User by ID

### Update User Profile

**Endpoint:** `PUT /users/{id}`

**Request:**
```json
{
  "first_name": "John",
  "last_name": "Smith",
  "email": "john.smith@example.com",
  "phone": "+1234567890"
}
```

**Response:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "first_name": "John",
  "last_name": "Smith",
  "email": "john.smith@example.com",
  "phone": "+1234567890",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T14:20:00Z"
}
```

### Get Public Profiles

**Endpoint:** `POST /users/profiles`

**Request:**
```json
{
  "user_ids": [
    "123e4567-e89b-12d3-a456-426614174000",
    "987fcdeb-51a2-43d1-9876-543210987654"
  ]
}
```

**Response:**
```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "first_name": "John",
    "last_name": "Doe"
  },
  {
    "id": "987fcdeb-51a2-43d1-9876-543210987654",
    "first_name": "Jane",
    "last_name": "Smith"
  }
]
```

## Validation Rules

### Profile Update Validation

- **First Name**: Required, 1-100 characters
- **Last Name**: Required, 1-100 characters
- **Email**: Valid email format, unique in system
- **Phone**: Valid phone number format

### Business Rules

- Email uniqueness is enforced across all users
- User ID cannot be changed after creation
- Password is not included in profile operations
- Timestamps are automatically managed

## Security Considerations

### Access Control

- Users can only access their own profile data
- Public profiles return limited information
- Authentication required for profile updates

### Data Privacy

- Password is never returned in API responses
- Sensitive data is properly validated and sanitized
- Input validation prevents injection attacks

### Audit Trail

- UpdatedAt timestamp tracks profile changes
- Profile changes should be logged for audit purposes

## Error Handling

### Common Errors

- `USER_NOT_FOUND`: User does not exist
- `EMAIL_ALREADY_EXISTS`: Email is already registered
- `INVALID_INPUT`: Validation failed
- `UNAUTHORIZED`: Access denied

### Error Response Format

```json
{
  "error": "User not found",
  "code": "USER_NOT_FOUND"
}
```

## Public Profile Features

### Limited Data Exposure

Public profiles only expose:
- User ID
- First name
- Last name

### Batch Retrieval

Multiple public profiles can be retrieved in a single request for efficiency.

### Use Cases

- User search results
- Group member listings
- Friend suggestions
- Travel plan participants

## Performance Considerations

### Database Optimization

- Email index for fast lookups
- UUID index for user ID queries
- Batch queries for multiple profiles

### Caching Strategy

- Cache frequently accessed profiles
- Invalidate cache on profile updates
- Use Redis for session-based caching

## Future Enhancements

### Profile Completeness

- Profile completion percentage
- Required vs optional fields
- Profile verification status

### Privacy Settings

- Profile visibility controls
- Contact information privacy
- Activity privacy settings

### Extended Profile Data

- Profile pictures
- Bio/description
- Social media links
- Preferences and settings