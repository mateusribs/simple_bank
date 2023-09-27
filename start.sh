#!/bin/sh

set -e

echo "run DB migration"
/app/migrate -path /app/migration -database "$DB_SOURCE" --verbose up

echo "start app"
exec "$@"