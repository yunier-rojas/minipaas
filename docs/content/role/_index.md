---
title: Role Overview
description: Provision servers and create a Docker Swarm cluster.
weight: 1
---

## Overview

The MiniPaaS Ansible role (`minipaas-role/`) prepares Linux hosts to run a **Docker Swarm cluster**.

The role provides:

- Docker installation  
- Swarm initialization and worker joining  
- nftables firewall with safe defaults  
- syslog-ng logging  
- Fail2Ban for SSH protection  
- Telegram monitoring and log alerts  
- `swarm-cronjob` on the primary manager  

This keeps the role **simple, transparent, and reversible**, suitable for small clusters, side projects, and environments where clarity matters more than heavy automation.

---

## Features

### Swarm initialization  
The role creates a Swarm cluster:

- one primary manager  
- any number of workers  

### Application Layer  
The role provisions infrastructure only. Applications, routing, and secrets are handled by the MiniPaaS CLI:

- **[CLI Overview](/cli/)**: build, deploy, and manage the stack
- **[Commands](/cli/commands/)**: deploy, code, secret, config, shell

### Remote Access  
The role keeps Docker bound to the local Unix socket. The CLI reaches a remote manager over SSH.

### Firewall  
A default-deny nftables setup, allowing:

- SSH  
- required Swarm ports  
- user-defined ports  

### Logging & Security  
- syslog-ng for system and Docker logging  
- Fail2Ban for SSH protection  
- continuous system health monitoring  
- Telegram notifications for logs and reports  

### Host-level Cron for Swarm Tasks  
`swarm-cronjob` is installed on the host, enabling cron-driven Swarm tasks once services define cron expressions.

---

## Motivation

The MiniPaaS role exists to:

- ensure servers are configured consistently  
- create a consistent Swarm cluster  
- establish safe security defaults  
- separate **infrastructure provisioning** from **application deployment**  

This simplicity makes it ideal for:

- prototyping  
- hobby clusters  
- homelabs  
- small production workloads  
- teams that want clear control over their infrastructure  

---

## Behavior Summary

After the role completes:

- Docker is installed  
- Swarm is initialized and joined  
- Firewall and logging are configured  
- `swarm-cronjob` runs on the primary manager  

---

## Start Here

To begin using the MiniPaaS role:

- Learn how to install and execute the role → **[Installation](/role/installation/)**  
- Configure variables such as firewall and Swarm settings → **[Configuration](/role/configuration/)**  
- Understand what the role does on each node → **[Swarm](/role/swarm/)**  
- Review firewall behavior → **[Firewall](/role/firewall/)**  
- Learn about system monitoring options → **[Monitoring](/role/monitoring/)**  
- Understand the inventory structure → **[Inventory](/role/inventory/)**  
- View all configurable variables → **[Variables](/role/variables/)**  

For deploying applications into the cluster, use the MiniPaaS CLI:

- **[MiniPaaS CLI Overview](/cli/)**  
