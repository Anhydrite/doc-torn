# Iterative Doc Scanner Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build `doc-torn-scan` Go binary for iterative feature-by-feature documentation with state persistence.

**Architecture:** Go CLI with 5 subcommands (`tree`, `scaffold`, `complete`, `meta`, `status`). Agent drives discovery, naming, and ordering. Script handles filesystem traversal, skeleton generation, state management, and meta-doc generation.

**Tech Stack:** Go 1.23, standard library only (no external deps), `encoding/json`, `os`, `path/filepath`, `flag`, `text/template`.

---

## File Structure

```
tools/doc-torn-scan/
├── go.mod
├── main.go                         # Entrypoint, subcommand dispatch via os.Args
├── state/
│   └── state.go                    # DocTornState, Feature structs, Read/Write/Complete
├── scan/
│   └── scan.go                     # WalkFilesystem() - recursive file list with metadata
├── generate/
│   ├── scaffold.go                 # ScaffoldFeature() - create L1/L2/L3 skeletons
│   └── meta.go                     # GenerateMeta() - L0, architecture, matrix, definitions, dev-process, AGENTS.md
└── cmd/
    ├── tree.go                     # runTree()
    ├── scaffold.go                 # runScaffold()
    ├── complete.go                 # runComplete()
    ├── meta.go                     # runMeta()
    └── status.go                   # runStatus()
```

### Task 1: Module Setup + State Package

**Files:**
- Create: `tools/doc-torn-scan/go.mod`
- Create: `tools/doc-torn-scan/state/state.go`

- [ ] **Step 1: Initialize Go module**

Run: `mkdir -p tools/doc-torn-scan && cd tools/doc-torn-scan && go mod init github.com/Anhydrite/doc-torn/tools/doc-torn-scan`

- [ ] **Step 2: Write state.go with DocTornState struct**

```go
package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type FeatureStatus string

const (
	StatusPending   FeatureStatus = "pending"
	StatusCompleted FeatureStatus = "completed"
)

type Feature struct {
	Name                  string   `json:"name"`
	Files                 []string `json:"files"`
	SubFeatures           []string `json:"sub_features,omitempty"`
	ImplementationDetails []string `json:"implementation_details,omitempty"`
	Dependencies          []string `json:"dependencies,omitempty"`
	Status                FeatureStatus `json:"status"`
}

type DocTornState struct {
	Version           int       `json:"version"`
	ProjectRoot       string    `json:"project_root"`
	Features          []Feature `json:"features"`
	Order             []string  `json:"order"`
	CurrentFeature    string    `json:"current_feature"`
	MetaDocsGenerated bool      `json:"meta_docs_generated"`
}

func StateFilePath() string {
	return filepath.Join(".doc-torn-state.json")
}

func Read() (*DocTornState, error) {
	data, err := os.ReadFile(StateFilePath())
	if err != nil {
		return nil, fmt.Errorf("reading state: %w", err)
	}
	var s DocTornState
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parsing state: %w", err)
	}
	return &s, nil
}

func (s *DocTornState) Write() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling state: %w", err)
	}
	return os.WriteFile(StateFilePath(), data, 0644)
}

func (s *DocTornState) Complete(name string) error {
	found := false
	for i, f := range s.Features {
		if f.Name == name {
			s.Features[i].Status = StatusCompleted
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("feature %q not found in state", name)
	}
	// Find next pending feature
	s.CurrentFeature = ""
	for _, orderName := range s.Order {
		for _, f := range s.Features {
			if f.Name == orderName && f.Status == StatusPending {
				s.CurrentFeature = orderName
				return s.Write()
			}
		}
	}
	return s.Write()
}

func (s *DocTornState) AllCompleted() bool {
	for _, f := range s.Features {
		if f.Status != StatusCompleted {
			return false
		}
	}
	return true
}
```

- [ ] **Step 3: Verify build**

Run: `cd tools/doc-torn-scan && go build ./state/`
Expected: no errors

