package entry

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// renderTree renders the entries in a tree-like format.
func renderTree(entries []Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		sb.WriteString(e.displayName)
		sb.WriteByte('\n')
	}
	return sb.String()
}

// addTreePrefixes adds tree-like prefixes to directory entries and collects subdirectory entries.
// If a single entry is not a directory, it is returned without a tree structure.
func addTreePrefixes(path string, entries []Entry, cfg Config, prefix string, depth int, c colorizer) ([]Entry, error) {
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
		result = append(result, NewEntry(fi, getDisplayName(fi, cfg, c), path, ""))
	}

	for i, e := range entries {
		isLast := i == len(entries)-1
		connector := "├── "
		if isLast {
			connector = "└── "
		}

		// Create new entry with colored name + uncolored prefix
		coloredEntry := NewEntry(e.FileInfo, prefix+connector+e.displayName, e.path, e.target)
		result = append(result, coloredEntry)

		if e.IsDir() {
			subPath := filepath.Join(e.path, e.Name())
			subEntries, err := ReadEntries(subPath, cfg, c)
			if err != nil {
				continue // Skip unreadable directories
			}
			subPrefix := prefix + "│   "
			if isLast {
				subPrefix = prefix + "    "
			}
			subTree, err := addTreePrefixes(subPath, subEntries, cfg, subPrefix, depth+1, c)
			if err != nil {
				continue // Skip problematic subdirectories
			}
			result = append(result, subTree...)
		}
	}
	return result, nil
}
