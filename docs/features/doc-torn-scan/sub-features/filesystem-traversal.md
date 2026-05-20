# doc-torn-scan — Filesystem Traversal

## Business Rules

- `WalkFilesystem()` walks the project root recursively, producing a flat sorted list of files.
- Hidden directories (names starting with `.`) are skipped, except the project root.
- Standard ignore directories are hardcoded: `.git`, `node_modules`, `vendor`, `__pycache__`, `.venv`, `venv`, `dist`, `build`, `.doc-torn`.
- Permission errors on individual files are silently skipped (not fatal).
- Line counting uses `bytes.Count(data, '\n') + 1` — efficient, no regex.

## Edge Cases

- **Binary files**: counted by line count (may be 0 or 1). No binary detection.
- **Symlinks**: `filepath.Walk` does not follow symlinks by default.
- **Very large files**: read entirely into memory for line counting — acceptable for codebases under typical sizes.