- [ ] **Step 4: Commit**

```bash
git add tools/doc-torn-scan/go.mod tools/doc-torn-scan/state/state.go
git commit -m "feat(doc-torn-scan): module setup and state package"
```

### Task 2: Scan Package (Filesystem Tree)

**Files:**
- Create: `tools/doc-torn-scan/scan/scan.go`

- [ ] **Step 1: Write scan.go**

```go
package scan

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FileEntry struct {
	Path      string `json:"path"`
	Dir       string `json:"dir"`
	Extension string `json:"extension"`
	SizeBytes int64  `json:"size_bytes"`
	Lines     int    `json:"lines"`
}

func isIgnored(path string, info os.FileInfo, ignoreDirs map[string]bool) bool {
	if info.IsDir() {
		base := info.Name()
		if ignoreDirs[base] || strings.HasPrefix(base, ".") {
			return true
		}
	}
	return false
}

func defaultIgnoreDirs() map[string]bool {
	return map[string]bool{
		".git":        true,
		"node_modules": true,
		"vendor":      true,
		"__pycache__": true,
		".venv":       true,
		"venv":        true,
		"dist":        true,
		"build":       true,
		".doc-torn":    true,
	}
}

func countLines(path string) int {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	if len(data) == 0 {
		return 0
	}
	n := 1
	for _, b := range data {
		if b == '\n' {
			n++
		}
	}
	return n
}

func WalkFilesystem(root string) ([]FileEntry, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	ignoreDirs := defaultIgnoreDirs()
	var entries []FileEntry

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		// Skip ignored directories
		if isIgnored(path, info, ignoreDirs) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() {
			return nil
		}
		entries = append(entries, FileEntry{
			Path:      rel,
			Dir:       filepath.Dir(rel),
			Extension: filepath.Ext(rel),
			SizeBytes: info.Size(),
			Lines:     countLines(path),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})
	return entries, nil
}
```

- [ ] **Step 2: Verify build**

Run: `cd tools/doc-torn-scan && go build ./scan/`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add tools/doc-torn-scan/scan/scan.go
git commit -m "feat(doc-torn-scan): filesystem tree walk with .gitignore-aware filtering"
```

### Task 3: Generate Package — Scaffold

**Files:**
- Create: `tools/doc-torn-scan/generate/scaffold.go`

- [ ] **Step 1: Write generate/scaffold.go**

```go
package generate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Anhydrite/doc-torn/tools/doc-torn-scan/state"
)

func docsFeatureDir(featureName string) string {
	return filepath.Join("docs", "features", featureName)
}

func ensureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func writeFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func ScaffoldFeature(feature *state.Feature) ([]string, error) {
	baseDir := docsFeatureDir(feature.Name)
	var created []string

	// L1: README.md
	l1Path := filepath.Join(baseDir, "README.md")
	l1Content := fmt.Sprintf(`# %s

## Functional Objective

<!-- EXPLANATION: What does this feature do? Why does it exist? -->

## Technical Logic

<!-- EXPLANATION: How does it work? Main flow. -->

## Dependencies

- **Upstream**: <!-- EXPLANATION: what this feature depends on -->
- **Downstream**: <!-- EXPLANATION: what depends on this feature -->

## API / Interface

<!-- EXPLANATION: How is this feature invoked or consumed? -->

## Key Files

%s
`, feature.Name, fileListMarkdown(feature.Files))
	if err := writeFile(l1Path, l1Content); err != nil {
		return nil, fmt.Errorf("writing L1: %w", err)
	}
	created = append(created, l1Path)

	// L2: Sub-features
	for _, sf := range feature.SubFeatures {
		sfPath := filepath.Join(baseDir, "sub-features", sf+".md")
		sfContent := fmt.Sprintf(`# %s — %s

## Business Rules

<!-- EXPLANATION: What business rules apply? -->

## Edge Cases

<!-- EXPLANATION: What edge cases exist? -->
`, feature.Name, sf)
		if err := writeFile(sfPath, sfContent); err != nil {
			return nil, fmt.Errorf("writing L2 %s: %w", sf, err)
		}
		created = append(created, sfPath)
	}

	// L3: Implementation details
	for _, id := range feature.ImplementationDetails {
		idPath := filepath.Join(baseDir, "implementation", id+".md")
		idContent := fmt.Sprintf(`# %s — %s

## Why This Approach

<!-- EXPLANATION: Why was this technical decision made? -->

## Tradeoffs

<!-- EXPLANATION: What alternatives were considered? -->
`, feature.Name, id)
		if err := writeFile(idPath, idContent); err != nil {
			return nil, fmt.Errorf("writing L3 %s: %w", id, err)
		}
		created = append(created, idPath)
	}

	return created, nil
}

