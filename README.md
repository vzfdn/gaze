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
go build -o gaze
```

Then run it with:

```
./gaze [flags] [path]
```

## Example

<pre>
> gaze -lah media/
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
</pre>

## Flags

- `-a, --all`: Show hidden entries (e.g., dot files on Unix)
- `-g, --grid`: Display entries in a grid layout (default)
- `-l, --long`: Use detailed listing format (permissions, owner, group, size, time, name)
- `-h, --header`: Include a header row in long format output
- `-r, --recursive`: Recursively list subdirectories

## Todo

- Symlinks: Support displaying symbolic links
- Sorting: Add `-s/--sort` flag for sorting by name, size, or time
- Media metadata: Add `-m/--media` flag to show metadata of media files (e.g., length)
- Colorization: Enable colored output
