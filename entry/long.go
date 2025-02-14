package entry

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type row struct {
	perms string
	user  string
	group string
	mod   string
	size  string
	name  string
}

type widths struct {
	perms int
	user  int
	group int
	mod   int
	size  int
}

// RenderLong renders a detailed view of file entries, aligned in columns.
// Returns a formatted string representing the file entries.
func RenderLong(entries []Entry, cfg Config) string {
	if len(entries) == 0 {
		return "0 File, 0B\n"
	}

	rows, w := processEntries(entries)
	var sb strings.Builder

	// Write summary line
	fmt.Fprintf(&sb, "%d File, %s\n", len(entries), TotalSize(entries))

	// write header if requested
	if cfg.Header {
		writeHeader(&sb, w)
	}

	// write rows
	for _, row := range rows {
		writeRow(&sb, row, w)
	}

	return sb.String()
}

// processEntries processes the file entries and calculates the maximum column widths.
// Returns a slice of rows and a struct containing the calculated widths.
func processEntries(entries []Entry) ([]row, widths) {
	rows := make([]row, len(entries))
	w := widths{
		perms: utf8.RuneCountInString("Permissions"),
		user:  utf8.RuneCountInString("User"),
		group: utf8.RuneCountInString("Group"),
		mod:   utf8.RuneCountInString("Date Modified"),
		size:  utf8.RuneCountInString("Size"),
	}

	for i, e := range entries {
		u, g := e.UserAndGroup()
		r := row{
			perms: e.Permission(),
			user:  u,
			group: g,
			mod:   e.Time(),
			size:  HumanReadableSize(e.Size()),
			name:  e.Name(),
		}
		rows[i] = r

		w.perms = max(w.perms, utf8.RuneCountInString(r.perms))
		w.user = max(w.user, utf8.RuneCountInString(r.user))
		w.group = max(w.group, utf8.RuneCountInString(r.group))
		w.mod = max(w.mod, utf8.RuneCountInString(r.mod))
		w.size = max(w.size, utf8.RuneCountInString(r.size))
	}

	return rows, w
}

// max returns the larger of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// writeHeader writes the header row to the strings.Builder.
// It uses the provided widths to align the column names.
func writeHeader(sb *strings.Builder, w widths) {
	fmt.Fprintf(sb, " %-*s  %-*s %-*s  %-*s  %-*s  %s\n",
		w.perms, "Permissions",
		w.user, "User",
		w.group, "Group",
		w.mod, "Date Modified",
		w.size, "Size",
		"Name",
	)
}

// writeRow writes a single file entry row to the strings.Builder.
// It uses the provided widths to align the row data.
func writeRow(sb *strings.Builder, r row, w widths) {
	fmt.Fprintf(sb, " %-*s  %-*s %-*s  %-*s  %-*s  %s\n",
		w.perms, r.perms,
		w.user, r.user,
		w.group, r.group,
		w.mod, r.mod,
		w.size, r.size,
		r.name,
	)
}
