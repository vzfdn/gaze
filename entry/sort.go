package entry

import (
	"cmp"
	"path/filepath"
	"slices"
)

// sortEntries sorts entries according to the given configuration.
func sortEntries(entries []Entry, cfg Config) {
	switch {
	case cfg.Size:
		sortBySize(entries)
	case cfg.Time:
		sortByTime(entries)
	case cfg.Kind:
		sortByKind(entries)
	case cfg.Ext:
		sortByExt(entries)
	}
	if cfg.Reverse {
		reverse(entries)
	}
}

// sortBySize sorts entries in descending order by file size.
func sortBySize(entries []Entry) {
	slices.SortStableFunc(entries, func(a, b Entry) int {
		return cmp.Compare(b.size(), a.size())
	})
}

// sortByTime sorts entries in descending order by modification time.
func sortByTime(entries []Entry) {
	slices.SortStableFunc(entries, func(a, b Entry) int {
		return cmp.Compare(b.info.ModTime().Unix(), a.info.ModTime().Unix())
	})
}

// sortByKind sorts entries by file type (directories first, then files).
func sortByKind(entries []Entry) {
	slices.SortStableFunc(entries, func(a, b Entry) int {
		if a.info.IsDir() && !b.info.IsDir() {
			return -1
		}
		if !a.info.IsDir() && b.info.IsDir() {
			return 1
		}
		return cmp.Compare(a.info.Name(), b.info.Name())
	})
}

// sortByExt sorts entries alphabetically by file extension.
func sortByExt(entries []Entry) {
	slices.SortStableFunc(entries, func(a, b Entry) int {
		return cmp.Compare(filepath.Ext(a.info.Name()), filepath.Ext(b.info.Name()))
	})
}

// reverse reverses the order of entries in-place.
func reverse(entries []Entry) {
	slices.Reverse(entries)
}
