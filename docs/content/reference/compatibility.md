---
title: Compatibility
description: Platforms and requirements for the MiniPaaS components.
weight: 5
---

## Overview

The two components have separate requirements. This page summarizes what each one targets.

---

## Role

- Linux hosts with SSH access and systemd
- Debian and Ubuntu families
- Python 3 and Ansible on the control machine
- Ansible collection `community.docker` and the `geerlingguy.docker` role

---

## CLI

- Go 1.25 or newer to build from source
- Prebuilt binaries for linux-amd64, darwin-arm64, and windows-amd64
- Docker CLI on the machine that runs the commands
- An SSH client and key for `ssh://` targets

---

## Cluster

- A Docker Swarm cluster; the CLI works with any reachable Docker endpoint
- `swarm-cronjob` for cron-triggered services, installed by the role (default `1.14.0`)

---

## Notes

- The role targets a single manager; each host in `managers` initializes its own cluster. Add capacity with workers.
- The CLI manages application artifacts and does not provision hosts.

---

## Next Steps

- **[Role Installation](/role/installation/)**: requirements and setup.
- **[CLI Installation](/cli/installation/)**: installation options.
- **[Variables](/reference/variables/)**: role configuration.
