# The Iron Law — No Code Before Docs

## Why Absolute?

The Iron Law is deliberately extreme because agents naturally default to "let me check the code." The prohibition extends to:

- Reading any source file (`.py`, `.js`, `.ts`, `.go`, `.rs`, etc.)
- Launching explore/grep/Task subagents that scan source code
- Running `grep`/`ripgrep` on source directories
- Any file outside `docs/`

## Violation Handling

- If caught during Phase 1-3: **STOP. Go back to Phase 1. Start over.**
- If caught during Phase 4-6: depends on context — if you opened code before exhausting docs, restart.

## Why Not "Check Both"?

Reading code and docs simultaneously builds a mixed mental model where you can't distinguish "what the code does" from "why it does it." The docs-first approach ensures you understand intent before implementation.
