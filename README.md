# rsstt — RSS to Telegram

Watches RSS feeds and delivers new posts to your Telegram chat. Single-user, self-hosted, runs as one binary.

## Features

- Polls RSS feeds on a configurable interval and sends new items to Telegram
- Web admin UI to add, pause, or delete feeds (no bot commands needed)
- Concurrent feed fetching with a configurable concurrency limit
- Retry logic and rate limiting for the Telegram API
- SQLite storage — no external database required

## Quick start

Create `config.json`:

```json
{
  "TELEGRAM_BOT_TOKEN": "your_bot_token",
  "TELEGRAM_ADMIN_ID": "your_telegram_chat_id",
  "DATABASE_PATH": "data/feed.db",
  "LOG_LEVEL": "info"
}
```

`TELEGRAM_ADMIN_ID` is your personal Telegram chat ID (send `/start` to [@userinfobot](https://t.me/userinfobot) to find it).

Run:

```bash
go run .
```

The web admin opens at **http://localhost:8080** — add your feeds there.

## Docker

```bash
docker compose up -d
```

`compose.yml` mounts `./data/feed.db` for persistence and `./config.json` read-only.

## Configuration

All settings can be set via `config.json` or environment variables. Env vars take precedence over the file.

| Key | Default | Description |
|-----|---------|-------------|
| `TELEGRAM_BOT_TOKEN` | — | Bot token from [@BotFather](https://t.me/BotFather) |
| `TELEGRAM_ADMIN_ID` | — | Your Telegram chat ID |
| `DATABASE_PATH` | `data/feed.db` | SQLite database path |
| `LOG_LEVEL` | `info` | `debug`, `info`, `warn`, `error` |
| `ADMIN_PORT` | `8080` | Web admin port |
| `FEED_UPDATE_INTERVAL` | `5` | Feed poll interval (minutes) |
| `SEND_TO_USERS_INTERVAL` | `5` | Send interval (minutes) |
| `HTTP_TIMEOUT` | `30` | HTTP request timeout (seconds) |
| `MAX_RETRIES` | `3` | Telegram send retries |
| `MAX_CONCURRENT_FEEDS` | `5` | Parallel feed fetches |
| `MAX_MESSAGES_PER_USER` | `10` | Max items sent per interval |

## Project structure

```
rsstt/
├── main.go          — startup, goroutines, graceful shutdown
├── config.go        — config loading (file + env vars)
├── models/          — GORM models and DB connection
├── repository/      — database queries
├── service/         — feed fetching and Telegram delivery
├── web/             — web admin HTTP handlers and HTML template
└── logger/          — structured logging
```
