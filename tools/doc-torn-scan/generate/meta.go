package generate

import (
	"fmt"
	"strings"

	"github.com/Anhydrite/doc-torn/tools/doc-torn-scan/state"
)

func Meta(s *state.DocTornState) ([]string, error) {
	var created []string

	if err := writeFile("docs/README.md", generateL0(s)); err != nil {
		return nil, err
	}
	created = append(created, "docs/README.md")

	if err := writeFile("docs/architecture/architecture.md", generateArchitecture(s)); err != nil {
		return nil, err
	}
	created = append(created, "docs/architecture/architecture.md")

	matrix := generateDependencyMatrix(s)
	if err := writeFile("docs/architecture/dependency-matrix.md", matrix); err != nil {
		return nil, err
	}
	created = append(created, "docs/architecture/dependency-matrix.md")

	defs := generateDefinitions()
	if err := writeFile("docs/user/definitions.md", defs); err != nil {
		return nil, err
	}
	created = append(created, "docs/user/definitions.md")

	dev := generateDevProcess()
	if err := writeFile("docs/user/dev-process.md", dev); err != nil {
		return nil, err
	}
	created = append(created, "docs/user/dev-process.md")

	agents := generateAgents(s)
	if err := writeFile("AGENTS.md", agents); err != nil {
		return nil, err
	}
	created = append(created, "AGENTS.md")

	return created, nil
}

func generateL0(s *state.DocTornState) string {
	var featureList strings.Builder
	for _, name := range s.Order {
		for _, f := range s.Features {
			if f.Name == name {
				featureList.WriteString(fmt.Sprintf("- [%s](features/%s/README.md)\n", f.Name, f.Name))
			}
		}
	}

	mermaid := generateMermaid(s)

	return fmt.Sprintf(`# Project Documentation

## In one line

<!-- EXPLANATION: One-line project summary -->

## Architecture

`+"```mermaid"+`
graph TD
%s
`+"```"+`

## Major Features

%s

## External Dependencies

<!-- EXPLANATION: List external dependencies -->
`, mermaid, featureList.String())
}

func generateArchitecture(s *state.DocTornState) string {
	mermaid := generateMermaid(s)
	return fmt.Sprintf(`# Functional Architecture

## Block Diagram

`+"```mermaid"+`
graph TD
%s
`+"```"+`

## Data Flows

<!-- EXPLANATION: Describe main data flows -->

## Key Boundaries

<!-- EXPLANATION: Describe architectural boundaries -->
`, mermaid)
}

func generateMermaid(s *state.DocTornState) string {
	var sb strings.Builder
	for _, name := range s.Order {
		for _, f := range s.Features {
			if f.Name == name {
				escaped := strings.ReplaceAll(name, "-", "_")
				sb.WriteString(fmt.Sprintf("    %s[%s]\n", escaped, f.Name))
				break
			}
		}
	}
	for _, name := range s.Order {
		for _, f := range s.Features {
			if f.Name == name {
				for _, dep := range f.Dependencies {
					fromName := dep
					toName := name
					fromEscaped := strings.ReplaceAll(fromName, "-", "_")
					toEscaped := strings.ReplaceAll(toName, "-", "_")
					sb.WriteString(fmt.Sprintf("    %s --> %s\n", fromEscaped, toEscaped))
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

func generateDefinitions() string {
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
	var featureList strings.Builder
	for _, name := range s.Order {
		for _, f := range s.Features {
			if f.Name == name {
				featureList.WriteString(fmt.Sprintf("- [%s](docs/features/%s/README.md)\n", f.Name, f.Name))
				break
			}
		}
	}

	mermaid := generateMermaid(s)

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
`, featureList.String(), mermaid)
}
