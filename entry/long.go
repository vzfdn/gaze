package entry

import (
	"fmt"
	"strings"
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
}

// widths holds the maximum column widths for long format rendering.
type widths struct {
	perms int
	user  int
	group int
	mod   int
	size  int
}

// renderLong renders a detailed view of file entries, aligned in columns.
func renderLong(entries []Entry, cfg Config) string {
	if len(entries) == 0 {
		return "total 0\n"
	}

	rows, w := processEntries(entries)

	var sb strings.Builder
	// Preallocate approximate capacity
	sb.Grow(20 + len(entries)*(w.perms+w.user+w.group+w.mod+w.size+20))

	// Summary with file count
	files := "Files"
	if len(entries) == 1 {
		files = "File"
	}
	fmt.Fprintf(&sb, "%d %s, %s\n", len(entries), files, totalSize(entries))

	// Header
	if cfg.Header {
		fmt.Fprintf(&sb, " %-*s %-*s %-*s %-*s %*s %s\n",
			w.perms, "Permissions",
			w.user, "User",
			w.group, "Group",
			w.mod, "Modified",
			w.size, "Size",
			"Name",
		)
	}

	// Rows
	for _, r := range rows {
		fmt.Fprintf(&sb, " %-*s %-*s %-*s %-*s %*s %s\n",
			w.perms, r.perms,
			w.user, r.user,
			w.group, r.group,
			w.mod, r.modTime,
			w.size, r.size,
			r.name,
		)
	}

	return sb.String()
}

// processEntries processes file entries and calculates maximum column widths.
func processEntries(entries []Entry) ([]row, widths) {
	rows := make([]row, 0, len(entries))
	w := widths{
		perms: utf8.RuneCountInString("Permissions"),
		user:  utf8.RuneCountInString("User"),
		group: utf8.RuneCountInString("Group"),
		mod:   utf8.RuneCountInString("Date Modified"),
		size:  utf8.RuneCountInString("Size"),
	}

	for _, e := range entries {
		u, g := e.UserAndGroup()
		r := row{
			perms:   e.Permission(),
			user:    u,
			group:   g,
			modTime: e.Time(),
			size:    humanReadableSize(e.Size()),
			name:    e.Name(),
		}
		w.perms = max(w.perms, utf8.RuneCountInString(r.perms))
		w.user = max(w.user, utf8.RuneCountInString(r.user))
		w.group = max(w.group, utf8.RuneCountInString(r.group))
		w.mod = max(w.mod, utf8.RuneCountInString(r.modTime))
		w.size = max(w.size, utf8.RuneCountInString(r.size))
		rows = append(rows, r)
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

// humanReadableSize converts a size in bytes to a human-readable string with units.
func humanReadableSize(size int64) string {
	const (
		_ = 1 << (iota * 10) // Ignore first value
		K                    // 1024
		M                    // 1024^2
		G                    // 1024^3
		T                    // 1024^4
		P                    // 1024^5
		E                    // 1024^6
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

// totalSize returns the total size of entries in a human-readable format.
func totalSize(entries []Entry) string {
	var t int64
	for _, e := range entries {
		t += e.info.Size()
	}
	return humanReadableSize(t)
}
