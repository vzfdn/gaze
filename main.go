package main

import (
	"fmt"
	"os"

	"github.com/vzfdn/gaze/entry"
)

func main() {
	cfg, f, err := entry.ParseConfig()
	if err != nil {
		fmt.Printf("Error parsing flags: %v", err)
		os.Exit(1)
	}

	path, err := entry.ResolvePath(f)
	if err != nil {
		fmt.Printf("Error resolving path: %v", err)
		os.Exit(1)
	}

	err = entry.PrintEntries(path, cfg)
	if err != nil {
		fmt.Printf("Error printing entries: %v", err)
		os.Exit(1)
	}
}

// TODO refactor renderGrid
// TODO improve godocs
// TODO symlinks
// TODO adding flags: -m --media, -s --sort
// TODO output colorization
