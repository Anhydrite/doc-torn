# Functional Architecture

## Block Diagram

```mermaid
graph TD
    subgraph Skills
        SD[structured-documentation]
        DDE[doc-driven-exploration]
        DC[documentation-consistency]
        TS[doc-torn-scan]
    end

    subgraph Project Target
        Docs[docs/]
        Code[Codebase]
        State[.doc-torn-state.json]
    end

    DDE -->|Phase 1-3: read before code| Docs
    SD -->|init: orchestrates| TS
    SD -->|update / verify| Docs
    TS -->|tree| Code
    TS -->|scaffold + meta| Docs
    TS -->|state| State
    DC -->|Step 2-6: audit + auto-fix| Docs
    DC -->|Step 3: compare| Code
    Code -.->|"real code > docs"| DC
```

## Data Flows

### Init flow (first-time documentation)

```mermaid
sequenceDiagram
    Agent->>TS: doc-torn-scan tree
    TS->>Codebase: walk filesystem
    TS-->>Agent: file tree JSON
    Agent->>State: write .doc-torn-state.json (features, order, deps)
    loop For each feature
        Agent->>TS: doc-torn-scan scaffold <feature>
        TS->>Docs: generate L1/L2/L3 skeletons
        Agent->>Docs: write explanations (the "why")
        Agent->>TS: doc-torn-scan complete <feature>
    end
    Agent->>TS: doc-torn-scan meta
    TS->>Docs: generate L0, architecture, matrix, definitions, dev-process, AGENTS.md
    Agent->>Docs: review and adjust
```

### Update flow (after a feature)

```mermaid
sequenceDiagram
    Agent->>SD: structured-documentation update
    SD->>Docs: update feature README (L1)
    SD->>Docs: create/update L2/L3
    SD->>Docs: recalculate dependency-matrix
    SD->>Docs: update definitions + dev-process
    SD->>Docs: update L0 if features changed
    SD->>Docs: update architecture if changed
    SD->>Root: update AGENTS.md
```

### Consistency audit flow

```mermaid
sequenceDiagram
    Agent->>DC: documentation-consistency
    DC->>Docs: load skeleton (4 files)
    DC->>Docs: scan every feature doc
    DC->>Codebase: scan for undocumented features
    DC->>Codebase: verify dependency-matrix against imports
    DC->>Docs: verify dependency-matrix against docs
    DC->>Codebase: verify definitions against entities
    DC->>Docs: verify definitions against docs
    DC->>Docs: auto-fix all gaps
    DC->>Agent: generate drift report
```

## Key Boundaries

- **Skills are Markdown, tool is Go** — `skills/<name>/SKILL.md` files are the skill definitions. `tools/doc-torn-scan/` is a Go binary for filesystem scanning and doc generation.
- **Examples are not features** — `examples/AGENTS.md` and `examples/hooks/` are project templates for consumers of doc-torn, not internal features.
- **doc-torn-scan is a CLI tool, not a skill** — it does not use the OpenCode skill system. It is invoked directly by the agent during `structured-documentation init`.
