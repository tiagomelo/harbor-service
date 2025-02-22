#!/bin/sh

DB_FILE="storage/sqlite/harbor.db"
MIGRATIONS_PATH="/root/storage/sqlite/migrations"

# Run migrations using golang-migrate
echo "Running database migrations..."
migrate -database "sqlite3://${DB_FILE}" -path "${MIGRATIONS_PATH}" up

# Start the service
exec ./harbor-service -p "${PORT:-8080}"
