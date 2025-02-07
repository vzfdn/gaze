package entry

import (
	"io/fs"
	"os"
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

// User returns the user and group name associated with the fileInfo inside the Entry.
func (e Entry) UserAndGroup() (string, string) {
	return FileUserGroup(e)
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

// Name returns the file name of the Entry, quoted if it contains special characters or whitespace.
func (e Entry) Name() string {
	name := e.info.Name()
	if strings.ContainsAny(name, " \t\n\v\f\r") ||
		strings.ContainsAny(name, "!@#$%^&*()[]{}<>?/|\\~`") {
		return "'" + name + "'"
	}
	return name
}

// ReadEntries reads and returns a []Entry from the specified path.
func ReadEntries(path string, showHidden bool) ([]Entry, error) {
    dirEntries, err := os.ReadDir(path)
    if err != nil {
        return nil, err
    }
	
    entries := make([]Entry, 0, len(dirEntries))
    for _, de := range dirEntries {
        info, err := de.Info()
        if err != nil {
            return nil, err
        }
        
        name := info.Name()
        if showHidden || name[0] != '.' {
            entries = append(entries, Entry{
                info: info,
                path: path,
            })
        }
    }
    return entries, nil
}

/*type VideoFile struct {
	Entry
	Duration   string
	Resolution string
}

func FormatVideoEntries(s []VideoFile) string {
	pad := 8
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%*sType%*sSize%*sDuration%*sFile name%*s \n",
		pad/4, "", pad+4, "", pad+1, "", pad, "", pad, ""))

	for _, v := range s {
		sb.WriteString(fmt.Sprintf("%v%*s%v%*s%v%*s%v\n",
			v.User, pad, "",
			FormatSize(v.Size), pad-1, "",
			v.Duration, pad, "",
			v.Name))
	}
	return sb.String()
}

type AudioFile struct {
	Entry
	Duration string
	Bitrate  string
}
*/
