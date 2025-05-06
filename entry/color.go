package entry

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

// Default colors for file types when LS_COLORS is not set.
var fallbackColors = map[string]string{
	typeDir:  "34", // blue
	typeLink: "36", // cyan
	typeExec: "32", // green
	typeFile: "0",  // default
}

// colorizer applies colors to filenames based on file type.
type colorizer struct {
	isTTY    bool
	lsColors map[string]string
}

// newColorizer creates a new instance that applies ANSI colors to filenames.
func newColorizer() colorizer {
	return colorizer{
		isTTY:    term.IsTerminal(int(os.Stdout.Fd())),
		lsColors: parseLSColors(),
	}
}

// colorize applies ANSI colors to a filename based on file type.
func (c colorizer) colorize(fileType, fileName string) string {
	if !c.isTTY {
		return fileName
	}
	return applyColor(fileName, lookupColorCode(fileType, fileName, c.lsColors))
}

// parseLSColors converts LS_COLORS environment variable to a color map.
func parseLSColors() map[string]string {
	colors := make(map[string]string)
	lsColors := os.Getenv("LS_COLORS")
	if lsColors == "" {
		return colors
	}
	for _, seq := range strings.Split(lsColors, ":") {
		if parts := strings.SplitN(seq, "=", 2); len(parts) == 2 {
			colors[parts[0]] = parts[1]
		}
	}
	return colors
}

// lookupColorCode finds the appropriate ANSI color code for a file.
// It checks file extensions first, then falls back to file type colors.
func lookupColorCode(fileType, fileName string, lsColors map[string]string) string {
	if ext := strings.ToLower(filepath.Ext(fileName)); ext != "" {
		if color, ok := lsColors["*"+ext]; ok {
			return color
		}
	}
	if color, ok := lsColors[fileType]; ok {
		return color
	}
	return fallbackColors[fileType]
}

// applyColor wraps text with ANSI color codes.
func applyColor(text, colorCode string) string {
	if colorCode == "" || colorCode == "0" {
		return text
	}
	return fmt.Sprintf("\x1b[%sm%s\x1b[0m", colorCode, text)
}
