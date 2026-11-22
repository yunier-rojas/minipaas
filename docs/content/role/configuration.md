---
title: Configuration
description: Variables used to configure the MiniPaaS Ansible roles.
weight: 3
---

## Overview

`minipaas-node/` controls:

- SSH hardening
- firewall ports
- Fail2Ban behavior
- Telegram monitoring and log alerts

`minipaas-swarm/` controls:

- the `swarm-cronjob` release installed on the primary manager

Set the variables in:

- `group_vars/managers.yml`
- `group_vars/workers.yml`
- `host_vars/<hostname>.yml`
- or directly in `inventory.ini`

---

# Node Role

### `node_ssh_harden`

Apply the SSH hardening drop-in. Defaults to `true`.

```yaml
node_ssh_harden: true
```

The node role disables password and keyboard-interactive authentication by default and keeps password root login disabled. Run `--check --tags node` first and confirm that key-based access works.

### `node_ssh_restrict` and `node_ssh_allowed_cidrs`

Restrict SSH to a list of IPv4 ranges:

```yaml
node_ssh_restrict: true
node_ssh_allowed_cidrs:
  - 203.0.113.10/32
  - 198.51.100.0/24
```

When `node_ssh_restrict` is true, the list must not be empty.

### `node_firewall_tcp_ports` and `node_firewall_udp_ports`

TCP and UDP ports to allow through the nftables firewall on every node.

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

Defaults to `[]`. The node role always allows the configured SSH port and local syslog traffic.

### `telegram_bot_token` / `telegram_chat_id`

Credentials for the Telegram bot that receives monitoring reports and log alerts. When both are set, the node role installs the monitoring script and forwards matching log entries to the chat.

```yaml
telegram_bot_token: "123:ABC"
telegram_chat_id: "987654321"
```

They default to the `TELEGRAM_TOKEN` and `TELEGRAM_CHAT` environment variables.

Set `node_require_telegram_alerts: true` to fail provisioning when the credentials are missing.

### Fail2Ban

| Variable | Default | Purpose |
| --- | --- | --- |
| `node_fail2ban_bantime` | `1h` | Initial ban duration |
| `node_fail2ban_findtime` | `10m` | Window used to count failures |
| `node_fail2ban_maxretry` | `4` | Failures before a ban |
| `node_fail2ban_bantime_maxtime` | `4w` | Maximum duration for repeated bans |
| `node_fail2ban_ignore_ips` | loopback ranges | Addresses that are never banned |

### Monitoring

| Variable | Purpose |
| --- | --- |
| `node_monitoring_script_url` | Source of `system-monitoring.sh` |
| `node_monitoring_flags` | Thresholds and checks passed to the script |

---

# Swarm Role

### `swarm_cronjob_version`

The `swarm-cronjob` release installed on the primary manager.

```yaml
swarm_cronjob_version: "1.14.0"
```

The default is the version tested with the swarm role.

---

# Example Configuration

`group_vars/managers.yml`:

```yaml
node_firewall_tcp_ports:
  - 2377
  - 7946
  - 80
  - 443
node_firewall_udp_ports:
  - 7946
  - 4789

telegram_bot_token: "123:ABC"
telegram_chat_id: "987654321"
```

`group_vars/workers.yml`:

```yaml
node_firewall_tcp_ports:
  - 7946
node_firewall_udp_ports:
  - 7946
  - 4789
```

---

# Best Practices

* Keep your inventory simple: one manager group, one worker group.
* Store sensitive variables in **Ansible Vault**.
* Re-run the roles to apply changes.
* Keep configurations in Git for reproducibility.
