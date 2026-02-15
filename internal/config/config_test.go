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
			Rules: []Rule{{Name: "Test", Actions: []Action{{Type: "move", Target: filepath.Join(os.TempDir(), "archive")}}}},
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
			Watch: []WatchFolder{{Path: filepath.Join(os.TempDir(), "non_existent_straw_test_path")}},
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

func TestConfig_Save(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "config_save_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	socketPath := filepath.Join(os.TempDir(), "straw_test.sock")
	configPath := filepath.Join(tmpDir, "config.toml")
	c := &Config{
		SocketPath: socketPath,
		Watch:      []WatchFolder{{Path: tmpDir}},
		Rules:      []Rule{{Name: "Test", Actions: []Action{{Type: "trash"}}}},
		TUI: TUIConfig{
			Theme: "catppuccin",
		},
	}

	if err := c.Save(configPath); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("config file was not created")
	}

	// Load it back and verify
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("failed to load saved config: %v", err)
	}

	if loaded.TUI.Theme != "catppuccin" {
		t.Errorf("expected theme catppuccin, got %s", loaded.TUI.Theme)
	}
	if loaded.SocketPath != socketPath {
		t.Errorf("expected socket path %s, got %s", socketPath, loaded.SocketPath)
	}
}
