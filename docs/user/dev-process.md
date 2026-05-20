# Development Process

## Build Commands

Build `doc-torn-scan` from source:

```bash
cd tools/doc-torn-scan && go build -o ~/.local/bin/doc-torn-scan .
```

## Validation Practices

### Pre-Commit Verification

The `examples/hooks/pre-commit` hook checks:
- `AGENTS.md` exists at project root
- `docs/README.md` (L0) exists
- `docs/architecture/architecture.md` or `docs/architecture/architecture-fonctionnelle.md` exists
- `docs/architecture/dependency-matrix.md` exists
- `docs/user/definitions.md` exists
- `docs/user/dev-process.md` exists
- Every feature under `docs/features/` has a `README.md`

If any check fails, the commit is blocked.

### Post-Commit Reminder

The `examples/hooks/post-commit` hook prints a documentation checklist reminder after every commit.

## Documentation Workflow

### First Time (init)
1. Load `structured-documentation` skill (init mode)
2. Phase 1: Run `doc-torn-scan tree`, study output, write `.doc-torn-state.json`
3. Phase 2: For each feature: `doc-torn-scan scaffold` → write "why" → `doc-torn-scan complete`
4. Phase 3: Run `doc-torn-scan meta`, review generated files
5. Run `documentation-consistency` to audit generated docs against real code

### Before a Feature (read)
1. Load `doc-driven-exploration` skill
2. Read Phase 1 skeleton (4 files)
3. Navigate to relevant feature docs
4. Read thoroughly (L1→L2→L3)
5. Only then read code if needed

### After a Feature (update)
1. Load `structured-documentation update` mode
2. Update feature L1, L2/L3
3. Recalculate dependency-matrix
4. Update definitions and dev-process if changed
5. Update L0 and architecture if scope changed
6. Update AGENTS.md

### Periodic Audit (consistency)
1. Load `documentation-consistency` skill
2. Full scan and auto-fix
3. Review drift report

## Technology Stack

- **Skills**: Markdown (SKILL.md), no runtime dependencies
- **Scanner**: Go 1.23+, zero external dependencies (stdlib only)
