# ─── Build Stage ─────────────────────────────────────────────────────────────
# Build context must be the repository root (not learning-resources/) because
# go.mod has a replace directive: github.com/learnbot/database => ../database
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /workspace

# Copy the database module first (required by replace directive)
COPY database/go.mod database/go.sum ./database/
COPY database/ ./database/

# Cache learning-resources dependencies
COPY learning-resources/go.mod learning-resources/go.sum ./learning-resources/
WORKDIR /workspace/learning-resources
RUN go mod download && go mod verify

# Copy learning-resources source and build
COPY learning-resources/ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo dev)" \
    -o /learning-resources \
    ./cmd/server/main.go

# ─── Runtime Stage ───────────────────────────────────────────────────────────
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata wget

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /learning-resources /usr/local/bin/learning-resources

USER appuser

EXPOSE 8082

HEALTHCHECK --interval=30s --timeout=5s --start-period=30s --retries=3 \
    CMD wget -qO- http://localhost:8082/health || exit 1

ENTRYPOINT ["learning-resources"]
