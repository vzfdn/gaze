
gaze
----

gaze is a simple ls-like command-line tool written in Go, designed to list
directory contents in a grid or detailed format. It works on both Unix and
Windows systems.

Usage
-----

gaze [flags] [path]

Examples:
```*
$ gaze -lah media/
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
Flags
-----

- -a, --all: show hidden entries (e.g., dot files on Unix)
- -g, --grid: display entries in a grid layout (default)
- -l, --long: use detailed listing format (permissions, owner, group, size,
  time, name)
- -h, --header: include a header row in long format output
- -R, --recursive: recursively list subdirectories

TODO
----
- Symlinks: support displaying symbolic links
- Sorting: add -s/--sort flag for sorting by name, size, or time
- Media Metadata: add -m/--media flag to show metadata of media files (e.g.,
  length)
- Colorization: enable colored output
