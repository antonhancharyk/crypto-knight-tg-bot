FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/tgbot ./cmd/tgbot

FROM gcr.io/distroless/static
COPY --from=builder /bin/tgbot /bin/tgbot
ENTRYPOINT ["/bin/tgbot"]