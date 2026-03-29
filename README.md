# crypto-knight-tg-bot

Telegram bot for interacting with the **Crypto Knight trading core**, providing trading signals, PnL reports, and system notifications in real-time.

---

## Badges / Status

![CI](https://github.com/antonhancharyk/crypto-knight-tg-bot/actions/workflows/ci.yml/badge.svg)
[![Codecov](https://codecov.io/gh/antonhancharyk/crypto-knight-tg-bot/branch/main/graph/badge.svg)](https://app.codecov.io/gh/antonhancharyk/crypto-knight-tg-bot)
[![Go Report Card](https://goreportcard.com/badge/github.com/antonhancharyk/crypto-knight-tg-bot)](https://goreportcard.com/report/github.com/antonhancharyk/crypto-knight-tg-bot)
![Go](https://img.shields.io/badge/go-1.26-blue)
![Docker Image Version](https://img.shields.io/docker/v/antgoncharik/crypto-knight-tgbot)
![Docker Image Size](https://img.shields.io/docker/image-size/antgoncharik/crypto-knight-tgbot/latest)
![Release](https://img.shields.io/github/v/release/antonhancharyk/crypto-knight-tg-bot)

---

## Architecture

The project follows a **clean/hexagonal architecture**:

- **`cmd/tgbot`** — entrypoint; loads config, wires dependencies, runs the app.
- **`internal/domain`** — core entities (e.g., `Report`).
- **`internal/ports`** — interfaces for external concerns: `Logger`, `ReportFetcher`.
- **`internal/usecase`** — business logic (e.g., report fetching and validation).
- **`internal/transport/telegram`** — Telegram bot handler (updates, callbacks, menus).
- **`internal/infra`** — implementations: HTTP client, RabbitMQ consumer, Zap logger, health server.

Use cases depend on **ports**, not concrete implementations → easy to test with mocks.

---

## Features

- Receive trading signals via Telegram
- PnL reports for your trading accounts
- Admin notifications
- Graceful shutdown and health monitoring
- Full test coverage with Codecov
- Docker-ready with healthchecks

---

## Configuration

Required environment variables:

| Variable           | Description |
|-------------------|-------------|
| `BOT_TOKEN`       | Telegram Bot API token |
| `ADMIN_USER_IDS`  | Comma-separated Telegram user IDs for admin notifications |
| `RABBITMQ_URL`    | RabbitMQ connection URL |

Optional variables:

| Variable                  | Default                     | Description |
|---------------------------|----------------------------|-------------|
| `API_BASE_URL`            | `http://localhost:8081`    | External API base URL |
| `HTTP_TIMEOUT`            | `10`                        | HTTP client timeout in seconds |
| `NOTIFICATION_GROUP_ID`   |                             | Telegram group ID for notifications |
| `TRADING_SIGNALS_QUEUE`   | `trading-signals-queue`     | Queue name for trading signals |
| `TRADING_SIGNALS_GROUP_ID`|                             | Consumer group ID |
| `PNL_REPORTS_QUEUE`       | `pnl-reports-queue`         | Queue name for PnL reports |
| `PNL_REPORTS_GROUP_ID`    |                             | Consumer group ID |
| `SYSTEM_QUEUE`            | `system-queue`              | Queue name for system messages |
| `SYSTEM_GROUP_ID`         |                             | Consumer group ID |
| `HEALTH_LISTEN_ADDR`      | `:8080`                     | Address for health endpoint |

---

## Installation

Clone the repository:

```bash
git clone https://github.com/antonhancharyk/crypto-knight-tg-bot.git
cd crypto-knight-tg-bot