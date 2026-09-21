---
title: Configuration
description: Variables used to configure the MiniPaaS Ansible role.
weight: 3
---

## Overview

The `minipaas-role/` exposes a small set of variables that control:

- Additional firewall ports
- Telegram monitoring and log alerts
- The `swarm-cronjob` release installed on the primary manager

Set them in:

- `group_vars/managers.yml`
- `group_vars/workers.yml`
- `host_vars/<hostname>.yml`
- or directly in `inventory.ini`

---

# Variables

### `minipaas_extra_ports`

Additional TCP ports to allow through the nftables firewall on every node.

```yaml
minipaas_extra_ports:
  - 8080
  - 9090
```

Defaults to `[]`. Set it when you publish services on host ports outside the Swarm and HTTP/HTTPS defaults.

---

### `telegram_bot_token` / `telegram_chat_id`

Credentials for the Telegram bot that receives monitoring reports and log alerts. When both are set, the role installs the monitoring script and forwards matching log entries to the chat.

```yaml
telegram_bot_token: "123:ABC"
telegram_chat_id: "987654321"
```

They default to the `TELEGRAM_TOKEN` and `TELEGRAM_CHAT` environment variables.

---

### `swarm_cronjob_version`

The `swarm-cronjob` release installed on the primary manager.

```yaml
swarm_cronjob_version: "1.14.0"
```

The default is the version tested with the role.

---

# Example Configuration

```yaml
minipaas_extra_ports:
  - 8080

telegram_bot_token: "123:ABC"
telegram_chat_id: "987654321"
```

---

# Best Practices

* Keep your inventory simple: one manager group, one worker group.
* Store sensitive variables in **Ansible Vault**.
* Re-run the role to apply changes.
* Keep configurations in Git for reproducibility.
