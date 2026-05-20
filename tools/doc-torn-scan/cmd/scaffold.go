package cmd

import (
	"fmt"

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
