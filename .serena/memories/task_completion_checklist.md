# Task Completion Checklist

## Sau khi hoàn thành một task, PHẢI chạy các commands sau:

### 1. Code Generation (nếu có thay đổi SQL/Proto)
```bash
make gen               # Generate cả proto và sqlc code
# hoặc riêng lẻ:
make proto            # Nếu có thay đổi .proto files
make sqlc             # Nếu có thay đổi SQL queries
```

### 2. Code Quality Checks
```bash
make fmt              # Format Go code
make vet              # Run go vet
make lint             # Run linter (nếu available)
```

### 3. Testing
```bash
make test             # Run all tests
make test-coverage    # Run tests với coverage report
```

### 4. Dependencies
```bash
make tidy             # Clean up go modules
```

### 5. Build Verification
```bash
make build            # Build service binary để verify
```

### 6. Health Check
```bash
make health           # Check all services healthy
```

## Migration Tasks (nếu có schema changes)
```bash
make migrate-up       # Apply new migrations
make db-reset         # Reset database nếu cần thiết
```

## Documentation Updates
- Update service-specific CLAUDE.md nếu có changes significant
- Update relevant files trong docs/ directory nếu cần

## IMPORTANT Notes
- KHÔNG commit trừ khi user yêu cầu explicitly
- LUÔN chạy linting và formatting trước khi submit code
- Verify service starts successfully sau changes
- Check logs không có errors: `make logs-user`