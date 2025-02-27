package entry

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds command-line configuration options for directory listing.
type Config struct {
	All     bool
	Grid    bool
	Long    bool
	Header  bool
	Recurse bool
	Size    bool
	Time    bool
	Kind    bool
	Ext     bool
	Reverse bool
}

// boolFlag defines a flag with short/long names and a description.
type boolFlag struct {
	ptr       *bool
	shortName string
	longName  string
	usage     string
}

// ParseConfig parses command-line flags and returns the configuration.
func ParseConfig() (Config, *flag.FlagSet, error) {
	cfg := Config{}
	fs := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ContinueOnError)

	// Define flags concisely using a slice
	flags := []boolFlag{
		{&cfg.All, "a", "all", "include hidden entries"},
		{&cfg.Grid, "g", "grid", "display as grid (default)"},
		{&cfg.Long, "l", "long", "detailed listing format"},
		{&cfg.Header, "h", "header", "show a header row for long format"},
		{&cfg.Recurse, "R", "recursive", "list subdirectories recursively"},
		{&cfg.Size, "s", "size", "sort by file size"},
		{&cfg.Time, "t", "time", "sort by modification time"},
		{&cfg.Kind, "k", "kind", "sort by file type (kind)"},
		{&cfg.Ext, "x", "extension", "sort by file extension"},
		{&cfg.Reverse, "r", "reverse", "reverse the sorting order"},
	}

	// Register flags efficiently
	for _, f := range flags {
		fs.BoolVar(f.ptr, f.shortName, false, f.usage)
		fs.BoolVar(f.ptr, f.longName, false, "alias for -"+f.shortName)
	}

	// Parse with optimized argument splitting
	args := expandShortFlags(os.Args[1:])
	if err := fs.Parse(args); err != nil {
		return Config{}, nil, err
	}

	// Set Grid as default if no format specified
	if !cfg.Long && !cfg.Grid {
		cfg.Grid = true
	}

	return cfg, fs, nil
}

// expandShortFlags splits combined short flags for shorthand input compatibility.
// For example, it converts "-al" to "-a -l" with optimized allocations.
func expandShortFlags(args []string) []string {
	// Estimate capacity: each arg could split into multiple flags
	result := make([]string, 0, len(args)*2)
	for _, arg := range args {
		if len(arg) < 2 || arg[0] != '-' || arg[1] == '-' {
			result = append(result, arg) // Pass through lone "-", long flags, or non-flags
			continue
		}
		// Split short flags (e.g., "-al" -> "-a", "-l")
		for i := 1; i < len(arg); i++ { // Start at 1 to skip initial "-"
			result = append(result, "-"+string(arg[i]))
		}
	}
	return result
}

// ResolvePath returns the first non-flag argument or the current directory.
func ResolvePath(f *flag.FlagSet) (string, error) {
	if f.NArg() > 0 {
		return f.Arg(0), nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot get current directory: %w", err)
	}
	return wd, nil
}
