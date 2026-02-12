package main

import (
	"fmt"
	"io"
	"strings"

	"straw/internal/config"
	"straw/internal/tui"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"
)

type RuleItem struct {
	rule config.Rule
}

func (i RuleItem) Title() string       { return i.rule.Name }
func (i RuleItem) Description() string { return i.rule.Description }
func (i RuleItem) FilterValue() string { return i.rule.Name }

type RuleDelegate struct {
	styles tui.Styles
	theme  tui.Theme
}

func NewRuleDelegate(styles tui.Styles, theme tui.Theme) RuleDelegate {
	return RuleDelegate{
		styles: styles,
		theme:  theme,
	}
}

func (d RuleDelegate) Height() int  { return 3 }
func (d RuleDelegate) Spacing() int { return 1 }

func (d RuleDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d RuleDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	rule, ok := item.(RuleItem)
	if !ok {
		return
	}

	title := rule.Title()
	desc := rule.Description()

	// Create subtitle with match criteria and actions
	var criteria []string
	if rule.rule.Match.Glob != "" {
		criteria = append(criteria, fmt.Sprintf("Glob: %s", rule.rule.Match.Glob))
	}
	if rule.rule.Match.Extension != "" {
		criteria = append(criteria, fmt.Sprintf("Ext: %s", rule.rule.Match.Extension))
	}
	if rule.rule.Match.Regex != "" {
		criteria = append(criteria, fmt.Sprintf("Reg: %s", rule.rule.Match.Regex))
	}

	matchStr := strings.Join(criteria, ", ")

	var actions []string
	for _, a := range rule.rule.Actions {
		actions = append(actions, strings.ToUpper(a.Type))
	}
	actionStr := strings.Join(actions, " -> ")

	// Styles
	titleStyle := d.styles.ListItem
	descStyle := d.styles.ListDim.PaddingLeft(2)
	metaStyle := d.styles.ListDim.PaddingLeft(2).Italic(true)

	if index == m.Index() {
		titleStyle = d.styles.ListSelected
		descStyle = descStyle.Foreground(d.theme.Foreground)
		metaStyle = metaStyle.Foreground(d.theme.Accent)
	}

	status := "●"
	if !rule.rule.Enabled {
		status = "○"
		titleStyle.Foreground(d.theme.Dim)
	}

	fmt.Fprintf(w, "%s %s\n%s\n%s",
		titleStyle.Render(status),
		titleStyle.Render(title),
		descStyle.Render(desc),
		metaStyle.Render(fmt.Sprintf("%s  •  %s", matchStr, actionStr)))
}
