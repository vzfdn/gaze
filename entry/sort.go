package entry

import (
	"cmp"
	"path/filepath"
	"slices"
	"strings"
)

// sortEntries sorts entries according to the given configuration.
func sortEntries(entries []Entry) {
	switch {
	case cfg.Size:
		sortBySize(entries)
	case cfg.Time:
		sortByTime(entries)
	case cfg.Kind:
		sortByKind(entries)
	case cfg.Ext:
		sortByExt(entries)
	default:
		sortByName(entries)
	}
	if cfg.Reverse {
		reverse(entries)
	}
}

// sortByExt sorts entries alphabetically by file name.
func sortByName(entries []Entry) {
	slices.SortStableFunc(entries, func(a, b Entry) int {
		return cmp.Compare(strings.ToLower(a.Name()), strings.ToLower(b.Name()))
	})
}

// sortBySize sorts entries in descending order by file size.
func sortBySize(entries []Entry) {
	slices.SortStableFunc(entries, func(a, b Entry) int {
		return cmp.Compare(b.Size(), a.Size())
	})
}

// sortByTime sorts entries in descending order by modification time.
func sortByTime(entries []Entry) {
	slices.SortStableFunc(entries, func(a, b Entry) int {
		return cmp.Compare(b.ModTime().Unix(), a.ModTime().Unix())
	})
}

// sortByKind sorts entries by file type (directories first, then files).
func sortByKind(entries []Entry) {
	slices.SortStableFunc(entries, func(a, b Entry) int {
		if a.IsDir() && !b.IsDir() {
			return -1
		}
		if !a.IsDir() && b.IsDir() {
			return 1
		}
		return cmp.Compare(strings.ToLower(a.Name()), strings.ToLower(b.Name()))
	})
}

// sortByExt sorts entries alphabetically by file extension.
func sortByExt(entries []Entry) {
	slices.SortStableFunc(entries, func(a, b Entry) int {
		extA, extB := filepath.Ext(a.Name()), filepath.Ext(b.Name())
		if extA != extB {
			return cmp.Compare(extA, extB)
		}
		// If extensions are equal, fallback to name comparison
		return cmp.Compare(strings.ToLower(a.Name()), strings.ToLower(b.Name()))
	})
}

// reverse reverses the order of entries in-place.
func reverse(entries []Entry) {
	slices.Reverse(entries)
}
