# Update Mode — Post-Feature Sync

## Business Rules

- Every feature change, no matter how small, triggers a full update cycle.
- Before editing, run `doc-torn-scan tree` to inspect the current filesystem state and verify the expected doc structure (README.md, sub-features/, implementation/).
- Update recalculates the dependency matrix — dependencies shift even when the feature itself doesn't change.
- AGENTS.md feature index must be refreshed to reflect any new or removed features.
- If the project scope evolves (new business stakes, changed audience), update AGENTS.md stakes section.
- After all edits, run `doc-torn-scan tree` again for structural verification against the expected doc template.

## Edge Cases

- **Feature removed**: remove its doc directory, update dependency matrix, update L0 feature list.
- **Multiple features changed**: treat each independently but produce a single update pass.
- **No changes after update cycle**: the verification step (doc-torn-scan tree) should confirm no drift.
