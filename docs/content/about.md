---
title: About
description: What MiniPaaS is, why it exists, and how the pieces fit together.
weight: 1
---

## The Problem

Running production services requires repeatable deployments, firewalling, monitoring, secrets management, and a way to coordinate background work.  
Teams often face a difficult choice:

- **Kubernetes**: powerful but heavy for small clusters; sizable operational overhead.
- **Docker Compose**: simple but locked to a single host and lacking production workflows.
- **Raw Docker Swarm**: good primitives, few conventions, and a limited toolkit.

The gap is especially visible for small teams who want infrastructure that is:

- Understandable end-to-end
- Automatable
- CI-friendly
- Portable across cloud providers and your own hardware

## The Idea

MiniPaaS provides a small set of open components that layer onto Docker Swarm:

- **CLI** (`minipaas-cli/`): Builds images, rolls out the stack, and manages Swarm secrets, configs, and Caddy routing from Compose files.
- **Ansible Role** (`minipaas-role/`): Turns plain Linux hosts into secure Swarm nodes with firewalling, monitoring, syslog, and swarm-cronjob.

Each component is standalone. You can adopt one or both.

## Why Build It?

MiniPaaS aims to reduce infrastructure friction while keeping everything transparent:

- Give small teams practical tooling for everyday deployments.
- Provide a consistent deployment workflow based on Compose files teams already use.
- Standardize Swarm node provisioning with predictable security defaults.
- Keep migration paths open for moving to another platform later.

## Principles

- **Small, composable parts**: adopt what you need.
- **Explicit over implicit**: configuration and state stay in visible files.
- **Self-hosted by default**: works offline, locally, and in CI.
- **Readable and auditable**: everything is plain Go, Ansible, and YAML.

## Learn More

- Infrastructure: [Ansible Role](/role/)
- Deployments: [CLI Overview](/cli/)  
