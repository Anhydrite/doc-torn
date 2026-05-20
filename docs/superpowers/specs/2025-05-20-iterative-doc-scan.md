# Iterative Documentation Scanner — Design Spec

## Problem

The current `structured-documentation init` mode performs an exhaustive codebase traversal in one pass, then writes ALL documentation at once. This saturates the model's context window, loses precise per-file details, and prevents the agent from regularly persisting progress.

## Solution

A Go binary `doc-torn-scan` that handles structure discovery, state persistence, and markdown scaffolding. The agent reads code and writes explanations ("why"). Together they iterate feature-by-feature with regular back-and-forth.

## Architecture

```
Project Root/
├── tools/doc-torn-scan/
│   ├── main.go              # Entrypoint, subcommand dispatch
│   ├── cmd/
│   │   ├── tree.go          # doc-torn-scan tree
│   │   ├── scaffold.go      # doc-torn-scan scaffold <feature>
│   │   ├── complete.go      # doc-torn-scan complete <feature>
│   │   ├── meta.go          # doc-torn-scan meta
│   │   └── status.go        # doc-torn-scan status
│   ├── scan/                # Filesystem traversal, file listing
│   ├── generate/            # Markdown scaffolding and meta-doc generation
│   └── state/               # JSON state file read/write
├── .doc-torn-state.json     # Persisted progress
```

**Key difference from initial design**: no `init` or `next` subcommands. The agent drives discovery and naming. The script provides raw data (`tree`) and executes on agent decisions (`scaffold`, `complete`, `meta`).

## State File (`.doc-torn-state.json`)

The state file is written by the **agent**, not by the script. The script only reads it (for scaffold/complete/meta) and writes it on `complete`.

```json
{
  "version": 1,
  "project_root": ".",
  "features": [
    {
      "name": "structured-documentation",
      "files": ["skills/structured-documentation/SKILL.md"],
      "sub_features": ["init-mode", "update-mode"],
      "implementation_details": ["exhaustive-traversal", "l0-l3-hierarchy"],
      "dependencies": [],
      "status": "pending"
    }
  ],
  "order": ["structured-documentation", "doc-driven-exploration", "documentation-consistency"],
  "current_feature": null,
  "meta_docs_generated": false
}
```

The agent populates this after running `doc-torn-scan tree` and studying the output.

## Subcommands

### `doc-torn-scan tree`

Walk the project root (respecting `.gitignore`), output a flat JSON list of all files with metadata:

```json
[
  {
    "path": "skills/structured-documentation/SKILL.md",
    "dir": "skills/structured-documentation",
    "extension": ".md",
    "size_bytes": 15200,
    "lines": 320
  }
]
```

No interpretation — just raw filesystem data. The agent decides what is a feature.

### `doc-torn-scan scaffold <feature>`

1. Read state file for feature metadata
2. Create directory tree:
   ```
   docs/features/<feature>/
     README.md                    # L1 skeleton
     sub-features/                # L2 skeletons
       <sub-feature-1>.md
       <sub-feature-2>.md
     implementation/              # L3 skeletons
       <detail-1>.md
       <detail-2>.md
   ```
3. Each skeleton file contains:
   - Auto-generated header (feature name, file list)
   - Section structure with `<!-- EXPLANATION -->` comments where the agent must write content
4. Output: list of created files

### `doc-torn-scan complete <feature>`

1. Read state file
2. Set feature `status: "completed"`
3. Set `current_feature` to next pending feature
4. Write state file

### `doc-torn-scan meta`

1. Read state file (all features must be completed)
2. Generate `docs/README.md` (L0):
   - One-line summary
   - Architecture diagram (generated from dependency graph)
   - Feature list with links
3. Generate `docs/architecture/architecture.md`:
   - Block diagram from dependency graph (Mermaid)
   - Data flow placeholders
   - Boundaries section
4. Generate `docs/architecture/dependency-matrix.md`:
   - Table from state file dependencies
5. Generate `docs/user/definitions.md`:
   - Standard glossary template
6. Generate `docs/user/dev-process.md`:
   - Standard template with conventions
7. Update `AGENTS.md`:
   - Feature index
   - Architecture diagram
   - Rules section

### `doc-torn-scan status`

1. Read state file
2. Print table: Feature | Files | Status | Docs Generated

## Modified SKILL.md Workflow (`init` mode)

```
### Mode `init` — Iterative documentation (AUTONOMOUS)

Phase 1 — Discovery (agent-driven)
  doc-torn-scan tree
  → Agent étudie l'arborescence, identifie les features,
    les nomme, les ordonne
  → Agent écrit .doc-torn-state.json

Phase 2 — Itération feature par feature
  Tant que des features "pending" existent:
    1. Lire les fichiers source de la feature courante
    2. Comprendre le code (structure, logique, edge cases)
    3. doc-torn-scan scaffold <feature>
    4. Écrire les explications dans les squelettes (le "why")
    5. doc-torn-scan complete <feature>
    6. Commit optionnel

Phase 3 — Méta-docs
  doc-torn-scan meta
  → génère L0, architecture, matrix, definitions, dev-process, AGENTS.md
  Agent révise et ajuste
```

## Responsibilities

### Agent
- **Étudier `tree` output** et identifier les features
- **Nommer les features** (business-meaningful names, not directory names)
- **Grouper les fichiers** sous les bonnes features
- **Ordonner les features** par dépendances
- **Lire le code** de chaque feature
- **Écrire le "why"** dans les squelettes générés
- **Réviser les méta-docs** après `meta`

### Script (`doc-torn-scan`)
- **Dumper l'arborescence** (`tree`)
- **Générer les squelettes markdown** (`scaffold`) — structure de dossiers, fichiers avec sections prêtes
- **Gérer l'état** (`complete`, état dans .doc-torn-state.json)
- **Générer les méta-docs** (`meta`) — L0, architecture, matrix, etc.
- **Output déterministe** — même arbre → même squelette

## Out of Scope (v1)

- AST-level deep code analysis
- Auto-generating explanations ("why") — always the agent's job
- Live file watching
- Integration with CI/CD
