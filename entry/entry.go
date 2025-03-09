// Package entry provides functionality for listing directory entries
// with customizable formatting and cross-platform support.
package entry

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Entry represents a file or directory with its metadata.
type Entry struct {
	info   fs.FileInfo
	path   string
	target string
}

// Permission returns the file permissions of the Entry as a string.
func (e Entry) Permission() string {
	return e.info.Mode().String()
}

// UserAndGroup returns the user and group names for the Entry's file info.
func (e Entry) UserAndGroup() (string, string) {
	return userGroup(e)
}

// Time returns the formatted modification time of the Entry.
// Uses "Jan 02 15:04" for current-year entries, "Jan 02  2006" otherwise.
func (e Entry) Time() string {
	mt := e.info.ModTime()
	if mt.Year() == time.Now().Year() {
		return mt.Format("Jan 02 15:04")
	}
	return mt.Format("Jan 02  2006")
}

// Size returns the size of the Entry in bytes.
func (e Entry) Size() int64 {
	return e.info.Size()
}

// Name returns the file name of the Entry,
// quoted if it has special characters or whitespace.
func (e Entry) Name() string {
	name := e.info.Name()
	if strings.ContainsAny(name, " \t\n\v\f\r") ||
		strings.ContainsAny(name, "!@#$%^&*()[]{}<>?/|\\~`") {
		return "'" + name + "'"
	}
	return name
}

// ReadEntries reads directory entries from path, applying Config
// rules for symlinks and hidden files. Returns entries or error.
func ReadEntries(path string, cfg Config) ([]Entry, error) {
	dirs, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("cannot access %s: %w", path, err)
	}

	entries := make([]Entry, 0, len(dirs))
	for i := range dirs {
		info, err := dirs[i].Info()
		if err != nil {
			return nil, fmt.Errorf("cannot access %s/%s: %w", path, dirs[i].Name(), err)
		}

		if cfg.All || info.Name()[0] != '.' {
			e := Entry{info: info, path: path}

			// Checks if the entry is a symlink
			if info.Mode()&os.ModeSymlink != 0 {
				targetPath := filepath.Join(path, info.Name())
				target, _ := os.Readlink(targetPath)

				if cfg.Dereference {
					// Dereference symlink and use target's info.
					if targetInfo, err := os.Stat(targetPath); err == nil {
						e.info = targetInfo
					} else {
						e.target = target // Mark as broken if target is invalid.
					}
				} else {
					e.target = target // Store symlink target without dereferencing.
				}
			}

			entries = append(entries, e)
		}
	}
	return entries, nil
}

// formatEntries generates output based on entries and configuration.
// It uses long format if -l is set, otherwise defaults to grid.
func formatEntries(entries []Entry, cfg Config) (string, error) {
	switch {
	case cfg.Size:
		sortBySize(entries)
	case cfg.Time:
		sortByTime(entries)
	case cfg.Kind:
		sortByKind(entries)
	case cfg.Ext:
		sortByExt(entries)
	}

	if cfg.Reverse {
		reverse(entries)
	}

	if cfg.Long && cfg.Grid {
		fmt.Fprintf(os.Stderr, "warning: -l and -g are mutually exclusive, using long format\n")
	}

	if cfg.Long {
		return renderLong(entries, cfg), nil
	}
	return renderGrid(entries)
}

// PrintEntries prints entries to stdout.
// It optionally recurses into subdirectories based on Config.Recurse.
func PrintEntries(path string, cfg Config) error {
	entries, err := ReadEntries(path, cfg)
	if err != nil {
		return err
	}

	if cfg.Recurse {
		fmt.Printf("%s:\n", path)
	}

	output, err := formatEntries(entries, cfg)
	if err != nil {
		return fmt.Errorf("%s: format error: %w", path, err)
	}
	fmt.Println(output)

	if cfg.Recurse {
		for i := range entries {
			if entries[i].info.IsDir() {
				subDir := filepath.Join(path, entries[i].info.Name())
				if err := PrintEntries(subDir, cfg); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
