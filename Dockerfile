# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ana ./cmd/server/main.go

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install necessary packages
RUN apk --no-cache add ca-certificates tzdata

# Set timezone
ENV TZ=UTC

# Copy the binary from the builder stage
COPY --from=builder /app/ana .

# Copy web assets
COPY web/ ./web/

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080
ENV GIN_MODE=release

# Create a non-root user to run the app
RUN adduser -D -g '' anauser
RUN chown -R anauser:anauser /app
USER anauser

# Run the application
CMD ["./ana"]

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:8080/health || exit 1

