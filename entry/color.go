package entry

import (
	"os"
	"path/filepath"
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
	maxAnsiSeqLen    = len(ansiEscapePrefix) + len(colorReadPerm) + 1 + len(resetCode)
)

// Default colors for file types when LS_COLORS is not set.
var fallbackColors = map[string]string{
	typeDir:  "34", // blue
	typeLink: "36", // cyan
	typeExec: "32", // green
	typeFile: "0",  // default
}

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
	return !c.isTTY || cfg.NoColor
}

func (c colorizer) colorCode(e Entry, fileName string) string {
	if ext := filepath.Ext(fileName); ext != "" {
		if color, ok := c.lsColors["*"+ext]; ok {
			return color
		}
	}
	fileType, _ := e.Classify()
	if color, ok := c.lsColors[fileType]; ok {
		return color
	}
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
	if c.disabled() {
		return mode.String()
	}
	// Preallocate buffer to avoid repeated allocations
	permStr := mode.String()
	var sb strings.Builder
	sb.Grow(len(permStr) * maxAnsiSeqLen)
	for i := range permStr {
		ch := permStr[i]
		var color string
		switch ch {
		case 'r':
			color = colorReadPerm
		case 'w':
			color = colorWritePerm
		case 'x', 's', 'S', 't', 'T':
			color = colorExecPerm
		default:
			color = colorPlaceholder
		}
		sb.WriteString(ansiEscapePrefix)
		sb.WriteString(color)
		sb.WriteByte('m')
		sb.WriteByte(ch)
		sb.WriteString(resetCode)
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
	if !strings.Contains(s, ansiEscapePrefix) {
		return utf8.RuneCountInString(s)
	}
	count := 0
	for i := 0; i < len(s); {
		// Detect ANSI sequence start
		if i+1 < len(s) && s[i] == ansiEscapePrefix[0] && s[i+1] == ansiEscapePrefix[1] {
			// Skip prefix and code until 'm'.
			i += len(ansiEscapePrefix)
			for i < len(s) && s[i] != 'm' {
				i++
			}
			if i < len(s) {
				i++ // Skip 'm'
			}
			continue
		}
		// Count visible rune and advance by its byte length.
		_, size := utf8.DecodeRuneInString(s[i:])
		count++
		i += size
	}
	return count
}
