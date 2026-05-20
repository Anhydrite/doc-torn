# doc-driven-exploration — Why Six Phases

## Why This Approach

The 6-phase structure mirrors the natural progression from "what does the project look like?" to "how does this specific feature work?" to "what did I learn?"

- **Phase 1 (Skeleton)** solves the blank-slate problem — without the 4 skeleton files, the agent has no map of the project.
- **Phase 2 (Navigate)** prevents reading irrelevant docs — keyword matching to feature directories is the fastest path.
- **Phase 3 (Deep read)** prevents surface-level understanding — L2/L3 contain edge cases and rationale that the README intentionally omits.
- **Phase 4 (Code last)** is the critical discipline — docs must be exhausted before code is reached.
- **Phase 5-6 (Update)** close the loop — every exploration improves the docs for the next explorer.

## Tradeoffs

- **Rigidity**: the 6-phase flow adds overhead for trivial lookups (e.g., "what line is X defined on?"). Mitigation: Phase 4 allows code access when docs truly can't answer.
- **Iron Law**: deliberately extreme. A softer rule ("prefer docs first") would be ignored in practice. Absolute prohibition creates a clear violation signal.
