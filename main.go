package main

import (
	"fmt"
	"os"

	"github.com/vzfdn/gaze/entry"
)

func main() {
	cfg, f, err := entry.ParseConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	path, err := entry.ResolvePath(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	err = entry.PrintEntries(path, cfg)
	if err != nil {
		os.Exit(2)
	}
	os.Exit(0)
}

// TODO symlinks
// TODO adding flags: -m --media, -s --sort
// TODO output colorization
