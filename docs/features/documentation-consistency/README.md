# Documentation Consistency

## Functional Objective

Perform a full systematic audit of all documentation against the current codebase, then auto-fix every gap found. Prevents documentation drift from accumulating until docs become misleading.

## Technical Logic

Seven-step workflow:

1. **Step 1** — Load doc skeleton (same 4 files as doc-driven-exploration Phase 1).
2. **Step 2** — Scan every feature doc: README exists? Key files on disk? Dependencies match imports? L2/L3 reference real code?
3. **Step 3** — Scan codebase for undocumented features: exhaustive traversal, then match modules to feature docs.
4. **Step 4** — Verify dependency-matrix against actual imports.
5. **Step 5** — Verify definitions.md against actual code entities.
6. **Step 6** — Auto-fix all gaps immediately (no to-do list).
7. **Step 7** — Generate drift report summarizing what was found and fixed.

The Golden Rule: **Real code > Documentation. Always.** When they disagree, update docs, not code.

## Dependencies

- **Upstream**: structured-documentation (audits its doc structure), doc-driven-exploration (benefits from consistent docs during Phase 1)
- **Downstream**: none

## API / Interface

- Declared in SKILL.md metadata (name, description, trigger conditions)
- The agent selects it after code changes, during reasoning when drift is noticed, or periodically
- Loading: `skill({ name: "documentation-consistency" })`

## Key Files

| File | Purpose |
|------|---------|
| `skills/documentation-consistency/SKILL.md` | Skill implementation — 7 steps, auto-fix rules, drift report |
| `docs/architecture/architecture.md` | Read in Step 1, compared against code layout |
| `docs/architecture/dependency-matrix.md` | Verified in Step 4 against actual imports |
| `docs/user/definitions.md` | Verified in Step 5 against code entities |
