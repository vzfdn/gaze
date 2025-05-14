package entry

import (
	"cmp"
	"path/filepath"
	"slices"
	"strings"
)

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
	slices.SortStableFunc(entries, func(a, b Entry) int {
		return cmp.Compare(strings.ToLower(a.Name()), strings.ToLower(b.Name()))
	})
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

func reverse(entries []Entry) {
	slices.Reverse(entries)
}
