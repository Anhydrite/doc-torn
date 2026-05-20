# Auto-Fix Rules

## Philosophy

Auto-fix is not optional. The whole point of the consistency skill is that "later never comes." Every gap found in Steps 2-5 must be fixed immediately in Step 6.

## Rule Table

| Condition | Action |
|-----------|--------|
| README missing for a feature | Create from code exploration of that feature |
| Key file in README does not exist on disk | Remove or update the reference |
| Sub-feature doc missing | Create L2 doc from code findings |
| Implementation doc missing | Create L3 doc from code findings |
| Dependency missing in matrix | Add row |
| Dependency removed | Remove or mark as legacy |
| Entity class without definition | Add entry to definitions.md |
| Obsolete definition | Update or mark deprecated |
| User-used term not defined | Add definition |
| Feature doc older than 90 days without changes | Flag for review |

## What NOT to Auto-Fix

- Architectural decisions that disagree with code: update the doc, don't change the matrix to "match" if the code has a bug.
- Intentional discrepancies (e.g., deprecated features kept in docs for reference).