func fileListMarkdown(files []string) string {
	if len(files) == 0 {
		return "_No source files listed._"
	}
	var result string
	for _, f := range files {
		result += fmt.Sprintf("- `%s`\n", f)
	}
	return result
}
```

- [ ] **Step 2: Verify build**

Run: `cd tools/doc-torn-scan && go build ./generate/`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add tools/doc-torn-scan/generate/scaffold.go
git commit -m "feat(doc-torn-scan): markdown scaffold generation for L1/L2/L3"
```

### Task 4: Generate Package — Meta Docs

**Files:**
- Create: `tools/doc-torn-scan/generate/meta.go`

- [ ] **Step 1: Write generate/meta.go**

```go
package generate

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Anhydrite/doc-torn/tools/doc-torn-scan/state"
)

func Meta(s *state.DocTornState) ([]string, error) {
	var created []string

	// docs/README.md (L0)
	l0, err := generateL0(s)
	if err != nil {
		return nil, fmt.Errorf("L0: %w", err)
	}
	if err := writeFile("docs/README.md", l0); err != nil {
		return nil, err
	}
	created = append(created, "docs/README.md")

	// docs/architecture/architecture.md
	arch, err := generateArchitecture(s)
	if err != nil {
		return nil, fmt.Errorf("architecture: %w", err)
	}
	if err := writeFile("docs/architecture/architecture.md", arch); err != nil {
		return nil, err
	}
	created = append(created, "docs/architecture/architecture.md")

	// docs/architecture/dependency-matrix.md
	matrix := generateDependencyMatrix(s)
	if err := writeFile("docs/architecture/dependency-matrix.md", matrix); err != nil {
		return nil, err
	}
	created = append(created, "docs/architecture/dependency-matrix.md")

	// docs/user/definitions.md
	defs := generateDefinitions(s)
	if err := writeFile("docs/user/definitions.md", defs); err != nil {
		return nil, err
	}
	created = append(created, "docs/user/definitions.md")

	// docs/user/dev-process.md
	dev := generateDevProcess()
	if err := writeFile("docs/user/dev-process.md", dev); err != nil {
		return nil, err
	}
	created = append(created, "docs/user/dev-process.md")

	// AGENTS.md
	agents := generateAgents(s)
	if err := writeFile("AGENTS.md", agents); err != nil {
		return nil, err
	}
	created = append(created, "AGENTS.md")

	return created, nil
}

func generateL0(s *state.DocTornState) (string, error) {
	var featureList string
	for _, name := range s.Order {
		for _, f := range s.Features {
			if f.Name == name {
				featureList += fmt.Sprintf("- [%s](features/%s/README.md)\n", f.Name, f.Name)
			}
		}
	}

	// Build Mermaid architecture diagram
	var diagram string
	diagram = "graph TD\n"
	for _, name := range s.Order {
		escaped := strings.ReplaceAll(name, "-", "_")
		diagram += fmt.Sprintf("    %s[%s]\n", escaped, name)
	}
	for _, name := range s.Order {
		for _, f := range s.Features {
			if f.Name == name {
				for _, dep := range f.Dependencies {
					depEscaped := strings.ReplaceAll(dep, "-", "_")
					nameEscaped := strings.ReplaceAll(name, "-", "_")
					diagram += fmt.Sprintf("    %s --> %s\n", depEscaped, nameEscaped)
				}
				break
			}
		}
	}

	return fmt.Sprintf(`# Project Documentation

