# Builder image
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the application
RUN go build -o main ./cmd/main.go

# Final image
FROM alpine:latest

WORKDIR /app
# Copy the .env file
COPY .env .env

# Copy the binary from the builder image
COPY --from=builder /app/main /app/main

# Expose the port the application listens on
EXPOSE 8080

# Run the application
CMD ["./main"]