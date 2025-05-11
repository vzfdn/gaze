package main

import (
	"fmt"
	"os"

	"github.com/vzfdn/gaze/entry"
)

func main() {
	flags, err := entry.ParseFlags()
	if err != nil {
		os.Exit(1)
	}
	path, err := entry.ResolvePath(flags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	err = entry.PrintEntries(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	os.Exit(0)
}
