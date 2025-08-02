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
if goose postgres "user=${DATABASE_USER} password=${DATABASE_PASSWORD} dbname=${DATABASE_NAME} host=${DATABASE_HOST} sslmode=${DATABASE_SSL_MODE}" -dir /app/internal/database/migrations up; then
    echo "✅ Migrations applied successfully"

    echo "Starting database seeding..."
    if ./bin/fluxend seed settings; then
        echo "✅ Seeding completed successfully"
    else
        echo "⚠️  Seeding failed, but continuing to start server..."
    fi
else
    echo "⚠️  Migrations failed, but continuing to start server..."
fi

# Start the application
echo "Starting server..."
exec ./bin/fluxend server