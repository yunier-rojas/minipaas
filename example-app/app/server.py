#!/usr/bin/env python3
import http.server
import json
import logging
from datetime import date, datetime, time as time_of_day, timedelta
from decimal import Decimal

import pgque
import psycopg

# Inline database configuration.
DB_HOST = "postgres"
DB_USER = "postgres"
DB_NAME = "postgres"
SSL_MODE = "disable"
SECRET_FILE = "/run/secrets/postgres_password"

# PgQue queue names, consumers and event types.
QUEUE_NAME = "example_queue"
QUEUE_CONSUMER = "example-worker"
QUEUE_EVENT_TYPE = "task"
STREAM_NAME = "example_stream"
STREAM_CONSUMER = "example-consumer"
STREAM_EVENT_TYPE = "event"

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


def jsonable(value):
    if isinstance(value, (datetime, date, time_of_day)):
        return value.isoformat()
    if isinstance(value, timedelta):
        return str(value)
    if isinstance(value, Decimal):
        return float(value)
    if isinstance(value, dict):
        return {k: jsonable(v) for k, v in value.items()}
    if isinstance(value, (list, tuple)):
        return [jsonable(v) for v in value]
    return value


def query(sql, params=None):
    """Run a single statement and return rows as JSON-friendly dicts."""
    with psycopg.connect(DB_CONN_STR, autocommit=True) as conn:
        cur = conn.execute(sql, params or ())
        if cur.description is None:
            return []
        columns = [desc.name for desc in cur.description]
        return [
            {name: jsonable(value) for name, value in zip(columns, row)}
            for row in cur.fetchall()
        ]


def publish(queue, event_type, payload):
    """Send one event to a PgQue queue and return its event id."""
    with pgque.connect(DB_CONN_STR, autocommit=True) as client:
        return client.send(queue, payload, type=event_type)


def queue_info(queue):
    return query("SELECT * FROM pgque.get_queue_info(%s);", (queue,))


def consumer_info(queue=None):
    if queue is None:
        return query("SELECT * FROM pgque.get_consumer_info();")
    return query("SELECT * FROM pgque.get_consumer_info(%s);", (queue,))


class RequestHandler(http.server.BaseHTTPRequestHandler):

    def do_GET(self):
        if self.path.startswith("/records"):
            self.handle_records_get()
        elif self.path.startswith("/queue"):
            self.handle_queue_get()
        elif self.path.startswith("/stream"):
            self.handle_stream_get()
        elif self.path.startswith("/consumers"):
            self.handle_consumers_get()
        elif self.path.startswith("/"):
            self.handle_health_get()
        else:
            self.send_error(404, "Resource not found")

    def do_POST(self):
        if self.path.startswith("/records"):
            self.handle_records_post()
        elif self.path.startswith("/queue"):
            self.handle_queue_post()
        elif self.path.startswith("/stream"):
            self.handle_stream_post()
        elif self.path.startswith("/error"):
            self.handle_error_post()
        else:
            self.send_error(404, "Resource not found")

    def send_json(self, payload, status=200):
        response_json = json.dumps(payload)
        self.send_response(status)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(response_json)))
        self.end_headers()
        self.wfile.write(response_json.encode("utf-8"))

    def read_json(self):
        try:
            content_length = int(self.headers.get("Content-Length", 0))
        except ValueError:
            content_length = 0
        body = self.rfile.read(content_length).decode("utf-8")
        return json.loads(body)

    # ----- /records Handlers -----
    def handle_records_get(self):
        try:
            records = query("SELECT id, data FROM records ORDER BY id;")
            self.send_json({"records": records})
        except Exception as e:
            self.send_error(500, str(e))

    def handle_records_post(self):
        try:
            payload = self.read_json()
        except Exception:
            self.send_error(400, "Invalid JSON")
            return
        if "data" not in payload:
            self.send_error(400, "Missing 'data' field")
            return
        try:
            rows = query(
                "INSERT INTO records (data) VALUES (%s) RETURNING id, data;",
                (payload["data"],),
            )
            self.send_json(rows[0])
        except Exception as e:
            self.send_error(500, str(e))

    # ----- /queue Handlers -----
    def handle_queue_get(self):
        try:
            self.send_json({
                "queue": queue_info(QUEUE_NAME),
                "consumers": consumer_info(QUEUE_NAME),
            })
        except Exception as e:
            self.send_error(500, str(e))

    def handle_queue_post(self):
        try:
            payload = self.read_json()
        except Exception:
            self.send_error(400, "Invalid JSON")
            return
        if "payload" not in payload:
            self.send_error(400, "Missing 'payload' field")
            return
        try:
            event_id = publish(QUEUE_NAME, QUEUE_EVENT_TYPE, payload["payload"])
            self.send_json({
                "id": event_id,
                "queue": QUEUE_NAME,
                "type": QUEUE_EVENT_TYPE,
                "payload": payload["payload"],
            })
        except Exception as e:
            self.send_error(500, str(e))

    # ----- /stream Handlers -----
    def handle_stream_get(self):
        try:
            self.send_json({
                "stream": queue_info(STREAM_NAME),
                "consumers": consumer_info(STREAM_NAME),
            })
        except Exception as e:
            self.send_error(500, str(e))

    def handle_stream_post(self):
        try:
            payload = self.read_json()
        except Exception:
            self.send_error(400, "Invalid JSON")
            return
        if "event" not in payload:
            self.send_error(400, "Missing 'event' field")
            return
        try:
            event_id = publish(STREAM_NAME, STREAM_EVENT_TYPE, payload["event"])
            self.send_json({
                "id": event_id,
                "queue": STREAM_NAME,
                "type": STREAM_EVENT_TYPE,
                "payload": payload["event"],
            })
        except Exception as e:
            self.send_error(500, str(e))

    # ----- /consumers Handlers -----
    def handle_consumers_get(self):
        try:
            self.send_json({
                "queue_consumer": QUEUE_CONSUMER,
                "stream_consumer": STREAM_CONSUMER,
                "consumers": consumer_info(),
            })
        except Exception as e:
            self.send_error(500, str(e))

    # ----- / Handlers (Health) -----
    def handle_health_get(self):
        try:
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.end_headers()
        except Exception as e:
            self.send_error(500, str(e))

    def handle_error_post(self):
        try:
            code = int(self.path.split("/")[-1])
            logging.error(f"Raising {code}!")
        except (ValueError, IndexError):
            code = 200

        try:
            self.send_response(code)
            self.send_header("Content-Type", "application/json")
            self.send_header("Content-Length", "0")
            self.end_headers()
        except Exception as e:
            self.send_error(500, str(e))


if __name__ == '__main__':
    port = 8080
    server_address = ("", port)
    try:
        conn = psycopg.connect(DB_CONN_STR, autocommit=True)
        conn.close()
        print("Database connection successful.")
    except Exception as e:
        print("Database not ready yet, starting anyway:", e)
    print(f"Starting server on port {port}...")
    httpd = http.server.HTTPServer(server_address, RequestHandler)
    httpd.serve_forever()
