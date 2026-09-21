---
title: minipaas.yaml
description: Environment configuration file for MiniPaaS.
weight: 3
---

## Overview

`minipaas.yaml` is the **environment descriptor**. Each environment directory has one:

```
project/
├── compose.yaml
├── compose.build.yaml
└── minipaas/
    ├── minipaas.yaml
    └── compose.apps.yaml
```

The file controls:

- Which Compose files belong to the environment
- How the CLI connects to Docker
- Which values are exported for Compose interpolation

Paths in `compose.files` are relative to the environment directory, so the same file can be committed for every environment. Local-dev files one level up are written as `../compose.yaml`.

---

# Schema

```yaml
compose:
  files:
    - compose.common.yml
    - compose.apps.yaml

docker:
  host: ssh://deploy@1.2.3.4   # optional; omit to use the local Docker socket

namespace: minipaas   # optional; defaults to minipaas

vars:
  MINIPAAS_DEPLOY_VERSION: v1
  MINIPAAS_IMAGE_PREFIX: ghcr.io/org
```

---

## Field Reference

### `compose.files`

Compose files for this environment, listed in deep-merge order (base first, overrides later):

```yaml
compose:
  files:
    - ../compose.yaml
    - ../compose.build.yaml
    - compose.common.yml
    - compose.apps.yaml
    - compose.postgres.yml
    - compose.caddy.yml
```

- Relative paths are resolved against the environment directory.
- Absolute paths are used as-is.
- The CLI deep-merges the files in order for build, rollout, and canary.
- `deploy secret` and `deploy config` patch the *last* file that defines each service, which is normally the environment's swarm override.

---

### `docker.host`

Docker endpoint used by the environment. Omit it to use the local Docker socket:

```yaml
compose:
  files:
    - compose.apps.yaml
```

For a remote Swarm manager, set an SSH target:

```yaml
docker:
  host: ssh://deploy@manager.example.com
```

The CLI exports this as `DOCKER_HOST`. Docker connects to the remote daemon over SSH. The target needs a reachable SSH key and a known host. Any options supported by `ssh` (user, non-default port, `~/.ssh/config` aliases) work in the URL, for example `ssh://deploy@manager.example.com:2222`.

---

### `namespace`

Deployment namespace for the environment. It is used as the stack name, the service name prefix (`<namespace>_<service>`), and the overlay network prefix (`<namespace>_network`, `<namespace>_public`):

```yaml
namespace: acme
```

Defaults to `minipaas` when omitted. `code init --namespace` records it. Changing it lets multiple environments share one Swarm without name collisions.

---

### `vars`

Map of environment variables exported while CLI commands run:

```yaml
vars:
  MINIPAAS_DEPLOY_VERSION: v2
  MINIPAAS_IMAGE_PREFIX: ghcr.io/org
  LOG_LEVEL: debug
```

Keys are used verbatim. Compose files can interpolate them, for example in image tags:

```yaml
image: ${MINIPAAS_IMAGE_PREFIX}/example:${MINIPAAS_DEPLOY_VERSION}
```

`MINIPAAS_IMAGE_PREFIX` is set by `minipaas code init --registry`. It names the external registry (and optional namespace) that `deploy build` pushes to and that rollout pulls from. When it is omitted, built images are not pushed. `MINIPAAS_DEPLOY_VERSION` tags the built images.

---

## Full Example

```yaml
compose:
  files:
    - compose.common.yml
    - compose.apps.yaml
    - compose.postgres.yml

docker:
  host: ssh://deploy@10.0.0.5

namespace: minipaas

vars:
  MINIPAAS_DEPLOY_VERSION: v2025-02-01
  MINIPAAS_IMAGE_PREFIX: ghcr.io/org
```

---

## Best Practices

- Use one environment directory per deployment target (`dev/`, `staging/`, `prod/`).
- Commit `minipaas.yaml` to Git: it is part of the environment definition.
- List Compose files relative to the environment directory so they survive being moved.
- Omit `docker.host` for local development.
- Set `docker.host` to an `ssh://` target to deploy against a remote Swarm manager.
- Set `namespace` to run several environments in the same Swarm without name collisions.
- Optionally set `MINIPAAS_IMAGE_PREFIX` to the registry that built images are pushed to.
- Authenticate with `docker login <registry-host>` before `deploy build` when a registry prefix is set.
- Use `vars` for interpolation values such as the deploy version tag.
- Do not store secrets here: use `minipaas deploy secret` / `deploy config` instead.
