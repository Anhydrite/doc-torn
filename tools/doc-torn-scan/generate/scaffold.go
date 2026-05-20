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
