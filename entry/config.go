package entry

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds command-line configuration options for directory listing.
type Config struct {
	All         bool
	Grid        bool
	Long        bool
	Header      bool
	Recurse     bool
	Tree        bool
	Classify    bool
	Dereference bool
	Size        bool
	Time        bool
	Kind        bool
	Ext         bool
	Reverse     bool
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
	flags := []boolFlag{
		{&cfg.All, "a", "all", "include hidden entries"},
		{&cfg.Grid, "g", "grid", "display as grid (default)"},
		{&cfg.Long, "l", "long", "detailed listing format"},
		{&cfg.Header, "h", "header", "show a header row for long format"},
		{&cfg.Classify, "F", "classify", "append file type indicators"},
		{&cfg.Recurse, "R", "recursive", "list subdirectories recursively"},
		{&cfg.Tree, "T", "tree", "recursively display directory contents as a tree"},
		{&cfg.Dereference, "L", "dereference", "show info for the target file, not the symlink"},
		{&cfg.Size, "s", "size", "sort by file size"},
		{&cfg.Time, "t", "time", "sort by modification time"},
		{&cfg.Kind, "k", "kind", "sort by file type (kind)"},
		{&cfg.Ext, "x", "extension", "sort by file extension"},
		{&cfg.Reverse, "r", "reverse", "reverse the sorting order"},
	}
	for _, f := range flags {
		fs.BoolVar(f.ptr, f.shortName, false, f.usage)
		fs.BoolVar(f.ptr, f.longName, false, "alias for -"+f.shortName)
	}
	args := expandShortFlags(os.Args[1:])
	if err := fs.Parse(args); err != nil {
		return Config{}, nil, err
	}
	if !cfg.Long && !cfg.Grid {
		cfg.Grid = true
	}
	return cfg, fs, nil
}

// expandShortFlags splits combined short flags (e.g., "-al" to "-a -l").
func expandShortFlags(args []string) []string {
	result := make([]string, 0, len(args)*2)
	for _, arg := range args {
		// Pass through non-flags, long flags, or single "-"
		if len(arg) < 2 || arg[0] != '-' || arg[1] == '-' {
			result = append(result, arg)
			continue
		}
		// Split short flags like "-al" into "-a", "-l"
		for _, flag := range arg[1:] {
			result = append(result, "-"+string(flag))
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
