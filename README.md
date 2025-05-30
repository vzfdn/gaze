# gaze

**gaze** is a simple ls-like command-line tool written in Go, designed to list directory contents in a grid or detailed format. It works on both Unix and Windows systems.

## Install

Ensure you have Go installed on your system. The simplest way is to run:

```
go install github.com/vzfdn/gaze@latest
```

This downloads, builds, and installs gaze to `$GOPATH/bin` (or `$HOME/go/bin`). Ensure it's in your `$PATH` to run gaze anywhere.
After installation, you can run the tool from anywhere simply by typing:

```
gaze [flags] [path]
```

Alternatively, to build manually from source:

```
git clone https://github.com/vzfdn/gaze.git
cd gaze
go build
```

Then run it with:

```
./gaze [flags] [path]
```

## Example

![gaze terminal demo](./img/demo.png)

## Flags

- `-a, --all`: Show hidden entries (e.g., dot files on Unix)
- `-g, --grid`: Display entries in a grid layout (default)
- `-l, --long`: Use detailed listing format (permissions, owner, group, size, time, name)
- `-h, --header`: Include a header row in long format output
- `-R, --recursive`: Recursively list subdirectories
- `-T, --tree`: Recursively display directory contents as a tree-like format
- `-L, --dereference`: Show info for the target file, not the symlink
- `-F, --classify`: Append file type indicators (e.g., / for directories, \* for executables, @ for symlinks)
- `-s, --size`: Sort entries by file size (largest first)
- `-t, --time`: Sort entries by modification time (newest first)
- `-k, --kind`: Sort entries by file type (directories first, then files)
- `-x, --extension`: Sort entries by file extension (alphabetically)
- `-r, --reverse`: Reverse the order of sorting
- `-U, --no-sort`: Do not sort entries
- `--no-color`: Do not colorize output

## TODO
- [ ] **Performance**: Replace slice buffering with stream processing for entries  
- [ ] **Performance**: Refactor tree view rendering to reduce memory usage  
- [ ] **Feature**: Add `-m/--media` flag to show file metadata (e.g., media length)  