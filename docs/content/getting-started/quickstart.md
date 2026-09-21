---
title: Quickstart
description: Provision a Swarm cluster and deploy your first service with MiniPaaS.
weight: 2
---

## Overview

This walkthrough goes from fresh Linux hosts to a running service reachable through Caddy:

1. Provision the cluster with the Ansible role.
2. Install the CLI.
3. Scaffold an environment.
4. Create a database secret.
5. Build images and roll out the stack.
6. Add and apply a Caddy route.

The examples use `manager.example.com` for the manager, `deploy` for the SSH user, `ghcr.io/org` for an optional registry, and `api:8080` for an application service. Replace them with your own values.

---

## 1. Provision the cluster

On your workstation, clone the repository and install the role dependencies:

```bash
git clone https://github.com/yunier-rojas/minipaas.git
cd minipaas
ansible-galaxy install -r minipaas-role/requirements.yml
ansible-galaxy role install geerlingguy.docker
```

Create `inventory.ini` and `playbook.yml`:

```ini
[managers]
manager1 ansible_host=1.2.3.4

[workers]
worker1 ansible_host=1.2.3.5
```

```yaml
- name: Install and Configure MiniPaaS
  hosts: all
  become: true
  roles:
    - ./minipaas-role
```

Run the playbook:

```bash
ansible-playbook -i inventory.ini playbook.yml
```

See **[Role Installation](/role/installation/)** for requirements and verification steps.

---

## 2. Install the CLI

```bash
go install github.com/yunier-rojas/minipaas/minipaas-cli/cmd/minipaas@main
```

See **[CLI Installation](/cli/installation/)** for prebuilt binaries.

---

## 3. Scaffold an environment

From your project directory, create the environment against the remote manager:

```bash
minipaas code init minipaas --local=false --host deploy@manager.example.com --registry ghcr.io/org
```

The command detects `compose.yaml` and `compose.build.yaml` in the project root and writes `minipaas.yaml`, the swarm Compose files, and `caddy.json` into `minipaas/`.

See **[minipaas.yaml](/cli/minipaas-file/)** and **[Base Stack Files](/reference/base-stack/)** for the generated files.

---

## 4. Create the database secret

```bash
echo "${POSTGRES_PASSWORD}" | minipaas deploy secret minipaas \
  --name postgres_password \
  --for postgres
```

See **[Secrets & Configs](/cli/secrets-configs/)** for configs and rotation.

---

## 5. Build and roll out

Authenticate against the registry when `--registry` was set, then build and deploy:

```bash
docker login ghcr.io
minipaas deploy build minipaas
minipaas deploy rollout minipaas
```

---

## 6. Route traffic through Caddy

```bash
minipaas code route minipaas https://app.example.com api:8080
minipaas deploy routing minipaas
```

`code route` updates `caddy.json` and marks the service as resilient; `deploy routing` reloads the running Caddy container. See **[Caddy and TLS](/guides/caddy-tls/)**.

---

## 7. Verify

```bash
docker -H ssh://deploy@manager.example.com service ls
curl https://app.example.com
```

---

## Next Steps

- **[Example App Walkthrough](/guides/example-app/)**: run the full example project.
- **[CLI Usage](/cli/usage/)**: environments, secrets, and deploy patterns.
- **[Troubleshooting](/troubleshooting/)**: diagnose common issues.
