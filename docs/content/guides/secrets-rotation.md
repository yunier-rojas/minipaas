---
title: Secrets Rotation
description: Replace a secret value without committing it to the repository.
weight: 3
---

## Overview

Secrets and configs are content-addressed as `<name>.<hash8>`. New content creates a new object, and the CLI patches the Compose reference for every service attached with `--for`. The previous object stays in Swarm until you remove it.

---

## Rotate a Value

1. Run the create command again with the new value:

```bash
echo "${NEW_POSTGRES_PASSWORD}" | minipaas deploy secret dev \
  --name postgres_password \
  --for postgres --for api
```

2. The CLI computes a new name and patches the last Compose file that defines each service.
3. Roll out the stack so services pick up the new reference:

```bash
minipaas deploy rollout dev
```

4. Remove the previous object once the rollout succeeds:

```bash
docker secret ls
docker secret rm postgres_password.<old-hash>
```

---

## Notes

- The mounted file name inside containers stays `<name>`, so applications do not change.
- `deploy config` follows the same model for non-sensitive files.
- Only services listed with `--for` are patched; update any other consumers explicitly.
- Running the command again with the same content reuses the existing object.

See **[Secrets & Configs](/cli/secrets-configs/)** for naming and CLI details.
