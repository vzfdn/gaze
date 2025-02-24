
# gaze

gaze is a simple `ls`-like command-line tool written in Go, designed to list directory contents in either a grid or detailed format. It runs seamlessly on both Unix and Windows systems.

## Usage

gaze [flags] [path]

**Example:** `gaze -lhR /path/to/dir`

## Flags

- `-a`, `--all`: show hidden entries (e.g., dot files on Unix).
- `-g`, `--grid`: display entries in a grid layout (default).
- `-l`, `--long`: use detailed listing format (permissions, owner, group, size, time, name).
- `-h`, `--header`: include a header row in long format output.
- `-R`, `--recursive`: recursively list subdirectories.

## TODO

- **Symlinks**: support displaying symbolic links.
- **Sorting**: add `-s`/`--sort` flag for sorting by name, size, or time.
- **Media Metadata**: add `-m`/`--media` flag to show metadata of media files (e.g., length).
- **Colorization**: enable colored output.
