package main

import (
	"fmt"
	"log"

	"github.com/vzfdn/gaze/entry"
)

func main() {
	cfg, f, err := entry.ParseConfig()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	path, err := entry.ResolvePath(f)
	if err != nil {
		log.Fatalf("Error resolving path: %v", err)
	}

	entries, err := entry.ReadEntries(path, cfg.All)
	if err != nil {
		log.Fatalf("Error reading entries: %v", err)
	}

	output, err := entry.Format(entries, cfg)
	if err != nil {
		log.Fatalf("Error formatting output: %v", err)
	}

	fmt.Print(output)
}

// TODO fix long position
// TODO improve error messages
// TODO adding flags: -m --media -s --sort
// TODO output colorization
