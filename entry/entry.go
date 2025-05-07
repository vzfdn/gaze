// Package entry provides functionality for listing directory entries
// with customizable formatting and cross-platform support.
package entry

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	typeDir  = "di"
	typeLink = "ln"
	typeExec = "ex"
	typeFile = "fi"

	symDir  = "/"
	symLink = "@"
	symExec = "*"

	execBits     = 0o111 // Executable permission bits
	specialChars = " \t\n\v\f\r!@#$%^&*()[]{}<>?/|\\~`"
)

// Entry represents a file or directory with its metadata.
// It embeds os.FileInfo to provide direct access to file information methods (Size, ModTime, IsDir, etc).
type Entry struct {
	os.FileInfo
	displayName string
	path        string
	target      string
}

// NewEntry creates a file/directory entry with metadata.
func NewEntry(fi os.FileInfo, name, path, target string) Entry {
	return Entry{
		FileInfo:    fi,
		displayName: name,
		path:        path,
		target:      target,
	}
}

// PrintEntries prints entries to stdout.
// It optionally recurses into subdirectories based on Config.Recurse.
func PrintEntries(path string, cfg Config) error {
	c := newColorizer()
	entries, err := ReadEntries(path, cfg, c)
	if err != nil {
		return err
	}

	if cfg.Tree {
		cfg.Recurse = false
		entries, err = addTreePrefixes(path, entries, cfg, "", 0, c)
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
			if e.IsDir() {
				subDir := filepath.Join(path, e.Name())
				if path == "." {
					subDir = "./" + e.Name()
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
func ReadEntries(path string, cfg Config, c colorizer) ([]Entry, error) {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}

	if !fileInfo.IsDir() {
		e, included, err := processEntry(path, fileInfo, cfg, c)
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
		e, included, err := processEntry(filepath.Join(path, fileInfo.Name()), fileInfo, cfg, c)
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
func processEntry(fullPath string, fileInfo os.FileInfo, cfg Config, c colorizer) (Entry, bool, error) {
	// Skip hidden files unless config.All is true.
	if !cfg.All && isHidden(fileInfo) {
		return Entry{}, false, nil
	}
	e := NewEntry(fileInfo, getDisplayName(fileInfo, cfg, c), filepath.Dir(fullPath), "")

	// Handle symlinks
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		linkTarget, err := os.Readlink(fullPath)
		if err != nil {
			return Entry{}, false, fmt.Errorf("reading symlink %s: %w", fullPath, err)
		}
		e.target = linkTarget
		if cfg.Dereference {
			// Replace entry info with dereferenced target info if available
			if targetInfo, err := os.Stat(fullPath); err == nil {
				e.FileInfo = targetInfo
				e.target = ""
			}
		}
	}
	return e, true, nil
}

// isHidden reports whether a file is hidden by its name starting with a dot.
// Returns false for empty filenames to avoid potential panics.
func isHidden(fileInfo os.FileInfo) bool {
	name := fileInfo.Name()
	return len(name) > 0 && name[0] == '.'
}

// fileType returns a short type code for the given file (e.g. "di", "ln").
func fileType(info os.FileInfo) string {
	mode := info.Mode()
	switch {
	case mode.IsDir():
		return typeDir
	case mode&os.ModeSymlink != 0:
		return typeLink
	case mode.Perm()&execBits != 0:
		return typeExec
	default:
		return typeFile
	}
}

// quoteSpecialChars checks if the name contains special characters and quotes it.
func quoteSpecialChars(name string) string {
	if strings.ContainsAny(name, specialChars) {
		return "'" + name + "'"
	}
	return name
}

// getDisplayName returns the formatted file name, quoting special characters
// and appending a classification symbol if enabled.
func getDisplayName(info os.FileInfo, cfg Config, c colorizer) string {
	name := quoteSpecialChars(info.Name())
	infoType := fileType(info)
	colored := c.colorize(infoType, name)
	if cfg.Classify {
		switch infoType {
		case typeDir:
			return colored + symDir
		case typeLink:
			if cfg.Grid {
				return colored + symLink
			}
		case typeExec:
			return colored + symExec
		}
	}
	return colored
}

// render generates output based on entries and configuration.
func render(entries []Entry, cfg Config) (string, error) {
	if cfg.Long {
		return renderLong(entries, cfg), nil
	}
	if cfg.Tree {
		return renderTree(entries), nil
	}
	return renderGrid(entries)
}
