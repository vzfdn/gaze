package entry

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

// renderGrid formats entries as a table-like grid string.
func renderGrid(entries []Entry) (string, error) {
	var names []string
	var maxLen int
	for _, e := range entries {
		entrylen := utf8.RuneCountInString(e.Name())
		if entrylen > maxLen {
			maxLen = entrylen
		}
		names = append(names, e.Name())
	}
	tw, _ := terminalWidth()
	columns, rows := getTableDimensions(tw, maxLen, len(entries))
	return generateTable(names, maxLen, columns, rows), nil
}

// generateTable generates a formatted table string from a slice of names.
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
		if x+1 == columns && rows != 1 {
			sb.WriteString("\n")
		} else {
			// Padding ensures columns are aligned by adding space between entries
			pad := strings.Repeat(" ", maxLen-len(nameStr)+2)
			sb.WriteString(pad)
		}
	}
	return sb.String()
}

// getTableDimensions calculates the optimal number of columns and rows for a table.
func getTableDimensions(width int, maxLen int, entriesLen int) (int, int) {
	cols := width / maxLen
	if cols == 0 {
		cols = 1
	}
	rows := (entriesLen + cols - 1) / cols
	return cols, rows
}

// terminalWidth returns the current terminal width, falling back to 80 if unavailable.
func terminalWidth() (int, error) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: cannot get terminal size: %v\n", err)
		return 80, nil
	}
	if width <= 0 {
		return 80, nil
	}
	return width, nil
}
