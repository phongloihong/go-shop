# User Service

The User Service is a microservice responsible for user management in the go-shop travel planning application. It handles user registration, authentication, profile management, and provides user data to other services in the system.

## Table of Contents

- [Quick Start](#quick-start)
- [Architecture](#architecture)
- [APIs](#apis)
- [Features](#features)
- [Setup](#setup)
- [Database](#database)

## Quick Start

1. Install dependencies: `go mod download`
2. Set up environment variables (see [setup/environment.md](setup/environment.md))
3. Run database migrations: `./scripts/migrate.sh`
4. Generate SQL code: `make gen-query`
5. Start the service: `go run cmd/main.go`

## Architecture

The User Service implements Clean Architecture with four distinct layers:

- **[Domain Layer](architecture/domain-layer.md)**: Business entities, value objects, and domain logic
- **[Application Layer](architecture/application-layer.md)**: Use cases and business orchestration
- **[Adapter Layer](architecture/adapter-layer.md)**: Connect-RPC handlers and external interfaces
- **[Infrastructure Layer](architecture/infrastructure-layer.md)**: Database, cache, and external service implementations

## APIs

- **[User Management](apis/user-management.md)**: User CRUD operations and profile management
- **[Authentication](apis/authentication.md)**: User authentication and authorization endpoints

## Features

- **[User Registration](features/user-registration.md)**: New user account creation with validation
- **[Profile Management](features/profile-management.md)**: User profile updates and data management
- **[Password Management](features/password-management.md)**: Secure password handling and updates

## Setup

- **[Environment Configuration](setup/environment.md)**: Required environment variables and configuration
- **[Database Setup](setup/database.md)**: PostgreSQL setup and migration instructions
- **[Development Setup](setup/development.md)**: Local development environment setup

## Database

- **[Schema Design](database/schema.md)**: Database table structure and relationships
- **[Queries](database/queries.md)**: Available SQL queries and operations
- **[Migrations](database/migrations.md)**: Database migration management