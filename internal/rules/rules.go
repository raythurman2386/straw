package rules

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"straw/internal/config"
	"straw/internal/pathutil"
	"straw/internal/watcher"

	"github.com/bmatcuk/doublestar/v4"
)

type Engine struct {
	rules []config.Rule
	mu    sync.RWMutex
}

func NewEngine(rules []config.Rule) *Engine {
	return &Engine{rules: rules}
}

func (e *Engine) SetRules(rules []config.Rule) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules = rules
}

// Evaluate checks if the event matches any rules and returns the actions to take.
func (e *Engine) Evaluate(event watcher.Event) []config.Action {
	var actions []config.Action

	// We generally only act on existing files.
	// If the file was removed, we can't inspect it.
	if event.Type == watcher.Remove {
		return nil
	}

	info, err := os.Stat(event.Path)
	if err != nil {
		return nil
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, rule := range e.rules {
		if !rule.Enabled {
			continue
		}

		// Skip invalid rules
		if err := rule.Validate(); err != nil {
			continue
		}

		if e.matches(rule.Match, event.Path, info) {
			// If the rule's primary move/copy outcome is already satisfied,
			// skip it to prevent infinite loops and duplicate shell notifications.
			if e.isRuleSatisfied(rule, event.Path) {
				continue
			}

			// Append actions from this rule
			// Currently we accumulate all actions from all matching rules.
			actions = append(actions, rule.Actions...)
		}
	}
	return actions
}

// isRuleSatisfied checks if a file is already in the target directory of the rule's move or copy action.
func (e *Engine) isRuleSatisfied(rule config.Rule, path string) bool {
	fileDir := filepath.Clean(filepath.Dir(path))

	for _, action := range rule.Actions {
		if action.Type == "move" || action.Type == "copy" {
			target := pathutil.ExpandPath(action.Target)
			target = filepath.Clean(target)
			if fileDir == target {
				return true
			}
		}
	}
	return false
}

func (e *Engine) matches(m config.Match, path string, info os.FileInfo) bool {
	name := filepath.Base(path)

	// Glob
	if m.Glob != "" {
		// If glob contains a path separator, match against full path
		target := name
		if strings.Contains(m.Glob, "/") || strings.Contains(m.Glob, string(filepath.Separator)) {
			target = path
		}

		matched, _ := doublestar.Match(m.Glob, target)
		if !matched {
			return false
		}
	}

	// Extension
	if m.Extension != "" {
		if filepath.Ext(name) != m.Extension {
			return false
		}
	}

	// Regex
	if m.Regex != "" {
		matched, _ := regexp.MatchString(m.Regex, name)
		if !matched {
			return false
		}
	}

	// Size
	if m.MinSize > 0 && info.Size() < m.MinSize {
		return false
	}
	if m.MaxSize > 0 && info.Size() > m.MaxSize {
		return false
	}

	// Age (Days)
	if m.MinAgeDays > 0 {
		age := time.Since(info.ModTime()).Hours() / 24
		if age < float64(m.MinAgeDays) {
			return false
		}
	}
	if m.MaxAgeDays > 0 {
		age := time.Since(info.ModTime()).Hours() / 24
		if age > float64(m.MaxAgeDays) {
			return false
		}
	}

	// FileType
	if m.FileType != "" {
		if m.FileType == "file" && info.IsDir() {
			return false
		}
		if m.FileType == "dir" && !info.IsDir() {
			return false
		}
	}

	// Hidden
	if m.Hidden != nil {
		if *m.Hidden != isHidden(name, path) {
			return false
		}
	}

	return true
}
