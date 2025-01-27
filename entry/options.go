package entry

import (
	"errors"
	"flag"
	"os"
	"strings"
)

type Flags struct {
	all  *bool
	grid *bool
	long *bool
}

var Options Flags

// setBoolFlag registers a flag and returns the pointer to its value.
func setBoolFlag(name, short, usage string) *bool {
	b := flag.Bool(name, false, usage)
	flag.BoolVar(b, short, false, usage) // Bind the short flag to the same variable.
	return b
}

// ParseFlags sets up the flags and manually handles any combined short flags.
func ParseFlags() error {
	Options.all = setBoolFlag("all", "a", "display all entries without ignoring . entries")
	Options.grid = setBoolFlag("grid", "g", "display entries as a grid (default)")
	Options.long = setBoolFlag("long", "l", "display entries as a long listing format")

	if err := parseCombinedFlags(); err != nil {
		return err
	}
	flag.Parse()
	return nil
}

// parseCombinedFlags processes any combined short flags (like -lg or -ag).
func parseCombinedFlags() error {
	for i, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") && len(arg) > 2 {
			for _, ch := range arg[1:] {
				switch ch {
				case 'a':
					*Options.all = true
				case 'g':
					*Options.grid = true
				case 'l':
					*Options.long = true
				default:
					return errors.New("invalid flag: " + string(ch))
				}
			}
			// Remove processed combined flag from os.Args to avoid re-parsing
			os.Args = append(os.Args[:i+1], os.Args[i+2:]...)
		}
	}
	return nil
}

// Format returns a formatted string of the directory entries based on the selected flags.
func Format(entries []Entry) string {
	switch {
	case *Options.long && *Options.all:
		return Long(entries)
	case *Options.all:
		return Grid(entries)
	case *Options.long:
		return Long(entries)
	default:
		return Grid(entries)
	}
}

// DirPath returns the directory path from the command-line arguments or defaults to the current working directory.
func DirPath() (string, error) {
	if flag.NArg() > 0 {
		return flag.Arg(0), nil
	}
	return os.Getwd()
}
