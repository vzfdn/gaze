package entry

import (
	"fmt"
	"strings"
	"time"
)

const (
	headerPerms   = "Permissions"
	headerUser    = "User"
	headerGroup   = "Group"
	headerModTime = "Date Modified"
	headerSize    = "Size"
	headerName    = "Name"

	labelFiles = "Files"
	labelFile  = "File"

	placeholderPerms    = "----------"
	placeholderField    = "-"
	placeholderNonexist = " [nonexist]"
	linkPrefix          = " -> "

	timeFormatCurrentYear = "Jan 02 15:04"
	timeFormatOlderYear   = "Jan 02  2006"

	fieldPadding   = 6
	coloredColumns = 4
	capacityBuffer = 50
)

const (
	_ = 1 << (iota * 10)
	kb
	mb
	gb
	tb
	pb
	eb
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
	
	rows, widths := buildTable(entries)
	summary := summaryLine(entries)

	var sb strings.Builder
	sb.Grow(estimateCapacity(rows, widths) + len(summary))
	sb.WriteString(summary)

	if cfg.Header {
		header := row{
			perms:   headerPerms,
			user:    headerUser,
			group:   headerGroup,
			size:    headerSize,
			modTime: headerModTime,
			name:    headerName,
		}
		writeRow(&sb, header, widths)
	}

	for _, r := range rows {
		writeRow(&sb, r, widths)
	}
	return sb.String()
}

func summaryLine(entries []Entry) string {
	label := labelFiles
	if len(entries) == 1 {
		label = labelFile
	}
	return fmt.Sprintf("%d %s, %s\n", len(entries), label, totalSize(entries))
}

func estimateCapacity(rows []row, widths columnWidths) (capacity int) {
	if cfg.Header {
		capacity += len(headerPerms+headerGroup+headerUser+headerSize+headerModTime+headerName) + fieldPadding
	}
	baseRowWidth := widths.perms + widths.user + widths.group + widths.size + widths.mod + fieldPadding
	capacity += len(rows) * baseRowWidth
	for _, r := range rows {
		// Add visible width of name and symlink target
		capacity += visibleWidth(r.name) + len(r.target)
	}
	if !cfg.NoColor {
		// Add capacity for colored output (permissions + metadata fields)
		capacity += len(rows) * (len(placeholderPerms) + coloredColumns) * ansiPerField
	}
	// Add 2% buffer to reduce chance of reallocation
	capacity += capacity / capacityBuffer
	return capacity
}

func writeRow(sb *strings.Builder, r row, widths columnWidths) {
	fmt.Fprintf(sb, "%s %s %s %s %s %s%s\n",
		padToWidth(r.perms, widths.perms, false),
		padToWidth(r.user, widths.user, false),
		padToWidth(r.group, widths.group, false),
		padToWidth(r.size, widths.size, true),
		padToWidth(r.modTime, widths.mod, false),
		r.name,
		r.target,
	)
}

func buildTable(entries []Entry) ([]row, columnWidths) {
	var widths columnWidths
	if cfg.Header {
		widths = columnWidths{
			perms: len(headerPerms),
			user:  len(headerUser),
			group: len(headerGroup),
			size:  len(headerSize),
			mod:   len(headerModTime),
		}
	}
	rows := make([]row, len(entries))
	for i, entry := range entries {
		rows[i] = makeRow(entry)
		widths.perms = max(widths.perms, visibleWidth(rows[i].perms))
		widths.user = max(widths.user, visibleWidth(rows[i].user))
		widths.group = max(widths.group, visibleWidth(rows[i].group))
		widths.size = max(widths.size, visibleWidth(rows[i].size))
		widths.mod = max(widths.mod, visibleWidth(rows[i].modTime))
	}
	return rows, widths
}

func makeRow(entry Entry) row {
	if entry.link != nil && cfg.Dereference {
		return row{
			perms:   color.placeholder(placeholderPerms),
			user:    color.placeholder(placeholderField),
			group:   color.placeholder(placeholderField),
			size:    color.placeholder(placeholderField),
			modTime: color.placeholder(placeholderField),
			name:    entry.DisplayName(),
			target:  placeholderNonexist,
		}
	}
	user, group := userGroup(entry)
	r := row{
		perms:   color.permissions(entry.Mode()),
		user:    color.user(user),
		group:   color.group(group),
		size:    formatSize(entry.Size()),
		modTime: color.modTime(formatModTime(entry.ModTime())),
		name:    entry.DisplayName(),
	}
	if entry.link != nil {
		r.target = linkPrefix + entry.link.target
		r.size = formatSize(int64(len(entry.link.target)))
	}
	return r
}

func padToWidth(s string, width int, rightAlign bool) string {
	paddingNeeded := width - visibleWidth(s)
	if paddingNeeded <= 0 {
		return s
	}
	padding := strings.Repeat(" ", paddingNeeded)
	if rightAlign {
		return padding + s
	}
	return s + padding
}

func formatSize(size int64) string {
	if size < 0 {
		size = 0
	}
	switch {
	case size < kb:
		return color.colorize(fmt.Sprintf("%d", size), colorSizeBytes)
	case size < mb:
		return color.colorize(fmt.Sprintf("%.1fK", float64(size)/kb), colorSizeKB)
	case size < gb:
		return color.colorize(fmt.Sprintf("%.1fM", float64(size)/mb), colorSizeMB)
	case size < tb:
		return color.colorize(fmt.Sprintf("%.1fG", float64(size)/gb), colorSizeGB)
	case size < pb:
		return color.colorize(fmt.Sprintf("%.1fT", float64(size)/tb), colorSizeTB)
	case size < eb:
		return color.colorize(fmt.Sprintf("%.1fP", float64(size)/pb), colorSizePB)
	default:
		return color.colorize(fmt.Sprintf("%.1fE", float64(size)/eb), colorSizeEB)
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
		return t.Format(timeFormatCurrentYear)
	}
	return t.Format(timeFormatOlderYear)
}
