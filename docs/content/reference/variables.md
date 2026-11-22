---
title: Variables
description: Configuration variables supported by the MiniPaaS Ansible roles.
weight: 2
---

## Overview

`minipaas-node/` and `minipaas-swarm/` read these variables. Set them in `inventory.ini`, `group_vars/`, or `host_vars/`.

### Node role

| Variable | Purpose | Default |
| --- | --- | --- |
| `node_ssh_harden` | Apply SSH hardening | `true` |
| `node_ssh_port` | SSH port opened by the firewall | `22` |
| `node_ssh_restrict` | Restrict SSH to allowed CIDRs | `false` |
| `node_ssh_allowed_cidrs` | IPv4 ranges allowed to reach SSH | `[]` |
| `node_firewall_tcp_ports` | Additional TCP ports for nftables | `[]` |
| `node_firewall_udp_ports` | Additional UDP ports for nftables | `[]` |
| `telegram_bot_token` | Telegram bot token for monitoring alerts | `TELEGRAM_TOKEN` env var |
| `telegram_chat_id` | Telegram chat ID for monitoring alerts | `TELEGRAM_CHAT` env var |
| `node_require_telegram_alerts` | Fail when Telegram credentials are missing | `false` |
| `node_monitoring_script_url` | Monitoring script source | pinned commit |
| `node_monitoring_flags` | Monitoring thresholds and checks | CPU 80, RAM 70, disk 90 |

### Swarm role

| Variable | Purpose | Default |
| --- | --- | --- |
| `swarm_cronjob_version` | `swarm-cronjob` release on the primary manager | `1.14.0` |

Examples and behavior live with the role docs; see **[Role Configuration](/role/configuration/)**.
