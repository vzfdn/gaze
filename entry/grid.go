package entry

import (
	"strings"
	"unicode/utf8"
)

// RenderGrid formats and returns a table-like grid string representation of the given entries
func RenderGrid(entries []Entry) string {
	var names []string
	var maxLen int
	for _, e := range entries {
		entrylen := utf8.RuneCountInString(e.Name())
		if entrylen > maxLen {
			maxLen = entrylen
		}
		names = append(names, e.Name())
	}
	tw := GetTerminalWidth()
	columns, rows := getTableDimensions(tw, maxLen, len(entries))
	return generateTable(names, maxLen, columns, rows)
}

func generateTable(names []string, maxLen, columns, rows int) string {
	var sb strings.Builder
	if columns == 1 {
		for _, str := range names {
			sb.WriteString(str)
		}
		return sb.String()
	}
	for i := 0; i < columns*rows; i++ {
		// Compute the column (x) and row (y) for the current grid position
		x, y := i%columns, i/columns
		// Convert the 2D grid coordinates (x, y) to a linear index in the names slice
		index := y + x*rows

		var nameStr string
		if index < len(names) {
			nameStr = names[index]
			sb.WriteString(nameStr)
		}
		// Print a line break if it's the last column
		if x+1 == columns {
			sb.WriteString("\n")
		} else {
			// Padding ensures columns are aligned by adding space between entries
			pad := strings.Repeat(" ", maxLen-len(nameStr)+2)
			sb.WriteString(pad)
		}
	}
	return sb.String()
}

func getTableDimensions(width int, maxLen int, entriesLen int) (int, int) {
	cols := width / maxLen
	if cols == 0 {
		cols = 1
	}
	rows := (entriesLen + cols - 1) / cols
	return cols, rows
}
