package actions

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestExecutor_ExpandPath(t *testing.T) {
	e := NewExecutor()
	home, _ := os.UserHomeDir()

	t.Run("Expands tilde", func(t *testing.T) {
		got := e.expandPath("~/Downloads")
		want := filepath.Join(home, "Downloads")
		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("Ignores non-tilde", func(t *testing.T) {
		path := filepath.Join(os.TempDir(), "test")
		if e.expandPath(path) != path {
			t.Error("should not modify absolute path")
		}
	})
}

func TestExecutor_CollisionHandling(t *testing.T) {
	e := NewExecutor()
	tmpDir, err := os.MkdirTemp("", "action_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	srcFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(srcFile, []byte("source"), 0644); err != nil {
		t.Fatal(err)
	}

	destDir := filepath.Join(tmpDir, "dest")
	if err := os.Mkdir(destDir, 0755); err != nil {
		t.Fatal(err)
	}

	existingFile := filepath.Join(destDir, "test.txt")
	if err := os.WriteFile(existingFile, []byte("existing"), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("Rename on collision", func(t *testing.T) {
		destPath := e.getDestPath(srcFile, destDir)
		if destPath == existingFile {
			t.Error("should have generated a new path name")
		}
		if !strings.Contains(destPath, "test_") {
			t.Errorf("expected timestamped name, got %s", destPath)
		}
	})
}

func TestExecutor_MoveCrossDevice(t *testing.T) {
	e := NewExecutor()
	tmpDir, err := os.MkdirTemp("", "move_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	src := filepath.Join(tmpDir, "src.txt")
	if err := os.WriteFile(src, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	destDir := filepath.Join(tmpDir, "out")

	err = e.move(src, destDir)
	if err != nil {
		t.Fatalf("move failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(destDir, "src.txt")); os.IsNotExist(err) {
		t.Error("file should exist in destination")
	}
	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Error("source file should have been removed")
	}
}

// TestExecutor_MoveSourceAlreadyMoved tests that calling move on a file
// that was already moved to the destination (e.g. by a duplicate event)
// succeeds without returning an error.
func TestExecutor_MoveSourceAlreadyMoved(t *testing.T) {
	e := NewExecutor()
	tmpDir, err := os.MkdirTemp("", "move_already_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	src := filepath.Join(tmpDir, "src.txt")
	if err := os.WriteFile(src, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	destDir := filepath.Join(tmpDir, "out")

	// First move should succeed normally
	err = e.move(src, destDir)
	if err != nil {
		t.Fatalf("first move failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(destDir, "src.txt")); os.IsNotExist(err) {
		t.Fatal("file should exist in destination after first move")
	}

	// Second move on the same (now missing) source path should not error,
	// because the file is already at the destination.
	err = e.move(src, destDir)
	if err != nil {
		t.Errorf("second move should succeed (file already at destination), got: %v", err)
	}
}

// TestExecutor_MoveSourceGone tests that calling move on a source file
// that no longer exists and is not at the destination returns an error.
func TestExecutor_MoveSourceGone(t *testing.T) {
	e := NewExecutor()
	tmpDir, err := os.MkdirTemp("", "move_gone_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	src := filepath.Join(tmpDir, "ghost.txt")
	destDir := filepath.Join(tmpDir, "out")

	// Source never existed — move should return an error
	err = e.move(src, destDir)
	if err == nil {
		t.Error("move should fail when source does not exist and is not at destination")
	}
}

func TestExecutor_Shell(t *testing.T) {
	e := NewExecutor()
	tmpDir, err := os.MkdirTemp("", "shell_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	src := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(src, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("Split command and substitute $FILE", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			// Use 'cmd /c dir' on Windows instead of 'ls'
			err := e.shell(src, "cmd", []string{"/c", "dir", "$FILE"})
			if err != nil {
				t.Errorf("shell failed: %v", err)
			}
		} else {
			err := e.shell(src, "ls -l $FILE", nil)
			if err != nil {
				t.Errorf("shell failed: %v", err)
			}
		}
	})

	t.Run("Full command as Target", func(t *testing.T) {
		touchedFile := filepath.Join(tmpDir, "touched.txt")
		if runtime.GOOS == "windows" {
			// Use 'cmd /c type nul >' equivalent via Go-friendly command
			err := e.shell(src, "cmd", []string{"/c", "copy", "nul", touchedFile})
			if err != nil {
				t.Errorf("shell failed: %v", err)
			}
		} else {
			err := e.shell(src, "touch "+touchedFile, nil)
			if err != nil {
				t.Errorf("shell failed: %v", err)
			}
		}
		if _, err := os.Stat(touchedFile); os.IsNotExist(err) {
			t.Error("shell command did not create file")
		}
	})
}
