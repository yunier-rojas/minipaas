---
title: Monitoring
description: Telegram-based monitoring and log alerts from the MiniPaaS role.
weight: 9
---

## Overview

The MiniPaaS role adds lightweight monitoring to each node:

- **syslog-ng** collects system logs and the Docker log stream, and forwards matching entries to Telegram
- a **system monitoring script** watches CPU, memory, disk, and Docker container usage, and alerts on threshold breaches
- **Fail2Ban** is installed for brute-force protection

Telegram is the alert channel. Set `telegram_bot_token` and `telegram_chat_id` to activate log forwarding and the boot report.

---

## Syslog-ng

The role installs syslog-ng on every node and deploys `/etc/syslog-ng/syslog-ng.conf`. The configuration:

- reads system and internal logs
- reads the Docker log stream on `udp/514`
- forwards container messages tagged `minipaas-*` that contain `ERROR`, `CRITICAL`, `FATAL`, or `PANIC`
- forwards error-level logs from any program

The Docker daemon writes container logs to `udp://127.0.0.1:514` with the tag `minipaas-<container>`. The tag is set in `/etc/docker/daemon.json`, so it applies to every container on the node, including deployments with a custom namespace, whose container names become `<namespace>_<service>`.

---

## Monitoring Script

With Telegram credentials set, the role installs:

- `/usr/local/bin/system-monitoring.sh`: the monitoring script
- `/usr/local/bin/server-monitoring.sh`: starts the script with server thresholds: CPU 80%, RAM 70%, disk 90%, plus load average, reboot, SSH login, and Docker container checks

A cron entry in `/etc/cron.d/monitoring-telegram` starts the script at every reboot, and the script keeps running. It samples CPU and memory every 60 seconds.

### Container Monitoring

The script runs with `--DOCKER-MONITOR` and checks every running container. Tune the thresholds per container with Docker labels:

| Label                   | Meaning                                                             | Default   |
|-------------------------|---------------------------------------------------------------------|-----------|
| `alert.cpu_threshold`   | CPU usage threshold in percent; a `%` suffix is optional            | `90`      |
| `alert.memory_threshold`| Memory usage threshold; accepts `B`, `KiB`, `MiB`, `GiB`, and similar units | `1000MiB` |

Add the labels to a service in `compose.apps.yaml`:

```yaml
services:
  example:
    labels:
      alert.cpu_threshold: "85"
      alert.memory_threshold: "1GiB"
```

A container alerts when usage is greater than its threshold. Alerts arrive as `DOCKER-CPU` and `DOCKER-MEMORY` messages with a 10-minute cooldown per container and metric. Missing or invalid labels fall back to the defaults.

`minipaas code init` writes both labels with the defaults above into `compose.apps.yaml`; values already set in the source Compose files are kept.

Labels are read from the running container, so redeploy the stack (`minipaas deploy rollout`) after changing them.

---

## Telegram Alerts

```yaml
telegram_bot_token: "123:ABC"
telegram_chat_id: "999111222"
```

Set both variables to enable alerts. Store them in Ansible Vault or the playbook environment.

---

## Fail2Ban

The role installs the Fail2Ban package and restarts the service. Jails follow the package defaults and any files you add under `/etc/fail2ban/`. Fail2Ban logs flow through syslog-ng like other system logs.

---

## Adjusting the Report

Change the thresholds in `minipaas-role/tasks/monitoring.yml` and re-run the playbook. Edit `/etc/cron.d/monitoring-telegram` on a node to change when the script runs.

---

## Summary

The role provides Telegram monitoring, syslog collection, and SSH protection for small clusters. Components are standard packages and files, easy to inspect and remove.
