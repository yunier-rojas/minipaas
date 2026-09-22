---
title: Secrets & Configs
description: Manage Docker secrets and configs with content-addressed names and Compose patching.
weight: 6
---

## Overview

MiniPaaS provides high-level commands for managing Docker **secrets** and **configs**:

- Creates Swarm secrets/configs with content-addressed names
- Patches the swarm Compose file declared in `minipaas.yaml`
- Attaches only to the services you mark with `--for`

Sensitive values stay **out of your repository**. The CLI stores content in Swarm and keeps the repository definitions explicit.

---

## Naming

Every object is content-addressed:

```
<name>.<hash8>
```

- `<name>` is `--name`, or the base name of the positional file.
- `<hash8>` is the first 8 hex characters of the SHA-256 of the content.

The object name changes when the content changes. Running the command again with new content creates a **new** object; it does not replace the existing value. The mounted file name is always `<name>`.

This keeps every version immutable. To rotate a value, run the command again, then update the consumer to the new name.

---

## Outputs

`minipaas.yaml` declares the Compose files for the environment, and `deploy secret` / `deploy config` patch the last file that defines each `--for <service>`. Because the environment lists its swarm overrides after the project's local-dev Compose file, references land in the swarm file and the local-dev file stays clean.

```bash
echo "${POSTGRES_PASSWORD}" | minipaas deploy secret dev \
  --name postgres_password \
  --for postgres \
  --for api \
  --for worker
```

The CLI creates the Swarm object and patches the swarm Compose file for each service.

---

## Creating Secrets

Use `deploy secret` to create a Swarm secret:

```bash
echo postgres | minipaas deploy secret dev \
  --name postgres_password \
  --for postgres \
  --for api
```

The CLI:

1. Reads the content from the positional file, or from STDIN when no file is given.
2. Computes the content-addressed name.
3. Creates the Swarm secret (skipped when it already exists).
4. Patches the Compose services.

Secrets are never stored in Compose files; only the secret reference is added.

---

## Creating Configs

Configs follow the same model and are intended for non-sensitive files.

```bash
minipaas deploy config dev ./configs/app.json \
  --name app.json \
  --for api \
  --for worker
```

---

## `--name` and files

`--name` is authoritative:

- With `--name`, the object and mounted file use that name. A positional file only supplies content.
- Without `--name`, a positional file supplies both the content and the base name.
- Reading from STDIN requires `--name`.

---

## Rotation

Because the name is derived from the content:

- Re-running the command with new content creates a new object name.
- Compose references are updated automatically by the CLI.

The old object remains in Swarm until it is removed manually.

---

## Best Practices

- Never commit secret values: commit only Compose references.
- Use one environment directory per target (`dev`, `staging`, `prod`).
- Re-run the command after a value changes, then update the consumer.
- Inspect cluster state with `docker secret ls` and `docker config ls`.
