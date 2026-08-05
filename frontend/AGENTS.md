# AGENTS.md (frontend/)

Scope: everything under `frontend/`. This file overrides the root `AGENTS.md`
for this subtree; anything not restated here still applies (tests mandatory,
conventional commits, no secrets hardcoded, etc.).

## Stack

React + TypeScript dashboard. Consumes the Backend API described in the root
docs (`Monitoring Agent -> HTTP API -> Backend API -> PostgreSQL -> Frontend
Dashboard`). Single page for now: **Dashboard**, showing registered servers,
their CPU/RAM/Disk metrics, and health status.

## Structure

```
frontend/
├── src/
│   ├── components/       # one component per file, colocated tests
│   ├── hooks/             # data-fetching hooks (useServers, useServerMetrics, ...)
│   ├── services/           # API client(s), fetch/axios wrappers
│   ├── types/                # shared TS types/interfaces (Server, Metric, ...)
│   ├── pages/                  # top-level routed views (Dashboard)
│   ├── App.tsx
│   └── main.tsx
├── public/
├── package.json
└── tsconfig.json
```

## Conventions

- **Functional components only** — no class components.
- **One component per file**, colocated test: `ServerCard.tsx` +
  `ServerCard.test.tsx` in the same directory.
- Components are presentational. **Data fetching lives in `hooks/` or
  `services/`, never inline inside a component.**

```tsx
// components/ServerCard.tsx
export function ServerCard({ server }: { server: Server }) {
  return <div className="server-card">{server.hostname}</div>;
}
```

```tsx
// hooks/useServers.ts
export function useServers() {
  const [servers, setServers] = useState<Server[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchServers()
      .then(setServers)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, []);

  return { servers, loading, error };
}
```

- API calls go through `services/` (a thin client wrapper), not `fetch`
  scattered across components/hooks.
- Types for API payloads (`Server`, `Metric`, health status, etc.) live in
  `types/` and should mirror the backend `model` structs — if the backend
  response shape changes, update `types/` in the same PR.
- No global state library unless a real cross-page need shows up — plain
  hooks/context are enough for a single-page dashboard.

## Data displayed (per Server, per Phase 6 spec)

- Server list: hostname, IP, OS, status, last seen.
- Per-server metrics: CPU, RAM, Disk (from `GET /api/servers/{id}/metrics`).
- Health status should be visually distinct (e.g. healthy / stale / down)
  based on `last_seen` — keep the status-derivation logic in a hook/util, not
  inline JSX.

## Testing (mandatory)

Every new or changed component and hook ships with a matching `*.test.tsx` /
`*.test.ts` colocated alongside it. No exceptions — enforced at PR review.

- Component tests: React Testing Library, assert rendered output/behavior,
  not implementation details.
- Hook tests: mock the `services/` API client, assert loading/error/success
  states.
- Don't hit the real backend in tests — mock `services/`.

Before committing:
```bash
npm run lint       # must be clean
npm run test        # must pass
npm run build         # must succeed (tsc + bundler)
```

## Config

API base URL and any other env-specific values come from environment
variables (`.env`, gitignored, modeled on `.env.example`). Never hardcode the
backend URL in `services/`.

## Do not

- Write class components.
- Put `fetch`/`axios` calls directly inside a component or `.tsx` JSX file.
- Merge a new/changed component or hook without a colocated test.
- Hardcode API URLs or secrets.
- Introduce a state-management library speculatively — keep it to
  hooks/context until there's a proven cross-page need.