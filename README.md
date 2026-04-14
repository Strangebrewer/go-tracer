# go-service-template

A GitHub template repository for Go REST services. Built as part of a larger portfolio project to establish consistent patterns across a multi-service architecture — and to deepen Go experience along the way.

Use the **Use this template** button on GitHub to create a new service repo.

---

## After Creating a New Service

1. Update the module name in `go.mod` (e.g. `github.com/Strangebrewer/go-auth`)
2. Find and replace all import paths from `go-service-template` to the new module name
3. Replace the `example/` package with your domain packages
4. Update `app/app.go` with your stores
5. Update `server/routes.go` to mount your routes

---

## Stack

- **Language**: Go
- **Router**: [chi](https://github.com/go-chi/chi)
- **Database**: Postgres via [pgx](https://github.com/jackc/pgx)
- **Query generation**: [sqlc](https://sqlc.dev) — SQL queries are written by hand and compiled to type-safe Go code
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Auth**: Stateless RSA JWT validation — tokens are issued by a dedicated auth service and verified independently per service, with no cross-service call per request
- **Logging**: `slog` with JSON output — Cloud Run ingests stdout directly into Cloud Logging

---

## Structure

```
cmd/
  server/      ← entry point: wiring only, no business logic
  migrate/     ← migration runner (up/down)
app/           ← aggregates domain stores into a single Application struct
server/        ← router setup and route registration
config/        ← environment variable loading
db_connection/ ← connection pool setup
db/
  schema.sql   ← source of truth for sqlc
  queries/     ← named SQL queries for sqlc to compile
  migrations/  ← golang-migrate up/down pairs
  generated/   ← sqlc output (committed)
health/        ← GET /health
middleware/    ← request ID, structured logging, JWT auth
example/       ← reference domain: model, store, handler, routes, integration test
```

Each domain follows the same four-file pattern — model, store, handler, routes — with a clean separation between HTTP concerns and data access. The store talks to the database; the handler talks to the store; the routes file wires them together. Auth is applied at the mount point, not buried inside individual route files.

---

## Running Locally

Copy `.env.example` to `.env.local` and fill in values.

```bash
# Start the server
go run ./cmd/server

# Run migrations
go run ./cmd/migrate up
go run ./cmd/migrate down   # rolls back one step

# Run tests
go test ./...
```

---

## Migrations

Migration files live in `db/migrations/` and follow golang-migrate naming:

```
000001_create_things.up.sql
000001_create_things.down.sql
```

---

## Database Queries (sqlc)

Write your schema in `db/schema.sql` and named queries in `db/queries/`. Then from the `db/` directory:

```bash
sqlc generate
```

This compiles your SQL into type-safe Go in `db/generated/`. The output is committed — no tooling required to build or run the service after generation.

---

## Testing

Integration tests spin up a real Postgres container via [testcontainers](https://testcontainers.com) — the database is never mocked. `TestMain` handles the container lifecycle; individual tests get a real store backed by a real schema.

The `example/` package includes a test skeleton with `t.Skip(...)` placeholders. Remove the skips once store methods are implemented.

Unit tests are written for pure logic with no database dependency — calculations, parsing utilities, and the like. Not for "does this handler call the store."

---

## Environment Variables

| Variable | Description |
|---|---|
| `PORT` | HTTP port (defaults to 8080) |
| `DATABASE_URL` | Postgres connection string |
| `JWT_PUBLIC_KEY` | RSA public key PEM for validating JWTs |
| `ALLOWED_ORIGINS` | Comma-separated list of allowed CORS origins |
