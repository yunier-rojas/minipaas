# Repository Guidelines

## Project Structure & Module Organization
- `minipaas-cli/`: Go CLI (entry: `cmd/minipaas/`), with `_test.go` unit tests.
- `minipaas-role/`: Ansible role (tasks, defaults, templates) for Swarm provisioning.
- `docs/`: Hugo site (Hextra theme). Config in `docs/hugo.toml`, content in `docs/content/`.
- `docs/i18n/`: UI string overrides for the docs site.
- `example-app/`: Minimal example app and Compose files used by documentation and CLI demos.
- `.github/`: CI workflows, issue templates.
- `mk/`: Makefile fragments (docs, build, lint targets).

---

## Build & Development Commands
- Build CLI:  
  `cd minipaas-cli/cmd/minipaas && go build -o minipaas`
- Run CLI tests:  
  `cd minipaas-cli && go test ./...`
- Format Go:  
  `gofmt -s -w .` (run inside `minipaas-cli`)
- Build docs:  
  `make build-docs` (runs `hugo --source docs --gc --minify`, output in `docs/public/`)
- Preview docs:  
  `make serve-docs` (runs `hugo server --source docs`)

---

## minipaas-cli Architecture

Clean/hexagonal layering. Dependencies point inward.

```
cmd/minipaas              CLI parsing (go-arg). Only allowed to import
                        internal/runtime and the logger framework.

internal/core/entities  Pure domain types + helpers. No external deps
                        except fmt/strings.
internal/core/usecases  Business logic (interactors) + ports. No framework
                        imports; only entities.
internal/frameworks/*   Concrete adapters: caddy, compose, config, deploy,
                        docker, files, logger, process, scaffold, shell.
internal/runtime        Wiring/composition root. Constructs interactors with
                        framework implementations and exposes use cases.
```

### Import rules are enforced

`qa/.import.yaml` defines per-folder allowlists and is validated by
`make imports`. Do not introduce imports outside the allowed list for a folder.
Notably:

- `core/usecases` may **not** import frameworks; if a use case needs I/O,
  define a port (interface) and implement it in a framework.
- `core/entities` stays dependency-free (only `fmt`, `strings`).
- `internal/runtime` is the only place that wires core + frameworks together.

### Conventions

- Package layering: when adding a feature, add the interactor in
  `internal/core/usecases`, the adapter in the relevant `internal/frameworks/*`
  package, and wire it in `internal/runtime`.
- Follow existing naming: use cases are `...Interactor`, services in frameworks
  are `...Service`, engines `...Engine`.
- Tests live next to the code as `*_test.go` in `internal/core/usecases`.
- Error handling uses `github.com/cockroachdb/errors`.
- Logging goes through `internal/frameworks/logger` (standard library only); do
  not log directly in core packages.
- YAML parsing uses `gopkg.in/yaml.v3`; Go templating uses
  `github.com/Masterminds/sprig/v3`.
- Do **not** add comments unless the surrounding code already documents that way
  or the user asks.
- Ask before adding any external dependency.

### Formatting

- Go files and Makefiles use tabs; YAML/JSON use 2 spaces (`.editorconfig`).
- Use `make qa` (`go fmt`) before committing; do not hand-format.

---

## Testing Guidelines
### CLI
- Tests live next to code as `<file>_test.go`.
- Prefer table-driven tests.
- Run via `go test ./...`.

### Role
- Add reproducible examples in `docs/` when relevant.
- Add molecule/db tests where feasible in their module directories.

---

## Documentation

### Structure
- Docs live in `docs/content/`; the Hugo site root is `docs/` (config in `docs/hugo.toml`).
- Section landing pages use `_index.md`; page order comes from the `weight` front matter field.
- Every page uses `title` and `description` front matter.
- Index pages (`docs/content/_index.md`, `docs/content/cli/_index.md`, `docs/content/role/_index.md`):
    - **No code blocks**.
    - Must contain a **`## Start Here`** section.
    - Must reference their source folder (`minipaas-cli/`, `minipaas-role/`).
    - Keep index pages high-level; details belong in subpages.

### Tone
- Describe **what components do**, not what they do not do.
- Keep explanations short and actionable.
- Base statements strictly on repository behavior—no speculation.
- Maintain consistent terminology (manager, worker, rollout, queue, stream, job, cron).

### Subpages
- Installation/usage pages may include code blocks.
- Use real commands only.
- Keep examples minimal and relevant.

### Linking
- Use site-absolute paths (`/role/installation/`) for internal links. Hugo's link render hook prefixes the base URL.
- `docs/layouts/_markup/render-link.html` fails the build when an internal link cannot be resolved to a page or asset.
- Navigation is defined by `hugo.toml` menus and page `weight` values.

### Consistency Rules for Agents & Contributors
- Update docs when behavior changes—CLI or role.
- Avoid inventing flags, commands, or variables.
- Prefer smaller pages over large ones.
- Follow the positive, capability-focused voice enforced project-wide.

---

## Docs

- User-facing docs live in `docs/` and are built with Hugo using the Hextra theme (`hugo.toml`).
- Update them when changing CLI commands or the `minipaas.yaml` format.
- Encouraging, accessible, clear, and empathetic.
- Use short sentences, active voice, concise language.
- Forbidden to use fluff, clichés, and corporate jargon.
- Keeps docs up to date with the code.
- Document features and behaviors, use facts, examples but not opinions.
- State limitations and caveats.
- Heavily reduce the usage of emojis.
- Forbid the use of — in the docs

---

## Commit & Pull Request Guidelines
- Use conventional commits with component scopes:
    - `feat(cli): ...`
    - `fix(role): ...`
    - `docs: ...`
- Keep changes focused; avoid mixing refactors with feature work.

---

## Security & Configuration Tips
- Never commit secrets or Docker contexts.
- Use `.gitignore` for local env files.
- Store sensitive config in secret managers or Ansible Vault.
- Document credential-related steps in `docs/` without exposing real artifacts.
