#!/usr/bin/env bash
set -euo pipefail

export PGPASSWORD="${POSTGRES_PASSWORD}"

until psql -v ON_ERROR_STOP=1 -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" -c 'SELECT 1' >/dev/null 2>&1; do
  echo "Waiting for PostgreSQL to accept connections..."
  sleep 2
done

for f in /migrations/*.sql; do
  echo "Applying $f"
  psql -v ON_ERROR_STOP=1 -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" -f "$f"
done
