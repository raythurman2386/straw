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

	t.Run("Invalid rule", func(t *testing.T) {
		c := &Config{
			Watch: []WatchFolder{{Path: tmpDir}},
			Rules: []Rule{{Name: "", Actions: []Action{}}}, // Missing name and actions
		}
		if err := c.Validate(); err == nil {
			t.Error("expected error for invalid rule")
		}
	})
}

func TestRule_Validate(t *testing.T) {
	tests := []struct {
		name    string
		rule    Rule
		wantErr bool
	}{
		{"valid rule", Rule{Name: "test", Actions: []Action{{Type: "trash"}}}, false},
		{"missing name", Rule{Name: "", Actions: []Action{{Type: "trash"}}}, true},
		{"missing actions", Rule{Name: "test", Actions: []Action{}}, true},
		{"unknown action", Rule{Name: "test", Actions: []Action{{Type: "foo"}}}, true},
		{"move missing target", Rule{Name: "test", Actions: []Action{{Type: "move"}}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rule.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Rule.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_LoadResiliency(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "config_load_resiliency")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.toml")
	content := `
socket_path = "/tmp/straw.sock"
[[watch]]
path = "` + filepath.ToSlash(tmpDir) + `"

[[rules]]
name = "bad_rule"
# Missing actions - this is invalid

[[rules]]
name = "good_rule"
actions = [{type = "trash"}]
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Load should succeed because watch folder is valid, even if rules are bad.
	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load should have succeeded with bad rules, but got error: %v", err)
	}

	if len(cfg.Rules) != 2 {
		t.Errorf("expected 2 rules to be loaded, got %d", len(cfg.Rules))
	}
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
