# AGENTS.md (backend/)

Scope: everything under `backend/`. This file overrides the root `AGENTS.md`
for this subtree; anything not restated here still applies (tests mandatory,
conventional commits, no secrets hardcoded, etc.).

## Stack

Go + [Chi router](https://github.com/go-chi/chi) + PostgreSQL. Part of the
`Monitoring Agent -> HTTP API -> Backend API -> PostgreSQL -> Frontend`
pipeline described in the root docs.

## Structure

```
backend/
├── cmd/
│   └── main.go          # entrypoint: config load, DB connect, router wiring, server start
├── internal/
│   ├── handler/          # HTTP layer: decode request, call service, encode response
│   ├── service/           # business logic
│   ├── repository/        # data access (SQL, PostgreSQL)
│   ├── model/              # structs: requests, responses, DB entities
│   └── db/                  # connection setup, migration runner
├── migrations/                # 00N_description.sql, additive-only
└── go.mod
```

## Layering rules

Strict one-way dependency: `handler -> service -> repository`.

- **handler**: parse/validate the HTTP request, call exactly one `service`
  method, write the response. No SQL, no business rules here.
- **service**: business logic and orchestration. No `net/http` types
  (`http.Request`, `http.ResponseWriter`) in this package.
- **repository**: all SQL lives here. Returns `model` structs or errors, never
  HTTP status codes.
- **model**: plain structs shared across layers (`ServerRegisterRequest`,
  `Server`, `Metric`, etc.).

```go
// handler/server.go
func (h *Handler) RegisterServer(w http.ResponseWriter, r *http.Request) {
    var req model.ServerRegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid body", http.StatusBadRequest)
        return
    }
    server, err := h.service.RegisterServer(r.Context(), req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(server)
}
```

## API & schema conventions

All routes live under `/api`, JSON in/out, and follow the same
`/api/<resource>` shape — go through `handler -> service -> repository`, don't
introduce a different pattern. Don't restate the current endpoint list or
route table here; that drifts out of sync with the code. If an agent needs
the current routes, read `internal/handler/` directly.

Schema changes go through `migrations/00N_description.sql`, additive-only —
never edit a migration that's already merged, add a new one instead. Don't
restate the current table/column layout here either; read the latest
migration files for the actual schema.

## Testing (mandatory)

Every new or changed function in `handler`, `service`, or `repository` ships
with a matching `*_test.go` in the same package. No exceptions — this is
enforced at PR review, not a suggestion.

- `handler` tests: use `httptest`, assert status code + response body, mock
  the `service` dependency (interface).
- `service` tests: mock the `repository` dependency (interface), test
  business logic in isolation.
- `repository` tests: prefer a real/test PostgreSQL instance (or testcontainers)
  over mocking SQL — this is the layer that should catch real query bugs.

Before committing:
```bash
gofmt -l .        # must print nothing
go vet ./...
go test ./...
```

## Config

DB connection string and other secrets come from environment variables
(loaded via `.env` locally, gitignored). Nothing sensitive hardcoded in
`cmd/main.go`, `internal/db/`, or anywhere else.

## Do not

- Put SQL in `handler/` or `service/`.
- Put `net/http` types in `service/` or `repository/`.
- Edit a migration file that's already merged.
- Merge a change to `handler`, `service`, or `repository` without a
  corresponding `*_test.go`.
- Hardcode DB credentials or connection strings.