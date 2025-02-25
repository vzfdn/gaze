# gaze

gaze is a simple ls-like command-line tool written in go, designed to list directory contents in a grid or detailed format. it works on both unix and windows systems.

## Install

Ensure you have go installed on your system. the simplest way is to run:
`go install github.com/vzfdn/gaze@latest`
this downloads, builds, and installs gaze to `$GOPATH/bin` (or `$HOME/go/bin`). ensure it’s in your `$PATH` to run gaze anywhere.

Alternatively, if you want to build it manually from source:
`git clone https://github.com/vzfdn/gaze.git`
`cd gaze`
`go build -o gaze`
after building manually, you can run it with:
`./gaze [flags] [path]`.

## Examples

```*
$ ./gaze -lah media/
  8 Files, 16.2M
  Permissions User Group Modified        Size Name
  -rw-rw-rw-  Rwin None  Oct 30  2024      40  .file
  -rw-rw-rw-  Rwin None  Aug 08  2024  144.8K 'Baldur`s_Gate_3.webp'
  -rw-rw-rw-  Rwin None  Apr 10  2024  121.9K  EldenRing.jpg
  -rw-rw-rw-  Rwin None  Nov 21  2023   12.5M 'Evanescence - Tourniquet.mp3'
  -rw-rw-rw-  Rwin None  Sep 10  2023    2.6M  aesthetic.mp4
  -rw-rw-rw-  Rwin None  Dec 15  2022  912.2K  cpumemory.pdf
  drwxrwxrwx  Rwin None  Feb 24 18:07    4.0K  notes
  -rw-rw-rw-  Rwin None  Feb 24 18:21    1.1K  test.txt
```
## Flags

- `-a, --all`: show hidden entries (e.g., dot files on unix)
- `-g, --grid`: display entries in a grid layout (default)
- `-l, --long`: use detailed listing format (permissions, owner, group, size, time, name)
- `-h, --header`: include a header row in long format output
- `-r, --recursive`: recursively list subdirectories

## Todo

- symlinks: support displaying symbolic links
- sorting: add `-s/--sort` flag for sorting by name, size, or time
- media metadata: add `-m/--media` flag to show metadata of media files (e.g., length)
- colorization: enable colored output
