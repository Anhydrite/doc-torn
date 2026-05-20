package scan

import (
	"bytes"
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

func isIgnored(info os.FileInfo, ignoreDirs map[string]bool) bool {
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
		".git":         true,
		"node_modules": true,
		"vendor":       true,
		"__pycache__":  true,
		".venv":        true,
		"venv":         true,
		"dist":         true,
		"build":        true,
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
	return bytes.Count(data, []byte{'\n'}) + 1
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
			// Permission denied etc. — skip, don't abort
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return nil
		}
		if rel == "." {
			return nil
		}
		if isIgnored(info, ignoreDirs) {
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
