# doc-torn-scan

## Functional Objective

Go CLI tool for iterative feature-by-feature documentation. Handles filesystem scanning, markdown scaffold generation, state persistence, and meta-doc generation. The agent drives discovery, naming, and explanation — the tool executes structural tasks deterministically.

## Technical Logic

The tool exposes 5 subcommands dispatched from a single entrypoint (`main.go`):

1. **`tree`** — walks the project filesystem, respecting `.gitignore` conventions, outputting a flat JSON list of all files with metadata (path, extension, size, line count).
2. **`scaffold <name>`** — reads `.doc-torn-state.json`, creates the full `docs/features/<name>/` directory tree with L1 (README.md), L2 (sub-features/), and L3 (implementation/) skeletons. Each skeleton contains section headers with `<!-- EXPLANATION -->` comments.
3. **`complete <name>`** — marks a feature as `completed` in the state file and advances `current_feature` to the next pending feature.
4. **`meta`** — when all features are completed, generates global docs: L0 (`docs/README.md`), architecture (`docs/architecture/architecture.md`), dependency matrix, definitions, dev-process, and `AGENTS.md`.
5. **`status`** — prints a table of features with their file count, status, and current pointer.

The state file (`.doc-torn-state.json`) is written by the agent during Phase 1 and read/updated by the tool on scaffold/complete/meta.

## Dependencies

- **Upstream**: Go 1.23+ stdlib — no external dependencies
- **Downstream**: structured-documentation (calls doc-torn-scan during init mode)

## API / Interface

Invoked from the command line:

```bash
doc-torn-scan tree                        # List all files as JSON
doc-torn-scan scaffold <feature-name>     # Generate doc skeletons
doc-torn-scan complete <feature-name>     # Mark feature done
doc-torn-scan meta                        # Generate global docs
doc-torn-scan status                      # Show progress
```

## Key Files

| File | Purpose |
|------|---------|
| `tools/doc-torn-scan/main.go` | Entrypoint — subcommand dispatch via `os.Args` |
| `tools/doc-torn-scan/state/state.go` | `DocTornState` struct, JSON read/write, Complete logic |
| `tools/doc-torn-scan/scan/scan.go` | `WalkFilesystem()` — recursive tree walk with .gitignore-aware filtering |
| `tools/doc-torn-scan/generate/scaffold.go` | `ScaffoldFeature()` — L1/L2/L3 markdown skeleton generation |
| `tools/doc-torn-scan/generate/meta.go` | `Meta()` — L0, architecture, matrix, definitions, dev-process, AGENTS.md |
| `tools/doc-torn-scan/cmd/tree.go` | `RunTree()` — CLI handler for `tree` |
| `tools/doc-torn-scan/cmd/scaffold.go` | `RunScaffold()` — CLI handler for `scaffold` |
| `tools/doc-torn-scan/cmd/complete.go` | `RunComplete()` — CLI handler for `complete` |
| `tools/doc-torn-scan/cmd/meta.go` | `RunMeta()` — CLI handler for `meta` |
| `tools/doc-torn-scan/cmd/status.go` | `RunStatus()` — CLI handler for `status` |
