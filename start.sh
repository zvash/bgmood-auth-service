#!/bin/sh

set -e

echo "migrate db"
/app/migrate -path db/migration -database "$DB_SOURCE" -verbose up
echo "start the app"
exec "$@"