# Doc-Driven Exploration

## Functional Objective

Enforce a rigid "docs before code" discipline for all agent operations. Any request involving the codebase (search, feature work, debugging, onboarding) must first exhaust documentation before opening a single code file.

## Technical Logic

Six-phase workflow:

1. **Phase 1 (Mandatory)** — Read all 4 skeleton files: architecture, dependency-matrix, definitions, dev-process. No code.
2. **Phase 2** — Navigate request keywords to `docs/features/<name>/` directories. Still no code.
3. **Phase 3** — Read feature docs thoroughly (L1→L2→L3). Still no code.
4. **Phase 4** — Code as last resort, only if docs are insufficient.
5. **Phase 5** — Update docs with findings (fill gaps, add L2/L3).
6. **Phase 6** — Handle special terms, update definitions.md.

The "Iron Law" forbids reading any `.py`, `.js`, `.ts`, `.go`, `.rs` file, launching explore subagents, or running grep on source directories until Phase 4.

## Dependencies

- **Upstream**: structured-documentation (reads docs it produces)
- **Downstream**: documentation-consistency (enforced doc-first discipline feeds accurate docs into audits)

## API / Interface

- Declared in SKILL.md metadata (name, description, trigger conditions)
- The agent selects it when the task involves codebase search or feature work
- Loading: `skill({ name: "doc-driven-exploration" })`

## Key Files

| File | Purpose |
|------|---------|
| `skills/doc-driven-exploration/SKILL.md` | Skill implementation — 6 phases, Iron Law, violation rules |
| `docs/architecture/architecture.md` | Read in Phase 1 |
| `docs/architecture/dependency-matrix.md` | Read in Phase 1 |
| `docs/user/definitions.md` | Read in Phase 1 |
| `docs/user/dev-process.md` | Read in Phase 1 |
