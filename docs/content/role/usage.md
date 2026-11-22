---
title: Usage
description: How to use the MiniPaaS roles to provision and manage Swarm nodes.
weight: 4
---

## Overview

This guide shows how to **run**, **re-run**, and **evolve** a MiniPaaS-powered Swarm cluster after installation.

The roles are designed to be:

- fully reproducible
- idempotent (safe to run repeatedly)
- easy to extend
- easy to undo
- suitable for small Swarm clusters and side projects

You describe hosts in your `inventory.ini`, adjust variables in `group_vars` / `host_vars`, and apply the roles using Ansible.

---

## Running the Roles

After creating `playbook.yml` and your inventory (see **[Installation](/role/installation/)**):

```bash
ansible-playbook -i inventory.ini playbook.yml
```

This performs:

* SSH hardening
* firewall setup (nftables)
* syslog-ng installation
* Fail2Ban installation
* Telegram monitoring and log alerts (when credentials are set)
* Docker installation
* Swarm initialization or joining
* swarm-cronjob installation on the primary manager

Re-running the playbook converges the nodes to the declared configuration.

---

## Re-running the Roles

You should re-run the playbook when you:

* add new nodes
* modify firewall rules
* modify SSH access
* update monitoring or alerting settings
* verify cluster health after changes

Example:

```bash
ansible-playbook -i inventory.ini playbook.yml
```

Re-running the roles restarts services whose configuration changes, such as syslog-ng, Fail2Ban, and nftables.

---

## Running One Layer

Use tags to run a single role:

```bash
ansible-playbook -i inventory.ini playbook.yml --tags node
ansible-playbook -i inventory.ini playbook.yml --tags swarm
ansible-playbook -i inventory.ini playbook.yml --tags docker
```

Package upgrades run only when requested:

```bash
ansible-playbook -i example-provisioning/inventory.ini example-provisioning/update.yml
```

---

## Multiple Managers

The swarm role initializes Swarm on each host in the `managers` group. MiniPaaS targets small clusters, so keep one host in `managers` and add capacity with `workers`. Each extra host in `managers` starts its own Swarm cluster.

---

## Adding a Worker Node

1. Add the host to `[workers]`
2. Add the worker firewall ports in `group_vars/workers.yml`
3. Apply the roles:

```bash
ansible-playbook -i inventory.ini playbook.yml
```

The worker joins the cluster with the worker join token from the primary manager.

---

## Connecting From Your Workstation

Docker stays bound to the local Unix socket. To use the MiniPaaS CLI or the `docker` client against a remote manager, connect over SSH and set an `ssh://` target:

```yaml
docker:
  host: ssh://deploy@manager.example.com
```

See **[minipaas.yaml](/cli/minipaas-file/)** for the configuration details.

---

## Deploying Applications

Applications, secrets, configs, and Caddy routing are managed by the MiniPaaS CLI. See **[Commands](/cli/commands/)**.

---

## Updating Firewall Rules

Modify the port lists for the group:

```yaml
node_firewall_tcp_ports:
  - 2377
  - 7946
  - 80
  - 443
  - 8080
node_firewall_udp_ports:
  - 7946
  - 4789
```

Then apply:

```bash
ansible-playbook -i inventory.ini playbook.yml --tags node
```

The nftables rules are regenerated and applied.

---

## Updating Monitoring or Alerts

Set the Telegram credentials:

```yaml
telegram_bot_token: "123:ABC"
telegram_chat_id: "999111222"
```

Monitoring activates when both are set. Apply changes with:

```bash
ansible-playbook -i inventory.ini playbook.yml --tags node
```

---

## Removing Components

Everything installed by the roles is:

* standard Debian/Ubuntu packages
* systemd units and configuration files

To remove a component:

* stop and disable its systemd unit
* uninstall its package
* remove its task from the matching role's `tasks/main.yml` so re-runs keep it removed

All files live at standard system paths.

---

## Cluster Lifecycle Examples

### Rebuild firewall after opening a new port

```yaml
node_firewall_tcp_ports:
  - 2377
  - 7946
  - 80
  - 443
  - 8080
```

```bash
ansible-playbook -i inventory.ini playbook.yml --tags node
```

### Re-run monitoring setup

```yaml
telegram_bot_token: "123:ABC"
telegram_chat_id: "999111222"
```

```bash
ansible-playbook -i inventory.ini playbook.yml --tags node
```

### Reset a worker node

1. Reinstall the OS or clean Docker
2. Add the node back to `[workers]`
3. Run the playbook
4. The worker rejoins automatically

---

## Summary

Using the MiniPaaS roles consists of:

* writing a clear inventory
* setting per-group or per-host variables
* running Ansible to converge the cluster
* re-running to apply changes safely

The roles provide a **minimal, transparent, reproducible Swarm environment**, ideal for small deployments and teams that value low operational overhead.
