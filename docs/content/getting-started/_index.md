---
title: Getting Started
description: Provision a cluster and deploy your first application with MiniPaaS.
weight: 1
---

## Overview

MiniPaaS combines three components:

- **`minipaas-node/`** hardens Linux hosts and configures firewall, logging, and monitoring.
- **`minipaas-swarm/`** installs Docker and creates a Docker Swarm cluster.
- **`minipaas-cli/`** builds images, rolls out the stack, and manages Swarm secrets, configs, and Caddy routing.

You can adopt the node role on its own. The quickstart uses all three to take fresh servers from provisioning to a service reachable through Caddy.

## Start Here

- **Deploy your first application**: [Quickstart](/getting-started/quickstart/)
- **Full example project**: [Example App Walkthrough](/guides/example-app/)
- **Provision servers only**: [Role Installation](/role/installation/)
- **Deploy from a workstation only**: [CLI Installation](/cli/installation/)
