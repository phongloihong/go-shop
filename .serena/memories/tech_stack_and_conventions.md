# Tech Stack và Code Conventions

## Tech Stack chính
- **Go 1.24.2**: Ngôn ngữ chính
- **PostgreSQL**: Database chính với pgx/v5 driver
- **Redis**: Cache layer
- **Docker & Docker Compose**: Development environment
- **SQLC**: Type-safe SQL code generation
- **golang-migrate**: Database migration tool
- **Connect-go**: HTTP/gRPC/WebRPC 3-in-1 API framework
- **Buf**: Protocol buffer management và validation
- **Viper**: Configuration management

## Code Style và Conventions

### Language Rules
- **Response**: Luôn trả lời bằng tiếng Việt
- **Comments**: Code comments luôn bằng tiếng Anh
- NO AI/Claude/Co-Author references trong code hoặc commit messages

### Clean Architecture Rules
- Dependencies chỉ flow từ outer layers vào inner layers
- Domain layer không phụ thuộc vào bất kỳ layer nào khác
- Repository interfaces định nghĩa trong domain, implement trong infrastructure
- Use cases orchestrate domain logic
- Value objects bất biến với validation tích hợp

### Naming Conventions
- Package names: lowercase, single word
- Interface names: thường kết thúc với "Repository", "Service"
- Value objects: camelCase với validation methods
- Entity constructors: NewEntityName()
- Error variables: ErrDescription

### File Organization Patterns
- Mỗi use case là một file riêng trong `application/usecase/`
- Repository interfaces trong domain layer
- Repository implementations trong infrastructure layer
- Value objects trong `internal/domain/valueObject/`
- Entities trong `internal/domain/entity/`

### Error Handling
- Domain-specific errors
- Wrap infrastructure errors appropriately  
- Meaningful error messages
- Proper logging levels