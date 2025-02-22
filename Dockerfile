# Stage 1: Build the Go binary.
FROM golang:1.24.0-alpine AS builder

WORKDIR /app

# Install GCC and dependencies for CGO.
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go binary.
RUN CGO_ENABLED=1 go build -o /harbor-service cmd/main.go

# Install golang-migrate for database migrations.
RUN go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Stage 2: Create a minimal runtime image.
FROM alpine:latest

WORKDIR /root/

# Install SQLite and migrations tool.
RUN apk add --no-cache sqlite

# Copy built binary, migrations, and entrypoint script.
COPY --from=builder /harbor-service .
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY storage/sqlite/migrations /root/storage/sqlite/migrations
COPY docker-entrypoint.sh /docker-entrypoint.sh

# Ensure the script has execution permissions.
RUN chmod +x /docker-entrypoint.sh

# Set the script as the entrypoint.
ENTRYPOINT ["/docker-entrypoint.sh"]
