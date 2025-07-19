# Adapter Layer

The Adapter Layer serves as the interface between external systems and the application core. It includes HTTP handlers, gRPC handlers, and other external interfaces that adapt external requests to internal use cases.

## Overview

The adapter layer is responsible for:

- **Protocol Translation**: Converting HTTP/gRPC requests to use case calls
- **Data Serialization**: Marshalling and unmarshalling request/response data
- **Error Handling**: Translating application errors to appropriate HTTP status codes
- **Authentication**: Applying security middleware and authorization
- **Routing**: Defining API endpoints and request routing

## Directory Structure

```
internal/delivery/
├── grpc/
│   ├── handler/
│   │   └── user_handler.go
│   └── proto/
│       └── user.proto
└── http/
    ├── handler/
    │   └── user_handler.go
    ├── middleware/
    │   └── auth.go
    └── router/
        └── router.go
```

## HTTP Interface

### User HTTP Handler

**Location:** `internal/delivery/http/handler/user_handler.go`

[TODO: Verify with team] - The current HTTP handler implementation appears to be a placeholder or incomplete.

### Expected HTTP Handler Structure

```go
type UserHandler struct {
    userUseCase usecase.UserUseCase
    logger      logger.Logger
}

func NewUserHandler(userUseCase usecase.UserUseCase, logger logger.Logger) *UserHandler {
    return &UserHandler{
        userUseCase: userUseCase,
        logger:      logger,
    }
}
```

### HTTP Endpoints

Based on the repository interface and domain model, the following endpoints should be implemented:

#### Create User

```go
func (h *UserHandler) CreateUser(c echo.Context) error {
    var req CreateUserRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{
            Error: "Invalid request format",
            Code:  "INVALID_REQUEST",
        })
    }
    
    user, err := h.userUseCase.CreateUser(c.Request().Context(), req)
    if err != nil {
        return h.handleError(c, err)
    }
    
    return c.JSON(http.StatusCreated, user)
}
```

#### Get User by ID

```go
func (h *UserHandler) GetUserByID(c echo.Context) error {
    id := c.Param("id")
    
    user, err := h.userUseCase.GetUserByID(c.Request().Context(), id)
    if err != nil {
        return h.handleError(c, err)
    }
    
    return c.JSON(http.StatusOK, user)
}
```

### Router Configuration

**Location:** `internal/delivery/http/router/router.go`

The router defines API endpoints and applies middleware:

```go
func SetupRoutes(e *echo.Echo, userHandler *handler.UserHandler) {
    api := e.Group("/api/v1")
    
    users := api.Group("/users")
    users.POST("", userHandler.CreateUser)
    users.GET("/:id", userHandler.GetUserByID)
    users.GET("/email/:email", userHandler.GetUserByEmail)
    users.PUT("/:id", userHandler.UpdateUser, middleware.AuthRequired())
    users.PUT("/:id/password", userHandler.UpdatePassword, middleware.AuthRequired())
    users.POST("/profiles", userHandler.GetPublicProfiles)
}
```

### Middleware

#### Authentication Middleware

**Location:** `internal/delivery/http/middleware/auth.go`

Provides JWT-based authentication for protected routes:

```go
func AuthRequired() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // JWT token validation logic
            return next(c)
        }
    }
}
```

## gRPC Interface

### User gRPC Handler

**Location:** `internal/delivery/grpc/handler/user_handler.go`

[TODO: Verify with team] - The current gRPC handler implementation appears to be a placeholder or incomplete.

### Protocol Buffer Definition

**Location:** `internal/delivery/grpc/proto/user.proto`

[TODO: Verify with team] - The proto file appears to be empty or incomplete.

### Expected gRPC Service Definition

```protobuf
syntax = "proto3";

package user;

service UserService {
  rpc CreateUser(CreateUserRequest) returns (UserResponse);
  rpc GetUser(GetUserRequest) returns (UserResponse);
  rpc UpdateUser(UpdateUserRequest) returns (UserResponse);
  rpc UpdatePassword(UpdatePasswordRequest) returns (EmptyResponse);
  rpc GetPublicProfiles(GetPublicProfilesRequest) returns (GetPublicProfilesResponse);
}

message CreateUserRequest {
  string first_name = 1;
  string last_name = 2;
  string email = 3;
  string phone = 4;
  string password = 5;
}

message UserResponse {
  string id = 1;
  string first_name = 2;
  string last_name = 3;
  string email = 4;
  string phone = 5;
  int64 created_at = 6;
  int64 updated_at = 7;
}
```

## Error Handling

### HTTP Error Response Format

```go
type ErrorResponse struct {
    Error string `json:"error"`
    Code  string `json:"code"`
}
```

### Error Translation

Convert application errors to appropriate HTTP status codes:

```go
func (h *UserHandler) handleError(c echo.Context, err error) error {
    switch {
    case errors.Is(err, usecase.ErrUserNotFound):
        return c.JSON(http.StatusNotFound, ErrorResponse{
            Error: "User not found",
            Code:  "USER_NOT_FOUND",
        })
    case errors.Is(err, usecase.ErrEmailAlreadyExists):
        return c.JSON(http.StatusConflict, ErrorResponse{
            Error: "Email already exists",
            Code:  "EMAIL_ALREADY_EXISTS",
        })
    case errors.Is(err, usecase.ErrInvalidInput):
        return c.JSON(http.StatusBadRequest, ErrorResponse{
            Error: err.Error(),
            Code:  "INVALID_INPUT",
        })
    default:
        h.logger.Error("Internal server error", err)
        return c.JSON(http.StatusInternalServerError, ErrorResponse{
            Error: "Internal server error",
            Code:  "INTERNAL_ERROR",
        })
    }
}
```

## Request Validation

Input validation using struct tags and validation middleware:

```go
type CreateUserRequest struct {
    FirstName string `json:"first_name" validate:"required,min=1,max=100"`
    LastName  string `json:"last_name" validate:"required,min=1,max=100"`
    Email     string `json:"email" validate:"required,email"`
    Phone     string `json:"phone" validate:"required,phone"`
    Password  string `json:"password" validate:"required,min=8"`
}
```

## Security Considerations

- **Input Sanitization**: All user input is validated and sanitized
- **Authentication**: Protected endpoints require valid JWT tokens
- **Password Security**: Passwords are never returned in responses
- **CORS**: Proper CORS configuration for web clients
- **Rate Limiting**: API rate limiting to prevent abuse

## Performance Considerations

- **Response Caching**: Cache frequently accessed user data
- **Database Connection Pooling**: Efficient database connection management
- **Request Timeouts**: Appropriate timeout handling
- **Graceful Shutdown**: Proper service shutdown handling

## Design Principles

### Single Responsibility

Each handler is responsible for one specific API endpoint or related group of endpoints.

### Dependency Injection

Handlers receive use cases and other dependencies through constructor injection.

### Protocol Independence

Business logic is isolated from protocol-specific details (HTTP vs gRPC).

### Error Boundaries

Errors are caught and translated at the adapter layer, preventing internal errors from leaking to clients.