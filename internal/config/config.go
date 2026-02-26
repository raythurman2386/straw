package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"straw/internal/pathutil"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	SocketPath string        `toml:"socket_path"`
	Watch      []WatchFolder `toml:"watch"`
	Rules      []Rule        `toml:"rules"`
	TUI        TUIConfig     `toml:"tui"`
}

type WatchFolder struct {
	Path      string `toml:"path"`
	Recursive bool   `toml:"recursive"`
}

type Rule struct {
	Name        string   `toml:"name"`
	Enabled     bool     `toml:"enabled"`
	Match       Match    `toml:"match"`
	Actions     []Action `toml:"actions"`
	Description string   `toml:"description"`
}

type Match struct {
	Glob       string `toml:"glob"`
	Regex      string `toml:"regex"`
	Extension  string `toml:"extension"`
	MinSize    int64  `toml:"min_size"`
	MaxSize    int64  `toml:"max_size"`
	MinAgeDays int    `toml:"min_age_days"`
	MaxAgeDays int    `toml:"max_age_days"`
	FileType   string `toml:"file_type"`
	Hidden     *bool  `toml:"hidden"`
}

type Action struct {
	Type   string   `toml:"type"`
	Target string   `toml:"target"`
	Args   []string `toml:"args"`
}

type TUIConfig struct {
	Theme string `toml:"theme"`
}

// DefaultConfigPath returns the platform-appropriate config file path.
// On Linux:   ~/.config/straw/config.toml
// On macOS:   ~/Library/Application Support/straw/config.toml
// On Windows: %APPDATA%\straw\config.toml
func DefaultConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "straw", "config.toml"), nil
}

// DefaultStateDir returns the platform-appropriate state directory.
// On Linux:   ~/.cache/straw/
// On macOS:   ~/Library/Caches/straw/
// On Windows: %LOCALAPPDATA%\straw\
func DefaultStateDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cacheDir, "straw"), nil
}

// DefaultSocketPath returns the platform-appropriate socket path.
// On Linux:   $XDG_RUNTIME_DIR/straw.sock or state dir fallback
// On macOS:   state dir
// On Windows: %LOCALAPPDATA%\straw\straw.sock or %TEMP%\straw.sock
func DefaultSocketPath() string {
	// Try XDG_RUNTIME_DIR first (Linux)
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir != "" {
		return filepath.Join(runtimeDir, "straw.sock")
	}

	// Fallback to state dir
	stateDir, err := DefaultStateDir()
	if err == nil {
		return filepath.Join(stateDir, "straw.sock")
	}

	// Last resort: platform temp directory
	return filepath.Join(os.TempDir(), "straw.sock")
}

func Load(path string) (*Config, error) {
	if path == "" {
		var err error
		path, err = DefaultConfigPath()
		if err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Expand tilde in all paths before validation
	cfg.expandPaths()

	// Validate critical components (Watch folders)
	if err := cfg.ValidateCritical(); err != nil {
		return nil, err
	}

	// Validate rules, but don't fail the whole load if some are bad.
	// We'll keep them in the config so they can be fixed in the TUI,
	// but the engine should probably be aware they are invalid.
	// For now, we just log (if we had a logger here, but we don't,
	// so the caller should handle it).
	// Actually, let's just let the engine handle invalid rules by ignoring them.

	return &cfg, nil
}

func (c *Config) Save(path string) error {
	if path == "" {
		var err error
		path, err = DefaultConfigPath()
		if err != nil {
			return err
		}
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := toml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// expandPaths normalizes all user-facing paths by expanding tilde prefixes.
// This must be called before Validate() so that os.Stat sees real paths.
func (c *Config) expandPaths() {
	for i := range c.Watch {
		c.Watch[i].Path = pathutil.ExpandPath(c.Watch[i].Path)
	}
	for i := range c.Rules {
		for j := range c.Rules[i].Actions {
			if c.Rules[i].Actions[j].Target != "" {
				c.Rules[i].Actions[j].Target = pathutil.ExpandPath(c.Rules[i].Actions[j].Target)
			}
		}
	}
}

func (c *Config) Validate() error {
	if err := c.ValidateCritical(); err != nil {
		return err
	}
	for i, r := range c.Rules {
		if err := r.Validate(); err != nil {
			return fmt.Errorf("rules[%d] (%s) is invalid: %w", i, r.Name, err)
		}
	}
	return nil
}

func (c *Config) ValidateCritical() error {
	if len(c.Watch) == 0 {
		return errors.New("config must include at least one watch folder")
	}
	for i, w := range c.Watch {
		if w.Path == "" {
			return fmt.Errorf("watch[%d] path is required", i)
		}
		// Check if path exists
		if info, err := os.Stat(w.Path); err != nil {
			return fmt.Errorf("watch[%d] path does not exist: %s", i, w.Path)
		} else if !info.IsDir() {
			return fmt.Errorf("watch[%d] path is not a directory: %s", i, w.Path)
		}
	}
	return nil
}

func (r Rule) Validate() error {
	if r.Name == "" {
		return errors.New("rule name is required")
	}
	if len(r.Actions) == 0 {
		return errors.New("at least one action is required")
	}
	for i, a := range r.Actions {
		if a.Type == "" {
			return fmt.Errorf("action[%d] type is required", i)
		}
		// Basic action validation
		switch a.Type {
		case "move", "copy":
			if a.Target == "" {
				return fmt.Errorf("action[%d] (%s) target is required", i, a.Type)
			}
		case "shell":
			// args or target (command) should be present?
			// Looking at actions.go might help.
		case "trash":
			// no target needed
		default:
			return fmt.Errorf("unknown action type: %s", a.Type)
		}
	}
	return nil
}
