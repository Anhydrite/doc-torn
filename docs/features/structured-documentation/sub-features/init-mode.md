# Init Mode — First-Time Documentation

## Business Rules

- Init is **fully autonomous**: the agent explores, decides, and produces everything without asking the user.
- The exhaustive traversal must visit **every nested file** — no depth limit, no early stop.
- Features are identified from the traversal results (user-facing modules, not nested implementations).
- The output must include all 9 deliverables: L0, architecture, dependency-matrix, L1 per feature, L2/L3 per feature, definitions, dev-process, AGENTS.md.

## Edge Cases

- **Empty codebase**: if no files found beyond the project skeleton, create a minimal L0 noting the project has no code yet.
- **Single-file project**: document it as a single feature with no sub-features.
- **Existing partial docs**: init should not overwrite; it should merge (detect existing docs, fill gaps only).
