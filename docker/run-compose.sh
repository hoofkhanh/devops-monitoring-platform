#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

docker compose \
  --env-file ../backend/.env \
  --env-file ../frontend/.env \
  up -d --build backend frontend
