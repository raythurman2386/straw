package main

import (
	"fmt"
	"io"
	"strings"

	"straw/internal/config"
	"straw/internal/tui"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type settingsItem struct {
	title       string
	description string
	value       string
	itemType    string // theme, etc.
}

func (i settingsItem) Title() string       { return i.title }
func (i settingsItem) Description() string { return i.description }
func (i settingsItem) FilterValue() string { return i.title }

type settingsDelegate struct {
	styles tui.Styles
	theme  tui.Theme
}

func (d settingsDelegate) Height() int  { return 2 }
func (d settingsDelegate) Spacing() int { return 1 }
func (d settingsDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
func (d settingsDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(settingsItem)
	if !ok {
		return
	}

	titleStyle := d.styles.ListItem
	descStyle := d.styles.ListDim.PaddingLeft(2)

	if index == m.Index() {
		titleStyle = d.styles.ListSelected
		descStyle = descStyle.Foreground(d.theme.Foreground)

		// If it's a theme, we can show the accent color
		if i.itemType == "theme" {
			titleStyle = titleStyle.Foreground(d.theme.Accent)
		}
	}

	fmt.Fprintf(w, "%s\n%s", titleStyle.Render(i.title), descStyle.Render(i.description))
}

type settingsModel struct {
	list       list.Model
	config     *config.Config
	configPath string
	styles     tui.Styles
	theme      tui.Theme
	updated    bool
	message    string
	width      int
	height     int
}

type settingsUpdatedMsg struct {
	config *config.Config
}

func newSettingsModel(cfg *config.Config, configPath string, styles tui.Styles, theme tui.Theme) settingsModel {
	var items []list.Item
	for _, t := range tui.AllThemes() {
		indicator := " [ ] "
		if cfg != nil && cfg.TUI.Theme == strings.ToLower(t.Name) {
			indicator = " [✓] "
		}
		items = append(items, settingsItem{
			title:       fmt.Sprintf("%sTheme: %s", indicator, t.Name),
			description: fmt.Sprintf("Set application color scheme to %s", t.Name),
			value:       strings.ToLower(t.Name),
			itemType:    "theme",
		})
	}

	d := settingsDelegate{styles: styles, theme: theme}
	l := list.New(items, d, 0, 0)
	l.Title = "Application Settings"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = styles.ListTitle
	l.SetShowHelp(false)

	return settingsModel{
		list:       l,
		config:     cfg,
		configPath: configPath,
		styles:     styles,
		theme:      theme,
	}
}

func (m settingsModel) Update(msg tea.Msg) (settingsModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateListSize()

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			item := m.list.SelectedItem().(settingsItem)
			if item.itemType == "theme" {
				m.config.TUI.Theme = item.value
				err := m.config.Save(m.configPath)
				if err != nil {
					m.message = fmt.Sprintf("Error saving config: %v", err)
				} else {
					m.message = fmt.Sprintf("Theme updated to %s", item.value)
					m.updated = true

					// Refresh items to update indicators
					var items []list.Item
					for _, t := range tui.AllThemes() {
						indicator := " [ ] "
						if m.config.TUI.Theme == strings.ToLower(t.Name) {
							indicator = " [✓] "
						}
						items = append(items, settingsItem{
							title:       fmt.Sprintf("%sTheme: %s", indicator, t.Name),
							description: fmt.Sprintf("Set application color scheme to %s", t.Name),
							value:       strings.ToLower(t.Name),
							itemType:    "theme",
						})
					}
					m.list.SetItems(items)
					m.updateListSize()

					return m, func() tea.Msg { return settingsUpdatedMsg{config: m.config} }
				}
			}
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *settingsModel) updateListSize() {
	h := m.height
	if m.message != "" {
		h -= 2 // Account for message and spacing
	}
	m.list.SetSize(m.width, h)
}

func (m settingsModel) View() string {
	if m.message == "" {
		return m.list.View()
	}

	color := m.theme.Success
	if strings.HasPrefix(m.message, "Error") {
		color = m.theme.Error
	}
	msgStr := lipgloss.NewStyle().
		Foreground(color).
		Padding(0, 2).
		Italic(true).
		Render(m.message)

	return lipgloss.JoinVertical(lipgloss.Left,
		m.list.View(),
		msgStr,
	)
}
