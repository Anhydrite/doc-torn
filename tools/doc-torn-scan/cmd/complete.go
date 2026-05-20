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
