---
title: Usage
description: Practical end-to-end usage of the MiniPaaS CLI.
weight: 4
---

## Overview

This page shows how the MiniPaaS CLI is used across common tasks:

- environment layout
- secret and config provisioning
- image builds and deployments
- CI pipelines

Each command takes the environment directory containing `minipaas.yaml` as its first positional argument.

---

# 1. Environment Layout

A MiniPaaS project keeps its local-dev Compose files at the project root and its swarm files in an environment directory:

```
project/
├── compose.yaml          # local development
├── compose.build.yaml    # build contexts
└── minipaas/             # swarm environment
    ├── minipaas.yaml
    ├── compose.apps.yaml
    └── ...
```

Scaffold the environment:

```bash
minipaas code init minipaas --registry ghcr.io/org
```

`code init` detects `compose.yaml`/`compose.yml` and `compose.build.yaml`/`compose.build.yml` in the project root. Pass `-c` to register additional Compose files. Pass `--registry` to push built images to an external registry; omit it to build images without pushing. Pass `--namespace` to deploy under a different stack, service, and network prefix; it defaults to `minipaas`.

The command writes `minipaas.yaml`, `compose.apps.yaml`, the base stack files (`compose.common.yml`, `compose.postgres.yml`, `compose.caddy.yml`), and `caddy.json`. It registers the detected files first in `minipaas.yaml`, so all commands deep-merge the local-dev files with the swarm overrides and `deploy build` can reach their `build` sections.

Each service in `compose.apps.yaml` also gets the `alert.cpu_threshold` and `alert.memory_threshold` labels used by the role's monitoring script. Labels already defined in the source Compose files keep their values.

Pass `--local=false --host deploy@manager.example.com` to target a remote Swarm manager over SSH; the CLI records `docker.host: ssh://deploy@manager.example.com` in `minipaas.yaml`.

You can also maintain `minipaas.yaml` by hand and list your Compose files in deep-merge order:

```yaml
compose:
  files:
    - ../compose.yaml
    - ../compose.build.yaml
    - compose.common.yml
    - compose.apps.yaml
    - compose.postgres.yml
    - compose.caddy.yml

namespace: minipaas

vars:
  MINIPAAS_DEPLOY_VERSION: v1
  MINIPAAS_IMAGE_PREFIX: ghcr.io/org
```

See the **[minipaas.yaml specification](/cli/minipaas-file/)** for every field.

---

# 2. Managing Secrets & Configs

Define secrets/configs and attach them to specific services.

## Secrets

```bash
echo postgres | minipaas deploy secret dev \
  --name postgres_password \
  --for postgres \
  --for api \
  --for migrate
```

## Configs

```bash
minipaas deploy config dev ./configs/app.json \
  --name app.json \
  --for api
```

Secrets/configs are created in Swarm and Compose references are patched in one run.

---

# 3. Routing

Add a Caddy route for a service in the environment's stack:

```bash
minipaas code route dev https://example.com/app api:8080
```

`code route` updates the environment's `caddy.json` and marks the service as resilient in `compose.apps.yaml` (healthcheck plus a replicated deploy block). Apply the routing file to the running Caddy container with:

```bash
minipaas deploy routing dev
```

---

# 4. Building and Deploying

Set the deploy version through `vars` so image tags resolve at build and rollout time. Optionally set a registry prefix to push built images there:

```yaml
vars:
  MINIPAAS_DEPLOY_VERSION: v1
  MINIPAAS_IMAGE_PREFIX: ghcr.io/org
```

Authenticate against the registry before building (only needed when a prefix is set):

```bash
docker login ghcr.io
```

Build every service that declares a `build` section:

```bash
minipaas deploy build dev
```

Built images are tagged `${MINIPAAS_IMAGE_PREFIX}/<name>:${MINIPAAS_DEPLOY_VERSION}` and pushed to that registry. When `MINIPAAS_IMAGE_PREFIX` is not set, images are built with plain tags and no push is performed. Rollout and canary pass registry credentials to Swarm with `--with-registry-auth`.

Roll out the environment's Compose files as the environment's stack (namespace):

```bash
minipaas deploy rollout dev
```

Release additional replicas for a subset of services:

```bash
minipaas deploy canary dev --replicas 1 api
```

Apply the environment's `caddy.json` to the running `minipaas_caddy` container:

```bash
minipaas deploy routing dev
```

---

# 5. Workers, Jobs, and Cron

`code worker`, `code job`, and `code cron` add the matching Swarm `deploy` block to services in `compose.apps.yaml`:

```bash
minipaas code worker minipaas example-worker
minipaas code job minipaas example-migration
minipaas code cron minipaas --cron "*/5 * * * *" example-cron
```

* `worker` → long-lived replicated service, restart on any failure
* `job` → one-shot task, restart on failure with bounded attempts
* `cron` → 0 replicas plus `swarm.cronjob.*` labels, triggered by `swarm-cronjob`

Install `swarm-cronjob` through the Ansible role so cron schedules take effect.

---

# 6. CI / CD Pipeline Pattern

A typical CI pipeline:

```bash
# Provision secrets/configs
echo "${POSTGRES_PASSWORD}" | minipaas deploy secret prod \
  --name postgres_password --for postgres --for api

# Build and roll out
minipaas deploy build prod
minipaas deploy rollout prod
minipaas deploy routing prod
```

---

# Summary

MiniPaaS usage combines:

* `code init` → scaffold an environment
* `code worker` / `code job` / `code cron` → shape service deploy blocks
* `deploy secret` / `deploy config` → payload & metadata management
* `deploy` → build, rollout, canary, routing
* `code route` → author Caddy routes
* `shell` → open a Docker-ready shell

Use this page as a pattern reference; for command-level details see the **Commands** page.
