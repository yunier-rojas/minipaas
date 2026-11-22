---
title: Installation
description: Install and run the MiniPaaS Ansible roles to harden hosts and create a minimal Docker Swarm cluster.
weight: 2
---

## Overview

MiniPaaS uses two roles:

- `minipaas-node/` prepares and hardens each host.
- `minipaas-swarm/` installs Docker and creates the Swarm cluster.

The node role installs:

- SSH hardening
- nftables firewall
- syslog-ng
- Fail2Ban
- Telegram monitoring and log alerts

The swarm role installs:

- Docker
- Swarm cluster (init on the primary manager, join on workers)
- `swarm-cronjob` on the primary manager

---

## Migrating from `minipaas-role`

The former `minipaas-role/` is now two roles. Replace the single role entry in your playbook with the node and swarm roles shown below.

`minipaas_extra_ports` is replaced by `node_firewall_tcp_ports` and `node_firewall_udp_ports`. Swarm ports must be declared explicitly. SSH hardening is enabled by default, so run `--check --tags node` and confirm key-based access first.

---

## Requirements

You need:

- Linux hosts with SSH access and key-based authentication
- Python 3 and Ansible installed locally
- An inventory file defining `managers` and `workers`
- (Optional) tokens/IDs for monitoring and alerting

Supported operating systems:

- Debian / Ubuntu family with systemd

---

## Install the Roles Locally

Clone the repository:

```bash
git clone https://github.com/yunier-rojas/minipaas.git
cd minipaas
```

Install the Galaxy dependencies used by the swarm role:

```bash
ansible-galaxy install -r minipaas-swarm/requirements.yml
ansible-galaxy role install geerlingguy.docker
```

---

## Prepare an Inventory

Create an `inventory.ini` listing your hosts:

```ini
[managers]
manager1 ansible_host=1.2.3.4

[workers]
worker1 ansible_host=1.2.3.5
```

The swarm role automatically:

* initializes Swarm on the manager node
* joins workers through the primary manager's worker join token

More details: **[Inventory](/role/inventory/)**

---

## Create a Playbook

Create `playbook.yml`:

```yaml
- name: Provision MiniPaaS
  hosts: all
  become: true
  roles:
    - role: ./minipaas-node
      tags: [node]
    - role: ./minipaas-swarm
      tags: [swarm]
```

Role paths are resolved relative to the playbook.

A complete example lives in `example-provisioning/`.

---

## Open the Required Ports

The node role opens only the ports you declare. Set them per group:

```yaml
# group_vars/managers.yml
node_firewall_tcp_ports: [2377, 7946, 80, 443]
node_firewall_udp_ports: [7946, 4789]
```

```yaml
# group_vars/workers.yml
node_firewall_tcp_ports: [7946]
node_firewall_udp_ports: [7946, 4789]
```

Add any service ports you publish on the host. More details: **[Firewall](/role/firewall/)**

---

## Telemetry and Alerts

Set both values to enable monitoring and log alerts:

```bash
export TELEGRAM_TOKEN="123:ABC"
export TELEGRAM_CHAT="999111222"
```

The variables default to the `TELEGRAM_TOKEN` and `TELEGRAM_CHAT` environment variables. More details: **[Monitoring](/role/monitoring/)**

---

## Remote Access

Docker binds to the local Unix socket (`unix:///var/run/docker.sock`). The MiniPaaS CLI reaches a remote manager over SSH.

Point `ssh` at a manager host and use an `ssh://` target in `minipaas.yaml`:

```yaml
docker:
  host: ssh://deploy@manager.example.com
```

---

## Run the Roles

Preview the node layer before changing SSH:

```bash
ansible-playbook -i inventory.ini playbook.yml --check --tags node
```

Provision the cluster:

```bash
ansible-playbook -i inventory.ini playbook.yml
```

The roles perform:

* SSH hardening
* nftables firewall configuration
* syslog-ng configuration
* Fail2Ban installation
* Telegram monitoring and log alerts (when credentials are set)
* Docker installation
* Swarm init on the primary manager
* Swarm join on workers
* `swarm-cronjob` installation on the primary manager

Re-running the playbook converges the nodes to the declared configuration.

SSH hardening is enabled by default. Run `--check --tags node` and confirm that key-based access works before the first real run, or set `node_ssh_harden: false`.

---

## Update OS Packages

Package upgrades are opt-in. Use the update playbook:

```bash
ansible-playbook -i example-provisioning/inventory.ini example-provisioning/update.yml
```

Or import the task file from your own playbook:

```yaml
- name: Update the system
  ansible.builtin.import_tasks: ../minipaas-node/tasks/system-update.yml
```

---

## Verifying Installation

### Check SSH and the firewall

```bash
sudo sshd -T
sudo nft list ruleset
```

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

Once the roles finish:

1. Machines are hardened and configured
2. The Swarm cluster is ready
3. The CLI can build and roll out the stack

To proceed:

* Use the MiniPaaS CLI → **[MiniPaaS CLI Overview](/cli/)**

---

## Summary

Use the MiniPaaS Ansible roles to:

* prepare and harden machines
* bootstrap Swarm
* secure the environment
* install host-level cron orchestration
* standardize logs and firewall
