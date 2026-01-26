# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sparkyfitness-mcp ./cmd/sparkyfitness-mcp

# Final stage - Alpine base image
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN adduser -D -u 1000 mcp

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/sparkyfitness-mcp .

# Use non-root user
USER mcp

# Run the server
ENTRYPOINT ["./sparkyfitness-mcp"]
