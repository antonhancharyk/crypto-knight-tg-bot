# syntax=docker/dockerfile:1.7
# ---------------- BUILD STAGE ----------------
FROM --platform=$BUILDPLATFORM golang:1.26.1-alpine3.23 AS builder
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=$(go env GOARCH) \
    go build \
        -trimpath \
        -ldflags="-s -w -buildid=" \
        -o app ./cmd/tgbot/main.go
# ---------------- RUNTIME STAGE ----------------
FROM alpine:3.23.3
RUN apk add --no-cache tini ca-certificates curl
RUN addgroup -S app && adduser -S app -G app
WORKDIR /app
COPY --from=builder --chown=app:app /src/app .
USER app
EXPOSE 8080
ENTRYPOINT ["/sbin/tini", "--", "/app/app"]