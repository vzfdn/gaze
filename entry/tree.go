package entry

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	connectorBranch = "├── "
	connectorLast   = "└── "
	subPrefix       = "│   "
	subPrefixLast   = "    "
)

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
	tree := make([]Entry, 0, estimatedCapacity)

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
		tree = append(tree, Entry{FileInfo: fi, path: path})
	}

	for i, e := range entries {
		isLast := i == len(entries)-1
		connector := connectorBranch
		if isLast {
			connector = connectorLast
		}

		e.treePrefix = color.treePrefix(prefix + connector)
		tree = append(tree, e)

		// collect subdirectory entries
		if e.IsDir() {
			subPath := filepath.Join(path, e.Name())
			subEntries, err := readEntries(subPath)
			if err != nil {
				continue // Skip unreadable directories
			}
			subPrefixNext := prefix + subPrefix
			if isLast {
				subPrefixNext = prefix + subPrefixLast
			}
			subTree, err := addTreePrefixes(subPath, subEntries, subPrefixNext, depth+1)
			if err != nil {
				continue // Skip problematic subdirectories
			}
			tree = append(tree, subTree...)
		}
	}

	return tree, nil
}
