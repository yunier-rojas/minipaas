---
title: Example App Walkthrough
description: Run the multi-service MiniPaaS example end to end.
weight: 2
---

## Overview

`example-app/` is a runnable multi-service project that exercises the full MiniPaaS workflow: local development, secrets, remote deployment over SSH, configuration, service orchestration, and background job execution.

It includes a web server, a queue worker, a stream consumer, and a cron service, all backed by PostgreSQL with PgQue.

---

## Prerequisites

- Docker and a running Swarm (`docker swarm init` for a local run)
- The MiniPaaS CLI: see **[CLI Installation](/cli/installation/)**
- `hurl` for the HTTP checks in `api.http`

---

## Run the Demo

From `example-app/`:

```bash
./demo-script.sh
```

The script initializes the `dev` environment with the `sample` namespace, routes the `example` service locally, shapes the worker, job, and cron services, creates the Postgres secret, builds and rolls out the stack, applies Caddy routing, and runs the Hurl tests.

For a remote environment, run the individual commands against an environment that targets a manager over SSH. See **[CLI Usage](/cli/usage/)**.

---

## What the Example Covers

- `example-worker` consumes `example_queue`
- `example-consumer` consumes `example_stream`
- `example-cron` runs `cleanup.py` on a schedule
- `example-migration` runs the PgQue and application migrations once
- HTTP endpoints: `POST /queue`, `GET /queue`, `POST /stream`, `GET /stream`, `GET /consumers`

---

## Next Steps

- **[Workers, Jobs, and Cron](/cli/usage/)**: shape service deploy blocks.
- **[Secrets Rotation](/guides/secrets-rotation/)**: replace values safely.
- **[Caddy and TLS](/guides/caddy-tls/)**: expose services on a domain.
