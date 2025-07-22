# Development Dockerfile with hot reload and proto generation support
FROM golang:1.24.2-alpine AS development

# Install development dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    make \
    curl \
    protobuf \
    protobuf-dev

# Install Go development tools
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest && \
    go install github.com/air-verse/air@latest && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest && \
    go install github.com/bufbuild/buf/cmd/buf@latest

# Install golang-migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/

# Set working directory
WORKDIR /app

# Copy go mod files for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Expose port for development
EXPOSE 8080

# Create tmp directory for air
RUN mkdir -p tmp

# Default command for development (can be overridden)
CMD ["air", "-c", ".air.toml"]