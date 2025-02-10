package entry

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// HumanReadableSize converts a size in bytes to a human-readable
// string with appropriate units.
func HumanReadableSize(size int64) string {
	const (
		_ = 1 << (iota * 10) // ignore first value
		K                    // 1024
		M                    // 1024^2
		G                    // 1024^3
		T                    // 1024^4
		P                    // 1024^5
		E                    // 1024^6
	)

	switch {
	case size < K:
		return fmt.Sprintf("%d", size) // Bytes
	case size < M:
		return fmt.Sprintf("%.1fK", float64(size)/K) // Kilobytes
	case size < G:
		return fmt.Sprintf("%.1fM", float64(size)/M) // Megabytes
	case size < T:
		return fmt.Sprintf("%.1fG", float64(size)/G) // Gigabytes
	case size < P:
		return fmt.Sprintf("%.1fT", float64(size)/T) // Terabytes
	case size < E:
		return fmt.Sprintf("%.1fP", float64(size)/P) // Petabytes
	default:
		return fmt.Sprintf("%.1fE", float64(size)/E) // Exabytes
	}
}

// TotalSize returns total size of entries.
func TotalSize(displayEntries []Entry) string {
	var t int64
	for _, e := range displayEntries {
		t += e.info.Size()
	}
	return HumanReadableSize(t)
}

// GetTerminalWidth returns the current terminal width,
// falling back to 80 if an error occurs or width is invalid.
func GetTerminalWidth() (int, error) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0, fmt.Errorf("cannot determine terminal width: %w", err)
	}
	if width <= 0 {
		width = 80
	}
	return width, nil
}
