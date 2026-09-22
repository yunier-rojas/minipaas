---
title: Installation
description: Install and run the MiniPaaS Ansible role to prepare servers for a minimal Docker Swarm cluster.
weight: 2
---

## Overview

The MiniPaaS role provisions servers so they are ready to run a **Docker Swarm cluster**.

The role installs:

- Docker
- Swarm cluster (init on the primary manager, join on workers)
- nftables firewall
- syslog-ng
- Fail2Ban
- Telegram monitoring and log alerts
- `swarm-cronjob` on the primary manager

---

## Requirements

You need:

- Linux hosts with SSH access
- Python 3 and Ansible installed locally
- An inventory file defining `managers` and `workers`
- (Optional) tokens/IDs for monitoring and alerting

Supported operating systems:

- Debian / Ubuntu family with systemd

---

## Install the Role Locally

Clone the repository:

```bash
git clone https://github.com/yunier-rojas/minipaas.git
cd minipaas
```

Install the Galaxy dependencies used by the role:

```bash
ansible-galaxy install -r minipaas-role/requirements.yml
ansible-galaxy role install geerlingguy.docker
```

---

## Prepare an Inventory

Create an `inventory.ini` in the repository root listing your hosts:

```ini
[managers]
manager1 ansible_host=1.2.3.4

[workers]
worker1 ansible_host=1.2.3.5
```

The role automatically:

* initializes Swarm on the manager node
* joins workers through the primary manager's worker join token

More details: **[Inventory](/role/inventory/)**

---

## Create a Playbook

Ansible applies the role through a playbook. Create `playbook.yml` in the repository root:

```yaml
- name: Install and Configure MiniPaaS
  hosts: all
  become: true
  roles:
    - ./minipaas-role
```

The `./minipaas-role` path is resolved relative to the playbook.

---

## Remote Access

The role binds Docker to the local Unix socket (`unix:///var/run/docker.sock`). The MiniPaaS CLI reaches a remote manager over SSH.

Point `ssh` at a manager host and use an `ssh://` target in `minipaas.yaml`:

```yaml
docker:
  host: ssh://deploy@manager.example.com
```

---

## Deploying Applications

After provisioning, the MiniPaaS CLI builds images, rolls out the stack, and authorizes Caddy routes. The Caddy service runs as part of the environment's stack.

See **[Commands](/cli/commands/)** for the deploy and code commands.

---

## Run the Role

Provision the entire cluster:

```bash
ansible-playbook -i inventory.ini playbook.yml
```

The role performs:

* Docker installation
* Swarm init on the primary manager
* Swarm join on workers
* nftables firewall configuration
* syslog-ng configuration
* Fail2Ban installation
* Telegram monitoring and log alerts (when credentials are set)
* `swarm-cronjob` installation on the primary manager

Re-running the playbook converges the nodes to the declared configuration.

---

## Verifying Installation

### Check Docker

```bash
docker info
```

### Check that Swarm is active

```bash
docker node ls
```

### Check that `swarm-cronjob` is running on the primary manager

```bash
systemctl status swarm-cronjob
```

### Check remote SSH access

From your workstation:

```bash
docker -H ssh://deploy@<manager-ip> info
```

---

## What Happens Next?

Once the role finishes:

1. Machines are hardened and configured
2. The Swarm cluster is ready
3. The CLI can build and roll out the stack

To proceed:

* Use the MiniPaaS CLI → **[MiniPaaS CLI Overview](/cli/)**

---

## Summary

Use the MiniPaaS Ansible role to:

* prepare machines
* bootstrap Swarm
* secure the environment
* install host-level cron orchestration
* standardize logs and firewall
