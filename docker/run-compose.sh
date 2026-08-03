#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"
docker compose --env-file ../backend/.env.backend up -d --build backend