## In one line

<!-- EXPLANATION: One-line project summary -->

## Architecture

`+"```mermaid"+`
%s
`+"```"+`

## Major Features

%s

## External Dependencies

<!-- EXPLANATION: List external dependencies -->
`, diagram, featureList), nil
}

func generateArchitecture(s *state.DocTornState) (string, error) {
	return `# Functional Architecture

## Block Diagram

` + "```mermaid" + `
graph TD
` + generateMermaidBlocks(s) + `
` + "```" + `

## Data Flows

<!-- EXPLANATION: Describe main data flows -->

## Key Boundaries

<!-- EXPLANATION: Describe architectural boundaries -->
`, nil
}

func generateMermaidBlocks(s *state.DocTornState) string {
	var sb strings.Builder
	for _, name := range s.Order {
		escaped := strings.ReplaceAll(name, "-", "_")
		sb.WriteString(fmt.Sprintf("    %s[%s]\n", escaped, name))
	}
	for _, name := range s.Order {
		for _, f := range s.Features {
			if f.Name == name {
				for _, dep := range f.Dependencies {
					depEscaped := strings.ReplaceAll(dep, "-", "_")
					nameEscaped := strings.ReplaceAll(name, "-", "_")
					sb.WriteString(fmt.Sprintf("    %s --> %s\n", depEscaped, nameEscaped))
				}
				break
			}
		}
	}
	return sb.String()
}

func generateDependencyMatrix(s *state.DocTornState) string {
	var sb strings.Builder
	sb.WriteString("# Dependency Matrix\n\n")
	sb.WriteString("| Feature | Depends On | Used By |\n")
	sb.WriteString("|---------|------------|--------|\n")

	for _, name := range s.Order {
		for _, f := range s.Features {
			if f.Name == name {
				deps := "—"
				if len(f.Dependencies) > 0 {
					deps = strings.Join(f.Dependencies, ", ")
				}
				usedBy := findUsedBy(s, name)
				sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", f.Name, deps, usedBy))
				break
			}
		}
	}
	return sb.String()
}

func findUsedBy(s *state.DocTornState, featureName string) string {
	var users []string
	for _, f := range s.Features {
		for _, dep := range f.Dependencies {
			if dep == featureName {
				users = append(users, f.Name)
			}
		}
	}
	if len(users) == 0 {
		return "—"
	}
	return strings.Join(users, ", ")
}

func generateDefinitions(s *state.DocTornState) string {
	return `# Definitions

## Core Concepts

| Term | Definition |
|------|------------|
| L0 | Highest-level documentation: project overview, architecture diagram, feature list |
| L1 | Per-feature documentation: objective, logic, dependencies, API, key files |
| L2 | Sub-feature documentation: edge cases, business rules, sub-flows |
| L3 | Implementation documentation: technical decisions, rationale |

<!-- EXPLANATION: Add project-specific terms -->
`
}

func generateDevProcess() string {
	return `# Development Process

## Build Commands

<!-- EXPLANATION: List build commands -->

## Test Commands

<!-- EXPLANATION: List test commands -->

## Documentation Workflow

### First Time (init)
1. Load structured-documentation skill
2. Follow iterative init workflow

### Before a Feature (read)
1. Load doc-driven-exploration skill
2. Read skeleton files
3. Read relevant feature docs

### After a Feature (update)
1. Load structured-documentation update mode
2. Update impacted docs

## Conventions

<!-- EXPLANATION: Code style, naming, commit conventions -->
`
}

func generateAgents(s *state.DocTornState) string {
	// Note: this generates a minimal AGENTS.md. The agent should review and enhance.
	var featureList string
	for _, name := range s.Order {
		for _, f := range s.Features {
			if f.Name == name {
				featureList += fmt.Sprintf("- [%s](docs/features/%s/README.md)\n", f.Name, f.Name)
				break
			}
		}
	}

	return fmt.Sprintf(`# Project — Agent Guide

