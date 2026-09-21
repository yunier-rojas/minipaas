# MiniPaaS

MiniPaaS is a lightweight toolkit for simple, repeatable infrastructure.  
It is built around **Docker Swarm** and **Ansible**, with a focus on clarity and practical workflows for small
deployments.

MiniPaaS consists of two independent components:

- **minipaas-role/** – An Ansible role that prepares servers for Docker Swarm: Docker installation, Swarm init/join,
  firewall, logging, monitoring, and host-level swarm-cronjob.
- **minipaas-cli/** – A Go CLI for building images, rolling out the stack, managing Swarm secrets and configs, and
  configuring Caddy routing.

You can adopt each component separately or use them together:

- Use the **role** to prepare nodes for Swarm,
- Use the **CLI** to build, deploy, and manage the stack.

MiniPaaS is designed for:

- side projects,
- small production clusters,
- self-hosted environments,
- teams that prefer explicit, scriptable tooling.

It scales from a single VM to a modest Swarm cluster.

---

# minipaas-role

## Purpose

A focused Ansible role that prepares plain Linux hosts to run Docker Swarm reliably and securely.

## What it sets up

- Docker CE installation
- Swarm initialization or join
- nftables firewall (allow + user-defined ports)
- syslog-ng for system and Docker logs
- fail2ban (SSH protection)
- Telegram monitoring and log alerts
- **swarm-cronjob** on the primary manager

The role performs **host provisioning**.  
Application deployments and routing are handled by the CLI and the Caddy service in the stack.

## Minimal Inventory Example

```ini
[managers]
manager1 ansible_host=192.0.2.10 ansible_user=root

[workers]
worker1 ansible_host=192.0.2.11 ansible_user=root
```

## Minimal Playbook Example

```yaml
- name: Install and Configure MiniPaaS
  hosts: all
  become: true
  roles:
    - ./minipaas-role
```

## Role Variables

| Variable                | Default                  | Purpose                                        |
|-------------------------|--------------------------|------------------------------------------------|
| `minipaas_extra_ports`  | `[]`                     | Additional TCP ports to allow through nftables |
| `telegram_bot_token`    | `TELEGRAM_TOKEN` env var | Token for Telegram monitoring alerts           |
| `telegram_chat_id`      | `TELEGRAM_CHAT` env var  | Chat ID for Telegram monitoring alerts         |
| `swarm_cronjob_version` | `1.14.0`                 | `swarm-cronjob` release on the primary manager |

Apply the role:

```bash
ansible-playbook -i inventory.ini playbook.yml
```

After the role runs, the Swarm cluster is ready and the CLI can deploy the stack.

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

Use conventional commits with component scopes (`feat(cli): ...`, `fix(role): ...`, `docs: ...`).
Include tests where appropriate.

---

# License

MIT License — see `LICENSE`.

