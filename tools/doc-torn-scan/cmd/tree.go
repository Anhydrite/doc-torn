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
