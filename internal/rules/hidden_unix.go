//go:build !windows

package rules

import "strings"

// isHidden reports whether a file is hidden on Unix systems.
// On Unix, files starting with a dot are considered hidden.
func isHidden(name, path string) bool {
	return strings.HasPrefix(name, ".")
}
