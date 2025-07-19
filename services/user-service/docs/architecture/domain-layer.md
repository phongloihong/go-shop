# Domain Layer

The Domain Layer is the core of the Clean Architecture implementation, containing the business logic, entities, and domain rules that are independent of external concerns.

## Overview

The domain layer follows Domain-Driven Design principles and contains:

- **Entities**: Core business objects with identity and behavior
- **Value Objects**: Immutable objects that represent domain concepts
- **Repository Interfaces**: Contracts for data persistence
- **Domain Services**: Cross-entity business logic

## Directory Structure

```
internal/domain/
├── entity/
│   ├── user.go
│   └── userPublicProfile.go
├── repository/
│   └── user_repository.go
├── service/
│   └── money.go
└── valueObject/
    ├── email.go
    ├── password.go
    ├── phone.go
    └── time.go
```

## Entities

### User Entity

**Location:** `internal/domain/entity/user.go`

The User entity represents the core user business object with the following attributes:

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

**Key Methods:**
- `NewUser()`: Factory method for creating new users with validation
- `UserFromDatabase()`: Factory method for reconstructing users from database
- `Validate()`: Business rule validation

**Business Rules:**
- User ID is automatically generated using UUID
- Email must be valid format
- Password must meet security requirements
- Phone number must be valid format
- Created/Updated timestamps are automatically managed

### User Public Profile

**Location:** `internal/domain/entity/userPublicProfile.go`

Represents public user information that can be shared with other services or users.

## Value Objects

Value objects encapsulate domain concepts and provide validation and behavior.

### Email Value Object

**Location:** `internal/domain/valueObject/email.go`

```go
type Email string
```

**Features:**
- Email format validation using `net/mail` package
- Immutable string-based representation
- Built-in validation method

### Password Value Object

**Location:** `internal/domain/valueObject/password.go`

```go
type Password string
```

**Features:**
- Minimum 8 character validation
- bcrypt hashing capability
- Secure password comparison
- Immutable representation

**Security Methods:**
- `Hash()`: Generates bcrypt hash
- `CompareHash()`: Verifies password against hash
- `Validate()`: Ensures password meets requirements

### Phone Value Object

**Location:** `internal/domain/valueObject/phone.go`

```go
type Phone string
```

**Features:**
- Phone number format validation
- Immutable string representation

### DateTime Value Object

**Location:** `internal/domain/valueObject/time.go`

```go
type DateTime struct
```

**Features:**
- Time handling and formatting
- Consistent timestamp management
- Conversion utilities

## Repository Interfaces

### User Repository Interface

**Location:** `internal/domain/repository/user_repository.go`

Defines the contract for user data persistence without specifying implementation details:

```go
type UserRepository interface {
    Create(user *entity.User) error
    GetByID(id string) (*entity.User, error)
    GetByEmail(email string) (*entity.User, error)
    Update(user *entity.User) error
    UpdatePassword(id string, password string) error
    GetPublicProfiles(ids []string) ([]*entity.UserPublicProfile, error)
}
```

## Domain Services

### Money Service

**Location:** `internal/domain/service/money.go`

Handles money-related business logic that spans across entities or requires special domain knowledge.

## Design Principles

### Dependency Rule

The domain layer:
- ✅ Can depend on standard Go libraries
- ✅ Can define interfaces for external dependencies  
- ❌ Cannot depend on infrastructure layer
- ❌ Cannot depend on application layer
- ❌ Cannot depend on delivery layer

### Business Logic Encapsulation

All business rules and validations are contained within the domain layer:
- Entity validation rules
- Value object constraints
- Domain service logic
- Repository contracts

### Immutability

Value objects are designed to be immutable, ensuring data consistency and preventing unintended modifications.

## Testing Strategy

The domain layer should be thoroughly tested with:
- Unit tests for entity business logic
- Value object validation tests
- Domain service behavior tests
- Mock implementations of repository interfaces