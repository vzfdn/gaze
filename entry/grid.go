package entry

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

// renderGrid formats entries into a grid, adapting to terminal width.
func renderGrid(entries []Entry, c colorizer) (string, error) {
	if len(entries) == 0 {
		return "", nil
	}
	maxLen := longestEntryName(entries)
	termWidth, _ := terminalWidth()
	cols := min(max(termWidth/(maxLen+2), 1), len(entries))
	rows := (len(entries) + cols - 1) / cols
	return buildGrid(entries, maxLen, cols, rows, c), nil
}

// longestEntryName returns the length of the longest name in runes.
func longestEntryName(entries []Entry) int {
	var maxLen int
	for _, e := range entries {
		if n := utf8.RuneCountInString(e.name); n > maxLen {
			maxLen = n
		}
	}
	return maxLen
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
func buildGrid(entries []Entry, maxLen, cols, rows int, c colorizer) string {
	var sb strings.Builder
	sb.Grow(rows * cols * (maxLen + 2)) // Rough capacity estimate

	for i, e := range entries {
		colored := c.colorize(classify(entries[i].FileInfo), e.name)
		sb.WriteString(colored)

		// Pad only if not at end of row and not last item
		if (i+1)%cols != 0 && i < len(entries)-1 {
			pad := maxLen - utf8.RuneCountInString(e.name) + 2
			sb.WriteString(strings.Repeat(" ", pad))
		} else {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}
