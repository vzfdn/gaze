package main

import (
	"fmt"
	"log"

	"github.com/vzfdn/gaze/entry"
)

func main() {
	flgs, err := entry.ParseFlags()
	if err != nil {
		log.Fatalf("Failed to parse flags: %v", err)
	}

	path, err := entry.ResolvePath()
	if err != nil {
		log.Fatalf("Failed to resolve path: %v", err)
	}

	entries, err := entry.ReadEntries(path, flgs.All)
	if err != nil {
		log.Fatalf("Failed to read entries: %v", err)
	}

	output, err := entry.Format(entries, flgs)
	if err != nil {
		log.Fatalf("Failed to format output: %v", err)
	}
	if _, err := fmt.Print(output); err != nil {
		log.Fatalf("Failed to write output: %v", err)
	}
}

// TODO adding flags: -m --media  -h --header -s --sort
// TODO output colorization
