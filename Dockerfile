# ---------- Stage 1: Build ----------
FROM golang:1.25-alpine AS builder

# Set the working directory
WORKDIR /app

# Ensure a portable, static-ish binary
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Install swag (Swagger generator)
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Generate Swagger JSON (swagger/swagger.json)
RUN swag init -g main.go -o ./swagger -ot json

# Build the Go application (strip debug info for smaller size)
RUN go build -trimpath -ldflags="-s -w" -o myapp .

# ---------- Stage 2: Final ----------
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Install runtime dependencies you actually need
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user for security
RUN addgroup -S appuser \
 && adduser -S -G appuser -H -s /sbin/nologin appuser

# Copy the binary and swagger docs, set ownership
COPY --from=builder --chown=appuser:appuser /app/myapp /app/myapp
COPY --from=builder --chown=appuser:appuser /app/swagger /app/swagger

# Run as non-root user
USER appuser

# Set the entrypoint command
ENTRYPOINT ["/app/myapp"]