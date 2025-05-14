package entry

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	HeaderPerms    = "Permissions"
	HeaderUser     = "User"
	HeaderGroup    = "Group"
	HeaderModified = "Date Modified"
	HeaderSize     = "Size"
	HeaderName     = "Name"

	LabelFiles = "Files"
	LabelFile  = "File"

	PlaceholderPerms    = "----------"
	PlaceholderField    = "-"
	PlaceholderNonexist = " [nonexist]"

	linkPrefix = " -> "

	TimeFormatCurrentYear = "Jan 02 15:04"
	TimeFormatOlderYear   = "Jan 02  2006"
)

type row struct {
	perms   string
	user    string
	group   string
	size    string
	modTime string
	name    string
	target  string
}

type columnWidths struct {
	perms, user, group, size, mod int
}

func renderLong(entries []Entry) string {
	if len(entries) == 0 {
		return ""
	}

	rows, cw := buildDisplayTable(entries)
	var sb strings.Builder
	sb.Grow(len(entries) * (cw.perms + cw.user + cw.group + cw.size + cw.mod + 20))

	files := LabelFiles
	if len(entries) == 1 {
		files = LabelFile
	}
	fmt.Fprintf(&sb, "%d %s, %s\n", len(entries), files, totalSize(entries))

	if cfg.Header {
		fmt.Fprintf(&sb, "%s %s %s %s %s %s\n",
			padToWidth(HeaderPerms, cw.perms, false),
			padToWidth(HeaderUser, cw.user, false),
			padToWidth(HeaderGroup, cw.group, false),
			padToWidth(HeaderSize, cw.size, true),
			padToWidth(HeaderModified, cw.mod, false),
			HeaderName,
		)
	}
	for _, r := range rows {
		formatRow(&sb, r, cw)
	}
	return sb.String()
}

func buildDisplayTable(entries []Entry) ([]row, columnWidths) {
	rows := make([]row, len(entries))
	cw := columnWidths{
		perms: utf8.RuneCountInString(HeaderPerms),
		user:  utf8.RuneCountInString(HeaderUser),
		group: utf8.RuneCountInString(HeaderGroup),
		size:  utf8.RuneCountInString(HeaderSize),
		mod:   utf8.RuneCountInString(HeaderModified),
	}

	for i, e := range entries {
		u, g := userGroup(e)
		rows[i] = row{
			perms:   color.permissions(e.Mode()),
			user:    color.user(u),
			group:   color.group(g),
			size:    formatSize(e.Size()),
			modTime: color.modTime(formatModTime(e.ModTime())),
			name:    e.DisplayName(),
		}

		if e.link != nil {
			if cfg.Dereference {
				rows[i] = row{
					perms:   color.placeholder(PlaceholderPerms),
					user:    color.placeholder(PlaceholderField),
					group:   color.placeholder(PlaceholderField),
					size:    color.placeholder(PlaceholderField),
					modTime: color.placeholder(PlaceholderField),
					name:    e.DisplayName(),
					target:  PlaceholderNonexist,
				}
				continue
			}
			rows[i].target = linkPrefix + e.link.target
			rows[i].size = formatSize(int64(len(e.link.target)))
		}

		r := rows[i]
		cw.perms = max(cw.perms, visibleWidth(r.perms))
		cw.user = max(cw.user, visibleWidth(r.user))
		cw.group = max(cw.group, visibleWidth(r.group))
		cw.size = max(cw.size, visibleWidth(r.size))
		cw.mod = max(cw.mod, visibleWidth(r.modTime))
	}
	return rows, cw
}

func formatRow(sb *strings.Builder, r row, cw columnWidths) {
	fmt.Fprintf(sb, "%s %s %s %s %s %s%s\n",
		padToWidth(r.perms, cw.perms, false),
		padToWidth(r.user, cw.user, false),
		padToWidth(r.group, cw.group, false),
		padToWidth(r.size, cw.size, true),
		padToWidth(r.modTime, cw.mod, false),
		r.name,
		r.target,
	)
}

func padToWidth(s string, width int, rightAlign bool) string {
	padding := strings.Repeat(" ", width-visibleWidth(s))
	if rightAlign {
		return padding + s
	}
	return s + padding
}

func formatSize(size int64) string {
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
		return color.colorize(fmt.Sprintf("%d", size), colorSizeBytes)
	case size < M:
		return color.colorize(fmt.Sprintf("%.1fK", float64(size)/K), colorSizeKB)
	case size < G:
		return color.colorize(fmt.Sprintf("%.1fM", float64(size)/M), colorSizeMB)
	case size < T:
		return color.colorize(fmt.Sprintf("%.1fG", float64(size)/G), colorSizeGB)
	case size < P:
		return color.colorize(fmt.Sprintf("%.1fT", float64(size)/T), colorSizeTB)
	case size < E:
		return color.colorize(fmt.Sprintf("%.1fP", float64(size)/P), colorSizePB)
	default:
		return color.colorize(fmt.Sprintf("%.1fE", float64(size)/E), colorSizePB)
	}
}

func totalSize(entries []Entry) string {
	var t int64
	for _, e := range entries {
		t += e.Size()
	}
	return formatSize(t)
}

func formatModTime(t time.Time) string {
	if t.Year() == time.Now().Year() {
		return t.Format(TimeFormatCurrentYear)
	}
	return t.Format(TimeFormatOlderYear)
}
