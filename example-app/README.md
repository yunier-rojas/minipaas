# Example App for MiniPaaS

This project demonstrates how to deploy a multi-service application using **MiniPaaS**, on a Docker Swarm environment. It covers local development, secrets, remote deployment over SSH, configuration management, service orchestration, and background job execution.

## Requirements

- Docker
- MiniPaaS CLI
  - ``cd ../minipaas-cli && go install ./cmd/minipaas``
- `hurl` utility

## Project Structure

```bash
example-app/
├── app/                      # Server, queue worker, stream consumer, ticker
├── db/                       # Dockerfile and SQL migrations for PostgreSQL
├── dev/                      # Production environment configuration
├── compose.yaml              # Base Compose file
├── compose.build.yaml        # Build-specific overrides
└── demo-script.sh            # Full automation for the demo
```

## Queues and Streams

The example uses [PgQue](https://github.com/NikolayS/pgque) to provide durable
queues and streams on PostgreSQL. The pinned PgQue schema (`v0.2.0`) is stored
in `db/migrations` and installed through dbmate migrations.

- `example_queue` carries work items. The `example-worker` service consumes them.
- `example_stream` carries events. The `example-consumer` service consumes them.
- `example-cron` runs `cleanup.py`, which ticks PgQue and runs maintenance.

HTTP endpoints:

- `POST /queue` sends a work item.
- `GET /queue` returns queue and consumer state.
- `POST /stream` publishes an event.
- `GET /stream` returns stream and consumer state.
- `GET /consumers` returns consumer cursors across queues.

## Infrastructure

```bash
docker swarm init
```

Creates a Docker Swarm in the local environment.

## Registry

Built images are pushed to an external registry. Set your registry prefix and authenticate before deploying:

```bash
export MINIPAAS_IMAGE_PREFIX=ghcr.io/your-user
docker login ghcr.io
```

`MINIPAAS_IMAGE_PREFIX` is stored in `dev/minipaas.yaml` and is exported while build, rollout, and canary commands run.

## MiniPaaS CLI

For more information, see the bash script `demo-script.sh`.
