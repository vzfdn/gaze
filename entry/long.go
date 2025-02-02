package entry

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Long formats and returns a detailed table-like string representation of the given entries.
func Long(entries []Entry) string {
	var sb strings.Builder

	// Calculate maximum widths for each column
	var maxUserLen, maxGroupLen, maxSizeLen, maxTimeLen int
	for _, e := range entries {
		user, group := e.UserAndGroup()
		sizeStr := HumanReadableSize(e.Size()) 
		timeStr := e.Time()                     

		if utf8.RuneCountInString(user) > maxUserLen {
			maxUserLen = utf8.RuneCountInString(user)
		}
		if utf8.RuneCountInString(group) > maxGroupLen {
			maxGroupLen = utf8.RuneCountInString(group)
		}
		if utf8.RuneCountInString(sizeStr) > maxSizeLen {
			maxSizeLen = utf8.RuneCountInString(sizeStr)
		}
		if utf8.RuneCountInString(timeStr) > maxTimeLen {
			maxTimeLen = utf8.RuneCountInString(timeStr)
		}
	}

	// Header line (total files and size)
	line := fmt.Sprintf("%d File, %s\n", len(entries), TotalSize(entries))
	sb.WriteString(line)

	for _, e := range entries {
		user, group := e.UserAndGroup()
		line = fmt.Sprintf("%s  %-*s  %-*s  %-*s  %-*s  %s\n",
			e.Permission(),
			maxUserLen, user,     
			maxGroupLen, group,   
			maxTimeLen, e.Time(),  
			maxSizeLen, HumanReadableSize(e.Size()),  
			e.Name(),
		)
		sb.WriteString(line)
	}

	return sb.String()
}
