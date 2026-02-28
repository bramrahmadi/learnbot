# ─── Build Stage ─────────────────────────────────────────────────────────────
# Build context must be the repository root (not api-gateway/) because
# go.mod has a replace directive: github.com/learnbot/resume-parser => ../resume-parser
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /workspace

# Copy the resume-parser module first (required by replace directive)
COPY resume-parser/go.mod resume-parser/go.sum ./resume-parser/
COPY resume-parser/ ./resume-parser/

# Cache api-gateway dependencies
COPY api-gateway/go.mod api-gateway/go.sum ./api-gateway/
WORKDIR /workspace/api-gateway
RUN go mod download && go mod verify

# Copy api-gateway source and build
COPY api-gateway/ ./
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

ENTRYPOINT ["/api-gateway"]
