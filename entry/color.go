package entry

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

const (
	colorReadPerm    = "38;5;189" // Light lavender
	colorWritePerm   = "38;5;140" // Medium purple
	colorExecPerm    = "38;5;98"  // Medium purple
	colorUser        = "38;5;97"  // Medium purple-blue
	colorGroup       = "38;5;147" // Light purple
	colorModTime     = "38;5;103" // Cool gray with purple tint
	colorSizeBytes   = "38;5;188" // Silver gray
	colorSizeKB      = "38;5;108" // Sage green
	colorSizeMB      = "38;5;175" // Mauve
	colorSizeGB      = "38;5;132" // Dusty rose
	colorSizeTB      = "38;5;61"  // Blue-purple
	colorSizePB      = "38;5;55"  // Deep purple-blue
	colorSizeEB      = "38;5;91"  // Rich purple
	colorPlaceholder = "38;5;146" // Light grayed purple
	colorTreePrefix  = "90"       // Gray
	ansiEscapePrefix = "\x1b["
	resetCode        = "\x1b[0m"
)

var (
	// Default colors for file types when LS_COLORS is not set.
	fallbackColors = map[string]string{
		typeDir:  "34", // blue
		typeLink: "36", // cyan
		typeExec: "32", // green
		typeFile: "0",  // default
	}
	ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`) // matches ANSI color codes
)

type colorizer struct {
	isTTY    bool
	lsColors map[string]string
}

func newColorizer() colorizer {
	return colorizer{
		isTTY:    term.IsTerminal(int(os.Stdout.Fd())),
		lsColors: parseLSColors(),
	}
}

func (c colorizer) disabled() bool {
	return !c.isTTY
}

func (c colorizer) colorCode(e Entry, fileName string) string {
	// Check for file extension colors first
	if ext := strings.ToLower(filepath.Ext(fileName)); ext != "" {
		if color, ok := c.lsColors["*"+ext]; ok {
			return color
		}
	}
	fileType, _ := e.Classify()
	// Fall back to file type colors
	if color, ok := c.lsColors[fileType]; ok {
		return color
	}
	// Use fallback colors as last resort
	return fallbackColors[fileType]
}

func (c colorizer) colorize(text, colorCode string) string {
	if c.disabled() {
		return text
	}
	return ansiEscapePrefix + colorCode + "m" + text + resetCode
}

func (c colorizer) fileName(e Entry, fileName string) string {
	if c.disabled() {
		return fileName
	}
	if code := c.colorCode(e, fileName); code != "" && code != "0" {
		return c.colorize(fileName, code)
	}
	return fileName
}

func (c colorizer) permissions(mode os.FileMode) string {
	permStr := mode.String()
	var sb strings.Builder
	sb.Grow(len(permStr) * 15)
	for _, perm := range permStr {
		var color string
		switch perm {
		case 'r':
			color = colorReadPerm
		case 'w':
			color = colorWritePerm
		case 'x':
			color = colorExecPerm
		default:
			color = colorPlaceholder
		}
		sb.WriteString(c.colorize(string(perm), color))
	}
	return sb.String()
}

func (c colorizer) user(text string) string        { return c.colorize(text, colorUser) }
func (c colorizer) group(text string) string       { return c.colorize(text, colorGroup) }
func (c colorizer) modTime(text string) string     { return c.colorize(text, colorModTime) }
func (c colorizer) placeholder(text string) string { return c.colorize(text, colorPlaceholder) }
func (c colorizer) treePrefix(text string) string  { return c.colorize(text, colorTreePrefix) }

func parseLSColors() map[string]string {
	ls := os.Getenv("LS_COLORS")
	if ls == "" {
		return nil
	}
	colors := make(map[string]string)
	for _, seq := range strings.Split(ls, ":") {
		if parts := strings.SplitN(seq, "=", 2); len(parts) == 2 {
			colors[parts[0]] = parts[1]
		}
	}
	return colors
}

func visibleWidth(s string) int {
	clean := ansiRegexp.ReplaceAllString(s, "")
	return utf8.RuneCountInString(clean)
}
