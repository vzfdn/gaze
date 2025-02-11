package entry

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type row struct {
	permission string
	user       string
	group      string
	time       string
	size       string
	name       string
}

// RenderLong renders a detailed view of file entries, aligned in columns.
func RenderLong(entries []Entry, showHeader bool) string {
	if len(entries) == 0 {
		return "0 File, 0B\n"
	}

	rows := make([]row, len(entries))
	var maxUser, maxGroup, maxTime, maxSize int

	for i, e := range entries {
		user, group := e.UserAndGroup()
		timeStr := e.Time()
		sizeStr := HumanReadableSize(e.Size())

		rows[i] = row{
			permission: e.Permission(),
			user:       user,
			group:      group,
			time:       timeStr,
			size:       sizeStr,
			name:       e.Name(),
		}

		maxUser = max(maxUser, utf8.RuneCountInString(user))
		maxGroup = max(maxGroup, utf8.RuneCountInString(group))
		maxTime = max(maxTime, utf8.RuneCountInString(timeStr))
		maxSize = max(maxSize, utf8.RuneCountInString(sizeStr))
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d File, %s\n", len(entries), TotalSize(entries))
	
	// Write header row if showHeader is true
	if showHeader {
		fmt.Fprintf(&sb, "%s  %-*s  %-*s  %-*s  %-*s  %s\n",
			"Permission",
			maxUser, "User",
			maxGroup, "Group",
			maxTime, "Date Modified",
			maxSize, "Size",
			"Name",
		)
	}

	for _, r := range rows {
		fmt.Fprintf(&sb, "%s  %-*s  %-*s  %-*s  %-*s  %s\n",
			r.permission,
			maxUser, r.user,
			maxGroup, r.group,
			maxTime, r.time,
			maxSize, r.size,
			r.name,
		)
	}

	return sb.String()
}

// max returns the larger of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
