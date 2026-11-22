---
title: Inventory
description: How to structure the Ansible inventory for MiniPaaS hosts.
weight: 5
---

## Overview

The MiniPaaS roles use a simple inventory layout.  
Hosts are grouped as **managers** or **workers**. Variables configure SSH, firewall ports, Telegram alerts, and the `swarm-cronjob` release; see **[Variables](/role/variables/)**.

This keeps the provisioning layer small, predictable, and easy to maintain.

---

## Basic Inventory Structure

A minimal `inventory.ini`:

```ini
[managers]
manager1 ansible_host=1.2.3.4

[workers]
worker1 ansible_host=1.2.3.5
```

The node role hardens each host. The swarm role performs `swarm init` or `swarm join` depending on the group.

---

## Group Variables

Group variables apply to all hosts in a group.
Examples:

### `group_vars/managers.yml`

```yaml
node_firewall_tcp_ports:
  - 2377
  - 7946
  - 80
  - 443
node_firewall_udp_ports:
  - 7946
  - 4789
```

### `group_vars/workers.yml`

```yaml
node_firewall_tcp_ports:
  - 7946
node_firewall_udp_ports:
  - 7946
  - 4789
```

---

## Host Variables

Use `host_vars/<hostname>.yml` to configure specific hosts.

Example: enabling Telegram alerts on a specific node:

```yaml
telegram_bot_token: "123:ABC"
telegram_chat_id: "999111222"
```

Example: restricting SSH to your office network:

```yaml
node_ssh_restrict: true
node_ssh_allowed_cidrs:
  - 203.0.113.10/32
```

---

## Inline Variables in the Inventory

For small setups, scalar variables can be set directly in `inventory.ini`:

```ini
[managers]
manager1 ansible_host=1.2.3.4

[workers]
worker1 ansible_host=1.2.3.5 ansible_user=deploy
```

Use `group_vars/` or `host_vars/` for list values such as `node_firewall_tcp_ports`.

Inline variables override defaults and group vars.

---

## Recommended Structure

```
inventory.ini
group_vars/
  managers.yml
  workers.yml
host_vars/
  manager1.yml
  worker1.yml
```

This approach keeps configurations organized and easy to scale.

---

## Summary

Your inventory defines:

* which hosts are managers
* which hosts are workers
* firewall ports
* optional SSH restrictions
* optional monitoring and alerting configuration

All other provisioning behavior is automatic and requires no additional settings.
