package entry

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

// renderGrid formats entries into a grid.
// It adapts to terminal width for a compact display.
func renderGrid(entries []Entry, cfg Config) (string, error) {
	names := make([]string, 0, len(entries))
	var maxNameLen int
	for _, entry := range entries {
		name := entry.Name(cfg)
		nameLen := utf8.RuneCountInString(name)
		if nameLen > maxNameLen {
			maxNameLen = nameLen
		}
		names = append(names, name)
	}

	termWidth, _ := terminalWidth()
	columns, rows := getTableDimensions(termWidth, maxNameLen, len(names))
	return generateTable(names, maxNameLen, columns, rows), nil
}

// generateTable builds a grid string from names using a row-major layout.
func generateTable(names []string, maxNameLen, columns, rows int) string {
	var sb strings.Builder
	// Preallocate capacity for names, padding, and newlines
	sb.Grow(len(names)*(maxNameLen+2) + rows)

	if columns <= 1 {
		for _, name := range names {
			sb.WriteString(name)
			sb.WriteByte('\n')
		}
		return sb.String()
	}

	for row := range rows {
		for col := range columns {
			// Map row,col to flat index in names
			index := row*columns + col
			if index >= len(names) {
				break
			}
			name := names[index]
			sb.WriteString(name)
			if col < columns-1 {
				// Pad to align next column
				padWidth := maxNameLen - utf8.RuneCountInString(name) + 2
				for range padWidth {
					sb.WriteByte(' ')
				}
			}
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}

// getTableDimensions computes the number of columns and rows for the grid.
func getTableDimensions(termWidth, maxNameLen, entryCount int) (columns, rows int) {
	// Divide terminal width by column width (name length + padding)
	columns = termWidth / (maxNameLen + 2)
	if columns < 1 {
		columns = 1
	}

	if columns > entryCount {
		columns = entryCount
	}

	rows = (entryCount + columns - 1) / columns
	return columns, rows
}

// terminalWidth retrieves the terminal width.
// It falls back to 80 if unavailable.
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
