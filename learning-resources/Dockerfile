# ─── Build Stage ─────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
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
