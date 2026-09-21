---
title: Multiple Environments
description: Run dev, staging, and production targets from one repository.
weight: 6
---

## Overview

Each environment is a directory with its own `minipaas.yaml`. The namespace in that file isolates the stack, service names, and networks inside a Swarm.

---

## Layout

```
project/
├── compose.yaml
├── compose.build.yaml
└── minipaas/
    ├── dev/
    │   └── minipaas.yaml
    └── prod/
        └── minipaas.yaml
```

Every command takes the environment directory as its first positional argument, so the target is always explicit.

---

## Namespaces

Set a distinct `namespace` per environment (`dev`, `staging`, `prod`). The CLI prefixes services as `<namespace>_<service>` and networks as `<namespace>_network` and `<namespace>_public`, allowing several environments to share one cluster.

---

## Targets

Each environment sets its own `docker.host`. Local development can omit it and use the local Docker socket; staging and production point at their managers over SSH.

---

## Next Steps

- **[minipaas.yaml](/cli/minipaas-file/)**: environment descriptor reference.
- **[CLI Usage](/cli/usage/)**: environment layout and deploy patterns.
- **[CI/CD Pipeline](/guides/ci-cd/)**: deploy each environment from CI.
