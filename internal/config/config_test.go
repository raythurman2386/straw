package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("Valid config", func(t *testing.T) {
		c := &Config{
			Watch: []WatchFolder{{Path: tmpDir, Recursive: true}},
			Rules: []Rule{{Name: "Test", Actions: []Action{{Type: "move", Target: "/tmp"}}}},
		}
		if err := c.Validate(); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("Missing watch", func(t *testing.T) {
		c := &Config{
			Watch: []WatchFolder{},
		}
		if err := c.Validate(); err == nil {
			t.Error("expected error for empty watch list")
		}
	})

	t.Run("Non-existent watch path", func(t *testing.T) {
		c := &Config{
			Watch: []WatchFolder{{Path: "/non/existent/path"}},
		}
		if err := c.Validate(); err == nil {
			t.Error("expected error for non-existent path")
		}
	})

	t.Run("Watch path is a file", func(t *testing.T) {
		tmpFile := filepath.Join(tmpDir, "file.txt")
		if err := os.WriteFile(tmpFile, []byte(""), 0644); err != nil {
			t.Fatal(err)
		}

		c := &Config{Watch: []WatchFolder{{Path: tmpFile}}}
		if err := c.Validate(); err == nil {
			t.Error("expected error for watch path being a file")
		}
	})
}
