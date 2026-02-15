//go:build !windows

package actions

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// trash moves a file to the FreeDesktop.org Trash on Linux/macOS.
// Uses the XDG Trash specification directory at ~/.local/share/Trash.
// Note: macOS has a native ~/.Trash directory, but this implementation
// uses the FreeDesktop.org spec for consistency across Unix platforms.
func (e *Executor) trash(src string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	trashBase := filepath.Join(home, ".local", "share", "Trash")
	filesDir := filepath.Join(trashBase, "files")
	infoDir := filepath.Join(trashBase, "info")

	if err := os.MkdirAll(filesDir, 0700); err != nil {
		return err
	}
	if err := os.MkdirAll(infoDir, 0700); err != nil {
		return err
	}

	absSrc, err := filepath.Abs(src)
	if err != nil {
		absSrc = src
	}

	destPath := e.getDestPath(src, filesDir)
	actualFilename := filepath.Base(destPath)

	// Create .trashinfo file per FreeDesktop.org Trash spec
	infoContent := fmt.Sprintf("[Trash Info]\nPath=%s\nDeletionDate=%s\n",
		absSrc,
		time.Now().Format("2006-01-02T15:04:05"))

	infoFilePath := filepath.Join(infoDir, actualFilename+".trashinfo")
	if err := os.WriteFile(infoFilePath, []byte(infoContent), 0600); err != nil {
		return fmt.Errorf("failed to create trash info: %w", err)
	}

	// Move file to trash
	if err := os.Rename(src, destPath); err != nil {
		// Fallback to copy+delete for cross-device moves
		if err := e.copyToPath(src, destPath); err != nil {
			// Cleanup info file if move fails
			os.Remove(infoFilePath)
			return fmt.Errorf("trash move failed: %w", err)
		}
		os.Remove(src)
	}

	return nil
}
