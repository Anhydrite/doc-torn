# L0→L3 Hierarchy — Why This Structure

## Problem

Flat documentation (a single README.md or a few cross-cutting files) collapses under growth. High-level vision, implementation details, and edge cases are mixed together, making it impossible to find what you need without reading everything.

## Solution

A 4-level hierarchy separates concerns by abstraction level:

- **L0** (`docs/README.md`): 5-minute overview. One-line summary, architecture diagram, feature list. For anyone asking "what is this project?"
- **L1** (`docs/features/<name>/README.md`): Per-feature objective, logic, dependencies, API, key files. For anyone asking "what does this feature do?"
- **L2** (`docs/features/<name>/sub-features/*.md`): Sub-features, edge cases, business rules. For anyone asking "what are the details?"
- **L3** (`docs/features/<name>/implementation/*.md`): Technical decisions, rationale, why-this-approach. For anyone asking "why was it done this way?"

## Why Not a Single Deep README

A single file that covers everything forces the reader to linear-scan irrelevant details. The hierarchy lets readers skip what they don't need: a new joiner reads L0, a feature developer reads L1+L2, a maintainer reads L3.
