# Multi-stage build for Go API

# Builder stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Install build tools and SQLite dependencies (required for CGO)
RUN apk add --no-cache git ca-certificates gcc musl-dev sqlite-dev && update-ca-certificates

# Cache modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary with CGO enabled (required for sqlite3)
RUN CGO_ENABLED=1 GOOS=linux go build -o server ./cmd/api

# Runtime stage
FROM alpine:3.20
WORKDIR /app

# Install SQLite runtime library (required for sqlite3 driver)
RUN apk add --no-cache ca-certificates tzdata sqlite && update-ca-certificates

# App binary
COPY --from=builder /app/server /app/server

# Data directory for sqlite db
RUN mkdir -p /data

# Default envs (can be overridden by compose)
ENV PORT=8080 \
    DB_PATH=/data/data.db

EXPOSE 8080

ENTRYPOINT ["/app/server"]

