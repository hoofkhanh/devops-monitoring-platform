# DevOps Monitoring Platform

Monitor a Linux machine’s CPU, memory, and disk usage with a lightweight agent, a Go API, PostgreSQL, and a React dashboard.

## Architecture

```text
Linux host
└── Monitoring agent
       │ POST /api/metrics
       ▼
┌──────────────────────────────────────────────┐
│ DevOps Monitoring Platform                   │
│                                              │
│  Nginx ──► Go backend ──► PostgreSQL         │
│                  │                           │
│                  └──────► React dashboard    │
└──────────────────────────────────────────────┘
```

- The monitoring agent collects CPU, memory, and disk utilization and posts it to the backend.
- The Go backend validates and stores metrics in PostgreSQL, then serves them to the dashboard.
- The React dashboard fetches the newest metric from the backend.
- Locally, Docker Compose starts PostgreSQL, pgAdmin, two backend containers, the frontend, and Nginx.
- For Kubernetes, the Helm chart deploys the frontend and backend, PostgreSQL StatefulSet, migration job, ConfigMaps, Secret, and Ingress.

For diagrams and deeper deployment context, see [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Run locally

### Prerequisites

- Docker Engine and Docker Compose v2
- Two local environment files: `backend/.env.docker` and `frontend/.env`

Set the backend environment variables in `backend/.env.docker`:

```dotenv
PORT=8080
DB_HOST=postgres
DB_PORT=5432
DB_NAME=devops_monitoring
DB_USER=your_database_user
DB_PASSWORD=your_database_password
APP_ENV=development

POSTGRES_DB=devops_monitoring
POSTGRES_USER=your_database_user
POSTGRES_PASSWORD=your_database_password
POSTGRES_PORT=5432

PGADMIN_DEFAULT_EMAIL=you@example.com
PGADMIN_DEFAULT_PASSWORD=choose_a_password
PGADMIN_PORT=5050
```

Set the frontend API URL in `frontend/.env`:

```dotenv
VITE_API_URL=http://localhost/api
VITE_FRONTEND_PORT=5173
```

in github/workflows/ci-pipeline.yaml
```
with:
  context: frontend
  push: true
  build-args: |
    VITE_API_URL=/api
```
Start the complete local stack from the repository root:

```bash
./docker/run-compose.sh
```

Open the services:

- Dashboard: `http://localhost:5173`
- API health check: `http://localhost/health`
- pgAdmin: `http://localhost:5050`

Stop the stack while retaining the database volumes:

```bash
cd docker
docker compose --env-file ../backend/.env.docker --env-file ../frontend/.env down
```

> The frontend reads `VITE_API_URL` at build time. The current Compose configuration does not pass it as a Docker build argument, so add it under `frontend.build.args` before using the containerized dashboard. The API stack can still be started and checked at `http://localhost/health`.

### Run the monitoring agent locally

After the local stack is running, configure [`monitoring-agent/collect-metrics.sh`](monitoring-agent/collect-metrics.sh) for your machine:

1. Change the `curl` URL to `http://localhost/api/metrics`.
2. Change the disk command’s `/dev/sdd` filter to the disk or mount point you want to monitor. Use `df -h` to identify it.

Then run the agent in the foreground:

```bash
chmod +x monitoring-agent/collect-metrics.sh
./monitoring-agent/collect-metrics.sh
```

It collects and submits a metric every five seconds. Keep the terminal open while it runs; press `Ctrl+C` to stop it. Confirm that the API is receiving data:

```bash
curl http://localhost/api/metrics
```

## Run automatically with scripts

These scripts deploy the Kubernetes version and install the monitoring agent as a `systemd` service. They require Linux, `sudo`, Helm 3, `kubectl`, and access to a Kubernetes cluster with an NGINX Ingress controller.

### Deploy and install the agent

```bash
chmod +x scripts/install.sh helm/run-helm-install.sh
./scripts/install.sh
```

The installer:

1. Creates and starts the `monitoring-agent` systemd service.
2. Prompts for the public backend hostname.
3. Creates database migration ConfigMaps when needed.
4. Installs or upgrades the Helm release named `dev` in the `devops-monitoring-platform` namespace.
5. Updates `/etc/hosts` and the monitoring-agent endpoint for the supplied hostname.

Check the agent after installation:

```bash
sudo systemctl status monitoring-agent
journalctl -u monitoring-agent -f
```

### Validate the Helm deployment

```bash
chmod +x helm/run-helm-test.sh
cd helm
./run-helm-test.sh
```

This runs Helm linting, renders the chart, and performs server-side dry-run validation against the cluster.

### Remove the deployment

```bash
chmod +x scripts/destroy.sh helm/run-helm-uninstall.sh
./scripts/destroy.sh
```

This stops the monitoring agent, removes the script-managed `/etc/hosts` entry, uninstalls the `dev` Helm release, and deletes the entire `devops-monitoring-platform` namespace.

> **Warning:** removing the namespace deletes all Kubernetes resources in it, including the PostgreSQL workload and its associated persistent-volume claim.
