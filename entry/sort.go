package entry

import (
	"cmp"
	"slices"
)

// sortBySizeDesc sorts entries in-place by Size in descending order.
func sortBySize(entries []Entry) {
	slices.SortFunc(entries, func(a, b Entry) int {
		return cmp.Compare(b.Size(), a.Size())
	})
}
