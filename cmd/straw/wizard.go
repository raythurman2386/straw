package main

import (
	"fmt"
	"strings"

	"straw/internal/config"
	"straw/internal/tui"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type wizardStep int

const (
	stepName wizardStep = iota
	stepMatch
	stepMinAge
	stepMaxAge
	stepActionType
	stepActionTarget
	stepConfirm
)

type wizardModel struct {
	step         wizardStep
	originalName string // If editing, stores the name before any changes
	nameInput    textinput.Model
	matchInput   textinput.Model
	minAgeInput  textinput.Model
	maxAgeInput  textinput.Model
	targetInput  textinput.Model
	actionType   string // move, copy, trash, shell
	styles       tui.Styles
	theme        tui.Theme
	width        int
	height       int
}

type wizardFinishedMsg struct {
	rule         *config.Rule
	originalName string // empty if new rule
}

type wizardCancelMsg struct{}

func newWizardModel(styles tui.Styles, theme tui.Theme, ruleToEdit *config.Rule) wizardModel {
	ni := textinput.New()
	ni.Placeholder = "Enter rule name (e.g. Clean PDFs)"

	mi := textinput.New()
	mi.Placeholder = "Enter extension (e.g. .pdf) or glob (e.g. *.txt)"

	minAi := textinput.New()
	minAi.Placeholder = "Min age in days (e.g. 7) - leave empty to skip"

	maxAi := textinput.New()
	maxAi.Placeholder = "Max age in days (e.g. 30) - leave empty to skip"

	ti := textinput.New()
	ti.Placeholder = "Enter destination directory path"

	actionType := "move"
	originalName := ""

	if ruleToEdit != nil {
		originalName = ruleToEdit.Name
		ni.SetValue(ruleToEdit.Name)

		if ruleToEdit.Match.Extension != "" {
			mi.SetValue(ruleToEdit.Match.Extension)
		} else if ruleToEdit.Match.Regex != "" {
			mi.SetValue(ruleToEdit.Match.Regex)
		} else {
			mi.SetValue(ruleToEdit.Match.Glob)
		}

		if ruleToEdit.Match.MinAgeDays > 0 {
			minAi.SetValue(fmt.Sprintf("%d", ruleToEdit.Match.MinAgeDays))
		}
		if ruleToEdit.Match.MaxAgeDays > 0 {
			maxAi.SetValue(fmt.Sprintf("%d", ruleToEdit.Match.MaxAgeDays))
		}

		if len(ruleToEdit.Actions) > 0 {
			actionType = ruleToEdit.Actions[0].Type
			ti.SetValue(ruleToEdit.Actions[0].Target)
		}
	}

	ni.Focus()

	return wizardModel{
		step:         stepName,
		originalName: originalName,
		nameInput:    ni,
		matchInput:   mi,
		minAgeInput:  minAi,
		maxAgeInput:  maxAi,
		targetInput:  ti,
		actionType:   actionType,
		styles:       styles,
		theme:        theme,
	}
}

func (m wizardModel) Update(msg tea.Msg) (wizardModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.nameInput.Width = m.width - 4
		m.matchInput.Width = m.width - 4
		m.minAgeInput.Width = m.width - 4
		m.maxAgeInput.Width = m.width - 4
		m.targetInput.Width = m.width - 4

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m.nextStep()
		case "esc":
			return m, func() tea.Msg { return wizardCancelMsg{} }
		case "m":
			if m.step == stepActionType {
				m.actionType = "move"
				return m.nextStep()
			}
		case "c":
			if m.step == stepActionType {
				m.actionType = "copy"
				return m.nextStep()
			}
		case "t":
			if m.step == stepActionType {
				m.actionType = "trash"
				return m.nextStep()
			}
		case "s":
			if m.step == stepActionType {
				m.actionType = "shell"
				return m.nextStep()
			}
		}
	}

	switch m.step {
	case stepName:
		m.nameInput, cmd = m.nameInput.Update(msg)
	case stepMatch:
		m.matchInput, cmd = m.matchInput.Update(msg)
	case stepMinAge:
		m.minAgeInput, cmd = m.minAgeInput.Update(msg)
	case stepMaxAge:
		m.maxAgeInput, cmd = m.maxAgeInput.Update(msg)
	case stepActionTarget:
		m.targetInput, cmd = m.targetInput.Update(msg)
	}

	return m, cmd
}

