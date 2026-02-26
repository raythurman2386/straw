package main

import (
	"fmt"
	"strings"

	"straw/internal/config"
	"straw/internal/tui"

	"github.com/charmbracelet/huh"
	tea "github.com/charmbracelet/bubbletea"
)

type wizardModel struct {
	form         *huh.Form
	originalName string // If editing, stores the name before any changes
	styles       tui.Styles
	theme        tui.Theme
	width        int
	height       int

	// Rule data captured by the form
	name       string
	match      string
	minAge     string
	maxAge     string
	actionType string
	target     string

	targetInput *huh.Input
}

type wizardFinishedMsg struct {
	rule         *config.Rule
	originalName string // empty if new rule
}

type wizardCancelMsg struct{}

func newWizardModel(styles tui.Styles, theme tui.Theme, ruleToEdit *config.Rule) wizardModel {
	m := wizardModel{
		originalName: "",
		styles:       styles,
		theme:        theme,
		actionType:   "move",
	}

	if ruleToEdit != nil {
		m.originalName = ruleToEdit.Name
		m.name = ruleToEdit.Name

		if ruleToEdit.Match.Extension != "" {
			m.match = ruleToEdit.Match.Extension
		} else if ruleToEdit.Match.Regex != "" {
			m.match = ruleToEdit.Match.Regex
		} else {
			m.match = ruleToEdit.Match.Glob
		}

		if ruleToEdit.Match.MinAgeDays > 0 {
			m.minAge = fmt.Sprintf("%d", ruleToEdit.Match.MinAgeDays)
		}
		if ruleToEdit.Match.MaxAgeDays > 0 {
			m.maxAge = fmt.Sprintf("%d", ruleToEdit.Match.MaxAgeDays)
		}

		if len(ruleToEdit.Actions) > 0 {
			m.actionType = ruleToEdit.Actions[0].Type
			m.target = ruleToEdit.Actions[0].Target
		}
	}

	m.targetInput = huh.NewInput().
		Key("target").
		Title("6. Enter destination or command:").
		Value(&m.target)

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("1. What should this rule be called?").
				Placeholder("e.g. Clean PDFs").
				Value(&m.name).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return fmt.Errorf("name cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Key("match").
				Title("2. What files should it match?").
				Placeholder("e.g. .pdf or *.txt").
				Description("Enter an extension like '.pdf' or a glob like 'data_*.csv'").
				Value(&m.match).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return fmt.Errorf("match pattern cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Key("minAge").
				Title("3. Minimum Age (days)?").
				Placeholder("e.g. 7").
				Description("Leave blank for no minimum age").
				Value(&m.minAge),
			huh.NewInput().
				Key("maxAge").
				Title("4. Maximum Age (days)?").
				Placeholder("e.g. 30").
				Description("Leave blank for no maximum age").
				Value(&m.maxAge),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("actionType").
				Title("5. Select an Action:").
				Options(
					huh.NewOption("MOVE", "move"),
					huh.NewOption("COPY", "copy"),
					huh.NewOption("TRASH", "trash"),
					huh.NewOption("SHELL", "shell"),
				).
				Value(&m.actionType),
		),
		huh.NewGroup(m.targetInput).WithHideFunc(func() bool { return m.actionType == "trash" }),
		huh.NewGroup(
			huh.NewNote().
				Title("7. Confirm and Save?").
				DescriptionFunc(func() string {
					targetStr := m.target
					if m.actionType == "trash" {
						targetStr = "System Trash"
					}
					return fmt.Sprintf("Name: %s\nMatch: %s\nAge: %s to %s days\nAction: %s -> %s",
						m.name, m.match, m.minAge, m.maxAge, strings.ToUpper(m.actionType), targetStr)
				}, &m.actionType),
			huh.NewConfirm().
				Title("Ready to save?").
				Affirmative("Yes, save it!").
				Negative("No, go back"),
		),
	).WithTheme(tui.GetHuhTheme(theme))

	return m
}

func (m wizardModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m wizardModel) Update(msg tea.Msg) (wizardModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.form = m.form.WithWidth(m.width)
	case tea.KeyMsg:
		if msg.String() == "esc" {
			return m, func() tea.Msg { return wizardCancelMsg{} }
		}
	}

	// Update the target input's title/description based on actionType
	if m.actionType != "trash" {
		if m.actionType == "shell" {
			m.targetInput.Title("6. Enter command to run:")
			m.targetInput.Description("Use $FILE for the file path")
		} else {
			m.targetInput.Title("6. Enter destination path:")
			m.targetInput.Description("The directory to move/copy files to")
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted {
		// Explicitly retrieve values from the form state to ensure they are captured
		name := m.form.GetString("name")
		match := m.form.GetString("match")
		minAge := m.form.GetString("minAge")
		maxAge := m.form.GetString("maxAge")
		actionType := m.form.GetString("actionType")
		target := m.form.GetString("target")

		rule := &config.Rule{
			Name:    name,
			Enabled: true,
			Actions: []config.Action{
				{Type: actionType, Target: target},
			},
		}

		if strings.HasPrefix(match, ".") {
			rule.Match.Extension = match
		} else {
			rule.Match.Glob = match
		}

		if minAge != "" {
			_, _ = fmt.Sscanf(minAge, "%d", &rule.Match.MinAgeDays)
		}
		if maxAge != "" {
			_, _ = fmt.Sscanf(maxAge, "%d", &rule.Match.MaxAgeDays)
		}

		return m, func() tea.Msg { return wizardFinishedMsg{rule: rule, originalName: m.originalName} }
	}

	if m.form.State == huh.StateAborted {
		return m, func() tea.Msg { return wizardCancelMsg{} }
	}

	return m, cmd
}

func (m wizardModel) View() string {
	var s strings.Builder

	titleText := "Rule Creation Wizard"
	if m.originalName != "" {
		titleText = "Edit Rule Wizard"
	}
	title := m.styles.ListTitle.Render(titleText)
	s.WriteString(title + "\n\n")

	s.WriteString(m.form.View())

	return s.String()
}
