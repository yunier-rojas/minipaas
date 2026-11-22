---
title: Swarm
description: How the MiniPaaS swarm role initializes and joins a Docker Swarm cluster.
weight: 7
---

## Overview

`minipaas-swarm/` provisions **Docker Swarm** in a minimal, infrastructure-only way:

- It installs Docker on each host.
- It initializes a Swarm manager on the primary node.
- It joins all worker nodes to the cluster.
- It installs `swarm-cronjob` **on the host** as a systemd service.

The result is a Swarm cluster ready for the CLI to deploy applications.

Host hardening, firewall, syslog, and monitoring belong to `minipaas-node/`. That separation keeps each role focused on one layer and leaves application behavior to the CLI and the Caddy service in the stack.

---

## What the Role Does

### 1. Installs Docker
All required packages are installed so the host can run containers and participate in a Swarm cluster. The swarm role also writes `/etc/docker/daemon.json` and the custom `docker.service` unit.

### 2. Initializes or Joins Swarm
- The primary manager performs `docker swarm init`.
- Worker nodes join with the worker join token.
- The cluster is left in a stable state.

Application services, stacks, and networks are created later by the CLI.

### 3. Installs swarm-cronjob on the Host
- `swarm-cronjob` is installed on the primary manager, outside the Swarm services.
- It runs as a systemd service.
- It communicates with the Docker Engine over the local Unix socket.
- It triggers the cron schedules that services declare in `swarm.cronjob.*` labels.

---

## Firewall Ports

The swarm role does not manage the firewall. Declare the Swarm ports in the node role configuration:

| Purpose | Port | Protocol | Group |
| --- | --- | --- | --- |
| Swarm management API | 2377 | TCP | managers |
| Gossip and node discovery | 7946 | TCP and UDP | managers and workers |
| Overlay networking | 4789 | UDP | managers and workers |
| HTTP and HTTPS | 80, 443 | TCP | managers |

See **[Firewall](/role/firewall/)** for a complete example.

---

## Cluster State After Provisioning

After both roles finish, your infrastructure looks like this:

### Managers
- SSH hardened
- nftables firewall active
- syslog-ng, Fail2Ban, and Telegram monitoring (when configured)
- Docker installed
- Swarm initialized
- `swarm-cronjob` installed on the primary manager

### Workers
- SSH hardened
- nftables firewall active
- syslog-ng, Fail2Ban, and Telegram monitoring (when configured)
- Docker installed
- Joined to the Swarm

The cluster is operational and ready for the CLI to deploy the stack.

---

## Deploying Into the Cluster

Once the Swarm is created by the swarm role, the MiniPaaS CLI handles:

- building images (`minipaas deploy build`)
- rolling out the stack (`minipaas deploy rollout`)
- releasing canaries (`minipaas deploy canary`)
- authoring and applying Caddy routing (`minipaas code route`, `minipaas deploy routing`)
- creating Swarm secrets and configs (`minipaas deploy secret`, `minipaas deploy config`)

The Caddy service runs inside the environment's stack and terminates HTTP routing for the services listed in `caddy.json`.

This ensures a clean split between:

- **Infrastructure provisioning** → roles
- **Application deployment and routing** → CLI + Caddy

---

## Best Practices

- Use a single manager; add capacity with workers.
- Deploy runtime components with the CLI after provisioning.
- Let the CLI create and maintain services, networks, and routes.
- Keep Docker bound to the local Unix socket and connect remotely over SSH.
- Treat `swarm-cronjob` as part of the infrastructure layer: it runs cron-triggered tasks for services deployed by the CLI.

---

## Summary

`minipaas-swarm/` prepares servers by:

- installing Docker
- creating a Swarm cluster
- installing the host-level cron scheduler for Swarm

Applications, routing, volumes, and databases run in the **stack** and are managed by the **CLI**.
