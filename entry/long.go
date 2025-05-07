package entry

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

// row represents a single file entry for long format rendering.
type row struct {
	perms   string
	user    string
	group   string
	modTime string
	size    string
	name    string
	target  string
}

// widths holds the maximum column widths for long format rendering.
type widths struct {
	perms, user, group, mod, size int
}

// renderLong renders a detailed view of file entries, aligned in columns.
func renderLong(entries []Entry, cfg Config) string {
	if len(entries) == 0 {
		return ""
	}

	rows, w := processEntries(entries, cfg)
	var sb strings.Builder
	sb.Grow(len(entries) * (w.perms + w.user + w.group + w.mod + w.size + 20))

	files := "Files"
	if len(entries) == 1 {
		files = "File"
	}
	fmt.Fprintf(&sb, "%d %s, %s\n", len(entries), files, totalSize(entries))

	if cfg.Header {
		fmt.Fprintf(&sb, "%-*s %-*s %-*s %-*s %*s %s\n",
			w.perms, "Permissions",
			w.user, "User",
			w.group, "Group",
			w.mod, "Modified",
			w.size, "Size",
			"Name",
		)
	}
	for _, r := range rows {
		formatRow(&sb, r, w)
	}
	return sb.String()
}

// processEntries builds rows and calculates column widths for long-format output.
func processEntries(entries []Entry, cfg Config) ([]row, widths) {
	rows := make([]row, len(entries))
	w := widths{
		perms: utf8.RuneCountInString("Permissions"),
		user:  utf8.RuneCountInString("User"),
		group: utf8.RuneCountInString("Group"),
		mod:   utf8.RuneCountInString("Modified"),
		size:  utf8.RuneCountInString("Size"),
	}

	for i, e := range entries {
		u, g := userGroup(e)
		rows[i] = row{
			perms:   e.Mode().String(),
			user:    u,
			group:   g,
			modTime: formatTime(e.ModTime()),
			size:    humanReadableSize(e.Size()),
			name:    e.displayName,
		}

		if e.target != "" {
			if cfg.Dereference {
				rows[i] = row{
					perms:   "----------",
					user:    "-",
					group:   "-",
					modTime: "-",
					size:    "-",
					name:    e.displayName,
					target:  " [nonexist]",
				}
				continue
			}
			rows[i].target = " -> " + e.target
			rows[i].size = humanReadableSize(int64(len(e.target)))
		}

		r := rows[i]
		w.perms = max(w.perms, utf8.RuneCountInString(r.perms))
		w.user = max(w.user, utf8.RuneCountInString(r.user))
		w.group = max(w.group, utf8.RuneCountInString(r.group))
		w.mod = max(w.mod, utf8.RuneCountInString(r.modTime))
		w.size = max(w.size, utf8.RuneCountInString(r.size))
	}
	return rows, w
}

// formatRow appends a formatted row with aligned columns to the builder.
func formatRow(sb *strings.Builder, r row, w widths) {
	fmt.Fprintf(sb, "%-*s %-*s %-*s %-*s %*s %s%s\n",
		w.perms, r.perms,
		w.user, r.user,
		w.group, r.group,
		w.mod, r.modTime,
		w.size, r.size,
		r.name, r.target,
	)
}

// humanReadableSize converts bytes to a human-readable string (e.g., "1.2M").
func humanReadableSize(size int64) string {
	const (
		_ = 1 << (iota * 10)
		K
		M
		G
		T
		P
		E
	)
	if size < 0 {
		size = 0
	}
	switch {
	case size < K:
		return fmt.Sprintf("%d", size)
	case size < M:
		return fmt.Sprintf("%.1fK", float64(size)/K)
	case size < G:
		return fmt.Sprintf("%.1fM", float64(size)/M)
	case size < T:
		return fmt.Sprintf("%.1fG", float64(size)/G)
	case size < P:
		return fmt.Sprintf("%.1fT", float64(size)/T)
	case size < E:
		return fmt.Sprintf("%.1fP", float64(size)/P)
	default:
		return fmt.Sprintf("%.1fE", float64(size)/E)
	}
}

// totalSize returns the total size of entries in human-readable format.
func totalSize(entries []Entry) string {
	var t int64
	for _, e := range entries {
		t += e.Size()
	}
	return humanReadableSize(t)
}

// formatTime returns the formatted time.
func formatTime(t time.Time) string {
	if t.Year() == time.Now().Year() {
		return t.Format("Jan 02 15:04")
	}
	return t.Format("Jan 02  2006")
}
