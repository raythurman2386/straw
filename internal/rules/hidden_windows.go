//go:build windows

package rules

import (
	"strings"
	"syscall"
)

// isHidden reports whether a file is hidden on Windows.
// Checks both the Windows FILE_ATTRIBUTE_HIDDEN attribute and the
// Unix-style dot prefix convention (which some cross-platform tools use).
func isHidden(name, path string) bool {
	// Check dot prefix (cross-platform convention)
	if strings.HasPrefix(name, ".") {
		return true
	}

	// Check Windows hidden file attribute
	p, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false
	}
	attrs, err := syscall.GetFileAttributes(p)
	if err != nil {
		return false
	}
	return attrs&syscall.FILE_ATTRIBUTE_HIDDEN != 0
}
