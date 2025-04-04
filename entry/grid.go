package entry

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

// renderGrid formats entries into a grid, adapting to terminal width.
func renderGrid(entries []Entry) (string, error) {
	if len(entries) == 0 {
		return "", nil
	}
	names, maxLen := extractNames(entries)
	termWidth, _ := terminalWidth()
	cols := min(max(termWidth/(maxLen+2), 1), len(names))
	rows := (len(names) + cols - 1) / cols
	return buildGrid(names, maxLen, cols, rows), nil
}

// extractNames pulls names from entries and finds the longest name length.
func extractNames(entries []Entry) ([]string, int) {
	names := make([]string, len(entries))
	maxLen := 0
	for i, e := range entries {
		names[i] = e.name
		if n := utf8.RuneCountInString(e.name); n > maxLen {
			maxLen = n
		}
	}
	return names, maxLen
}

// terminalWidth gets the terminal width, defaulting to 80 if unavailable.
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

// buildGrid constructs the grid string from names.
func buildGrid(names []string, maxLen, cols, rows int) string {
	var sb strings.Builder
	sb.Grow(rows * cols * (maxLen + 2)) // Rough capacity estimate
	for i, name := range names {
		sb.WriteString(name)
		// Pad only if not at end of row and not last item
		if (i+1)%cols != 0 && i < len(names)-1 {
			pad := maxLen - utf8.RuneCountInString(name) + 2
			sb.WriteString(strings.Repeat(" ", pad))
		} else {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}
