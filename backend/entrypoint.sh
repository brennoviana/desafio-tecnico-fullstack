#!/bin/sh

set -e

host="postgres"
port=5432

echo "Waiting for PostgreSQL at $host:$port..."

while ! nc -z $host $port; do
  sleep 1
done

echo "PostgreSQL is up - executing migrations"

goose -dir ./migrations postgres "postgres://admin:4cZZr2SVFqfDoLQg@postgres:5432/main_db?sslmode=disable" up

echo "Starting the app..."

exec ./app