## Business Stakes

<!-- EXPLANATION: Why this project exists, business impact -->

## Features

%s
## Architecture in 30s

`+"```mermaid"+`
graph TD
%s
`+"```"+`

See docs/architecture/architecture.md for the full diagram.

## Agent Rules

### Every search → doc-driven-exploration
Load the doc skeleton before opening any code.

### After each feature → structured-documentation update
Create/update feature doc, sub-features, dependency matrix, definitions.

### Code modified/tested/completed → documentation-consistency
Full audit of all docs against real code.

## Key Definitions

See docs/user/definitions.md for the full glossary.
`, featureList, generateMermaidBlocks(s))
}
```

- [ ] **Step 2: Verify build**

Run: `cd tools/doc-torn-scan && go build ./generate/`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add tools/doc-torn-scan/generate/meta.go
git commit -m "feat(doc-torn-scan): meta-doc generation for L0, architecture, matrix, definitions, dev-process, AGENTS.md"
```

### Task 5: Command Handlers

**Files:**
- Create: `tools/doc-torn-scan/cmd/tree.go`
- Create: `tools/doc-torn-scan/cmd/scaffold.go`
- Create: `tools/doc-torn-scan/cmd/complete.go`
- Create: `tools/doc-torn-scan/cmd/meta.go`
- Create: `tools/doc-torn-scan/cmd/status.go`

- [ ] **Step 1: Write cmd/tree.go**

```go
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Anhydrite/doc-torn/tools/doc-torn-scan/scan"
)

func RunTree() error {
	entries, err := scan.WalkFilesystem(".")
	if err != nil {
		return fmt.Errorf("walking filesystem: %w", err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}
```

- [ ] **Step 2: Write cmd/scaffold.go**

```go
package cmd

import (
	"fmt"
	"os"

	"github.com/Anhydrite/doc-torn/tools/doc-torn-scan/generate"
	"github.com/Anhydrite/doc-torn/tools/doc-torn-scan/state"
)

func RunScaffold(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: doc-torn-scan scaffold <feature-name>")
	}
	featureName := args[0]

	s, err := state.Read()
	if err != nil {
		return fmt.Errorf("reading state: %w", err)
	}

	var feature *state.Feature
	for i := range s.Features {
		if s.Features[i].Name == featureName {
			feature = &s.Features[i]
			break
		}
	}
	if feature == nil {
		return fmt.Errorf("feature %q not found in state", featureName)
	}

	created, err := generate.ScaffoldFeature(feature)
	if err != nil {
		return fmt.Errorf("scaffolding: %w", err)
	}

	fmt.Println("Created:")
	for _, path := range created {
		fmt.Printf("  %s\n", path)
	}
	return nil
}
```

- [ ] **Step 3: Write cmd/complete.go**

```go
package cmd

import (
	"fmt"

	"github.com/Anhydrite/doc-torn/tools/doc-torn-scan/state"
)

func RunComplete(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: doc-torn-scan complete <feature-name>")
	}
	featureName := args[0]

	s, err := state.Read()
	if err != nil {
		return fmt.Errorf("reading state: %w", err)
	}

	if err := s.Complete(featureName); err != nil {
		return fmt.Errorf("completing feature: %w", err)
	}

	fmt.Printf("Feature %q marked as completed.\n", featureName)
	if s.CurrentFeature != "" {
		fmt.Printf("Next feature: %s\n", s.CurrentFeature)
	} else {
		fmt.Println("All features completed. Run `doc-torn-scan meta` to generate global docs.")
	}
	return nil
}
```

- [ ] **Step 4: Write cmd/meta.go**

```go
package cmd

