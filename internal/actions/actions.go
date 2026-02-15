package actions

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"straw/internal/config"
	"straw/internal/pathutil"
)

type Executor struct{}

func NewExecutor() *Executor {
	return &Executor{}
}

func (e *Executor) Execute(action config.Action, sourcePath string) error {
	switch action.Type {
	case "move":
		return e.move(sourcePath, action.Target)
	case "copy":
		return e.copy(sourcePath, action.Target)
	case "trash":
		return e.trash(sourcePath)
	case "shell":
		return e.shell(sourcePath, action.Target, action.Args)
	default:
		return fmt.Errorf("unknown action type: %s", action.Type)
	}
}

func (e *Executor) move(src, destDir string) error {
	destDir = e.expandPath(destDir)
	destDir = filepath.Clean(destDir)

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	destPath := e.getDestPath(src, destDir)

	// Ensure we aren't moving a file into itself or sensitive areas if we can help it
	if destPath == src {
		return fmt.Errorf("source and destination are the same: %s", src)
	}

	// Attempt standard rename
	err := os.Rename(src, destPath)
	if err == nil {
		return nil
	}

	// If rename fails (e.g. cross-device), fallback to copy + delete
	slog.Debug("Standard rename failed, attempting copy+delete fallback", "src", src, "dest", destPath, "error", err)

	if err := e.copyToPath(src, destPath); err != nil {
		return fmt.Errorf("fallback copy failed: %w", err)
	}

	return os.Remove(src)
}

func (e *Executor) copy(src, destDir string) error {
	destDir = e.expandPath(destDir)
	destDir = filepath.Clean(destDir)

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	destPath := e.getDestPath(src, destDir)
	if destPath == src {
		return fmt.Errorf("source and destination are the same: %s", src)
	}

	return e.copyToPath(src, destPath)
}

func (e *Executor) expandPath(path string) string {
	return pathutil.ExpandPath(path)
}

func (e *Executor) copyToPath(src, dest string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func (e *Executor) getDestPath(src, destDir string) string {
	filename := filepath.Base(src)
	destPath := filepath.Join(destDir, filename)

	// If destination exists, append timestamp to prevent overwrite
	if _, err := os.Stat(destPath); err == nil {
		ext := filepath.Ext(filename)
		base := strings.TrimSuffix(filename, ext)
		newName := fmt.Sprintf("%s_%d%s", base, time.Now().Unix(), ext)
		destPath = filepath.Join(destDir, newName)
		slog.Info("Collision detected, renamed destination", "original", filename, "new", newName)
	}

	return destPath
}

func (e *Executor) shell(src, command string, args []string) error {
	// SECURITY WARNING: $FILE substitution can be dangerous if the command
	// is executed through a shell (e.g. bash -c "echo $FILE").
	// We use exec.Command which avoids shell interpretation by default,
	// but users should still be cautious of malicious filenames.
	var finalCmd string
	var finalArgs []string

	if len(args) == 0 {
		// If no args provided, split the command string
		parts := strings.Fields(command)
		if len(parts) == 0 {
			return fmt.Errorf("empty shell command")
		}
		finalCmd = parts[0]
		finalArgs = parts[1:]
	} else {
		finalCmd = command
		finalArgs = args
	}

	expandedArgs := make([]string, len(finalArgs))
	for i, arg := range finalArgs {
		expandedArgs[i] = strings.ReplaceAll(arg, "$FILE", src)
	}

	cmd := exec.Command(finalCmd, expandedArgs...)
	cmd.Env = append(os.Environ(), "STRAW_FILE="+src)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command failed: %s, output: %s", err, out)
	}
	return nil
}
