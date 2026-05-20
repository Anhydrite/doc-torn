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
