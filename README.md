# crypto-knight-tg-bot

Telegram bot for managing and receiving information from the crypto-knight trading core.

## Coverage

![CI](https://github.com/antonhancharyk/crypto-knight-tg-bot/actions/workflows/ci.yml/badge.svg)
[![codecov](https://codecov.io/gh/antonhancharyk/crypto-knight-tg-bot/branch/main/graph/badge.svg)](https://app.codecov.io/gh/antonhancharyk/crypto-knight-tg-bot)
[![Go Report Card](https://goreportcard.com/badge/github.com/antonhancharyk/crypto-knight-tg-bot)](https://goreportcard.com/report/github.com/antonhancharyk/crypto-knight-tg-bot)


## Architecture

The project follows **clean/hexagonal** layout:

- **`cmd/tgbot`** — entrypoint; loads config, wires dependencies, runs the app.
- **`internal/domain`** — core entities (e.g. `Report`).
- **`internal/ports`** — interfaces for external concerns: `Logger`, `ReportFetcher`.
- **`internal/usecase`** — business logic (e.g. report date validation and fetching).
- **`internal/transport/telegram`** — Telegram bot handler (updates, callbacks, menus).
- **`internal/infra`** — implementations: HTTP client (reports API), RabbitMQ consumer, Zap logger.

Use cases depend on **ports** (e.g. `ReportFetcher`), not on concrete HTTP clients, so they are easy to test with mocks.

## Configuration

Required env vars: `BOT_TOKEN`, `ADMIN_USER_IDS` (comma-separated), `RABBITMQ_URL`.

Optional: `API_BASE_URL` (default `http://localhost:8081`), `HTTP_TIMEOUT`, `NOTIFICATION_GROUP_ID`. Queue names and group IDs for RabbitMQ consumers can be overridden via `TRADING_SIGNALS_QUEUE`, `TRADING_SIGNALS_GROUP_ID`, `PNL_REPORTS_QUEUE`, `PNL_REPORTS_GROUP_ID`, `SYSTEM_QUEUE`, `SYSTEM_GROUP_ID`.

## Development

```bash
# Run tests
make test
# or
go test ./...

# Lint (requires golangci-lint)
make lint

# Build
make build
```

## Deployment

Docker image is built and pushed on push to `main`; see `.github/workflows/ci.yml`. The workflow runs tests and vet before building the image.
