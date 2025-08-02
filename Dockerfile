# Build stage
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache curl

WORKDIR /app
COPY . .

# Install dependencies for Go application
RUN go mod tidy

# Download Goose binary
RUN GOOSE_VERSION=v3.22.1 && \
    curl -fsSL https://github.com/pressly/goose/releases/download/${GOOSE_VERSION}/goose_linux_x86_64 -o /tmp/goose && \
    chmod +x /tmp/goose

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o bin/fluxend cmd/main.go

# Final stage
FROM alpine:latest

RUN apk add --no-cache docker ca-certificates

WORKDIR /app

# Copy binaries from builder stage
COPY --from=builder /tmp/goose /usr/local/bin/goose
COPY --from=builder /app/bin/fluxend ./bin/fluxend
COPY --from=builder /app/internal/database/migrations ./internal/database/migrations
COPY --from=builder /app/internal/database/seeders/client ./internal/database/seeders/client
COPY --from=builder /app/scripts/run.sh ./run.sh

# Make entrypoint executable
RUN chmod +x /app/run.sh

# Expose the port Echo is running on (change if needed)
EXPOSE 8080

# Run the entrypoint script
CMD ["/app/run.sh"]