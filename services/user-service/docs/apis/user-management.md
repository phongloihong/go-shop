# User Management API

This document describes the user management endpoints available in the User Service.

## Base URL

The service runs on the configured port (default: 8080) with base path `/api/v1`

## Endpoints

### Create User

Create a new user account.

**Endpoint:** `POST /users`

**Request Body:**
```json
{
  "first_name": "string",
  "last_name": "string", 
  "email": "string",
  "phone": "string",
  "password": "string"
}
```

**Response:**
```json
{
  "id": "uuid",
  "first_name": "string",
  "last_name": "string",
  "email": "string", 
  "phone": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

**Validation Rules:**
- Email must be valid format
- Password must be at least 8 characters
- Phone number must be valid format
- First name and last name are required

### Get User by ID

Retrieve user information by user ID.

**Endpoint:** `GET /users/{id}`

**Path Parameters:**
- `id` (string): User UUID

**Response:**
```json
{
  "id": "uuid",
  "first_name": "string",
  "last_name": "string",
  "email": "string",
  "phone": "string", 
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Get User by Email

Retrieve user information by email address.

**Endpoint:** `GET /users/email/{email}`

**Path Parameters:**
- `email` (string): User email address

**Response:**
```json
{
  "id": "uuid",
  "first_name": "string",
  "last_name": "string",
  "email": "string",
  "phone": "string",
  "created_at": "timestamp", 
  "updated_at": "timestamp"
}
```

### Update User

Update user profile information.

**Endpoint:** `PUT /users/{id}`

**Path Parameters:**
- `id` (string): User UUID

**Request Body:**
```json
{
  "first_name": "string",
  "last_name": "string",
  "email": "string",
  "phone": "string"
}
```

**Response:**
```json
{
  "message": "User updated successfully"
}
```

### Update User Password

Update user password.

**Endpoint:** `PUT /users/{id}/password`

**Path Parameters:**
- `id` (string): User UUID

**Request Body:**
```json
{
  "password": "string"
}
```

**Response:**
```json
{
  "message": "Password updated successfully"
}
```

**Validation:**
- Password must be at least 8 characters long
- Password is automatically hashed using bcrypt

### Get Public Profiles

Retrieve public profile information for multiple users.

**Endpoint:** `POST /users/profiles`

**Request Body:**
```json
{
  "user_ids": ["uuid1", "uuid2", "uuid3"]
}
```

**Response:**
```json
[
  {
    "id": "uuid",
    "first_name": "string", 
    "last_name": "string"
  }
]
```

## Error Responses

All endpoints return standard HTTP status codes:

- `400 Bad Request`: Invalid request data or validation errors
- `404 Not Found`: User not found
- `500 Internal Server Error`: Server error

Error response format:
```json
{
  "error": "Error message description",
  "code": "error_code"
}
```

## Authentication

[TODO: Verify with team] - Authentication mechanism not yet implemented in the current codebase.