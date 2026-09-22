---
title: Alert Labels
description: Per-container thresholds read by the role's monitoring script.
weight: 3
---

## Overview

The monitoring script installed by `minipaas-role/` checks every running container when `--DOCKER-MONITOR` is active. Two labels tune the thresholds.

| Label | Meaning | Default |
| --- | --- | --- |
| `alert.cpu_threshold` | CPU usage threshold in percent; a `%` suffix is optional | `90` |
| `alert.memory_threshold` | Memory usage threshold; accepts `B`, `KiB`, `MiB`, `GiB`, and similar units | `1000MiB` |

Add them to a service in `compose.apps.yaml`:

```yaml
services:
  example:
    labels:
      alert.cpu_threshold: "85"
      alert.memory_threshold: "1GiB"
```

---

## Behavior

- A container alerts when usage is greater than its threshold.
- Alerts arrive as `DOCKER-CPU` and `DOCKER-MEMORY` messages with a 10-minute cooldown per container and metric.
- Missing or invalid labels fall back to the defaults.
- `minipaas code init` writes both labels with the defaults; values already set in the source Compose files are kept.
- Labels are read from the running container, so redeploy the stack (`minipaas deploy rollout`) after changing them.

See **[Monitoring](/role/monitoring/)** for the full alerting setup.
