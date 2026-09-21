---
title: CLI Overview
description: High-level overview of the MiniPaaS CLI.
weight: 1
---

## Overview

The MiniPaaS CLI provides helpers for managing Swarm secrets and configs while staying fully Compose-first.  
It lives in the source folder **`minipaas-cli/`** and can be used independently of the Ansible role.

The CLI takes your existing Compose files and adds an environment model for secret and config management.

## Features

- **Secrets and configs** with automatic multi-file service lookup
- **Deploy commands** to build images, roll out the stack, release canaries, and apply Caddy routing
- **Environment scaffolding** that builds `compose.apps.yaml` and the base stack from your Compose files
- **Deploy shaping** that marks services as workers, jobs, or cron jobs
- **Route authoring** that adds or updates Caddy routes in `caddy.json`
- **Environment-aware execution** via `minipaas.yaml`
- **Local or SSH Docker targets** for running commands against a development daemon or a remote host
- **Shell helpers** for exporting the Docker context variables

## Motivation

The CLI exists to simplify deployments for teams that want:

- A Compose-first workflow that still supports multi-node clusters
- A single tool that handles secrets and configs
- CI-friendly behavior with no external services
- Minimal operational overhead

## Behavior

At a high level, the CLI:

- Reads an **environment directory** and its `minipaas.yaml` file
- Resolves services across multiple Compose files when modifying secrets/configs
- Builds images, rolls out the stack, releases canaries, and reloads Caddy routing from the same environment

The CLI is intentionally explicit: all changes are visible in Compose files or versioned artifacts.

## Start Here

- **Installation**: [CLI Installation](/cli/installation/)
- **Configuration**: [minipaas.yaml specification](/cli/minipaas-file/)
- **Using the CLI**: [Usage guide](/cli/usage/)
- **Command reference**: [Commands](/cli/commands/)

For related components:

- [Ansible Role Overview](/role/)
