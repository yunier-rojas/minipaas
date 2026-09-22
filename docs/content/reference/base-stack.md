---
title: Base Stack Files
description: Files that minipaas code init generates inside an environment.
weight: 4
---

## Overview

`minipaas code init` writes the environment descriptor and a small base stack next to it.

| File | Contents |
| --- | --- |
| `minipaas.yaml` | Compose file list, Docker target, namespace, and vars |
| `compose.apps.yaml` | Services collected from the source Compose files with swarm overrides |
| `compose.common.yml` | `<namespace>_network` (internal) and `<namespace>_public` overlay networks |
| `compose.postgres.yml` | `postgres:17.4` with a named volume and `POSTGRES_PASSWORD_FILE` |
| `compose.caddy.yml` | `caddy:2.9.1-alpine` on a manager with data and config volumes |
| `caddy.json` | Empty Caddy server that `code route` fills in |

---

## Details

- `compose.common.yml` defines the two overlay networks; services in `compose.apps.yaml` are moved to `<namespace>_network`.
- `compose.postgres.yml` expects the password secret mounted at `/run/secrets/postgres_password`; create it with `deploy secret` and attach it to `postgres`.
- `compose.caddy.yml` constrains Caddy to a manager node and publishes container port `80` on host port `8000`; see **[Caddy and TLS](/guides/caddy-tls/)**.
- Image tags can be pinned with the `POSTGRES_IMAGE_TAG` and `CADDY_IMAGE_TAG` vars in `minipaas.yaml`.
- The files are ordinary Compose files: edit them or leave them out of `compose.files` when they are not needed.

---

## Next Steps

- **[minipaas.yaml](/cli/minipaas-file/)**: full descriptor reference.
- **[CLI Usage](/cli/usage/)**: environment layout and commands.
- **[Reference](/reference/)**: variables and compatibility.
