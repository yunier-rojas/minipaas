#!/bin/sh
set -u

export PASSWORD=$(cat /run/secrets/postgres_password)

export DATABASE_URL="postgres://postgres:${PASSWORD}@postgres/postgres?sslmode=disable"

attempts=15
until dbmate up; do
    attempts=$((attempts - 1))
    if [ "$attempts" -le 0 ]; then
        echo "migrations failed after retries"
        exit 1
    fi
    echo "waiting for database..."
    sleep 2
done
