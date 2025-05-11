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
func (c colorizer) colorize(fi os.FileInfo, displayName string) string {
	if !c.isTTY {
		return displayName
	}
	colorCode := c.lookupColorCode(fi,  displayName)
	if colorCode == "" || colorCode == "0" {
		return displayName
	}
	return fmt.Sprintf("\x1b[%sm%s\x1b[0m", colorCode, displayName)
}

// lookupColorCode finds the appropriate ANSI color code for a file.
func (c colorizer) lookupColorCode(fi os.FileInfo, fileName string) string {
	// Check for file extension colors first
	if ext := strings.ToLower(filepath.Ext(fileName)); ext != "" {
		if color, ok := c.lsColors["*"+ext]; ok {
			return color
		}
	}
	fileType, _ := fileType(fi)  
	// Fall back to file type colors
	if color, ok := c.lsColors[fileType]; ok {
		return color
	}
	// Use fallback colors as last resort
	return fallbackColors[fileType]
}

// parseLSColors converts LS_COLORS environment variable to a color map.
func parseLSColors() map[string]string {
	colors := make(map[string]string)
	if lsColors := os.Getenv("LS_COLORS"); lsColors != "" {
		for _, seq := range strings.Split(lsColors, ":") {
			if parts := strings.SplitN(seq, "=", 2); len(parts) == 2 {
				colors[parts[0]] = parts[1]
			}
		}
	}
	return colors
}