import (
	"fmt"

	"github.com/Anhydrite/doc-torn/tools/doc-torn-scan/generate"
	"github.com/Anhydrite/doc-torn/tools/doc-torn-scan/state"
)

func RunMeta() error {
	s, err := state.Read()
	if err != nil {
		return fmt.Errorf("reading state: %w", err)
	}

	if !s.AllCompleted() {
		return fmt.Errorf("not all features are completed yet. Run `doc-torn-scan status` to check progress")
	}

	created, err := generate.Meta(s)
	if err != nil {
		return fmt.Errorf("generating meta docs: %w", err)
	}

	s.MetaDocsGenerated = true
	if err := s.Write(); err != nil {
		return fmt.Errorf("updating state: %w", err)
	}

	fmt.Println("Meta-docs generated:")
	for _, path := range created {
		fmt.Printf("  %s\n", path)
	}
	fmt.Println("\nPlease review and adjust the generated files.")
	return nil
}
```

- [ ] **Step 5: Write cmd/status.go**

```go
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Anhydrite/doc-torn/tools/doc-torn-scan/state"
)

func RunStatus() error {
	s, err := state.Read()
	if err != nil {
		return fmt.Errorf("reading state: %w", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "Feature\tFiles\tStatus\tCurrent")
	fmt.Fprintln(w, "-------\t-----\t------\t-------")

	for _, name := range s.Order {
		for _, f := range s.Features {
			if f.Name == name {
				current := ""
				if s.CurrentFeature == name {
					current = ">"
				}
				fmt.Fprintf(w, "%s\t%d\t%s\t%s\n", f.Name, len(f.Files), f.Status, current)
				break
			}
		}
	}
	w.Flush()

	if s.AllCompleted() {
		fmt.Println("\nAll features completed!")
	} else if s.CurrentFeature != "" {
		fmt.Printf("\nNext to document: %s\n", s.CurrentFeature)
	}
	return nil
}
```

- [ ] **Step 6: Verify build**

Run: `cd tools/doc-torn-scan && go build ./cmd/`
Expected: no errors

- [ ] **Step 7: Commit**

```bash
git add tools/doc-torn-scan/cmd/
git commit -m "feat(doc-torn-scan): all 5 command handlers (tree, scaffold, complete, meta, status)"
```

### Task 6: Main Entrypoint + Build

**Files:**
- Create: `tools/doc-torn-scan/main.go`

- [ ] **Step 1: Write main.go**

```go
package main

import (
	"fmt"
	"os"

	"github.com/Anhydrite/doc-torn/tools/doc-torn-scan/cmd"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	subcommand := os.Args[1]
	args := os.Args[2:]

	var err error
	switch subcommand {
	case "tree":
		err = cmd.RunTree()
	case "scaffold":
		err = cmd.RunScaffold(args)
	case "complete":
		err = cmd.RunComplete(args)
	case "meta":
		err = cmd.RunMeta()
	case "status":
		err = cmd.RunStatus()
	case "help", "--help", "-h":
		printUsage()
		return
	default:
		err = fmt.Errorf("unknown subcommand: %s", subcommand)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`doc-torn-scan — Iterative documentation scanner

Usage:
  doc-torn-scan tree              List all project files as JSON
  doc-torn-scan scaffold <name>    Generate markdown skeleton for a feature
  doc-torn-scan complete <name>    Mark a feature as completed
  doc-torn-scan meta               Generate global documentation (L0, arch, matrix, etc.)
  doc-torn-scan status             Show documentation progress`)
}
```

- [ ] **Step 2: Build and verify**

Run: `cd tools/doc-torn-scan && go build -o doc-torn-scan .`
Expected: binary created at `tools/doc-torn-scan/doc-torn-scan`

- [ ] **Step 3: Quick smoke test**

Run: `cd tools/doc-torn-scan && ./doc-torn-scan tree 2>&1 | head -5`
Expected: JSON array of files

- [ ] **Step 4: Commit**

```bash
git add tools/doc-torn-scan/main.go
git commit -m "feat(doc-torn-scan): main entrypoint with subcommand dispatch"
# Add the binary to .gitignore first
echo "tools/doc-torn-scan/doc-torn-scan" >> .gitignore
git add .gitignore
git commit -m "chore: add doc-torn-scan binary to gitignore"
```

### Task 7: Update structured-documentation SKILL.md

**Files:**
- Modify: `skills/structured-documentation/SKILL.md`

- [ ] **Step 1: Replace init mode workflow with iterative version**

Replace the current `init` workflow in SKILL.md with the iterative workflow that references `doc-torn-scan`.

The `init` section should become:

```markdown
### Mode `init` — Iterative documentation (AUTONOMOUS)

**Prerequisite:** Build `doc-torn-scan` from `tools/doc-torn-scan/`:
```bash
cd tools/doc-torn-scan && go build -o doc-torn-scan .
```

**Phase 1 — Discovery (agent-driven)**

Run `./tools/doc-torn-scan/doc-torn-scan tree` to get the full file tree as JSON.
Study the output to identify features — a feature is a logical unit of functionality,
not necessarily a directory. Name features with business-meaningful names.

Write `.doc-torn-state.json` with:
- Feature name, source files, sub-features, implementation details
- Dependencies between features
- Topological order (dependencies first, dependents after)

**Phase 2 — Itération feature par feature**

Repeat until all features are completed:

1. Read the source files for the current feature
2. Understand the code (structure, logic, edge cases, business rules)
3. Run `./tools/doc-torn-scan/doc-torn-scan scaffold <feature-name>`
4. Open each generated file and replace `<!-- EXPLANATION -->` sections with real content (the "why")
5. Run `./tools/doc-torn-scan/doc-torn-scan complete <feature-name>`
6. Commit if desired

**Phase 3 — Méta-docs**

When all features are complete:
1. Run `./tools/doc-torn-scan/doc-torn-scan meta`
2. Review and adjust generated files:
   - docs/README.md (L0) — verify one-line summary, architecture diagram, feature list
   - docs/architecture/architecture.md — verify block diagram, add data flows
   - docs/architecture/dependency-matrix.md — verify dependencies
   - docs/user/definitions.md — add missing terms
   - docs/user/dev-process.md — add build/test commands
   - AGENTS.md — verify stakes, rules, architecture
3. Run `./tools/doc-torn-scan/doc-torn-scan status` to confirm completion
```

- [ ] **Step 2: Verify and commit**

```bash
git add skills/structured-documentation/SKILL.md
git commit -m "feat(sd): update init mode with iterative doc-torn-scan workflow"
```

### Task 8: Push

- [ ] **Step 1: Push everything**

```bash
git push --force-with-lease origin main
```

## Spec Coverage Map

| Spec Requirement | Task |
|-----------------|------|
| `doc-torn-scan tree` subcommand | Task 2 (scan pkg) + Task 5 (cmd/tree.go) |
| `doc-torn-scan scaffold` subcommand | Task 3 (scaffold gen) + Task 5 (cmd/scaffold.go) |
| `doc-torn-scan complete` subcommand | Task 1 (state pkg) + Task 5 (cmd/complete.go) |
| `doc-torn-scan meta` subcommand | Task 4 (meta gen) + Task 5 (cmd/meta.go) |
| `doc-torn-scan status` subcommand | Task 5 (cmd/status.go) |
| State file (`.doc-torn-state.json`) read/write | Task 1 (state/state.go) |
| Filesystem tree walk with ignore | Task 2 (scan/scan.go) |
| Agent-driven feature identification | Task 7 (SKILL.md — agent studies `tree` output) |
| Agent writes state file | Task 7 (SKILL.md — agent writes .doc-torn-state.json) |
| Skeleton generation with EXPLANATION comments | Task 3 (generate/scaffold.go) |
| Meta-doc generation (L0, arch, matrix, definitions, dev-process, AGENTS.md) | Task 4 (generate/meta.go) |
