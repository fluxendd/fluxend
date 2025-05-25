# Use the official Go image as a base
FROM golang:1.23-alpine

RUN apk add --no-cache docker

WORKDIR /app
COPY . .

# Install dependencies for Go application
RUN go mod tidy

# Build the Go app using your specific build command. This builds for Linux amd64 architecture.
# We don't use CROSS platform binary to reduce overall size
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o bin/fluxend cmd/main.go

# Expose the port Echo is running on (change if needed)
EXPOSE 8080

# Run the Go binary
CMD ["./bin/fluxton", "server"]
