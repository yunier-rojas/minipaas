---
title: Swarm
description: How the MiniPaaS role initializes and joins a Docker Swarm cluster.
weight: 7
---

## Overview

The MiniPaaS role provisions **Docker Swarm**, but does so in a minimal, infrastructure-only way:

- It installs Docker on each host.  
- It initializes a Swarm manager on the primary node.  
- It joins all worker nodes to the cluster.  
- It installs `swarm-cronjob` **on the host** as a systemd service.

The result is a Swarm cluster ready for the CLI to deploy applications.

This separation keeps the role focused on repeatable host provisioning and leaves application behavior to the CLI and the Caddy service in the stack.

---

## What the Role Does

### 1. Installs Docker
All required packages are installed so the host can run containers and participate in a Swarm cluster.

### 2. Initializes or Joins Swarm
- The primary manager performs `docker swarm init`.
- Worker nodes join with the worker join token.
- The cluster is left in a stable state.

Application services, stacks, and networks are created later by the CLI.

### 3. Configures Swarm Ports
The role ensures nftables allows required Swarm traffic:

- TCP 2377 (Swarm management)
- TCP/UDP 7946 (gossip)
- UDP 4789 (overlay networking)

### 4. Installs swarm-cronjob on the Host
- `swarm-cronjob` is installed on the primary manager, outside the Swarm services.
- It runs as a systemd service.
- It communicates with the Docker Engine over the local Unix socket.
- It triggers the cron schedules that services declare in `swarm.cronjob.*` labels.

---

## Cluster State After Provisioning

After the role finishes, your infrastructure looks like this:

### Managers
- Docker installed  
- Swarm initialized  
- nftables firewall active  
- syslog-ng, Fail2Ban, and Telegram monitoring (when configured)  
- `swarm-cronjob` installed on the primary manager  

### Workers
- Docker installed  
- Joined to the Swarm  
- nftables rules applied  
- syslog-ng installed  

The cluster is operational and ready for the CLI to deploy the stack.

---

## Deploying Into the Cluster

Once the Swarm is created by the role, the MiniPaaS CLI handles:

- building images (`minipaas deploy build`)  
- rolling out the stack (`minipaas deploy rollout`)  
- releasing canaries (`minipaas deploy canary`)  
- authoring and applying Caddy routing (`minipaas code route`, `minipaas deploy routing`)  
- creating Swarm secrets and configs (`minipaas deploy secret`, `minipaas deploy config`)  

The Caddy service runs inside the environment's stack and terminates HTTP routing for the services listed in `caddy.json`.

This ensures a clean split between:

- **Infrastructure provisioning** → role  
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

The MiniPaaS role prepares servers by:

- installing Docker  
- creating a Swarm cluster  
- configuring firewall and logging  
- installing the host-level cron scheduler for Swarm  

Applications, routing, volumes, and databases run in the **stack** and are managed by the **CLI**.

