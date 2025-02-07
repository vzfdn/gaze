package entry

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Flags struct {
	All  bool
	Grid bool
	Long bool
}

// ParseFlags parses command-line flags, supporting short and long options.
func ParseFlags() (Flags, error) {
	var flg Flags
	f := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ContinueOnError)

	f.BoolVar(&flg.All, "a", false, "include hidden entries")
	f.BoolVar(&flg.All, "all", flg.All, "alias for -a")
	f.BoolVar(&flg.Grid, "g", false, "display as grid (default)")
	f.BoolVar(&flg.Grid, "grid", flg.Grid, "alias for -g")
	f.BoolVar(&flg.Long, "l", false, "detailed listing format")
	f.BoolVar(&flg.Long, "long", flg.Long, "alias for -l")

	args := make([]string, 0, len(os.Args[1:]))
	for _, arg := range os.Args[1:] {
		if len(arg) > 2 && arg[0] == '-' && arg[1] != '-' {
			// Split combined short flags (e.g., "-al" → "-a", "-l")
			for _, c := range arg[1:] {
				args = append(args, "-"+string(c))
			}
		} else {
			args = append(args, arg)
		}
	}

	if err := f.Parse(args); err != nil {
		return Flags{}, fmt.Errorf("flags: %w", err)
	}

	return flg, nil
}

// ResolvePath returns the first non-flag argument or the current directory.
func ResolvePath() (string, error) {
	if flag.NArg() > 0 {
		return flag.Arg(0), nil
	}
	return os.Getwd()
}

// Format generates output based on entries and configuration.
func Format(entries []Entry, flg Flags) string {
	if flg.Long {
		return RenderLong(entries)
	}
	return RenderGrid(entries)
}
