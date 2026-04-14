# go-service-template — Claude Context

## What This Repo Is

A GitHub template repository for Go REST services in the personal-enterprise project. It provides the full structural skeleton — chi router, JWT auth middleware, structured logging, Postgres connection, sqlc setup, golang-migrate, testcontainers integration tests, Dockerfile, and CI — so new services can be created without rebuilding boilerplate from scratch.

When a concrete service is created from this template, the first things to change are:
1. The module name in `go.mod` (e.g. `github.com/Strangebrewer/go-auth`) and all import paths throughout the codebase
2. The `example/` package — replace with real domain packages
3. `app/app.go` — replace `ExampleStore` with real stores
4. `server/routes.go` — replace example route registration with real routes

---

## Architecture

```
cmd/
  server/main.go     ← wiring only: config, DB, stores, server.New(), graceful shutdown
  migrate/main.go    ← runs golang-migrate up/down against DATABASE_URL
app/
  app.go             ← Application struct aggregating all domain stores
server/
  server.go          ← chi router, global middleware, returns *Server
  routes.go          ← mounts domain routes with auth middleware
config/
  config.go          ← loads .env.local + env vars into Config struct
db_connection/
  db.go              ← creates and pings a pgxpool.Pool
db/
  schema.sql         ← sqlc reads this for type generation
  sqlc.yaml          ← sqlc config: engine, query/schema paths, output package
  queries/           ← named SQL queries (.sql files) for sqlc to generate from
  migrations/        ← golang-migrate migration files (up + down pairs)
  generated/         ← sqlc output — committed, not gitignored
health/
  handler.go         ← GET /health → 200 {"status":"ok"}
middleware/
  auth.go            ← RequireAuth(pemString) parses RSA public key once, returns middleware
  logging.go         ← structured slog request logging with status code capture
  requestid.go       ← generates/propagates X-Request-ID, injects into context
example/
  example_model.go   ← Example struct + CreateExampleRequest/UpdateExampleRequest DTOs
  example_store.go   ← Store wrapping pgxpool.Pool; stubbed CRUD showing sqlc pattern
  example_handler.go ← Handler with GetAll/GetOne/Create/Update/Delete methods
  example_routes.go  ← Routes(store) returns chi.Router with all endpoints registered
  example_test.go    ← testcontainers integration test skeleton (tests skip until implemented)
```

---

## Key Patterns

### Adding a Domain

Follow the `example/` package as the structural template:

1. Create `<domain>/<domain>_model.go`, `_store.go`, `_handler.go`, `_routes.go`
2. Add `<Domain>Store *<domain>.Store` to `app/app.go`
3. Instantiate the store in `cmd/server/main.go` and add to `Application`
4. Mount routes in `server/routes.go`

### Auth

`middleware.RequireAuth(cfg.JWTPublicKey)` parses the RSA public key PEM once at startup and returns a `func(http.Handler) http.Handler`. It validates Bearer JWTs and injects the user ID (string UUID from the `sub` claim) into context. Retrieve it in handlers via `middleware.UserIDFromContext(r.Context())`.

Auth is applied at the mount point in `server/routes.go`, not inside individual route files.

### Database

- `db_connection.NewPool(cfg.DatabaseURL)` — call once in main, pass pool to stores
- `sqlc generate` — run from the `db/` directory after modifying schema or queries. Output goes to `db/generated/` and is committed.
- `go run ./cmd/migrate [up|down]` — run from the repo root. `down` rolls back one step.
- Migration files follow golang-migrate naming: `000001_create_things.up.sql` / `000001_create_things.down.sql`
- The migrate command rewrites `postgres://` → `pgx5://` internally — `DATABASE_URL` format in `.env.local` doesn't need to change.

### Logging

`slog.SetDefault(logger)` is called in `main.go` before anything else. After that, all packages can call `slog.Info/Error/etc.` directly without receiving a logger as a parameter. The logger writes JSON to stdout, which Cloud Run ingests into Cloud Logging automatically.

### Testing

Integration tests use testcontainers to spin up a real Postgres container. `TestMain` handles container lifecycle — start, apply schema, create pool, tear down. Individual test functions use `t.Skip(...)` as placeholders until store methods are implemented. Remove the skip calls when implementing a concrete service.

---

## Conventions

- File naming: `<domain>_handler.go`, `<domain>_store.go`, etc. — not just `handler.go`
- Default exports for nothing — this is Go, everything is package-scoped
- Receiver names: `h` for handlers, `s` for stores
- Store methods: `GetAll`, `GetOne`, `Create`, `Update`, `Delete`
- Handler methods follow the same verbs
- Routes function signature: `Routes(store *Store) chi.Router`
- Errors: log with `slog.Error` server-side, return generic message to client — never leak internal details
- Monetary values as integers (cents) if this service handles money

---

## Environment Variables

| Variable | Description |
|---|---|
| `PORT` | HTTP port (defaults to 8080) |
| `DATABASE_URL` | Postgres connection string (`postgres://user:pass@host/db`) |
| `JWT_PUBLIC_KEY` | RSA public key PEM for validating JWTs issued by go-auth |
| `ALLOWED_ORIGINS` | Comma-separated list of allowed CORS origins |

Copy `.env.example` to `.env.local` for local dev. Never commit `.env.local`.
