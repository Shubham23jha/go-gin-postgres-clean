# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build all three services
RUN go build -o api-server ./cmd/server/main.go
RUN go build -o publisher ./cmd/publisher/main.go
RUN go build -o worker ./cmd/worker/main.go
RUN go build -o dashboard ./cmd/dashboard/main.go

# Final Stage (Standard Alpine for a tiny image)
FROM alpine:latest

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/api-server .
COPY --from=builder /app/publisher .
COPY --from=builder /app/worker .
COPY --from=builder /app/dashboard .
COPY --from=builder /app/cmd/dashboard/web ./web-dashboard
COPY --from=builder /app/cmd/server/web ./web-server
COPY --from=builder /app/migrations ./migrations

# We don't set a default CMD here because we will override it in K8s/Docker-Compose
# But for safety, we can default to the API server
CMD ["./api-server"]
