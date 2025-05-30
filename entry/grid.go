package entry

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func renderGrid(entries []Entry) (string, error) {
	if len(entries) == 0 {
		return "", nil
	}
	maxLen := longestEntryName(entries)
	termWidth, _ := terminalWidth()
	cols := min(max(termWidth/(maxLen+2), 1), len(entries))
	rows := (len(entries) + cols - 1) / cols
	return buildGrid(entries, maxLen, cols, rows), nil
}

func longestEntryName(entries []Entry) int {
	var maxLen int
	for _, e := range entries {
		displayName := e.DisplayName()
		if n := visibleWidth(displayName); n > maxLen {
			maxLen = n
		}
	}
	return maxLen
}

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

func buildGrid(entries []Entry, maxLen, cols, rows int) string {
	var sb strings.Builder
	sb.Grow(rows * cols * (maxLen + 2)) // Rough capacity estimate
	for i, e := range entries {
		displayName := e.DisplayName()
		sb.WriteString(displayName)
		// Pad only if not at end of row and not last item
		if (i+1)%cols != 0 && i < len(entries)-1 {
			// Calculate visible width of the displayName
			visibleWidth := visibleWidth(displayName)
			pad := maxLen - visibleWidth + 2
			sb.WriteString(strings.Repeat(" ", pad))
		} else {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}
