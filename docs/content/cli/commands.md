---
title: Commands
description: Reference for all MiniPaaS CLI commands.
weight: 5
---

# Overview

MiniPaaS CLI commands are organized into functional groups:

- **deploy** → build images, roll out the stack, release canaries, apply Caddy routing, and manage Swarm secrets and configs  
- **code** → scaffold an environment, shape worker/job/cron services, and author Caddy routes  
- **shell** → open a Docker-ready shell for local or remote contexts  

All commands take the environment directory as their first positional argument: `<env>`, a directory that contains `minipaas.yaml`.

---

# Common Flags

- `<env>`: environment directory containing `minipaas.yaml` (first positional argument, all commands)  
- `--verbose`: detailed output (all commands)  

Command-specific flags:

- `--name <name>`: object and mounted file name (`deploy secret`, `deploy config`); defaults to the positional file's base name  
- `--for <value>`: Compose service to attach the secret or config to (`deploy secret`, `deploy config`); repeatable  
- `--replicas <n>`: replica count (`deploy canary`); defaults to `1`  
- `--cron <schedule>`: cron schedule (`code cron`); defaults to `* * * * *`  
- `-c, --compose-file <file>`: extra Compose file to register (`code init`); repeatable  
- `--local <bool>`: use the local Docker daemon (`code init`); defaults to `true`  
- `--host <user@host>`: SSH target used when `--local=false` (`code init`)  
- `--registry <prefix>`: image prefix for built images (`code init`), for example `ghcr.io/org`  
- `--namespace <name>`: deployment namespace for the stack, services, and networks (`code init`); defaults to `minipaas`  

---

# `minipaas deploy`: Build, rollout, canary, routing, secrets, configs

These commands load the environment's Compose files listed in `compose.files`.

### Build images

```bash
minipaas deploy build dev
```

Runs `docker build` for every service that declares a `build` section, tagging each image with its `image` value. When `MINIPAAS_IMAGE_PREFIX` is set, each built image is also pushed to that registry; authenticate first with `docker login <registry-host>`. When it is not set, images are built only and no push is performed. Images are built with host networking.

### Roll out the stack

```bash
minipaas deploy rollout dev
```

Runs `docker stack deploy --with-registry-auth` against the environment's Compose files using the environment's namespace (default `minipaas`) as the stack name.

### Release a canary

```bash
minipaas deploy canary dev --replicas 1 api worker
```

Runs `docker service update --with-registry-auth --replicas <n> --image <image>` for each named service.

### Apply routing

```bash
minipaas deploy routing dev
```

Posts the environment's `caddy.json` to the `minipaas_caddy` container admin endpoint, reloading its configuration.

### Create a secret

```bash
echo postgres | minipaas deploy secret dev \
  --name postgres_password \
  --for postgres --for api --for migrate
```

Creates a content-addressed Swarm secret named `<name>.<hash8>` and patches the last Compose file that defines each service. That is the environment's swarm override, so the project's local-dev Compose file stays free of Swarm references.

### Create a config

```bash
minipaas deploy config dev ./configs/app.json \
  --name app.json \
  --for api --for worker
```

Creates a content-addressed Swarm config and patches the last Compose file that defines each service.

---

# `minipaas code`: Compose helpers

### Initialize an environment

```bash
minipaas code init minipaas --registry ghcr.io/org
```

`code init` auto-detects the project's local-dev Compose files in the project root (the parent of `<env>`): `compose.yaml`/`compose.yml` and `compose.build.yaml`/`compose.build.yml`. Pass `-c` to register additional Compose files. Pass `--registry` to push built images to an external registry; it is stored as `MINIPAAS_IMAGE_PREFIX` in `minipaas.yaml`. Omit it to build images without pushing them.

Pass `--namespace` to deploy under a different stack name. The namespace is used for the stack, every service (`<namespace>_<service>`), and the overlay networks (`<namespace>_network`, `<namespace>_public`). It defaults to `minipaas` and is recorded in `minipaas.yaml`.

Writes an environment directory containing:

- `minipaas.yaml` listing the detected files, any `-c` files, and the generated base stack, in deep-merge order
- `compose.apps.yaml` built from the source files, with services moved to `minipaas_network`, built images tagged `${MINIPAAS_IMAGE_PREFIX}/<name>:${MINIPAAS_DEPLOY_VERSION}`, and the monitoring labels `alert.cpu_threshold` and `alert.memory_threshold` on each service
- `compose.common.yml`, `compose.postgres.yml`, `compose.caddy.yml`, and `caddy.json`

The CLI loads all registered files in order and deep-merges them, so the local-dev file at the project root and the swarm overrides in the environment directory form one project.

Pass `--local=false --host deploy@manager.example.com` to target a remote Docker daemon over SSH. The CLI records `docker.host: ssh://deploy@manager.example.com` in `minipaas.yaml`.

### Mark a service as a worker

```bash
minipaas code worker minipaas example-worker
```

Adds a replicated `deploy` block (1 replica, `start-first` rolling update, restart on any failure) to each service in `compose.apps.yaml`.

### Mark a service as a job

```bash
minipaas code job minipaas example-migration
```

Adds a one-shot `deploy` block (1 replica, `stop-first`, restart on failure up to 10 attempts) to each service in `compose.apps.yaml`.

### Schedule a service as a cron job

```bash
minipaas code cron minipaas --cron "*/5 * * * *" example-cron
```

Adds a `swarm-cronjob` `deploy` block (0 replicas plus `swarm.cronjob.*` labels) to each service in `compose.apps.yaml`. `--cron` defaults to `* * * * *`. The cluster must run `swarm-cronjob`; the Ansible role installs it on the host.

### Route a service through Caddy

```bash
minipaas code route dev https://example.com/app api:8080
```

Arguments:

1. Public URL used to reach the service. Pass a plain path such as `/web` to match every host: the host constraint is omitted. Hostless URLs use plain HTTP on port `80`, since TLS requires a domain; add a port to bind elsewhere, for example `http://:8000/web`.
2. `service[:port]` inside the environment's stack; the port defaults to `80`.

Adds or replaces the matching host/path route in the environment's `caddy.json`, pointing at `minipaas_<service>`. It also marks the service as resilient in `compose.apps.yaml`: a healthcheck against the target port plus a replicated `deploy` block (2 replicas, `start-first` rolling update, rollback on failure). Run `minipaas deploy routing` afterwards to apply the routing file.

---

# `minipaas shell`: Docker-ready shell

Open a shell with `DOCKER_HOST` set from the environment:

```bash
minipaas shell dev
```

Useful for debugging or CI environments.

---

# Notes

* Object names are content-addressed as `<name>.<hash8>`; new content creates a new name.
* All compose modifications are **idempotent**.
* Secrets/configs patch only the last compose file where each service appears.
* Secret/config and deploy commands expect the Docker engine to already be reachable.
* Deploy commands interpolate `vars` from `minipaas.yaml`, so `${MINIPAAS_DEPLOY_VERSION}` resolves during build, rollout, and canary.
* `code route` updates `caddy.json` and the target service's deploy block in `compose.apps.yaml`; `deploy routing` pushes the routing file to the running Caddy container.
* The CLI never overwrites unrelated sections inside compose files.
