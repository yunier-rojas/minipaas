---
title: Variables
description: Configuration variables supported by the MiniPaaS Ansible roles.
weight: 6
---

## Overview

`minipaas-node/` and `minipaas-swarm/` read these variables.

### Node role

| Variable | Purpose | Default |
| --- | --- | --- |
| `node_ssh_harden` | Apply SSH hardening | `true` |
| `node_ssh_port` | SSH port opened by the firewall | `22` |
| `node_ssh_password_authentication` | Allow password authentication | `false` |
| `node_ssh_permit_root_login` | Root login policy | `prohibit-password` |
| `node_ssh_restrict` | Restrict SSH to `node_ssh_allowed_cidrs` | `false` |
| `node_ssh_allowed_cidrs` | IPv4 ranges allowed to reach SSH | `[]` |
| `node_firewall_tcp_ports` | Additional TCP ports for the firewall | `[]` |
| `node_firewall_udp_ports` | Additional UDP ports for the firewall | `[]` |
| `telegram_bot_token` | Telegram bot token for monitoring alerts | `TELEGRAM_TOKEN` env var |
| `telegram_chat_id` | Telegram chat ID for monitoring alerts | `TELEGRAM_CHAT` env var |
| `node_require_telegram_alerts` | Fail when Telegram credentials are missing | `false` |
| `node_fail2ban_bantime` | Initial ban duration | `1h` |
| `node_fail2ban_findtime` | Failure counting window | `10m` |
| `node_fail2ban_maxretry` | Failures before a ban | `4` |
| `node_fail2ban_bantime_maxtime` | Maximum ban duration | `4w` |
| `node_monitoring_script_url` | Monitoring script source | pinned commit |
| `node_monitoring_flags` | Monitoring thresholds and checks | CPU 80, RAM 70, disk 90 |

### Swarm role

| Variable | Purpose | Default |
| --- | --- | --- |
| `swarm_cronjob_version` | `swarm-cronjob` release on the primary manager | `1.14.0` |

Define them in `inventory.ini`, `group_vars/`, or `host_vars/`.

---

## Variable Details

### SSH

```yaml
node_ssh_harden: true
node_ssh_restrict: true
node_ssh_allowed_cidrs:
  - 203.0.113.10/32
```

### Firewall

A list of extra ports to allow on every node:

```yaml
node_firewall_tcp_ports:
  - 8000
  - 9090
node_firewall_udp_ports:
  - 51820
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

These variables customize SSH, the firewall, monitoring, and the `swarm-cronjob` release. All other provisioning tasks use built-in defaults.
