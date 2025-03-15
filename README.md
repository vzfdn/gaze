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

<pre>
> gaze -lash media/
8 Files, 16.1M
Permissions User Group Modified       Size Name
-rwxrwxrwx  rwin rwin  Nov 21  2023  12.5M 'Evanescence - Tourniquet.mp3'
-rwxrwxrwx  rwin rwin  Sep 10  2023   2.6M aesthetic.mp4
-rwxrwxrwx  rwin rwin  Dec 15  2022 912.2K cpumemory.pdf
-rwxrwxrwx  rwin rwin  Aug 08  2024 144.8K 'Baldur`s_Gate_3.webp'
drwxrwxrwx  rwin rwin  Feb 24 18:07   4.0K notes
-rwxrwxrwx  rwin rwin  Feb 24 18:21   1.1K test.txt
-rwxrwxrwx  rwin rwin  Oct 30  2024     40 .file
Lrwxrwxrwx  rwin rwin  Mar 04 15:52     13 'link to pdf' -> cpumemory.pdf
</pre>

## Flags

- `-a, --all`: Show hidden entries (e.g., dot files on Unix)
- `-g, --grid`: Display entries in a grid layout (default)
- `-l, --long`: Use detailed listing format (permissions, owner, group, size, time, name)
- `-h, --header`: Include a header row in long format output
- `-R, --recursive`: Recursively list subdirectories
- `-T, --tree`: Recursively display directory contents as a tree-like format
- `-L, --dereference`: Show info for the target file, not the symlink
- `-F, --classify`: Append file type indicators (e.g., / for directories, * for executables, @ for symlinks)
- `-s, --size`: Sort entries by file size (largest first)
- `-t, --time`: Sort entries by modification time (newest first)
- `-k, --kind`: Sort entries by file type (directories first, then files)
- `-x, --extension`: Sort entries by file extension (alphabetically)
- `-r, --reverse`: Reverse the order of sorting

## Todo

- Tree: Add `-T --tree` flag to Tree-like recursive view
- Media metadata: Add `-m/--media` flag to show metadata of media files (e.g., length)
- Colorization: Enable colored output
