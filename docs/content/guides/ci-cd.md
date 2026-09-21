---
title: CI/CD Pipeline
description: Run MiniPaaS deployments from a CI job.
weight: 5
---

## Overview

The CLI is stateless: it reads the environment directory, talks to Docker over the socket or SSH, and writes only files in the repository. That makes it suitable for CI jobs with no external services.

---

## Pipeline Outline

1. Check out the repository.
2. Install the CLI.
3. Point the environment at the target: commit `docker.host` in `minipaas.yaml`, or create the environment during the job.
4. Create secrets and configs from CI secret variables.
5. Build and push images.
6. Roll out the stack.
7. Apply Caddy routing.

```bash
echo "${POSTGRES_PASSWORD}" | minipaas deploy secret prod \
  --name postgres_password --for postgres --for api

docker login ghcr.io
minipaas deploy build prod
minipaas deploy rollout prod
minipaas deploy routing prod
```

---

## Requirements

- Registry credentials stored as CI secrets. `deploy build` pushes when `MINIPAAS_IMAGE_PREFIX` is set; rollout and canary pass credentials to Swarm with `--with-registry-auth`.
- An SSH key with access to the manager when the environment uses an `ssh://` target.

---

## Next Steps

- **[CLI Usage](/cli/usage/)**: environment layout and deploy commands.
- **[Secrets & Configs](/cli/secrets-configs/)**: secret and config lifecycle.
- **[Multiple Environments](/guides/multiple-environments/)**: isolate dev, staging, and production.
