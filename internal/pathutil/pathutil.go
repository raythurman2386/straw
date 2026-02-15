package pathutil

import (
	"os"
	"path/filepath"
	"strings"
)

// ExpandPath expands a leading ~/ in a path to the user's home directory.
// On all platforms, os.UserHomeDir() is used which returns:
//   - Linux:   $HOME
//   - macOS:   $HOME
//   - Windows: %USERPROFILE%
//
// If the path does not start with ~/ or the home directory cannot be
// determined, the path is returned unchanged.
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~\\") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}
