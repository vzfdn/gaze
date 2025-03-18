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
	name   string
	path   string
	target string
}

// permission returns the file permissions of the Entry as a string.
func (e Entry) permission() string {
	return e.info.Mode().String()
}

// userAndGroup returns the user and group names for the Entry's file info.
func (e Entry) userAndGroup() (string, string) {
	return userGroup(e)
}

// time returns the formatted modification time of the Entry.
// Uses "Jan 02 15:04" for current-year entries, "Jan 02  2006" otherwise.
func (e Entry) time() string {
	mt := e.info.ModTime()
	if mt.Year() == time.Now().Year() {
		return mt.Format("Jan 02 15:04")
	}
	return mt.Format("Jan 02  2006")
}

// size returns the size of the Entry in bytes.
func (e Entry) size() int64 {
	return e.info.Size()
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

// ReadEntries reads directory entries from path, applying Config rules.
// Returns a slice of entries or an error if reading fails.
func ReadEntries(path string, cfg Config) ([]Entry, error) {
	dirEntry, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("cannot access %s: %w", path, err)
	}
	entries := make([]Entry, 0, len(dirEntry))
	for _, de := range dirEntry {
		info, err := de.Info()
		if err != nil {
			return nil, fmt.Errorf("cannot access %s: %w", filepath.Join(path, de.Name()), err)
		}
		// Skip hidden files unless cfg.All is true.
		if !cfg.All && strings.HasPrefix(info.Name(), ".") {
			continue
		}
		e := Entry{
			info: info,
			name: formatName(info, cfg),
			path: path,
		}
		// Process symbolic links.
		if info.Mode()&os.ModeSymlink != 0 {
			linkPath := filepath.Join(path, info.Name())
			target, err := os.Readlink(linkPath)
			if err != nil {
				continue
			}
			e.target = target
			if cfg.Dereference {
				// Override with dereferenced info if available.
				if targetInfo, err := os.Stat(linkPath); err == nil {
					e.info = targetInfo
					e.target = "" // Clear target when successfully dereferenced.
				}
			}
		}
		entries = append(entries, e)
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

// formatName formats the file name based on the file type and configuration.
// It quotes the name if it contains special characters or whitespace,
// and appends a classification symbol ("/", "*", or "@") if cfg.Classify is enabled.
func formatName(info fs.FileInfo, cfg Config) string {
	name := info.Name()
	var indicator string
	switch {
	case info.IsDir():
		indicator = "/"
	case info.Mode()&os.ModeSymlink != 0:
		indicator = "@"
	case info.Mode()&0o111 != 0:
		indicator = "*"
	}
	const specialChars = " \t\n\v\f\r!@#$%^&*()[]{}<>?/|\\~`"
	if strings.ContainsAny(name, specialChars) {
		name = "'" + name + "'"
	}
	if cfg.Classify {
		name += indicator
	}
	return name
}
