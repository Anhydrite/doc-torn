# Definitions

## Core Concepts

| Term | Definition |
|------|------------|
| L0 | Highest-level documentation: project overview, architecture diagram, feature list. 5-minute read. |
| L1 | Per-feature documentation: objective, logic, dependencies, API, key files. |
| L2 | Sub-feature documentation: edge cases, business rules, sub-flows. |
| L3 | Implementation documentation: technical decisions, rationale, why-this-approach. |
| Structured Documentation | The doc-torn methodology: hierarchical docs (L0→L3) with explicit dependency matrix and lifecycle (init/read/update/verify). |
| Doc-Driven Exploration | The practice of reading documentation before opening any code file. Enforced by the doc-driven-exploration skill. |
| Documentation Consistency | The practice of auditing all docs against real code and auto-fixing discrepancies. Enforced by the documentation-consistency skill. |
| Iron Law | The rule that no code file may be opened before Phase 4 of doc-driven-exploration. Violation means restarting from Phase 1. |
| Golden Rule | Real code always takes precedence over documentation. When they disagree, update docs, not code. |

## Skill Names

| Term | Definition |
|------|------------|
| structured-documentation | Skill implementing the full doc lifecycle: init, read, update, verify. |
| doc-driven-exploration | Skill enforcing doc-first exploration before any code search. |
| documentation-consistency | Skill performing full doc vs code audit with auto-fix. |

## Tools

| Term | Definition |
|------|------------|
| doc-torn-scan | Go CLI tool for iterative feature-by-feature documentation. 5 subcommands: `tree`, `scaffold`, `complete`, `meta`, `status`. |
| .doc-torn-state.json | State file persisted by the agent and read/written by doc-torn-scan. Contains feature list, order, status, and metadata. |

## Project Artifacts

| Term | Definition |
|------|------------|
| AGENTS.md | Agent cheat sheet at project root: stakes, feature index, architecture, workflow rules. |
| dependency-matrix.md | Table of dependencies between features: what each imports and who imports it. |
| definitions.md | This file — evolving business glossary. |
| dev-process.md | Development conventions, build commands, validation practices. |
| SKILL.md | OpenCode skill definition file containing metadata (YAML frontmatter) and full instruction body. |
| SPEC.md | Design specification document at project root containing original spec, architecture, principles, and TDD results. |
