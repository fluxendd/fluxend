#!/bin/sh
set -e

echo "🚀 Starting Fluxend application..."

# Wait for database to be ready
echo "⏳ Waiting for database to be ready..."
until goose postgres "user=${DATABASE_USER} password=${DATABASE_PASSWORD} dbname=${DATABASE_NAME} host=${DATABASE_HOST} sslmode=${DATABASE_SSL_MODE}" -dir /app/internal/database/migrations status > /dev/null 2>&1; do
  echo "Database is unavailable - sleeping"
  sleep 2
done

echo "✅ Database is ready!"

# Run migrations
echo "📊 Running database migrations..."
goose postgres "user=${DATABASE_USER} password=${DATABASE_PASSWORD} dbname=${DATABASE_NAME} host=${DATABASE_HOST} sslmode=${DATABASE_SSL_MODE}" -dir /app/internal/database/migrations up

echo "✅ Migrations completed successfully!"

# Start the application
echo "🎯 Starting Fluxend server..."
exec ./bin/fluxend server