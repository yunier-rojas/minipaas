---
title: Firewall
description: How the MiniPaaS node role configures nftables for Docker Swarm nodes.
weight: 8
---

## Overview

`minipaas-node/` configures an **nftables firewall** with a simple, secure default:

- **Default-deny** for inbound traffic
- explicit TCP and UDP port lists
- optional SSH allowlist
- **Fail2Ban** for SSH protection

The goal is to provide a **safe baseline** for small Swarm clusters while keeping all rules readable, inspectable, and easy to remove.

The swarm role does not open ports. Declare the ports each group needs.

---

## Default Behavior

The node role writes `/etc/nftables.conf` and enables the `nftables` service. The rules:

- Drop all inbound traffic except the configured ports
- Allow all outbound traffic
- Allow related and established connections

The node role always allows the configured SSH port and local syslog traffic on `udp/514`.

---

## Ports to Declare

### All nodes

| Purpose | Port(s) | Protocol |
| --- | --- | --- |
| SSH | `node_ssh_port` (22) | TCP |
| Gossip and node discovery | 7946 | TCP and UDP |
| Overlay networking (VXLAN) | 4789 | UDP |

### Managers

| Purpose | Port(s) | Protocol |
| --- | --- | --- |
| Swarm management API | 2377 | TCP |
| HTTP and HTTPS ingress (Caddy) | 80, 443 | TCP |

Keep these ports open while the cluster is running.

---

## Example

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
```

`group_vars/workers.yml`:

```yaml
node_firewall_tcp_ports:
  - 7946
node_firewall_udp_ports:
  - 7946
  - 4789
```

Add any published service ports to the matching list.

---

## Restricting SSH

Set `node_ssh_restrict: true` and list IPv4 ranges:

```yaml
node_ssh_restrict: true
node_ssh_allowed_cidrs:
  - 203.0.113.10/32
  - 198.51.100.0/24
```

An ISP-wide range reduces exposure, but it does not authenticate the network. Key-only authentication and Fail2Ban remain in force.

---

## Fail2Ban

The node role installs Fail2Ban and writes an sshd jail under `/etc/fail2ban/jail.d/zz-minipaas.local`. The jail uses the `systemd` backend and the `nftables-multiport` ban action. It verifies that the jail starts.

---

## Inspecting Firewall Rules

After the playbook runs, inspect the live rules:

```bash
sudo nft list ruleset
```

Or the node role managed nftables config (`/etc/nftables.conf`).

All rules are plain nftables syntax, fully transparent and editable.

---

## Skipping Firewall Configuration

The node role has no toggle for the firewall. To manage firewalling elsewhere, remove the **Configure firewall** task from the node role's `tasks/main.yml`. This fits environments that:

* use cloud provider firewall rules
* route through an external network appliance
* run on a trusted private network

---

## Best Practices

* Expose HTTP services through the stack's **Caddy** service instead of ad-hoc ports.
* Keep the Swarm coordination ports open.
* Protect SSH with Fail2Ban and network-level rules.
* Re-run the node role to apply changes.
