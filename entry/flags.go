package entry

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the command-line configuration options.
type Config struct {
	All    bool
	Grid   bool
	Long   bool
	Header bool
}

// ParseConfig parses command-line flags and returns the configuration.
func ParseConfig() (Config, *flag.FlagSet, error) {
	var cfg Config
	f := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ContinueOnError)

	f.BoolVar(&cfg.All, "a", false, "include hidden entries")
	f.BoolVar(&cfg.All, "all", cfg.All, "alias for -a")
	f.BoolVar(&cfg.Grid, "g", false, "display as grid (default)")
	f.BoolVar(&cfg.Grid, "grid", cfg.Grid, "alias for -g")
	f.BoolVar(&cfg.Long, "l", false, "detailed listing format")
	f.BoolVar(&cfg.Long, "long", cfg.Long, "alias for -l")
	f.BoolVar(&cfg.Header, "h", false, "display a header row for each column")
	f.BoolVar(&cfg.Header, "header", cfg.Header, "alias for -h")

	// Split combined short flags (e.g., "-al" → "-a", "-l")
	args := splitCombinedFlags(os.Args[1:])

	// Parse flags
	if err := f.Parse(args); err != nil {
		return Config{}, nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	return cfg, f, nil
}

// splitCombinedFlags splits combined short flags into individual flags.
func splitCombinedFlags(args []string) []string {
	var result []string
	for _, arg := range args {
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
	return os.Getwd()
}

// Format generates output based on entries and configuration.
func Format(entries []Entry, cfg Config) (string, error) {
	switch {
	case cfg.Long && cfg.Header:
		return RenderLong(entries, true), nil
	case cfg.Long:
		return RenderLong(entries, false), nil
	default:
		return RenderGrid(entries)
	}
}
