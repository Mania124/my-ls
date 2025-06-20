package permissions

import (
	"fmt"
	"os"
)

// GetPermissions converts file mode to human-readable string
func GetPermissions(mode os.FileMode) string {
	perm := mode.String()
	if mode&os.ModeSymlink != 0 {
		perm = "l" + perm[1:] // Ensure lowercase 'l' for symbolic links
	}
	return perm
}

// HumanReadableSize converts bytes to human-readable sizes (-h flag)
func HumanReadableSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%dB", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%c", float64(size)/float64(div), "KMGTPE"[exp])
}
