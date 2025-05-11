package entry

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

// renderGrid formats entries into a grid, adapting to terminal width.
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

// longestEntryName returns the length of the longest name in runes.
func longestEntryName(entries []Entry) int {
	var maxLen int
	for _, e := range entries {
		displayName := e.DisplayName()
		if n := getVisibleWidth(displayName); n > maxLen {
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

// getVisibleWidth returns the visible width of a string by removing ANSI color codes
func getVisibleWidth(s string) int {
	ansiPattern := regexp.MustCompile("\x1b\\[[0-9;]*m")
	cleaned := ansiPattern.ReplaceAllString(s, "")
	return utf8.RuneCountInString(cleaned)
}

// buildGrid constructs the grid string from names.
func buildGrid(entries []Entry, maxLen, cols, rows int) string {
	var sb strings.Builder
	sb.Grow(rows * cols * (maxLen + 2)) // Rough capacity estimate
	for i, e := range entries {
		displayName := e.DisplayName()
		sb.WriteString(displayName)
		// Pad only if not at end of row and not last item
		if (i+1)%cols != 0 && i < len(entries)-1 {
			// Calculate visible width of the displayName
			visibleWidth := getVisibleWidth(displayName)
			pad := maxLen - visibleWidth + 2
			sb.WriteString(strings.Repeat(" ", pad))
		} else {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}
