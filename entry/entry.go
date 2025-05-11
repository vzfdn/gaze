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

	symbolDir  = "/"
	symbolLink = "@"
	symbolExec = "*"

	execBits     = 0o111 // Executable permission bits
	specialChars = " \t\n\v\f\r!@#$%^&*()[]{}<>?/|\\~`"
)

var (
	cfg   Config
	color = newColorizer()
)

// Entry represents a file or directory with its metadata.
// It embeds os.FileInfo to provide direct access to file information methods (Size, ModTime, IsDir, etc).
type Entry struct {
	os.FileInfo
	path       string
	target     string
	treePrefix string
}

// isHidden reports whether a file is hidden by its name starting with a dot.
// Returns false for empty filenames to avoid potential panics.
func (e Entry) isHidden() bool {
	name := e.Name()
	return len(name) > 0 && name[0] == '.'
}

// DisplayName formats the basename of an entry and returns the formatted string.
func (e Entry) DisplayName() string {
	name := e.Name()
	// Quote if needed
	if strings.ContainsAny(name, specialChars) {
		name = "'" + name + "'"
	}
	colored := color.colorize(e.FileInfo, name)
	if cfg.Classify {
		_, symbol := fileType(e.FileInfo)
		colored += symbol
	}
	if e.treePrefix != "" {
		colored = e.treePrefix + colored
	}
	return colored
}

// PrintEntries prints entries to stdout and, if cfg.Recurse is true,
// recurses into subdirectories.
func PrintEntries(path string) error {
	entries, err := readEntries(path)
	if err != nil {
		return err
	}

	if cfg.Tree {
		cfg.Recurse = false
		entries, err = addTreePrefixes(path, entries, "", 0)
		if err != nil {
			return fmt.Errorf("tree error: %w", err)
		}
	}

	output, err := render(entries)
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
				if err := PrintEntries(subDir); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// readEntries lists entries in path, applying filters from Config.
// If path is a file, it returns a single-entry slice or nil if skipped.
func readEntries(path string) ([]Entry, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}

	// Handle single file case
	if !fi.IsDir() {
		if e, included, err := processEntry(path, fi); err != nil {
			return nil, err
		} else if included {
			return []Entry{e}, nil
		}
		return nil, nil
	}

	// Process directory case
	dirEntries, err := readDir(path)
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(dirEntries))
	for _, de := range dirEntries {
		fi, err := de.Info()
		if err != nil {
			continue // Skip unreadable entries (e.g., permission denied).
		}
		if e, included, err := processEntry(filepath.Join(path, fi.Name()), fi); err != nil {
			continue // Skip problematic entries (e.g., broken symlinks).
		} else if included {
			entries = append(entries, e)
		}
	}

	if len(entries) > 1 {
		sortEntries(entries)
	}
	return entries, nil
}

// processEntry creates an Entry, filters hidden files, and resolves symlinks.
// Returns the Entry, a boolean indicating inclusion, and any error.
func processEntry(path string, fi os.FileInfo) (Entry, bool, error) {
	e := Entry{
		FileInfo: fi,
		path:     path,
	}
	// Skip hidden files unless cfg.All is true.
	if !cfg.All && e.isHidden() {
		return Entry{}, false, nil
	}
	// Handle symlinks
	if fi.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(path)
		if err != nil {
			return Entry{}, false, fmt.Errorf("read symlink %s: %w", path, err)
		}
		e.target = target
		if cfg.Dereference {
			// Replace entry info with dereferenced target info if available
			if targetInfo, err := os.Stat(path); err == nil {
				e.FileInfo = targetInfo
				e.target = ""
			}
		}
	}
	return e, true, nil
}

// readDir reads directory entries, avoiding the extra sort done by os.ReadDir.
func readDir(path string) ([]os.DirEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.ReadDir(-1)
}

// fileType classifies a file and returns its type identifier and display symbol.
// Returns two strings: (typeIdentifier, displaySymbol)
func fileType(fi os.FileInfo) (string, string) {
	switch {
	case fi.IsDir():
		return typeDir, symbolDir
	case fi.Mode()&os.ModeSymlink != 0:
		return typeLink, symbolLink
	case fi.Mode().Perm()&execBits != 0:
		return typeExec, symbolExec
	default:
		return typeFile, ""
	}
}

// render generates a string representation of the entries based on configuration settings,
func render(entries []Entry) (string, error) {
	if cfg.Long {
		return renderLong(entries), nil
	}
	if cfg.Tree {
		return renderTree(entries), nil
	}
	return renderGrid(entries)
}