func (m wizardModel) nextStep() (wizardModel, tea.Cmd) {
	switch m.step {
	case stepName:
		if m.nameInput.Value() == "" {
			return m, nil
		}
		m.step = stepMatch
		m.matchInput.Focus()
	case stepMatch:
		if m.matchInput.Value() == "" {
			return m, nil
		}
		m.step = stepMinAge
		m.minAgeInput.Focus()
	case stepMinAge:
		m.step = stepMaxAge
		m.maxAgeInput.Focus()
	case stepMaxAge:
		m.step = stepActionType
	case stepActionType:
		if m.actionType == "trash" {
			m.step = stepConfirm
		} else {
			m.step = stepActionTarget
			if m.actionType == "shell" {
				m.targetInput.Placeholder = "Enter command (use $FILE for path)"
			} else {
				m.targetInput.Placeholder = "Enter destination directory path"
			}
			m.targetInput.Focus()
		}
	case stepActionTarget:
		if m.targetInput.Value() == "" {
			return m, nil
		}
		m.step = stepConfirm
	case stepConfirm:
		rule := &config.Rule{
			Name:    m.nameInput.Value(),
			Enabled: true,
			Actions: []config.Action{
				{Type: m.actionType, Target: m.targetInput.Value()},
			},
		}

		matchVal := m.matchInput.Value()
		if strings.HasPrefix(matchVal, ".") {
			rule.Match.Extension = matchVal
		} else {
			rule.Match.Glob = matchVal
		}

		if m.minAgeInput.Value() != "" {
			_, _ = fmt.Sscanf(m.minAgeInput.Value(), "%d", &rule.Match.MinAgeDays)
		}
		if m.maxAgeInput.Value() != "" {
			_, _ = fmt.Sscanf(m.maxAgeInput.Value(), "%d", &rule.Match.MaxAgeDays)
		}

		return m, func() tea.Msg { return wizardFinishedMsg{rule: rule, originalName: m.originalName} }
	}
	return m, nil
}

