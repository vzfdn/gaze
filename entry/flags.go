package entry

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds command-line configuration options.
type Config struct {
	All     bool
	Grid    bool
	Long    bool
	Header  bool
	Recurse bool
}

// ParseConfig parses command-line flags and returns the configuration.
func ParseConfig() (Config, *flag.FlagSet, error) {
	var cfg Config
	f := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ContinueOnError)
	f.BoolVar(&cfg.All, "a", false, "include hidden entries")
	f.BoolVar(&cfg.All, "all", false, "alias for -a")
	f.BoolVar(&cfg.Grid, "g", false, "display as grid (default)")
	f.BoolVar(&cfg.Grid, "grid", false, "alias for -g")
	f.BoolVar(&cfg.Long, "l", false, "detailed listing format")
	f.BoolVar(&cfg.Long, "long", false, "alias for -l")
	f.BoolVar(&cfg.Header, "h", false, "show a header row for long format columns")
	f.BoolVar(&cfg.Header, "header", false, "alias for -h")
	f.BoolVar(&cfg.Recurse, "R", false, "list subdirectories recursively")
	f.BoolVar(&cfg.Recurse, "recursive", false, "alias for -R")

	// splitCombinedFlags splits combined short flags (e.g., "-al" to "-a -l").
	args := splitCombinedFlags(os.Args[1:])

	if err := f.Parse(args); err != nil {
		return Config{}, nil, err
	}

	// Set grid as default if no format flags are provided.
	if !cfg.Long && !cfg.Grid {
		cfg.Grid = true
	}

	return cfg, f, nil
}

// splitCombinedFlags splits combined short flags into individual flags.
// Preserves long flags (starting with --) and skips lone -.
func splitCombinedFlags(args []string) []string {
	var result []string
	for _, arg := range args {
		if arg == "-" {
			continue // Skip lone -
		}
		if len(arg) > 2 && arg[0] == '-' && arg[1] != '-' {
			for _, c := range arg[1:] {
				result = append(result, "-"+string(c))
			}
		} else {
			result = append(result, arg)
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
