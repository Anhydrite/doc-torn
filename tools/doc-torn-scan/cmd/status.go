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
