package main

import (
	"fmt"
	"os"

	"github.com/vzfdn/gaze/entry"
)

func main() {
	err := entry.ParseFlags()
	if err != nil {
		fmt.Printf("cannot parse flags: '%v'\n", err)
		os.Exit(1)
	}

	path, err := entry.DirPath()
	if err != nil {
		fmt.Printf("cannot open directory: '%v'\n", err)
		os.Exit(1)
	}

	entries, err := entry.ReadEntries(path)
	if err != nil {
		fmt.Printf("cannot read entries: '%v'\n", err)
		os.Exit(1)
	}

	output := entry.Format(entries)
	if _, err = fmt.Fprint(os.Stdout, output); err != nil {
		fmt.Printf("cannot write to stdout: '%v'\n", err)
		os.Exit(1)
	}
}

// TODO combine functionality
// TODO listing files
// TODO adding flags: -m --media  -h --header -s --sort
// TODO error handling
// TODO documentation
// TODO output colorization
// TODO windows support
