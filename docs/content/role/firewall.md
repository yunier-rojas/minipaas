---
title: Firewall
description: How the MiniPaaS role configures nftables for Docker Swarm nodes.
weight: 8
---

## Overview

The MiniPaaS role configures an **nftables firewall** with a simple, secure default:

- **Default-deny** for inbound traffic  
- Required **Docker Swarm ports opened automatically**  
- Custom TCP ports allowed via configuration  
- **Fail2Ban** for SSH protection

The goal is to provide a **safe baseline** for small Swarm clusters while keeping all rules readable, inspectable, and easy to remove.

The role provides the rules a Swarm cluster needs: SSH, Swarm coordination ports, HTTP/HTTPS on managers, and any extra ports you define.

---

## Default Behavior

The role writes `/etc/nftables.conf` and enables the `nftables` service. The rules:

- Drop all inbound traffic except the ports listed below
- Allow all outbound traffic
- Allow related and established connections

Managers and workers receive the shared rules; managers additionally open the management and ingress ports.

---

## Ports Opened Automatically

### All nodes

| Purpose                      | Port(s)      | Protocol |
|------------------------------|--------------|----------|
| SSH                          | 22           | TCP      |
| Gossip / node discovery      | 7946         | TCP/UDP  |
| Overlay networking (VXLAN)   | 4789         | UDP      |

### Managers

| Purpose                      | Port(s)      | Protocol |
|------------------------------|--------------|----------|
| Swarm management API         | 2377         | TCP      |
| HTTP/HTTPS ingress (Caddy)   | 80, 443      | TCP      |

Keep these ports open while the cluster is running.

---

## Custom Allowed Ports

To open additional inbound TCP ports, set `minipaas_extra_ports`:

```yaml
minipaas_extra_ports:
  - 8080
  - 9090
```

These are merged into the nftables configuration on all hosts.

Useful for:

* exposing services on additional ports
* internal dashboards
* custom TCP workloads

If you publish a service on an additional host port, you must open the port here.

---

## Fail2Ban

The role installs the Fail2Ban package and restarts the service. Jails follow the package defaults and any files you add under `/etc/fail2ban/`.

---

## Inspecting Firewall Rules

After the playbook runs, you can inspect the live rules:

```bash
sudo nft list ruleset
```

Or the role-managed nftables config (`/etc/nftables.conf`).

All rules are plain nftables syntax, fully transparent and editable.

---

## Skipping Firewall Configuration

The role has no toggle for the firewall. To manage firewalling elsewhere, remove the **Configure firewall** task from the role’s `tasks/main.yml`. This fits environments that:

* use cloud provider firewall rules
* route through an external network appliance
* run on a trusted private network

---

## Best Practices

* Expose HTTP services through the stack's **Caddy** service instead of ad-hoc ports.
* Keep the Swarm coordination ports open.
* Protect SSH with Fail2Ban or network-level rules.
* Re-run the role to apply changes.
