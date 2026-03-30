#!/bin/sh
set -eu

export PGHOST="$DB_HOST"
export PGPORT="$DB_PORT"
export PGUSER="$DB_USER"
export PGDATABASE="$DB_NAME"
export PGPASSWORD="$DB_PASSWORD"
export PGSSLMODE="${DB_SSLMODE:-disable}"

until psql -c "SELECT 1" >/dev/null 2>&1; do
  echo "Waiting for postgres..."
  sleep 2
done

psql -v ON_ERROR_STOP=1 -c "
  CREATE TABLE IF NOT EXISTS schema_migrations (
    version TEXT PRIMARY KEY,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
  );
"

for file in /migrations/*.up.sql; do
  version="$(basename "$file")"

  if [ "$(psql -tAc "SELECT 1 FROM schema_migrations WHERE version='$version' LIMIT 1;")" = "1" ]; then
    echo "Skipping migration: $version"
    continue
  fi

  echo "Applying migration: $version"
  psql -v ON_ERROR_STOP=1 -f "$file"
  psql -v ON_ERROR_STOP=1 -c "INSERT INTO schema_migrations(version) VALUES ('$version');"
done