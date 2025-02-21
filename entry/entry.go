package entry

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Entry struct {
	info fs.FileInfo
	path string
	// 	basename string
	// ext      string
}

// Permission returns the file permissions of the Entry as a string.
func (e Entry) Permission() string {
	return e.info.Mode().String()
}

// UserAndGroup returns the user and group name associated with the fileInfo inside the Entry.
func (e Entry) UserAndGroup() (string, string) {
	return fileUserGroup(e)
}

// Time returns the formatted modification time of the Entry.
// It uses "Jan 02 15:04" format for entries modified in the current year,
// or "Jan 02  2006" format for entries modified in other years.
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
// quoted if it contains special characters or whitespace.
func (e Entry) Name() string {
	name := e.info.Name()
	if strings.ContainsAny(name, " \t\n\v\f\r") ||
		strings.ContainsAny(name, "!@#$%^&*()[]{}<>?/|\\~`") {
		return "'" + name + "'"
	}
	return name
}

// ReadEntries reads and returns a []Entry from the specified path.
func ReadEntries(path string, cfg Config) ([]Entry, error) {
	dirs, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("cannot access %s: %w", path, err)
	}

	entries := make([]Entry, 0, len(dirs))
	for i := range dirs {
		info, err := dirs[i].Info()
		if err != nil {
			return nil, fmt.Errorf("cannot stat %s/%s: %w", path, dirs[i].Name(), err)
		}

		if cfg.All || info.Name()[0] != '.' {
			entries = append(entries, Entry{
				info: info,
				path: path,
			})
		}
	}

	return entries, nil
}

// formatEntries generates output based on entries and configuration.
// Uses long format if -l is set, otherwise defaults to grid format.
func formatEntries(entries []Entry, cfg Config) (string, error) {
	if cfg.Long && cfg.Grid {
		fmt.Fprintf(os.Stderr, "warning: -l and -g are mutually exclusive, using long format\n")
		return renderLong(entries, cfg), nil
	}

	if cfg.Long {
		return renderLong(entries, cfg), nil
	}

	return renderGrid(entries)
}

// PrintEntries prints entries (and recurses if needed) to stdout.
func PrintEntries(path string, cfg Config) error {
	entries, err := ReadEntries(path, cfg)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
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
					return fmt.Errorf("%s: %w", subDir, err)
				}
			}
		}
	}

	return nil
}
