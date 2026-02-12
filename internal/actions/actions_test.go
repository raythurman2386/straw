package actions

import (
	"os"
	"path/filepath"
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
		path := "/tmp/test"
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

func TestExecutor_Shell(t *testing.T) {
	e := NewExecutor()
	tmpDir, err := os.MkdirTemp("", "shell_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	src := filepath.Join(tmpDir, "test.sh")
	if err := os.WriteFile(src, []byte("echo hello"), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("Split command and substitute $FILE", func(t *testing.T) {
		// Using 'ls' as a safe command that should exist
		err := e.shell(src, "ls -l $FILE", nil)
		if err != nil {
			t.Errorf("shell failed: %v", err)
		}
	})

	t.Run("Full command as Target", func(t *testing.T) {
		err := e.shell(src, "touch "+filepath.Join(tmpDir, "touched.txt"), nil)
		if err != nil {
			t.Errorf("shell failed: %v", err)
		}
		if _, err := os.Stat(filepath.Join(tmpDir, "touched.txt")); os.IsNotExist(err) {
			t.Error("shell command did not create file")
		}
	})
}
