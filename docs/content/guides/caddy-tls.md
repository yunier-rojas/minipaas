---
title: Caddy and TLS
description: Expose services through the environment's Caddy instance.
weight: 4
---

## Overview

Caddy runs inside the environment's stack and serves the routes in `caddy.json`. `code route` authors routes; `deploy routing` posts the file to the running Caddy container.

Route behavior:

- A route with a hostname keeps automatic HTTPS enabled, so Caddy obtains certificates through ACME.
- A route with a plain path and no host matches every host and uses plain HTTP on port `80`.
- Add a port to bind elsewhere, for example `http://:8000/web`.

---

## Add a Route

```bash
minipaas code route dev https://app.example.com api:8080
minipaas deploy routing dev
```

`code route` updates `caddy.json` and marks the target service as resilient in `compose.apps.yaml`: a healthcheck on the target port plus a replicated deploy block. Nothing changes on the running cluster until `deploy routing` reloads Caddy.

---

## The Caddy Service

`code init` writes `compose.caddy.yml`, which runs `caddy:2.9.1-alpine` on a manager node, mounts the `caddy_data` and `caddy_config` volumes, and publishes container port `80` on host port `8000`.

For hostname routes that terminate TLS, publish container port `443` as well. The role's firewall opens HTTP and HTTPS on managers by default; see **[Firewall](/role/firewall/)**.

---

## Verify

- Inspect the generated routes in the environment's `caddy.json`.
- Confirm the service is running: `docker service ls`.
- Request the public URL from a client.

---

## Next Steps

- **[Commands](/cli/commands/)**: full command reference.
- **[Firewall](/role/firewall/)**: ports opened by the role.
- **[Base Stack Files](/reference/base-stack/)**: what `code init` writes.
