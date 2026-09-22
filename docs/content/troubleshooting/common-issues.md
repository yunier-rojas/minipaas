---
title: Common Issues
description: Symptoms, causes, and fixes for recurring MiniPaaS problems.
weight: 2
---

## Overview

Each section lists a symptom, the usual cause, and the fix.

---

## The CLI Cannot Reach Docker

Deploy and secret commands expect the Docker engine to be reachable. Check `docker.host` in `minipaas.yaml`; for SSH targets, confirm the key and known host, then test directly:

```bash
docker -H ssh://deploy@manager.example.com info
```

See **[minipaas.yaml](/cli/minipaas-file/)**.

---

## A New Secret Value Is Not Used

Secret and config names are content-addressed, so new content creates a new object and patches the Compose references. Redeploy the stack to apply the new reference:

```bash
minipaas deploy rollout dev
```

Remove the previous object afterwards. See **[Secrets Rotation](/guides/secrets-rotation/)**.

---

## Images Fail to Pull During Rollout

`deploy build` pushes only when `MINIPAAS_IMAGE_PREFIX` is set, and the build machine must be logged in. Rollout and canary pass credentials to Swarm with `--with-registry-auth`.

```bash
docker login ghcr.io
```

See **[CLI Usage](/cli/usage/)**.

---

## A Route Returns 404 or the Old Page

`code route` only writes `caddy.json`. Apply it to the running Caddy container:

```bash
minipaas deploy routing dev
```

Confirm the target is `service[:port]` inside the environment's namespace and that the service is running. See **[Caddy and TLS](/guides/caddy-tls/)**.

---

## Cron Services Never Run

`code cron` sets replicas to `0` and adds `swarm.cronjob.*` labels. The cluster needs `swarm-cronjob`, which the role installs on the primary manager:

```bash
systemctl status swarm-cronjob
```

See **[Swarm](/role/swarm/)**.

---

## Monitoring Alerts Are Missing

Both `telegram_bot_token` and `telegram_chat_id` must be set. Per-container thresholds are read from running containers, so redeploy after changing labels. See **[Monitoring](/role/monitoring/)**.

---

## A Worker Does Not Join

Swarm requires TCP 2377, TCP/UDP 7946, and UDP 4789 between nodes. Check that nftables allows them and that the worker can reach the primary manager. See **[Firewall](/role/firewall/)**.

---

## Each New Manager Forms Its Own Cluster

The role runs `swarm init` on every host in the `managers` group. Keep one host there and add capacity with `workers`. See **[Role Usage](/role/usage/)**.
