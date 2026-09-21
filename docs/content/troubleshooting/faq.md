---
title: FAQ
description: Short answers about MiniPaaS components and workflows.
weight: 3
---

## Do I need both components?

No. `minipaas-role/` and `minipaas-cli/` are independent. Use the role to prepare hosts and the CLI to deploy; either can be adopted alone. See **[MiniPaaS Overview](/)**.

## Can I use an existing Swarm cluster?

Yes. The CLI works against any reachable Docker endpoint. Point `docker.host` at the target, or omit it to use the local socket. See **[minipaas.yaml](/cli/minipaas-file/)**.

## Where is MiniPaaS state stored?

In files: `minipaas.yaml`, the Compose files, and `caddy.json` inside the environment directory. The CLI uses no external service. See **[CLI Overview](/cli/)**.

## Can several environments share one Swarm?

Yes. Give each environment a distinct `namespace` so stacks, services, and networks do not collide. See **[Multiple Environments](/guides/multiple-environments/)**.

## Does the role support multiple managers?

Keep one host in the `managers` group; every host in that group initializes its own cluster. Scale with workers. See **[Role Usage](/role/usage/)**.

## How do I remove the role?

Stop and disable the systemd units, uninstall the packages, and remove the task from `minipaas-role/tasks/main.yml`. See **[Role Usage](/role/usage/)**.

## Which registry can I use?

Any registry reachable from the build machine and the cluster. `code init --registry` stores the prefix as `MINIPAAS_IMAGE_PREFIX`; authenticate with `docker login`. See **[minipaas.yaml](/cli/minipaas-file/)**.

## Does the CLI support Windows?

Prebuilt binaries exist for windows-amd64. The role targets Debian and Ubuntu hosts. See **[Compatibility](/reference/compatibility/)**.
