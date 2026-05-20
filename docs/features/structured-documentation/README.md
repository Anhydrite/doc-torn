# Structured Documentation

## Functional Objective

Create and maintain a hierarchical documentation tree (L0→L3) that stays permanently in sync with code. Provides the core lifecycle: init (first-time) and update (after features). Reading and verification are delegated to companion skills.

## Technical Logic

The skill operates in two modes:

1. **init** — Iterative feature-by-feature documentation using `doc-torn-scan`. Three phases: discovery (`tree`), feature iteration (`scaffold` + `complete`), and meta-docs (`meta`).
2. **update** — After a feature is complete, updates the feature L1, L2/L3, recalculates dependencies, refreshes definitions and AGENTS.md.

> **Note**: The `read` and `verify` modes have been delegated to companion skills. See [doc-driven-exploration](../doc-driven-exploration/README.md) for pre-feature reading and [documentation-consistency](../documentation-consistency/README.md) for pre-commit drift detection.

## Dependencies

- **Upstream**: doc-torn-scan (invoked during init and update modes)
- **Downstream**: doc-driven-exploration (reads docs produced here), documentation-consistency (audits docs produced here)

## API / Interface

The skill is invoked via the agent's `skill` tool:
- Declared in SKILL.md metadata (name, description, trigger conditions)
- The agent selects it when the task matches the description
- Loading: `skill({ name: "structured-documentation" })`

## Key Files

| File | Purpose |
|------|---------|
| `skills/structured-documentation/SKILL.md` | Skill implementation — 2 modes (init, update), workflows, templates |
| `docs/README.md` | L0 output |
| `docs/architecture/architecture.md` | Architecture output |
| `docs/architecture/dependency-matrix.md` | Dependency matrix output |
| `AGENTS.md` | Agent guide output |
