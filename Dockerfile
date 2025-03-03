# Use the official Go image as a base
FROM golang:1.23-alpine

WORKDIR /app
COPY . .

# Install dependencies
RUN go mod tidy

# Build the Go app using your specific build command
RUN go build -o bin/fluxton cmd/*.go

# Expose the port Echo is running on (change if needed)
EXPOSE 80

# Run the Go binary
CMD ["./bin/fluxton"]
