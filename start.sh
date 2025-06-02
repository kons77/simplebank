#!/bin/sh

set -e

echo "=== Debug: Checking migration files ==="
echo "Contents of /app/migration/:"
ls -la /app/migration/ || echo "Migration directory not found"
echo ""

echo "Contents of /app/:"
ls -la /app/
echo ""

echo "=== Checking migrate binary ==="
/app/migrate --version
echo ""

echo "run db migration"
source /app/app.env
echo "DB_SOURCE: $DB_SOURCE"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"