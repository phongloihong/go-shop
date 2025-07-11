# **Clean Architecture Design for project**

## **Clean architecture layers**

- Presentation layer (HTTP, CLI, Websocket, etc)
- Application layer (Use cases, DTOs, Application services)
- Domain layer (Entities, Value Objects, Domain services)
- Infrastructure layer (Database, External services, File system)

## Key points

HTTP Request → Handler → Use Case → Entity → Repository Interface
                                                      ↓
HTTP Response ← Handler ← Use Case ← Entity ← Repository Implementation

- Dependency Rule: Dependencies only go from outer layter to inner layer.
  - Domain layer do not know anything about database or framework
  - Application layer do not know anything about HTTP or MongoDB
  - Adapters know about application and infrastructure layers

- Layer independence:
  - Domain layer: pure business logic, no dependencies
  - Application layer: use cases orchestrate domain logic,
  - Adapter layer: convert between formats
  - Infrastructure layer: Framework and driver details

- Benefits:
  - Testable: Easy to mock interfaces
  - Flexible: Change database or framework without affecting core logic
  - Maintainable: Clear separation of concerns
  - Scalable: Add new features without breaking existing code

- Domain-Driven Design Integration:
  - Entities with rich domain models
  - Value objects for immutable data
  - Domain services for cross-entity logic
  - Repository pattern for persistence abstraction

- Use case centric:
  - Each use case is a single file
  - Clear input/output DTOs
  - Independent and reusable

## Project flow

1. Create a plan
    - Plan setting (example: decisions need confirm from all members or leader only)
    - Add  polls to voting where, when to go and what to do
    - Invite friends to the plan
    - Voting and confirm

2. Design where to go and voting
    - Who will join who not
    - Add small transactions in one place (example: restaurant, hotel, etc)
    - Split the bill base on who joined

## Project structure

```
travel-planning-clean/
├── cmd/
│   ├── api/
│   │   └── main.go                 # HTTP API entry point
│   └── worker/
│       └── main.go                 # Background worker entry point
│
├── internal/
│   ├── domain/                     # Enterprise Business Rules (innermost layer)
│   │   ├── entity/
│   │   │   ├── plan.go          # Product entity
│   │   │   ├── user.go             # User entity
│   │   │   ├── location.go            # Order entity
│   │   │   └── transaction.go             # Cart entity
│   │   │
│   │   ├── valueobject/            # object do not have identity, only value
│   │   │   ├── money.go            # Money value object
│   │   │   ├── email.go            # Email value object
│   │   │
│   │   ├── repository/             # only interfaces, implementation in the infrastructure
│   │   │   ├── plan.go
│   │   │   ├── user.go
│   │   │   └── activity.go
│   │   │
│   │   ├── service/                # do not belong to any entity or belong many entities
│   │   │   ├── voting.go           # how to count vote and return the result 
│   │   │   └── pricing.go          # how to calculate price
│   │   │
│   │   └── event/                  # CDC events and event publisher
│   │       ├── events.go           # Event definitions
│   │       └── publisher.go        # Event publisher interface
│   │
│   ├── application/                # Application Business Rules
│   │   ├── usecase/                # Use cases (one per file)
│   │   │   ├── create_product.go
│   │   │   ├── list_products.go
│   │   │   ├── add_to_cart.go
│   │   │   ├── checkout.go
│   │   │   ├── create_order.go
│   │   │   └── process_payment.go
│   │   │
│   │   ├── dto/                    # Data Transfer Objects
│   │   │   ├── product_dto.go
│   │   │   ├── order_dto.go
│   │   │   └── user_dto.go
│   │   │
│   │   ├── mapper/                 # Entity <-> DTO mappers
│   │   │   └── mapper.go
│   │   │
│   │   └── port/                   # Secondary ports (interfaces for infra)
│   │       ├── email_service.go
│   │       ├── payment_gateway.go
│   │       └── notification.go
│   │
│   ├── adapter/                    # Interface Adapters
│   │   ├── inbound/                # Primary adapters (driving)
│   │   │   ├── http/
│   │   │   │   ├── handler/
│   │   │   │   │   ├── product.go
│   │   │   │   │   ├── order.go
│   │   │   │   │   └── user.go
│   │   │   │   ├── middleware/
│   │   │   │   │   ├── auth.go
│   │   │   │   │   └── cors.go
│   │   │   │   ├── request/       # HTTP request structs
│   │   │   │   ├── response/      # HTTP response structs
│   │   │   │   └── router.go
│   │   │   │
│   │   │   ├── grpc/              # gRPC handlers
│   │   │   └── graphql/           # GraphQL resolvers
│   │   │
│   │   └── outbound/              # Secondary adapters (driven)
│   │       ├── persistence/
│   │       │   ├── mongodb/
│   │       │   │   ├── product_repo.go
│   │       │   │   ├── user_repo.go
│   │       │   │   └── order_repo.go
│   │       │   └── redis/
│   │       │       └── cache.go
│   │       │
│   │       ├── email/
│   │       │   └── sendgrid.go
│   │       │
│   │       └── payment/
│   │           ├── stripe.go
│   │           └── paypal.go
│   │
│   └── infrastructure/            # Frameworks & Drivers (outermost layer)
│       ├── config/
│       │   └── config.go
│       ├── database/
│       │   └── mongodb.go
│       └── server/
│           └── http.go
│
└── pkg/                          # Shared utilities
    ├── errors/
    ├── logger/
    └── validator/
```
