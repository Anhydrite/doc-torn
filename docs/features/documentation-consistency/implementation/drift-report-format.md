# Drift Report Format

## Why a Report

The drift report serves two purposes:
1. **Accountability**: the agent and user both know what changed.
2. **Review**: if an auto-fix was wrong, the report makes it findable.

## Standard Template

```markdown
## Documentation Consistency Report

**Date**: <date>
**Duration**: <time elapsed>

### Fixed
- Feature X README: updated key files list (removed 2 obsolete, added 1)
- Dependency matrix: added Y→Z dependency
- definitions.md: added 3 new terms
- New feature doc created: <name>
- L2 sub-feature added: <name>

### Requires Attention
- Architecture.md: the diagram doesn't show the new CouchDB sync flow
- Feature A: sub-features directory exists but is empty

### Clean
- 12/15 features fully consistent
- definitions.md complete
```

## What to Include

- **Every fix**: even trivial ones (typo correction).
- **Requires attention**: items the agent cannot auto-fix (needs human judgment).
- **Clean count**: how many features passed without changes.
