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
		sb.WriteString(e.DisplayName())
		sb.WriteByte('\n')
	}
	return sb.String()
}

// addTreePrefixes adds tree-like prefixes to directory entries and collects subdirectory entries.
func addTreePrefixes(path string, entries []Entry, prefix string, depth int) ([]Entry, error) {
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
		result = append(result, Entry{FileInfo: fi, path: path})
	}

	for i, e := range entries {
		isLast := i == len(entries)-1
		connector := "├── "
		if isLast {
			connector = "└── "
		}

		e.treePrefix = prefix + connector
		result = append(result, e)

		if e.IsDir() {
			subPath := filepath.Join(e.path, e.Name())
			subEntries, err := readEntries(subPath)
			if err != nil {
				continue // Skip unreadable directories
			}
			subPrefix := prefix + "│   "
			if isLast {
				subPrefix = prefix + "    "
			}
			subTree, err := addTreePrefixes(subPath, subEntries, subPrefix, depth+1)
			if err != nil {
				continue // Skip problematic subdirectories
			}
			result = append(result, subTree...)
		}
	}
	return result, nil
}
