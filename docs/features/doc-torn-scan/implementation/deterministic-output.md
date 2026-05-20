# doc-torn-scan — Deterministic Output

## Why This Approach

The tool is designed to be **deterministic** — given the same filesystem state and state file, it produces identical output. This is critical because:

1. **Agent reproducibility**: the agent can re-run scaffolding after edits without unexpected changes.
2. **Diff-friendly**: generated files only change when the input changes, not due to randomness or timestamps.
3. **Verifiable**: meta-doc generation should produce the same Mermaid diagram for the same dependency graph.

## Tradeoffs

- **No AST analysis**: the tool treats all files as opaque. Deep code understanding is the agent's job.
- **No explanation generation**: `<!-- EXPLANATION -->` comments force the agent to write "why" content. The tool never guesses intent.
- **No live file watching**: v1 is purely command-based. File watching is deferred to a future version.
