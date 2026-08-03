# AGENTS.md

Root-level instructions for AI coding agents working in this repo. Nested
`AGENTS.md` files (if added under `backend/`, `frontend/`, etc.) override this
one for their subtree.

## Project

`devops-monitoring-platform` — agent-based server monitoring. Flow:
`Monitoring Agent -> HTTP API -> Backend API -> PostgreSQL -> Frontend Dashboard`.
Shipped via `GitHub Actions -> Docker Hub -> Kubernetes (Helm)`.

## Structure

```
backend/            Go + Chi API (cmd/, internal/{handler,service,repository,model,db}, migrations/)
frontend/            React + TypeScript dashboard
monitoring-agent/    Collects CPU/RAM/Disk every 30s, installs as systemd service
docker/              Dockerfiles, docker-compose.yml
kubernetes/          Raw manifests (namespace, frontend, backend, database, ingress, config)
helm/                Chart.yaml, values.yaml, templates/
scripts/             install.sh, deploy.sh, backup.sh, health-check.sh, cleanup.sh
.github/workflows/   CI/CD
```

## Conventions

### backend/ — Go + Chi API
- Layered: `handler` (parse request/write response) → `service` (business logic)
  → `repository` (data access). No business logic in handlers.
```go
// handler/server.go
func (h *Handler) RegisterServer(w http.ResponseWriter, r *http.Request) {
    var req model.ServerRegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid body", http.StatusBadRequest)
        return
    }
    server, err := h.service.RegisterServer(r.Context(), req)
    ...
}
```
- **Tests are mandatory, not optional.** Every new/changed `handler`, `service`,
  and `repository` function must ship with a matching `*_test.go` in the same
  package. No PR merges with backend logic that has zero test coverage.
- Run `gofmt` and `go vet` clean, and `go test ./...` green, before committing.
- REST API under `/api`, JSON in/out. Existing endpoints:
  `GET /health`, `POST /api/servers/register`, `POST /api/metrics`,
  `GET /api/servers`, `GET /api/servers/{id}/metrics`. New endpoints follow the
  same `/api/<resource>` shape.
- Schema changes are additive-only migrations in `backend/migrations/`, named
  `00N_description.sql`. Never edit a migration that's already been merged —
  add a new one instead.

### frontend/ — React + TypeScript dashboard
- Functional components only, one component per file, colocated tests
  (`ServerCard.tsx` + `ServerCard.test.tsx`).
```tsx
// components/ServerCard.tsx
export function ServerCard({ server }: { server: Server }) {
  return <div className="server-card">{server.hostname}</div>;
}
```
- Data fetching lives in hooks/services, not inside components.

### monitoring-agent/ — metrics collector
- Collects CPU/RAM/Disk on a fixed interval (30s) and posts to the backend
  HTTP API; keep the collect/send/sleep loop simple and side-effect-free
  outside that loop.
- Ships as a `systemd` service (`monitor-agent.service`, `Restart=always`).
  Any change to the binary's flags/env vars must be reflected in `install.sh`
  and the unit file together.

### docker/ — Dockerfiles & compose
- Multi-stage builds, minimal final images, no build tools in the runtime stage.
- `docker-compose.yml` is the local dev source of truth: services are
  `frontend`, `backend`, `postgres`, `nginx` — keep it runnable with a single
  `docker compose up`.

### kubernetes/ — raw manifests
- One resource kind per file, grouped by app (`frontend/`, `backend/`,
  `database/`, `ingress/`, `config/`).
- PostgreSQL is a `StatefulSet` with a `PVC` — never a plain `Deployment`.
- Config in `ConfigMap`, secrets in `Secret`. Nothing sensitive hardcoded in
  any manifest.

### helm/ — chart
- `values.yaml` is the single source of truth for env-specific values (image
  tag, replicas, resource limits). Templates read from values — don't
  hardcode environment-specific settings in `templates/`.
- Any manifest added under `kubernetes/` should have a corresponding template
  here; don't let the two drift apart.

### scripts/ — automation
- `install.sh`, `deploy.sh`, `backup.sh`, `health-check.sh`, `cleanup.sh` must
  each start with:
```bash
#!/usr/bin/env bash
set -euo pipefail
```
- Idempotent where practical — safe to re-run without side effects.

### .github/workflows/ — CI/CD
- Pipeline order: tests → security scan (Trivy + SonarQube) → build → push to
  Docker Hub → deploy via Helm. Don't reorder so that build/push happens
  before tests and scans pass.
- Never weaken or remove the Trivy/SonarQube steps to make a build pass.

## Config & Secrets

Local config goes through a `.env` file, gitignored, modeled on a committed
`.env.example`. Nothing sensitive is ever hardcoded in code or YAML.

## Commit / PR

- Conventional commits: `feat:`, `fix:`, `chore:`, `docs:`.
- PRs must pass: tests (including new backend tests), Trivy scan, SonarQube.
- Squash merge only.

## Best Practices — do not

- Commit `.env` files, kubeconfig, or any credential/secret value.
- Merge backend changes without accompanying tests.
- Edit an already-merged migration file.
- Put business logic in HTTP handlers — it belongs in `service/`.
- Expose services via `NodePort` in manifests meant for production.