# Go Shop - Project Overview

## Mục đích dự án
Go Shop là một ứng dụng **travel planning** (lập kế hoạch du lịch) được xây dựng theo kiến trúc microservices với Clean Architecture. Ứng dụng cho phép:

1. **Tạo kế hoạch du lịch**: Người dùng tạo plans với cơ chế voting
2. **Lập kế hoạch hợp tác**: Bạn bè tham gia plans và vote cho điểm đến/hoạt động
3. **Quản lý chi phí**: Theo dõi và chia bills dựa trên việc tham gia
4. **Quản lý địa điểm**: Thiết kế hành trình với voting cho nơi tham quan

## Kiến trúc tổng quan
Dự án tuân theo **Clean Architecture** với 4 tầng:
1. **Domain Layer** (innermost): Business logic thuần túy
2. **Application Layer**: Use cases và DTOs
3. **Adapter Layer**: HTTP/gRPC handlers và repository implementations
4. **Infrastructure Layer** (outermost): Database, external services

## Tech Stack
- **Language**: Go 1.24.2
- **Database**: PostgreSQL với pgx/v5 driver  
- **Cache**: Redis
- **Code Generation**: sqlc cho type-safe SQL queries
- **API**: connect-go (HTTP + gRPC + web_rpc 3-in-1)
- **Migration**: golang-migrate
- **Config**: Viper với YAML và environment variables
- **Security**: bcrypt, UUID

## Services hiện tại
- **user-service**: Service quản lý người dùng (đang triển khai)

## Business Domain chính
- Plan creation với configurable decision-making
- Voting system cho destinations, timing, activities
- Transaction management cho shared expenses  
- Bill splitting dựa trên participation