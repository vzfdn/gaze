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

type boolFlag struct {
	ptr       *bool
	shortName string
	longName  string
	usage     string
}

func boolFlags() []boolFlag {
	return []boolFlag{
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
}

// ParseFlags parses command-line flags and returns the configuration.
func ParseFlags() (*flag.FlagSet, error) {
	f := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ContinueOnError)
	for _, bf := range boolFlags() {
		f.BoolVar(bf.ptr, bf.shortName, false, bf.usage)
		f.BoolVar(bf.ptr, bf.longName, false, "alias for -"+bf.shortName)
	}

	args := expandShortFlags(os.Args[1:])
	if err := f.Parse(args); err != nil {
		return nil, err
	}
	if !cfg.Long && !cfg.Grid {
		cfg.Grid = true
	}
	return f, nil
}

// ResolvePath returns the first non-flag argument as a cleaned path,
// or "." if none is provided. It returns an error if the path is inaccessible.
func ResolvePath(f *flag.FlagSet) (string, error) {
	if f.NArg() == 0 {
		return ".", nil
	}
	path := filepath.Clean(f.Arg(0))
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("%q: no such file or directory", path)
		}
		if os.IsPermission(err) {
			return "", fmt.Errorf("%q: permission denied", path)
		}
		return "", fmt.Errorf("%q: %v", path, err)
	}
	return path, nil
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
