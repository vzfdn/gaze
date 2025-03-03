package entry

import (
	"cmp"
	"slices"
	"strings"
)

// sortBySize sorts entries in descending order by file size.
func sortBySize(entries []Entry) {
	slices.SortStableFunc(entries, func(a, b Entry) int {
		return cmp.Compare(b.Size(), a.Size())
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
		return cmp.Compare(a.Name(), b.Name())
	})
}

// sortByExt sorts entries alphabetically by file extension.
func sortByExt(entries []Entry) {
	slices.SortStableFunc(entries, func(a, b Entry) int {
		return cmp.Compare(getExt(a.Name()), getExt(b.Name()))
	})
}

// getExt returns the file extension or an empty string if none.
func getExt(name string) string {
	if dot := strings.LastIndex(name, "."); dot >= 0 {
		return name[dot:]
	}
	return ""
}

// reverse reverses the order of entries in-place.
func reverse(entries []Entry) {
	slices.Reverse(entries)
}
