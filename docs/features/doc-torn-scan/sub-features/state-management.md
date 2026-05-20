# doc-torn-scan — State Management

## Business Rules

- The state file (`.doc-torn-state.json`) is written by the **agent**, not the tool — the agent populates it after studying `tree` output.
- The tool reads the state on `scaffold`, `complete`, `meta`, `status`.
- The tool writes the state only on `complete` (status transition, current feature advance) and `meta` (sets `meta_docs_generated = true`).
- Feature ordering is topological — dependencies must be documented before dependents.

## Edge Cases

- **Invalid status**: the `validate()` function rejects features with unknown status values.
- **Feature not found**: scaffold/complete return a clear error if the feature name is not in state.
- **Meta before all features done**: `RunMeta()` blocks with an error listing remaining features.
- **Empty project**: `tree` returns an empty JSON array.