func (m wizardModel) View() string {
	var s strings.Builder

	titleText := "Rule Creation Wizard"
	if m.originalName != "" {
		titleText = "Edit Rule Wizard"
	}
	title := m.styles.ListTitle.Render(titleText)
	s.WriteString(title + "\n\n")

	successStyle := m.styles.ListSelected.Border(lipgloss.HiddenBorder())

	switch m.step {
	case stepName:
		s.WriteString("1. What should this rule be called?\n")
		s.WriteString(m.nameInput.View())
	case stepMatch:
		s.WriteString(fmt.Sprintf("Rule: %s\n\n", successStyle.Render(m.nameInput.Value())))
		s.WriteString("2. What files should it match?\n")
		s.WriteString(m.styles.LogMessage.Render("Enter an extension like '.pdf' or a glob like 'data_*.csv'"))
		s.WriteString("\n" + m.matchInput.View())
	case stepMinAge:
		s.WriteString(fmt.Sprintf("Rule: %s\n", successStyle.Render(m.nameInput.Value())))
		s.WriteString(fmt.Sprintf("Match: %s\n\n", successStyle.Render(m.matchInput.Value())))
		s.WriteString("3. Minimum Age (days)?\n")
		s.WriteString(m.styles.LogMessage.Render("Leave blank for no minimum age"))
		s.WriteString("\n" + m.minAgeInput.View())
	case stepMaxAge:
		s.WriteString(fmt.Sprintf("Rule: %s\n", successStyle.Render(m.nameInput.Value())))
		s.WriteString(fmt.Sprintf("Match: %s\n", successStyle.Render(m.matchInput.Value())))
		s.WriteString(fmt.Sprintf("Min Age: %s days\n\n", successStyle.Render(m.minAgeInput.Value())))
		s.WriteString("4. Maximum Age (days)?\n")
		s.WriteString(m.styles.LogMessage.Render("Leave blank for no maximum age"))
		s.WriteString("\n" + m.maxAgeInput.View())
	case stepActionType:
		s.WriteString(fmt.Sprintf("Rule: %s\n", successStyle.Render(m.nameInput.Value())))
		s.WriteString(fmt.Sprintf("Match: %s\n", successStyle.Render(m.matchInput.Value())))
		s.WriteString(fmt.Sprintf("Age: %s to %s days\n\n", successStyle.Render(m.minAgeInput.Value()), successStyle.Render(m.maxAgeInput.Value())))
		s.WriteString("5. Select an Action:\n\n")

		s.WriteString(m.renderActionOption("move", "m"))
		s.WriteString(m.renderActionOption("copy", "c"))
		s.WriteString(m.renderActionOption("trash", "t"))
		s.WriteString(m.renderActionOption("shell", "s"))

		s.WriteString("\n" + m.styles.LogMessage.Render("Press the letter to select and continue"))
	case stepActionTarget:
		s.WriteString(fmt.Sprintf("Rule: %s\n", successStyle.Render(m.nameInput.Value())))
		s.WriteString(fmt.Sprintf("Match: %s\n", successStyle.Render(m.matchInput.Value())))

		actionColor := m.theme.ActionColor(m.actionType)
		actionStr := lipgloss.NewStyle().Foreground(actionColor).Bold(true).Render(strings.ToUpper(m.actionType))
		s.WriteString(fmt.Sprintf("Action: %s\n\n", actionStr))

		if m.actionType == "shell" {
			s.WriteString("6. Enter command to run:\n")
		} else {
			s.WriteString("6. Enter destination path:\n")
		}
		s.WriteString(m.targetInput.View())
	case stepConfirm:
		s.WriteString("7. Confirm and Save?\n\n")
		s.WriteString(fmt.Sprintf(" Name:    %s\n", successStyle.Render(m.nameInput.Value())))
		s.WriteString(fmt.Sprintf(" Match:   %s\n", successStyle.Render(m.matchInput.Value())))
		s.WriteString(fmt.Sprintf(" Age:    %s to %s days\n", successStyle.Render(m.minAgeInput.Value()), successStyle.Render(m.maxAgeInput.Value())))

		actionColor := m.theme.ActionColor(m.actionType)
		actionStr := lipgloss.NewStyle().Foreground(actionColor).Bold(true).Render(strings.ToUpper(m.actionType))
		targetStr := m.targetInput.Value()
		if m.actionType == "trash" {
			targetStr = "System Trash"
		}
		s.WriteString(fmt.Sprintf(" Action:  %s -> %s\n", actionStr, targetStr))
		s.WriteString("\nPress ENTER to Save or ESC to Cancel")
	}

	s.WriteString("\n\n" + m.styles.Desc.Render("ENTER: Next • ESC: Cancel"))

	return s.String()
}

func (m wizardModel) renderActionOption(action, key string) string {
	color := m.theme.ActionColor(action)
	style := lipgloss.NewStyle().Foreground(color)

	if m.actionType == action {
		style = style.Bold(true).Underline(true)
		bullet := lipgloss.NewStyle().Foreground(m.theme.Success).Render("●")
		return fmt.Sprintf(" %s [%s] %s (Selected)\n", bullet, m.styles.Key.Render(key), style.Render(strings.ToUpper(action)))
	}

	return fmt.Sprintf(" ○ [%s] %s\n", m.styles.Key.Render(key), style.Render(strings.ToUpper(action)))
}
