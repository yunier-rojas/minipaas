---
title: MiniPaaS Overview
description: A minimal toolkit for provisioning servers and deploying applications.
cascade:
  type: docs
---

## Overview

MiniPaaS is a collection of small, focused tools for simple and repeatable service deployments.  
Each component is independent and can be adopted incrementally:

- **`minipaas-role/`** provisions servers and prepares them to run a Docker Swarm cluster.  
- **`minipaas-cli/`** builds images, rolls out the stack, and manages Swarm secrets, configs, and Caddy routing.  

The project targets side projects, small clusters, and teams that value explicit control over their infrastructure.

---

## Architecture at a Glance

MiniPaaS follows a clear separation of concerns:

- The **role** prepares machines: Docker, Swarm initialization, firewall, and monitoring.  
- The **CLI** builds images, rolls out the stack, and manages secrets, configs, and Caddy routing.  

The role changes rarely; the CLI evolves on its own.

---

## When MiniPaaS Helps

MiniPaaS is useful when you want:

- a simple and reproducible way to prepare servers  
- straightforward application deployment workflows  
- background workers, jobs, and cron tasks on a small cluster  
- tools that are easy to adopt and easy to remove  

MiniPaaS is built for developers who prefer **small, predictable, scriptable components**.

---

## Components

### Provisioning: `minipaas-role/`
An Ansible role that prepares servers for Docker Swarm with secure defaults, logging, firewall rules, and host-level cron orchestration.

### Deployment: `minipaas-cli/`
A command-line interface for building images, rolling out the stack, and managing secrets, configs, and Caddy routing across Compose files.

---

## Start Here

- Deploy your first application → **[Quickstart](/getting-started/quickstart/)**
- Server provisioning → **[Role Overview](/role/)**
- Deploying applications → **[CLI Overview](/cli/)**
- Workflow walkthroughs → **[Guides](/guides/)**

For more context:

- Role and CLI reference → **[Reference](/reference/)**
- Common problems → **[Troubleshooting](/troubleshooting/)**
- Project background → **[About](/about/)**
- Contact links → **[Contact](/contact/)**  
