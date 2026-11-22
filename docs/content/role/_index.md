---
title: Roles Overview
description: Harden Linux hosts and create a Docker Swarm cluster with two focused Ansible roles.
weight: 1
---

## Overview

MiniPaaS provides two Ansible roles:

- **`minipaas-node/`** hardens and monitors Linux hosts.
- **`minipaas-swarm/`** installs Docker and creates the Swarm cluster.

The split keeps host concerns separate from container orchestration. You can run the node role on its own, then add the Swarm role when you are ready to form the cluster.

The node role provides:

- SSH hardening
- nftables firewall with explicit TCP and UDP port lists
- Fail2Ban for SSH protection
- syslog-ng for system and Docker logs
- Telegram monitoring and log alerts

The swarm role provides:

- Docker CE installation
- Docker daemon configuration
- Swarm initialization and worker joining
- `swarm-cronjob` on the primary manager

This keeps both roles **simple, transparent, and reversible**, suitable for small clusters, side projects, and environments where clarity matters more than heavy automation.

---

## Features

### Host layer
`minipaas-node/` applies safe defaults to every host and exposes the ports you declare.

### Swarm layer
`minipaas-swarm/` creates the cluster:

- one primary manager
- any number of workers

### Application layer
The roles provision infrastructure only. Applications, routing, and secrets are handled by the MiniPaaS CLI:

- **[CLI Overview](/cli/)**: build, deploy, and manage the stack
- **[Commands](/cli/commands/)**: deploy, code, secret, config, shell

### Remote access
Docker stays bound to the local Unix socket. The CLI reaches a remote manager over SSH.

### Logging and security
- syslog-ng for system, Fail2Ban, and Docker logging
- Fail2Ban for SSH protection
- continuous system health monitoring
- Telegram notifications for logs and reports

### Host-level cron for Swarm tasks
`swarm-cronjob` runs on the primary manager, enabling cron-driven Swarm tasks once services define cron expressions.

---

## Motivation

The roles exist to:

- ensure servers are configured consistently
- separate host hardening from cluster orchestration
- create a consistent Swarm cluster
- establish safe security defaults
- separate **infrastructure provisioning** from **application deployment**

This simplicity makes them ideal for:

- prototyping
- hobby clusters
- homelabs
- small production workloads
- teams that want clear control over their infrastructure

---

## Behavior Summary

After both roles complete:

- SSH is hardened when `node_ssh_harden` is true
- the firewall allows only declared ports
- syslog-ng and Fail2Ban are running
- Telegram monitoring is active when credentials are set
- Docker is installed
- Swarm is initialized and joined
- `swarm-cronjob` runs on the primary manager

---

## Start Here

To begin using the roles:

- Learn how to install and execute them → **[Installation](/role/installation/)**
- Configure firewall, SSH, and Swarm variables → **[Configuration](/role/configuration/)**
- Understand what the swarm role does on each node → **[Swarm](/role/swarm/)**
- Review firewall behavior → **[Firewall](/role/firewall/)**
- Learn about system monitoring options → **[Monitoring](/role/monitoring/)**
- Understand the inventory structure → **[Inventory](/role/inventory/)**
- View all configurable variables → **[Variables](/role/variables/)**

For deploying applications into the cluster, use the MiniPaaS CLI:

- **[MiniPaaS CLI Overview](/cli/)**
