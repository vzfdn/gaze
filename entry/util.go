package entry

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

// HumanReadableSize converts a size in bytes to a human-readable string with appropriate units.
func HumanReadableSize(size int64) string {
	const (
		_     = 1 << (iota * 10) // ignore first value  
		K                        // 1024
		M                        // 1024^2
		G                        // 1024^3
		T                        // 1024^4
		P                        // 1024^5
		E                        // 1024^6
		width = 6
	)
	var sizeStr string
	switch {
	case size < K:
		sizeStr = fmt.Sprintf("%d", size)
	case size < M:
		sizeStr = fmt.Sprintf("%.1fK", float64(size)/K)
	case size < G:
		sizeStr = fmt.Sprintf("%.1fM", float64(size)/M)
	case size < T:
		sizeStr = fmt.Sprintf("%.1fG", float64(size)/G)
	case size < P:
		sizeStr = fmt.Sprintf("%.1fT", float64(size)/T)
	case size < E:
		sizeStr = fmt.Sprintf("%.1fP", float64(size)/P)
	default:
		sizeStr = fmt.Sprintf("%.1fE", float64(size)/E)
	}
	return fmt.Sprintf("%*s", width, sizeStr)
}

// totalSize returns total size of entries.
func TotalSize(displayEntries []Entry) string {
	var t int64
	for _, e := range displayEntries {
		t += e.info.Size()
	}
	return HumanReadableSize(t)
}

// GetTerminalWidth returns the current terminal width, falling back to 80 if an error occurs or width is invalid.
func GetTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Println("Error getting terminal width:", err)
		return 80 
	}
	if width <= 0 {
		width = 80
	}
	return width
}
