---
title: Usage
description: How to use the MiniPaaS role to provision and manage Swarm nodes.
weight: 4
---

## Overview

This guide shows how to **run**, **re-run**, and **evolve** a MiniPaaS-powered Swarm cluster after installation.

The MiniPaaS role is designed to be:

- fully reproducible
- idempotent (safe to run repeatedly)
- easy to extend
- easy to undo
- suitable for small Swarm clusters and side projects

You describe hosts in your `inventory.ini`, adjust variables in `group_vars` / `host_vars`, and apply the role using
Ansible.

---

## Running the Role

After creating `playbook.yml` and your inventory (see **[Installation](/role/installation/)**):

```bash
ansible-playbook -i inventory.ini playbook.yml
```

This performs:

* Docker installation
* Swarm initialization or joining
* Firewall setup (nftables)
* syslog-ng installation
* Fail2Ban installation
* Telegram monitoring and log alerts (when credentials are set)
* swarm-cronjob installation on the primary manager

Re-running the playbook converges the nodes to the declared configuration.

---

## Re-running the Role

You should re-run the playbook when you:

* add new nodes
* modify firewall rules
* update monitoring or alerting settings
* verify cluster health after changes

Example:

```bash
ansible-playbook -i inventory.ini playbook.yml
```

Re-running the role restarts services whose configuration changes, such as syslog-ng and Fail2Ban.

---

## Multiple Managers

The role initializes Swarm on each host in the `managers` group. MiniPaaS targets small clusters, so keep one host in `managers` and add capacity with `workers`. Each extra host in `managers` starts its own Swarm cluster.

---

## Adding a Worker Node

1. Add the host to `[workers]`
2. (Optional) add `minipaas_extra_ports` or other overrides
3. Apply the role:

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

Modify allowed TCP ports:

```yaml
minipaas_extra_ports:
  - 8080
  - 9090
```

Then apply:

```bash
ansible-playbook -i inventory.ini playbook.yml
```

The nftables rules will be regenerated and applied safely.

---

## Updating Monitoring or Alerts

Set the Telegram credentials:

```yaml
telegram_bot_token: "123:ABC"
telegram_chat_id: "999111222"
```

Monitoring activates when both are set. Apply changes with:

```bash
ansible-playbook -i inventory.ini playbook.yml
```

---

## Removing Components

Everything installed by the role is:

* standard Debian/Ubuntu packages
* systemd units and configuration files

To remove a component:

* stop and disable its systemd unit
* uninstall its package
* remove its task from `minipaas-role/tasks/main.yml` so re-runs keep it removed

All files live at standard system paths.

---

## Cluster Lifecycle Examples

### Rebuild firewall after opening a new port

```yaml
minipaas_extra_ports:
  - 8080
```

```bash
ansible-playbook -i inventory.ini playbook.yml
```

### Re-run monitoring setup

```yaml
telegram_bot_token: "123:ABC"
telegram_chat_id: "999111222"
```

```bash
ansible-playbook -i inventory.ini playbook.yml
```

### Reset a worker node

1. Reinstall OS or clean Docker
2. Add node back to `[workers]`
3. Run the playbook
4. Worker rejoins automatically

---

## Summary

Using the MiniPaaS role consists of:

* Writing a clear inventory
* Setting per-group or per-host variables
* Running Ansible to converge the cluster
* Re-running to apply changes safely

The role provides a **minimal, transparent, reproducible Swarm environment**, ideal for small deployments and teams that
value low operational overhead.
