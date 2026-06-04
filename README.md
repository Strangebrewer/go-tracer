# go-tracer

A lightweight distributed tracing collector built as part of my portfolio project. Services emit spans fire-and-forget after each operation; the frontend polls this service to retrieve the full trace and render it as a live request timeline.

---

## Role in the Architecture

Every Go and NestJS service in the stack has the option to document its operations by sending a span to go-tracer after each request or Pub/Sub handler completes. The frontend generates a `traceId` (`crypto.randomUUID()`) before initiating an action, passes it downstream as an `X-Trace-ID` header, and polls `GET /traces/:traceId` once the action resolves.

Two separate auth mechanisms keep internal traffic and browser traffic isolated:

- `POST /spans` — requires an `X-Service-Key` header. Used by backend services only; never called from the browser.
- `GET /traces/{traceId}` — requires a JWT Bearer token. Called by the frontend dashboard after an action completes.

Spans auto-expire via a MongoDB TTL index on `startTime`. The default retention window is 1 hour, configurable via `SPAN_TTL_DAYS`.

---

## Stack

- **Language**: Go
- **Router**: [chi](https://github.com/go-chi/chi)
- **Database**: MongoDB Atlas via [mongo-driver v2](https://github.com/mongodb/mongo-go-driver) — no ORM
- **Auth**: dual-mode — RSA JWT validation for browser clients; shared service key for internal service calls
- **Logging**: `slog` with JSON output — Cloud Run ingests stdout directly into Cloud Logging

---

## Structure

```
cmd/
  server/        ← entry point: config, DB, middleware, server wiring
config/          ← environment variable loading
db_connection/   ← MongoDB connect, ping, index creation
health/          ← GET /health
middleware/      ← JWT auth, service key auth, request ID, structured logging
server/          ← chi router setup and route registration
span/            ← model, store, handler, routes, integration test
```

Indexes are created at startup — no migration step required.

---

## API

| Method | Path                | Auth                   | Description                                            |
| ------ | ------------------- | ---------------------- | ------------------------------------------------------ |
| `POST` | `/spans`            | `X-Service-Key` header | Persist a new span                                     |
| `GET`  | `/traces/{traceId}` | JWT Bearer             | Retrieve all spans for a trace, ordered by `startTime` |
| `GET`  | `/health`           | none                   | Health check                                           |

### Span shape

```json
{
  "traceId": "string",
  "spanId": "string",
  "parentSpanId": "string (optional)",
  "service": "string",
  "operation": "string",
  "status": "string",
  "error": "string (optional)",
  "startTime": "ISO 8601",
  "endTime": "ISO 8601",
  "metadata": {}
}
```

---

## Running Locally

Copy `.env.example` to `.env.local` and fill in values.

```bash
go run ./cmd/server
```

No migration step — indexes are created automatically on first connect.

```bash
go test ./...
```

---

## Testing

Integration tests spin up a real MongoDB container via [testcontainers](https://testcontainers.com) — the database is never mocked. `TestMain` handles the container lifecycle; individual tests operate against a real collection with real indexes.

---

## Environment Variables

| Variable          | Description                                                        |
| ----------------- | ------------------------------------------------------------------ |
| `PORT`            | HTTP port (defaults to `8080`)                                     |
| `MONGODB_URI`     | MongoDB connection string                                          |
| `DB_NAME`         | MongoDB database name (defaults to `tracer`)                       |
| `JWT_PUBLIC_KEY`  | RSA public key PEM for validating JWTs issued by go-auth           |
| `SERVICE_KEY`     | Shared secret required on `POST /spans` via `X-Service-Key` header |
| `ALLOWED_ORIGINS` | Comma-separated list of allowed CORS origins                       |
