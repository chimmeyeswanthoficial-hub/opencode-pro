# Multi-stage Dockerfile for OpenCode Pro & Headless Daemon
FROM golang:1.24-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o /app/opencode .

# Runtime Stage
FROM alpine:3.20

WORKDIR /workspace
RUN apk add --no-cache git bash ripgrep fzf curl ca-certificates

COPY --from=builder /app/opencode /usr/local/bin/opencode

EXPOSE 8080

ENTRYPOINT ["opencode"]
CMD ["daemon", "--port", "8080", "--host", "0.0.0.0"]
