package rules

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"straw/internal/config"
	"straw/internal/watcher"
)

func TestEngine_Evaluate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "rules_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	rules := []config.Rule{
		{
			Name:    "Match Text",
			Enabled: true,
			Match: config.Match{
				Extension: ".txt",
			},
			Actions: []config.Action{
				{Type: "move", Target: "/tmp/archive"},
			},
		},
		{
			Name:    "Match Large",
			Enabled: true,
			Match: config.Match{
				MinSize: 100,
			},
			Actions: []config.Action{
				{Type: "trash"},
			},
		},
	}

	engine := NewEngine(rules)

	t.Run("Match extension", func(t *testing.T) {
		event := watcher.Event{Path: testFile, Type: watcher.Create}
		actions := engine.Evaluate(event)
		if len(actions) != 1 {
			t.Errorf("expected 1 action, got %d", len(actions))
		}
		if actions[0].Type != "move" {
			t.Errorf("expected move action, got %s", actions[0].Type)
		}
	})

	t.Run("No match size", func(t *testing.T) {
		event := watcher.Event{Path: testFile, Type: watcher.Create}
		// File is 5 bytes, rule requires 100
		actions := engine.Evaluate(event)
		// Only the first rule should match
		foundTrash := false
		for _, a := range actions {
			if a.Type == "trash" {
				foundTrash = true
			}
		}
		if foundTrash {
			t.Error("expected no trash action for small file")
		}
	})

	t.Run("Match age", func(t *testing.T) {
		oldFile := filepath.Join(tmpDir, "old.log")
		if err := os.WriteFile(oldFile, []byte("old"), 0644); err != nil {
			t.Fatal(err)
		}
		// Set mod time to 10 days ago
		oldTime := time.Now().AddDate(0, 0, -10)
		if err := os.Chtimes(oldFile, oldTime, oldTime); err != nil {
			t.Fatal(err)
		}

		engine.SetRules([]config.Rule{
			{
				Name:    "Old logs",
				Enabled: true,
				Match: config.Match{
					MinAgeDays: 7,
					Extension:  ".log",
				},
				Actions: []config.Action{{Type: "trash"}},
			},
		})

		event := watcher.Event{Path: oldFile, Type: watcher.Create}
		actions := engine.Evaluate(event)
		if len(actions) != 1 {
			t.Errorf("expected 1 action for old file, got %d", len(actions))
		}
	})
}
