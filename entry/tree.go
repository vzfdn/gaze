package entry

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// renderTree renders the entries in a tree-like format.
func renderTree(entries []Entry, c colorizer) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "%s\n", c.colorize(classify(e.FileInfo), e.name))
	}
	return sb.String()
}

// addTreePrefixes adds tree-like prefixes to directory entries and collects subdirectory entries.
// If a single entry is not a directory, it is returned without a tree structure.
func addTreePrefixes(path string, entries []Entry, cfg Config, prefix string, depth int) ([]Entry, error) {
	estimatedCapacity := len(entries)
	if depth == 0 {
		estimatedCapacity++ // For root directory
	}
	result := make([]Entry, 0, estimatedCapacity)

	if depth == 0 {
		fi, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("accessing path %s: %w", path, err)
		}
		// If `path` is a file, return it directly (no tree formatting)
		if !fi.IsDir() {
			return entries, nil
		}
		// If `path` is a directory, include it as the root
		result = append(result, NewEntry(fi, formatName(fi, cfg), path, ""))
	}

	for i, e := range entries {
		isLast := i == len(entries)-1
		connector := "├── "
		if isLast {
			connector = "└── "
		}
		e.name = prefix + connector + e.name
		result = append(result, e)

		if e.IsDir() {
			subPath := filepath.Join(e.path, e.Name())
			subEntries, err := ReadEntries(subPath, cfg)
			if err != nil {
				continue // Skip unreadable directories
			}
			subPrefix := prefix + "│   "
			if isLast {
				subPrefix = prefix + "    "
			}
			subTree, err := addTreePrefixes(subPath, subEntries, cfg, subPrefix, depth+1)
			if err != nil {
				continue // Skip problematic subdirectories
			}
			result = append(result, subTree...)
		}
	}
	return result, nil
}
