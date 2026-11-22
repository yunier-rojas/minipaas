# MiniPaaS

MiniPaaS is a lightweight toolkit for simple, repeatable infrastructure.  
It is built around **Docker Swarm** and **Ansible**, with a focus on clarity and practical workflows for small
deployments.

MiniPaaS consists of three components:

- **minipaas-node/** – An Ansible role that hardens Linux hosts: SSH, nftables, fail2ban, syslog-ng, and Telegram
  monitoring.
- **minipaas-swarm/** – An Ansible role that installs Docker, forms the Swarm cluster, and runs swarm-cronjob.
- **minipaas-cli/** – A Go CLI for building images, rolling out the stack, managing Swarm secrets and configs, and
  configuring Caddy routing.

You can adopt each component separately or use them together:

- Use the **node** role to harden hosts,
- use the **swarm** role to create the cluster,
- use the **CLI** to build, deploy, and manage the stack.

MiniPaaS is designed for:

- side projects,
- small production clusters,
- self-hosted environments,
- teams that prefer explicit, scriptable tooling.

It scales from a single VM to a modest Swarm cluster.

---

# minipaas-node and minipaas-swarm

## Purpose

Two focused Ansible roles that separate host management from Swarm orchestration.

## What minipaas-node sets up

- SSH hardening with configurable access rules
- nftables firewall with explicit TCP and UDP port lists
- fail2ban with an sshd jail
- syslog-ng for system, fail2ban, and Docker logs
- Telegram monitoring and log alerts
- OS package updates on demand through a separate playbook

## What minipaas-swarm sets up

- Docker CE installation
- Docker daemon configuration
- Swarm initialization or join
- **swarm-cronjob** on the primary manager

The roles perform **host provisioning**.  
Application deployments and routing are handled by the CLI and the Caddy service in the stack.

## Minimal Inventory Example

```ini
[managers]
manager1 ansible_host=192.0.2.10 ansible_user=deploy

[workers]
worker1 ansible_host=192.0.2.11 ansible_user=deploy
```

## Minimal Playbook Example

```yaml
- name: Provision MiniPaaS
  hosts: all
  become: true
  roles:
    - role: ./minipaas-node
      tags: [node]
    - role: ./minipaas-swarm
      tags: [swarm]
```

Swarm traffic needs explicit firewall ports. Set them per group:

```yaml
# group_vars/managers.yml
node_firewall_tcp_ports: [2377, 7946, 80, 443]
node_firewall_udp_ports: [7946, 4789]
```

```yaml
# group_vars/workers.yml
node_firewall_tcp_ports: [7946]
node_firewall_udp_ports: [7946, 4789]
```

A complete example lives in `example-provisioning/`.

## Role Variables

| Variable                       | Role    | Default                  | Purpose                                        |
|--------------------------------|---------|--------------------------|------------------------------------------------|
| `node_ssh_harden`              | node    | `true`                   | Apply SSH hardening                            |
| `node_ssh_restrict`            | node    | `false`                  | Restrict SSH to `node_ssh_allowed_cidrs`       |
| `node_firewall_tcp_ports`      | node    | `[]`                     | Additional TCP ports to allow through nftables |
| `node_firewall_udp_ports`      | node    | `[]`                     | Additional UDP ports to allow through nftables |
| `telegram_bot_token`           | node    | `TELEGRAM_TOKEN` env var | Token for Telegram monitoring alerts           |
| `telegram_chat_id`             | node    | `TELEGRAM_CHAT` env var  | Chat ID for Telegram monitoring alerts         |
| `swarm_cronjob_version`        | swarm   | `1.14.0`                 | `swarm-cronjob` release on the primary manager |

Apply the playbook:

```bash
ansible-playbook -i inventory.ini playbook.yml
```

After the roles run, the Swarm cluster is ready and the CLI can deploy the stack.

---

# minipaas-cli

## Purpose

A single binary that builds images, rolls out the stack, manages secrets and configs, authors Caddy routes, and provides
a Docker-ready shell.

Install:

```bash
go install github.com/yunier-rojas/minipaas/minipaas-cli/cmd/minipaas@main
```

## Major Features

* **Deploy**: build images, roll out the stack, release canaries, apply Caddy routing
* **Secrets & configs**: create and auto-patch Compose
* **Routing**: author Caddy routes into the environment's `caddy.json`
* **Remote targets**: run against the local daemon or a remote host over SSH
* **Shell environment**: load the environment's `DOCKER_HOST` and vars

## Command Groups

### `deploy`

Build images, roll out the stack, release canaries, apply Caddy routing, and manage Swarm secrets and configs (
`deploy secret`, `deploy config`).

### `code`

Scaffold an environment, shape service deploy blocks (`worker`, `job`, `cron`), and author routes into the environment's
`caddy.json`.

### `shell`

Start a shell preloaded with the environment's Docker variables.

---

# Contributing

Use conventional commits with component scopes (`feat(cli): ...`, `fix(node): ...`, `fix(swarm): ...`, `docs: ...`).
Include tests where appropriate.

---

# License

MIT License. See `LICENSE`.
