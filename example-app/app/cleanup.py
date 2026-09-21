#!/usr/bin/env python3
import sys
import time

import pgque

# Inline database configuration.
DB_HOST = "postgres"
DB_USER = "postgres"
DB_NAME = "postgres"
SSL_MODE = "disable"
SECRET_FILE = "/run/secrets/postgres_password"

# The cron service runs this script every minute. It ticks PgQue for a bit
# under a minute, then exits; the next cron run starts a fresh loop.
TICK_INTERVAL_SECONDS = 0.1
MAINT_INTERVAL_SECONDS = 30
RUN_SECONDS = 55

# Read the password from the secret file.
try:
    with open(SECRET_FILE, "r") as f:
        DB_PASSWORD = f.read().strip()
except Exception as e:
    raise Exception("Error reading secret file: " + str(e))

# Build the connection string inline.
DB_CONN_STR = (
    f"host={DB_HOST} user={DB_USER} password={DB_PASSWORD} "
    f"dbname={DB_NAME} sslmode={SSL_MODE}"
)


def run_maintenance(client):
    # Each call must commit on its own before the next tick.
    client.conn.execute("SELECT pgque.maint_retry_events();")
    client.conn.execute("SELECT pgque.maint();")
    print("PgQue maintenance complete", flush=True)


def run_ticker():
    client = pgque.connect(DB_CONN_STR, autocommit=True)
    started = time.monotonic()
    last_maintenance = 0.0
    try:
        while time.monotonic() - started < RUN_SECONDS:
            client.ticker_all()
            now = time.monotonic()
            if now - last_maintenance >= MAINT_INTERVAL_SECONDS:
                run_maintenance(client)
                last_maintenance = now
            time.sleep(TICK_INTERVAL_SECONDS)
    finally:
        client.close()
    print("PgQue ticker run complete", flush=True)


if __name__ == '__main__':
    run_ticker()
