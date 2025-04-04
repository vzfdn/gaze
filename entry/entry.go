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

// NewEntry creates a file/directory entry with metadata.
func NewEntry(fi fs.FileInfo, name, path, target string) Entry {
	return Entry{
		info:   fi,
		name:   name,
		path:   path,
		target: target,
	}
}

// size returns the size of the Entry in bytes.
func (e Entry) size() int64 {
	return e.info.Size()
}

// permission returns the file permissions of the Entry as a string.
func (e Entry) permission() string {
	return e.info.Mode().String()
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

// userAndGroup returns the user and group names for the Entry's file info.
func (e Entry) userAndGroup() (string, string) {
	return userGroup(e)
}

// PrintEntries prints entries to stdout.
// It optionally recurses into subdirectories based on Config.Recurse.
func PrintEntries(path string, cfg Config) error {
	entries, err := ReadEntries(path, cfg)
	if err != nil {
		return err
	}
	if cfg.Tree {
		cfg.Recurse = false
		entries, err = addTreePrefixes(path, entries, cfg, "", 0)
		if err != nil {
			return fmt.Errorf("tree error: %w", err)
		}
	}
	output, err := render(entries, cfg)
	if err != nil {
		return fmt.Errorf("render error: %w", err)
	}
	fmt.Fprint(os.Stdout, output)
	if cfg.Recurse {
		for _, e := range entries {
			if e.info.IsDir() {
				subDir := filepath.Join(path, e.info.Name())
				if path == "." {
					subDir = "./" + e.info.Name()
				}
				fmt.Printf("\n%s:\n", subDir)
				if err := PrintEntries(subDir, cfg); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// ReadEntries lists entries in path, applying filters from Config.
// If path is a file, it returns a single-entry slice or nil if skipped.
func ReadEntries(path string, cfg Config) ([]Entry, error) {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !fileInfo.IsDir() {
		e, included, err := processEntry(path, fileInfo, cfg)
		if err != nil {
			return nil, err
		}
		if included {
			return []Entry{e}, nil
		}
		return nil, nil
	}
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	entries := make([]Entry, 0, len(dirEntries))
	for _, de := range dirEntries {
		fileInfo, err := de.Info()
		if err != nil {
			continue // Skip unreadable entries (e.g., permission denied).
		}
		e, included, err := processEntry(filepath.Join(path, fileInfo.Name()), fileInfo, cfg)
		if err != nil {
			continue // Skip problematic entries (e.g., broken symlinks).
		}
		if included {
			entries = append(entries, e)
		}
	}
	if len(entries) > 1 {
		sortEntries(entries, cfg)
	}
	return entries, nil
}

// processEntry creates an Entry while applying filters for hidden files and handling symlinks.
// Returns the Entry and true if it should be included, false if skipped.
func processEntry(fullPath string, fileInfo fs.FileInfo, cfg Config) (Entry, bool, error) {
	// Skip hidden files unless config.All is true.
	if !cfg.All && isHidden(fileInfo) {
		return Entry{}, false, nil
	}
	e := NewEntry(fileInfo, formatName(fileInfo, cfg), filepath.Dir(fullPath), "")
	// Handle symlinks
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		linkTarget, err := os.Readlink(fullPath)
		if err != nil {
			return Entry{}, false, nil
		}
		e.target = linkTarget
		if cfg.Dereference {
			// Replace entry info with dereferenced target info if available
			if targetInfo, err := os.Stat(fullPath); err == nil {
				e.info = targetInfo
				e.target = ""
			}
		}
	}
	return e, true, nil
}

// isHidden reports whether a file is hidden by its name starting with a dot.
func isHidden(fileInfo fs.FileInfo) bool {
	return fileInfo.Name()[0] == '.'
}

// formatName formats the file name based on the file type and configuration.
// It quotes the name if it contains special characters or whitespace,
// and appends a classification symbol ("/", "*", or "@") if cfg.Classify is enabled.
func formatName(info fs.FileInfo, cfg Config) string {
	name := info.Name()
	if strings.ContainsAny(name, " \t\n\v\f\r!@#$%^&*()[]{}<>?/|\\~`") {
		name = fmt.Sprintf("'%s'", name)
	}
	if cfg.Classify {
		switch {
		case info.IsDir():
			name += "/"
		case info.Mode()&os.ModeSymlink != 0:
			name += "@"
		case info.Mode()&0o111 != 0:
			name += "*"
		}
	}
	return name
}

// addTreePrefixes adds tree-like prefixes to directory entries and collects subdirectory entries.
// If a single entry is not a directory, it is returned without a tree structure.
func addTreePrefixes(path string, entries []Entry, cfg Config, prefix string, depth int) ([]Entry, error) {
	var result []Entry
	if depth == 0 {
		fi, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		// If `path` is a file, return it directly (no tree formatting)
		if !fi.IsDir() {
			return entries, nil
		}
		// If `path` is a directory, include it as the root
		result = append(result, NewEntry(fi, formatName(fi, cfg), path, ""))
	}
	for i, e := range entries {
		isLast := i == len(entries)-1
		connector := "├── "
		if isLast {
			connector = "└── "
		}
		e.name = prefix + connector + e.name
		result = append(result, e)
		if e.info.IsDir() {
			subPath := filepath.Join(e.path, e.info.Name())
			subEntries, err := ReadEntries(subPath, cfg)
			if err != nil {
				return nil, err
			}
			subPrefix := prefix + "│   "
			if isLast {
				subPrefix = prefix + "    "
			}
			subTree, err := addTreePrefixes(subPath, subEntries, cfg, subPrefix, depth+1)
			if err != nil {
				return nil, err
			}
			result = append(result, subTree...)
		}
	}
	return result, nil
}

// render generates output based on entries and configuration.
// It uses long format if -l is set, otherwise defaults to grid.
func render(entries []Entry, cfg Config) (string, error) {
	if cfg.Long {
		return renderLong(entries, cfg), nil
	}
	if cfg.Tree {
		return renderTree(entries), nil
	}
	return renderGrid(entries)
}

// renderTree renders the entries in a tree-like format.
func renderTree(entries []Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "%s\n", e.name)
	}
	return sb.String()
}
