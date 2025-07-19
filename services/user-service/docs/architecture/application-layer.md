# Application Layer

The Application Layer orchestrates business use cases and coordinates between the domain layer and external interfaces. It contains application-specific business rules and use case implementations.

## Overview

The application layer is responsible for:

- **Use Case Implementation**: Specific business operations and workflows
- **Data Transformation**: Converting between domain objects and DTOs
- **Business Flow Orchestration**: Coordinating domain services and repositories
- **Transaction Management**: Ensuring data consistency across operations

## Directory Structure

```
internal/usecase/
└── user_usecase.go
```

## Use Cases

### User Use Case

**Location:** `internal/usecase/user_usecase.go`

The User Use Case implements business operations related to user management.

[TODO: Verify with team] - The current use case implementation appears to be a placeholder or incomplete.

### Expected Use Case Operations

Based on the domain model and repository interfaces, the following use cases should be implemented:

#### Create User Use Case

**Purpose**: Handle new user registration with business validation

**Flow**:
1. Validate input data
2. Check if user already exists (by email)
3. Create domain User entity with validation
4. Hash password securely
5. Persist user to database
6. Return user response (without password)

#### Get User Use Case

**Purpose**: Retrieve user information by ID or email

**Flow**:
1. Validate input parameters
2. Query repository for user data
3. Convert domain entity to response DTO
4. Return user information

#### Update User Use Case

**Purpose**: Update user profile information

**Flow**:
1. Validate input data
2. Retrieve existing user
3. Apply updates with domain validation
4. Persist changes to database
5. Return success confirmation

#### Update Password Use Case

**Purpose**: Securely update user password

**Flow**:
1. Validate new password requirements
2. Hash new password using bcrypt
3. Update password in database
4. Return success confirmation

#### Get Public Profiles Use Case

**Purpose**: Retrieve public user information for multiple users

**Flow**:
1. Validate user ID list
2. Query repository for public profiles
3. Return filtered public information

## Use Case Structure

Each use case should follow this pattern:

```go
type UserUseCase struct {
    userRepo domain.UserRepository
    // other dependencies
}

func NewUserUseCase(userRepo domain.UserRepository) *UserUseCase {
    return &UserUseCase{
        userRepo: userRepo,
    }
}

func (uc *UserUseCase) CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
    // Use case implementation
}
```

## Data Transfer Objects (DTOs)

### Request DTOs

Request DTOs represent incoming data from external layers:

```go
type CreateUserRequest struct {
    FirstName string `json:"first_name" validate:"required"`
    LastName  string `json:"last_name" validate:"required"`
    Email     string `json:"email" validate:"required,email"`
    Phone     string `json:"phone" validate:"required"`
    Password  string `json:"password" validate:"required,min=8"`
}

type UpdateUserRequest struct {
    FirstName string `json:"first_name" validate:"required"`
    LastName  string `json:"last_name" validate:"required"`
    Email     string `json:"email" validate:"required,email"`
    Phone     string `json:"phone" validate:"required"`
}

type UpdatePasswordRequest struct {
    Password string `json:"password" validate:"required,min=8"`
}
```

### Response DTOs

Response DTOs represent outgoing data to external layers:

```go
type UserResponse struct {
    ID        string    `json:"id"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Email     string    `json:"email"`
    Phone     string    `json:"phone"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type PublicProfileResponse struct {
    ID        string `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
}
```

## Error Handling

Use cases should handle and transform domain errors into application-specific errors:

```go
type UseCaseError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

func (e UseCaseError) Error() string {
    return e.Message
}
```

Common error scenarios:
- `USER_NOT_FOUND`: User does not exist
- `EMAIL_ALREADY_EXISTS`: Email is already registered
- `INVALID_INPUT`: Request validation failed
- `INTERNAL_ERROR`: Unexpected system error

## Dependency Injection

Use cases receive their dependencies through constructor injection:

```go
type Dependencies struct {
    UserRepo   domain.UserRepository
    Cache      infrastructure.CacheService
    Logger     infrastructure.Logger
}

func NewUserUseCase(deps Dependencies) *UserUseCase {
    return &UserUseCase{
        userRepo: deps.UserRepo,
        cache:    deps.Cache,
        logger:   deps.Logger,
    }
}
```

## Transaction Management

For operations requiring multiple database operations, use cases should coordinate transactions:

```go
func (uc *UserUseCase) ComplexOperation(ctx context.Context) error {
    return uc.dbManager.WithTransaction(ctx, func(tx Transaction) error {
        // Multiple operations within transaction
        return nil
    })
}
```

## Design Principles

### Single Responsibility

Each use case handles one specific business operation with a clear purpose and boundary.

### Dependency Inversion

Use cases depend on domain interfaces, not concrete implementations from infrastructure layer.

### Input Validation

All external input is validated at the use case level before being passed to domain entities.

### Error Translation

Domain errors are translated into application-specific errors appropriate for external consumers.