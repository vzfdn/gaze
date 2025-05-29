package entry

import (
	"cmp"
	"path/filepath"
	"slices"
	"strings"
)

type sortableEntry struct {
	entry    Entry
	lower    string
	extLower string
}

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

func sortByName(entries []Entry) {
	if len(entries) <= 1 {
		return
	}
	sorted := make([]sortableEntry, len(entries))
	for i, e := range entries {
		sorted[i] = sortableEntry{entry: e, lower: strings.ToLower(e.Name())}
	}
	slices.SortStableFunc(sorted, func(a, b sortableEntry) int {
		return cmp.Compare(a.lower, b.lower)
	})
	for i := range entries {
		entries[i] = sorted[i].entry
	}
}

func sortBySize(entries []Entry) {
	slices.SortStableFunc(entries, func(a, b Entry) int {
		return cmp.Compare(b.Size(), a.Size())
	})
}

func sortByTime(entries []Entry) {
	slices.SortStableFunc(entries, func(a, b Entry) int {
		return cmp.Compare(b.ModTime().Unix(), a.ModTime().Unix())
	})
}

func sortByKind(entries []Entry) {
	if len(entries) <= 1 {
		return
	}
	sorted := make([]sortableEntry, len(entries))
	for i, e := range entries {
		sorted[i] = sortableEntry{entry: e, lower: strings.ToLower(e.Name())}
	}
	slices.SortStableFunc(sorted, func(a, b sortableEntry) int {
		aDir, bDir := a.entry.IsDir(), b.entry.IsDir()
		if aDir && !bDir {
			return -1
		}
		if !aDir && bDir {
			return 1
		}
		return cmp.Compare(a.lower, b.lower)
	})
	for i := range entries {
		entries[i] = sorted[i].entry
	}
}

func sortByExt(entries []Entry) {
	if len(entries) <= 1 {
		return
	}
	sorted := make([]sortableEntry, len(entries))
	for i, e := range entries {
		name := e.Name()
		sorted[i] = sortableEntry{
			entry:    e,
			lower:    strings.ToLower(name),
			extLower: strings.ToLower(filepath.Ext(name)),
		}
	}
	slices.SortStableFunc(sorted, func(a, b sortableEntry) int {
		if a.extLower != b.extLower {
			return cmp.Compare(a.extLower, b.extLower)
		}
		return cmp.Compare(a.lower, b.lower)
	})
	for i := range entries {
		entries[i] = sorted[i].entry
	}
}

func reverse(entries []Entry) {
	slices.Reverse(entries)
}
