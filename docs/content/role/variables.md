---
title: Variables
description: Configuration variables supported by the MiniPaaS Ansible role.
weight: 6
---

## Overview

The MiniPaaS Ansible role (`minipaas-role/`) reads these variables:

| Variable                | Purpose                                        | Default                  |
|-------------------------|------------------------------------------------|--------------------------|
| `minipaas_extra_ports`  | Additional TCP ports for the firewall          | `[]`                     |
| `telegram_bot_token`    | Telegram bot token for monitoring alerts       | `TELEGRAM_TOKEN` env var |
| `telegram_chat_id`      | Telegram chat ID for monitoring alerts         | `TELEGRAM_CHAT` env var  |
| `swarm_cronjob_version` | `swarm-cronjob` release on the primary manager | `1.14.0`                 |

Define them in `inventory.ini`, `group_vars/`, or `host_vars/`.

---

## Variable Details

### Firewall

A list of extra TCP ports to allow on every node:

```yaml
minipaas_extra_ports:
  - 8000
  - 9090
```

### Telegram

Set both values to receive monitoring reports and alerts:

```yaml
telegram_bot_token: "123:ABC"
telegram_chat_id: "999111222"
```

### swarm-cronjob

`swarm_cronjob_version` selects the release downloaded on the primary manager.

---

## Summary

These variables customize the firewall, monitoring, and the `swarm-cronjob` release. All other provisioning tasks use built-in defaults.
