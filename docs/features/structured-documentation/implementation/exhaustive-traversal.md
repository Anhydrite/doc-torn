# Exhaustive Traversal — Why Every File Matters

## Problem

When documenting a codebase for the first time, agents tend to skip "unimportant" files — config files, tests, examples, nested directories. This guarantees incomplete documentation.

## Solution

The todo list drives discovery, not the agent's intuition:

1. List root directory → add a `todowrite` task for every file (read) and every directory (list contents).
2. When a directory is listed → add new tasks for everything found inside.
3. Repeat until no tasks remain.

This guarantees 100% coverage. Every file is at least seen and classified, even if the final decision is "document this under parent feature X" rather than "make this a separate feature."

## Tradeoff

- **Cost**: exhaustive traversal of large monorepos can take many steps.
- **Mitigation**: group trivial files (e.g., "read all `.gitkeep` files in tests/") under a single task.
- **When to skip**: only if the user explicitly says "ignore directory X."
