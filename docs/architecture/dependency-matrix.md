# Dependency Matrix

## Feature-to-Feature Dependencies

| Feature / Tool | Depends On | Used By | Import/Call |
|----------------|-----------|---------|-------------|
| structured-documentation | doc-torn-scan (invokes during init) | doc-driven-exploration, documentation-consistency | sd init calls `tree`, `scaffold`, `complete`, `meta`; dde Phase 1 reads sd output; dc Step 1 reads sd output |
| doc-driven-exploration | structured-documentation (reads its doc structure) | documentation-consistency | Phase 1 reads docs produced by sd; dc benefits from consistent docs |
| documentation-consistency | structured-documentation (reads its doc structure), doc-driven-exploration (benefits from consistent docs) | — | Step 1 reads skeleton files, Step 4 checks matrix |
| doc-torn-scan | — | structured-documentation (init mode) | sd init calls `tree`, `scaffold`, `complete`, `meta` |

## External Dependencies

| Dependency | Used By | Purpose |
|-----------|---------|---------|
| OpenCode | All skills | Runtime environment for skill execution |
| superpowers (optional) | All skills | Automatic skill discovery and loading |
| Go 1.23+ | doc-torn-scan | Build and run the scanner binary |
| Git | pre-commit hook | Detect staged changes for doc verification |
