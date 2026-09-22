---
title: Variables
description: Configuration variables supported by the MiniPaaS Ansible role.
weight: 2
---

## Overview

The `minipaas-role/` reads four variables. Set them in `inventory.ini`, `group_vars/`, or `host_vars/`.

| Variable | Purpose | Default |
| --- | --- | --- |
| `minipaas_extra_ports` | Additional TCP ports allowed through nftables | `[]` |
| `telegram_bot_token` | Telegram bot token for monitoring alerts | `TELEGRAM_TOKEN` env var |
| `telegram_chat_id` | Telegram chat ID for monitoring alerts | `TELEGRAM_CHAT` env var |
| `swarm_cronjob_version` | `swarm-cronjob` release on the primary manager | `1.14.0` |

Examples and behavior live with the role docs; see **[Role Configuration](/role/configuration/)**.
