package entry

import (
	"fmt"
	"strings"
)

func Long(displayEntries []Entry) string {
	var sb strings.Builder
	line := fmt.Sprintf("%d File,%s\n", len(displayEntries), TotalSize(displayEntries))
	sb.WriteString(line)
	const pad = 2
	for _, e := range displayEntries {
		line = fmt.Sprintf("%s%*s%s%*s%s%*s%s%*s%s%*s%s\n",
			e.Permission(), pad, "",
			e.User(), pad, "",
			e.Group(), pad, "",
			e.Time(), pad, "",
			HumanReadableSize((e.Size())), pad, "",
			e.Name(),
		)
		sb.WriteString(line)
	}
	return sb.String()
}
