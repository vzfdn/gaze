package main

import (
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

	err = entry.PrintEntries(path, cfg)
	if err != nil {
		log.Fatalf("Error printing entries: %v", err)
	}
}

// TODO renderLong: total, newline problem
// TODO improve godocs, error messages
// TODO symlinks
// TODO adding flags: -m --media, -s --sort  
// TODO output colorization
