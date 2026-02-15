package pathutil

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestExpandPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("could not get home dir: %v", err)
	}

	t.Run("Expands tilde with forward slash", func(t *testing.T) {
		got := ExpandPath("~/Downloads")
		want := filepath.Join(home, "Downloads")
		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("Expands tilde with backslash", func(t *testing.T) {
		got := ExpandPath("~\\Downloads")
		want := filepath.Join(home, "Downloads")
		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("Ignores absolute path", func(t *testing.T) {
		var path string
		if runtime.GOOS == "windows" {
			path = "C:\\Users\\test"
		} else {
			path = "/home/test"
		}
		if ExpandPath(path) != path {
			t.Error("should not modify absolute path")
		}
	})

	t.Run("Ignores empty path", func(t *testing.T) {
		if ExpandPath("") != "" {
			t.Error("should return empty string unchanged")
		}
	})

	t.Run("Returns tilde-only unchanged", func(t *testing.T) {
		if ExpandPath("~") != "~" {
			t.Error("bare tilde without separator should be returned unchanged")
		}
	})
}
