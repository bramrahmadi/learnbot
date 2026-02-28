# ─── Build Stage ─────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo dev)" \
    -o /api-gateway \
    ./cmd/server/main.go

# ─── Runtime Stage ───────────────────────────────────────────────────────────
FROM scratch

# Import CA certificates and timezone data from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /api-gateway /api-gateway

# Non-root user (numeric UID for scratch images)
USER 65534:65534

EXPOSE 8090

HEALTHCHECK --interval=30s --timeout=5s --start-period=30s --retries=3 \
    CMD ["/api-gateway", "-health-check"]

ENTRYPOINT ["/api-gateway"]
