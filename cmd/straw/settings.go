package main

import (
	"strings"

	"straw/internal/config"
	"straw/internal/tui"

	"github.com/charmbracelet/huh"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type settingsModel struct {
	form       *huh.Form
	config     *config.Config
	configPath string
	styles     tui.Styles
	theme      tui.Theme
	updated    bool
	message    string
	width      int
	height     int

	selectedTheme string
}

type settingsUpdatedMsg struct {
	config *config.Config
}

func newSettingsModel(cfg *config.Config, configPath string, styles tui.Styles, theme tui.Theme) settingsModel {
	m := settingsModel{
		config:     cfg,
		configPath: configPath,
		styles:     styles,
		theme:      theme,
	}

	if cfg != nil {
		m.selectedTheme = cfg.TUI.Theme
	}

	var options []huh.Option[string]
	for _, t := range tui.AllThemes() {
		options = append(options, huh.NewOption(t.Name, strings.ToLower(t.Name)))
	}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("theme").
				Title("Appearance").
				Description("Choose your preferred color scheme").
				Options(options...).
				Value(&m.selectedTheme),
		),
	).WithTheme(tui.GetHuhTheme(theme))

	return m
}

func (m settingsModel) Update(msg tea.Msg) (settingsModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.form = m.form.WithWidth(m.width)
	}

	prevTheme := m.selectedTheme
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	// If form completed (user hit enter on the last/only field), reset it to keep it visible
	if m.form.State == huh.StateCompleted {
		m.form.State = huh.StateNormal
	}

	// Explicitly get the theme from the form if it changed
	m.selectedTheme = m.form.GetString("theme")

	// If theme changed, save it
	if m.selectedTheme != "" && m.selectedTheme != prevTheme && m.config != nil {
		m.config.TUI.Theme = m.selectedTheme
		_ = m.config.Save(m.configPath)
		m.updated = true
		return m, func() tea.Msg { return settingsUpdatedMsg{config: m.config} }
	}

	return m, cmd
}

func (m settingsModel) View() string {
	var s strings.Builder

	s.WriteString(m.styles.ListTitle.Render("Application Settings") + "\n\n")
	s.WriteString(m.form.View())

	if m.message != "" {
		color := m.theme.Success
		if strings.HasPrefix(m.message, "Error") {
			color = m.theme.Error
		}
		s.WriteString("\n\n" + lipgloss.NewStyle().Foreground(color).Padding(0, 2).Italic(true).Render(m.message))
	}

	return s.String()
}

func (m settingsModel) Init() tea.Cmd {
	return m.form.Init()
}
