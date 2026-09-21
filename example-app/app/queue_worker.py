#!/usr/bin/env python3
import json
import logging
import sys

import pgque

# Inline database configuration.
DB_HOST = "postgres"
DB_USER = "postgres"
DB_NAME = "postgres"
SSL_MODE = "disable"
SECRET_FILE = "/run/secrets/postgres_password"

QUEUE_NAME = "example_queue"
CONSUMER_NAME = "example-worker"
EVENT_TYPE = "task"

# Read the password from the secret file.
try:
    with open(SECRET_FILE, "r") as f:
        DB_PASSWORD = f.read().strip()
except Exception as e:
    sys.exit("Error reading secret file: " + str(e))

# Build the connection string.
DB_CONN_STR = (
    f"host={DB_HOST} user={DB_USER} password={DB_PASSWORD} "
    f"dbname={DB_NAME} sslmode={SSL_MODE}"
)


def parse_payload(raw):
    if isinstance(raw, str):
        try:
            return json.loads(raw)
        except json.JSONDecodeError:
            return raw
    return raw


def process_message(message):
    """Process one queue message. Extend with your business logic."""
    print(
        f"Processing job id {message.msg_id} "
        f"with payload: {parse_payload(message.payload)}",
        flush=True,
    )


def main():
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s %(levelname)s %(name)s %(message)s",
    )

    consumer = pgque.Consumer(
        dsn=DB_CONN_STR,
        queue=QUEUE_NAME,
        name=CONSUMER_NAME,
        poll_interval=5,
    )

    @consumer.on(EVENT_TYPE)
    def handle(message):
        process_message(message)

    consumer.start()


if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        print("\nExiting queue consumer.", flush=True)
        sys.exit(0)
