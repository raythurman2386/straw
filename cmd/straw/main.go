package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"straw/internal/config"
	"straw/internal/ipc"
	"straw/internal/logging"
	"straw/internal/tui"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type tab int

const (
	tabActivity tab = iota
	tabRules
	tabCreate
	tabLogFile
	tabSettings
)

type model struct {
	client          *ipc.Client
	connected       bool
	err             error
	socketPath      string
	configPath      string
	config          *config.Config
	logs            []string
	viewport        viewport.Model
	logFileViewport viewport.Model
	rulesList       list.Model
	wizard          wizardModel
	settings        settingsModel
	activeTab       tab
	ready           bool
	theme           tui.Theme
	styles          tui.Styles
	width           int
	height          int
	logFilePath     string

	// Activity Stats
	statsProcessed int
	statsErrors    int
	lastActivity   time.Time
}

func initialModel(socketPath string, configPath string, cfg *config.Config, logFilePath string) model {
	themeName := "everforest"
	if cfg != nil && cfg.TUI.Theme != "" {
		themeName = cfg.TUI.Theme
	}
	theme := tui.GetTheme(themeName)
	styles := tui.GetStyles(theme)

	delegate := NewRuleDelegate(styles, theme)
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Active Rules"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = styles.ListTitle
	l.Styles.HelpStyle = styles.Desc
	l.Styles.NoItems = styles.ListDim.MarginLeft(2)
	l.Styles.PaginationStyle = styles.ListDim
	// Remove default border
	l.Styles.TitleBar = lipgloss.NewStyle()
	l.SetShowHelp(false)

	return model{
		socketPath:  socketPath,
		configPath:  configPath,
		config:      cfg,
		logFilePath: logFilePath,
		logs:        []string{},
		theme:       theme,
		styles:      styles,
		rulesList:   l,
		wizard:      newWizardModel(styles, theme, nil),
		settings:    newSettingsModel(cfg, configPath, styles, theme),
		activeTab:   tabActivity,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		connectCmd(m.socketPath),
		tea.EnterAltScreen,
	)
}

type connectedMsg *ipc.Client
type eventMsg ipc.Event
type rulesMsg []config.Rule
type errorMsg error
type logFileMsg string

func connectCmd(socketPath string) tea.Cmd {
	return func() tea.Msg {
		client := ipc.NewClient(socketPath)
		if err := client.Connect(); err != nil {
			return errorMsg(err)
		}
		return connectedMsg(client)
	}
}

func readLogFileCmd(path string) tea.Cmd {
	return func() tea.Msg {
		if path == "" {
			return logFileMsg("Log file path not configured")
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return logFileMsg(fmt.Sprintf("Error reading log file: %v", err))
		}
		return logFileMsg(string(content))
	}
}

func fetchRulesCmd(client *ipc.Client) tea.Cmd {
	return func() tea.Msg {
		res, err := client.Call(ipc.MethodGetRules, nil)
		if err != nil {
			return errorMsg(err)
		}
		var rules []config.Rule
		if err := json.Unmarshal(res, &rules); err != nil {
			return errorMsg(err)
		}
		return rulesMsg(rules)
	}
}

func reloadDaemonCmd(client *ipc.Client) tea.Cmd {
	return func() tea.Msg {
		_, err := client.Call(ipc.MethodTriggerReload, nil)
		if err != nil {
			return errorMsg(err)
		}
		// After triggering reload, wait a bit and fetch rules again
		time.Sleep(100 * time.Millisecond)
		return fetchRulesCmd(client)()
	}
}

func addRuleCmd(client *ipc.Client, rule *config.Rule) tea.Cmd {
	return func() tea.Msg {
		_, err := client.Call(ipc.MethodAddRule, rule)
		if err != nil {
			return errorMsg(err)
		}
		return fetchRulesCmd(client)()
	}
}

func updateRuleCmd(client *ipc.Client, originalName string, rule *config.Rule) tea.Cmd {
	return func() tea.Msg {
		params := ipc.UpdateRuleParams{
			OriginalName: originalName,
			Rule:         *rule,
		}
		_, err := client.Call(ipc.MethodUpdateRule, params)
		if err != nil {
			return errorMsg(err)
		}
		return fetchRulesCmd(client)()
	}
}

func waitForEvent(client *ipc.Client) tea.Cmd {
	return func() tea.Msg {
		event := <-client.Events()
		return eventMsg(event)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.activeTab != tabCreate {
			switch msg.String() {
			case "ctrl+c", "q":
				if m.client != nil {
					m.client.Close()
				}
				return m, tea.Quit
			case "1":
				m.activeTab = tabActivity
				return m, nil
			case "2":
				m.activeTab = tabRules
				if m.connected {
					cmds = append(cmds, fetchRulesCmd(m.client))
				}
				return m, tea.Batch(cmds...)
			case "3":
				m.activeTab = tabCreate
				m.wizard = newWizardModel(m.styles, m.theme, nil)
				return m, m.wizard.Init()
			case "4":
				m.activeTab = tabLogFile
				return m, readLogFileCmd(m.logFilePath)
			case "5":
				m.activeTab = tabSettings
				return m, m.settings.Init()
			case "r":
				if m.connected {
					m.addSystemLog("Triggering daemon reload...", true)
					cmds = append(cmds, reloadDaemonCmd(m.client))
				}
				return m, tea.Batch(cmds...)
			case "e":
				if m.activeTab == tabRules && m.rulesList.SelectedItem() != nil {
					ruleItem := m.rulesList.SelectedItem().(RuleItem)
					m.activeTab = tabCreate
					m.wizard = newWizardModel(m.styles, m.theme, &ruleItem.rule)
					return m, m.wizard.Init()
				}
			}
		} else {
			if msg.String() == "ctrl+c" {
				if m.client != nil {
					m.client.Close()
				}
				return m, tea.Quit
			}
		}

	case wizardCancelMsg:
		m.activeTab = tabRules

	case wizardFinishedMsg:
		m.activeTab = tabRules
		if m.connected {
			if msg.originalName != "" {
				cmds = append(cmds, updateRuleCmd(m.client, msg.originalName, msg.rule))
			} else {
				cmds = append(cmds, addRuleCmd(m.client, msg.rule))
			}
		}

	case settingsUpdatedMsg:
		m.config = msg.config
		m.theme = tui.GetTheme(m.config.TUI.Theme)
		m.styles = tui.GetStyles(m.theme)

		// Update rules list styles
		delegate := NewRuleDelegate(m.styles, m.theme)
		m.rulesList.SetDelegate(delegate)
		m.rulesList.Styles.Title = m.styles.ListTitle
		m.rulesList.Styles.HelpStyle = m.styles.Desc

		// Update wizard styles
		m.wizard.styles = m.styles
		m.wizard.theme = m.theme
		m.wizard.form = m.wizard.form.WithTheme(tui.GetHuhTheme(m.theme))

		// Update settings styles
		m.settings.styles = m.styles
		m.settings.theme = m.theme
		m.settings.form = m.settings.form.WithTheme(tui.GetHuhTheme(m.theme))

		m.addSystemLog(fmt.Sprintf("Theme updated to %s", m.config.TUI.Theme), true)

	case logFileMsg:
		m.logFileViewport.SetContent(string(msg))
		m.logFileViewport.GotoBottom()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		headerHeight := 4
		footerHeight := 3
		verticalMargin := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width-6, msg.Height-verticalMargin-2)
			m.viewport.YPosition = headerHeight

			m.logFileViewport = viewport.New(msg.Width-6, msg.Height-verticalMargin-2)
			m.logFileViewport.YPosition = headerHeight

			m.ready = true
		} else {
			m.viewport.Width = msg.Width - 6
			m.viewport.Height = msg.Height - verticalMargin - 2

			m.logFileViewport.Width = msg.Width - 6
			m.logFileViewport.Height = msg.Height - verticalMargin - 2
		}
		m.rulesList.SetSize(msg.Width-6, msg.Height-verticalMargin-2)

		// Send WindowSizeMsgs to sub-models so they can manage their own layout
		subModelMsg := tea.WindowSizeMsg{
			Width:  msg.Width - 6,
			Height: msg.Height - verticalMargin - 2,
		}
		m.settings, _ = m.settings.Update(subModelMsg)
		m.wizard, _ = m.wizard.Update(subModelMsg)

	case connectedMsg:
		m.connected = true
		m.client = msg
		m.addSystemLog("Connected to daemon", true)
		cmds = append(cmds, waitForEvent(m.client), fetchRulesCmd(m.client))

	case rulesMsg:
		items := make([]list.Item, len(msg))
		for i, r := range msg {
			items[i] = RuleItem{rule: r}
		}
		m.rulesList.SetItems(items)

	case errorMsg:
		m.err = msg
		m.addSystemLog(fmt.Sprintf("Error: %v", msg), false)

	case eventMsg:
		if msg.Type == ipc.EventNotification {
			m.handleNotification(msg.Payload)
		} else {
			// Fallback
			p, _ := json.Marshal(msg.Payload)
			m.addSystemLog(fmt.Sprintf("%s: %s", msg.Type, string(p)), false)
		}
		cmds = append(cmds, waitForEvent(m.client))
	}

	if m.activeTab == tabRules {
		m.rulesList, cmd = m.rulesList.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.activeTab == tabCreate {
		m.wizard, cmd = m.wizard.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.activeTab == tabLogFile {
		m.logFileViewport, cmd = m.logFileViewport.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.activeTab == tabSettings {
		m.settings, cmd = m.settings.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) handleNotification(payloadBytes []byte) {
	var payload struct {
		File   string `json:"file"`
		Action string `json:"action"`
		Status string `json:"status"`
		Error  string `json:"error"`
	}

	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		m.addSystemLog("Invalid event payload", false)
		return
	}

	m.statsProcessed++
	m.lastActivity = time.Now()
	if payload.Status == "error" {
		m.statsErrors++
	}

	timestamp := m.styles.LogTime.Render(time.Now().Format("15:04:05"))

	actionColor := m.theme.ActionColor(payload.Action)
	if payload.Status == "error" {
		actionColor = m.theme.Error
	}

	caser := cases.Title(language.English)
	actionStr := m.styles.LogAction.
		Foreground(actionColor).
		Render(caser.String(payload.Action))

	fileName := filepath.Base(payload.File)
	fileStr := m.styles.LogFile.Render(fileName)

	msgStr := ""
	if payload.Status == "error" {
		msgStr = m.styles.LogMessage.Foreground(m.theme.Error).Render(payload.Error)
	} else {
		dir := filepath.Dir(payload.File)
		if len(dir) > 20 {
			dir = "..." + dir[len(dir)-20:]
		}
		msgStr = m.styles.LogMessage.Render(dir)
	}

	line := lipgloss.JoinHorizontal(lipgloss.Left, timestamp, actionStr, fileStr, msgStr)
	m.addRawLog(line)
}

func (m *model) addSystemLog(msg string, success bool) {
	timestamp := m.styles.LogTime.Render(time.Now().Format("15:04:05"))

	color := m.theme.Warning
	if success {
		color = m.theme.Success
	} else if strings.HasPrefix(msg, "Error") {
		color = m.theme.Error
	}

	actionStr := m.styles.LogAction.
		Foreground(color).
		Render("SYSTEM")

	msgStr := m.styles.LogMessage.Render(msg)

	line := lipgloss.JoinHorizontal(lipgloss.Left, timestamp, actionStr, msgStr)
	m.addRawLog(line)
}

func (m *model) addRawLog(line string) {
	m.logs = append(m.logs, line)
	m.viewport.SetContent(strings.Join(m.logs, "\n"))
	m.viewport.GotoBottom()
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	// 1. Header
	title := m.styles.HeaderTitle.Render("STRAW")

	statusColor := m.theme.Error
	statusText := "OFFLINE"
	if m.connected {
		statusColor = m.theme.Success
		statusText = "ONLINE"
	}
	status := m.styles.HeaderStatus.Foreground(statusColor).Render(statusText)

	// Justify header: [TITLE STATUS]
	headerContent := lipgloss.JoinHorizontal(lipgloss.Center, title, status)
	header := m.styles.Header.Width(m.width - 4).Render(headerContent)

	// 2. Content Container
	contentStyle := m.styles.LogContainer.
		Width(m.width - 4).
		Height(m.height - 10)

	var content string
	if m.activeTab == tabRules {
		content = contentStyle.Padding(0, 1).Render(m.rulesList.View())
	} else if m.activeTab == tabCreate {
		content = contentStyle.Padding(1, 2).Render(m.wizard.View())
	} else if m.activeTab == tabLogFile {
		content = contentStyle.Padding(0, 1).Render(m.logFileViewport.View())
	} else if m.activeTab == tabSettings {
		content = contentStyle.Padding(1, 2).Render(m.settings.View())
	} else {
		content = m.renderActivity(contentStyle)
	}

	// 3. Footer
	keys := []string{
		m.renderKey("1", "Activity", m.activeTab == tabActivity),
		m.renderKey("2", "Rules", m.activeTab == tabRules),
		m.renderKey("3", "New", m.activeTab == tabCreate),
		m.renderKey("4", "Logs", m.activeTab == tabLogFile),
		m.renderKey("5", "Settings", m.activeTab == tabSettings),
		m.renderKey("r", "Reload", false),
	}

	if m.activeTab == tabRules {
		keys = append(keys, m.renderKey("e", "Edit", false))
	}

	keys = append(keys, m.renderKey("q", "Quit", false))

	footer := m.styles.Footer.Width(m.width - 4).Render(strings.Join(keys, "  "))

	return m.styles.App.Render(lipgloss.JoinVertical(lipgloss.Left,
		header,
		content,
		footer,
	))
}

func (m model) renderActivity(style lipgloss.Style) string {
	// Dashboard
	statsStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Tertiary).
		Padding(0, 2).
		MarginBottom(1).
		Width(style.GetWidth() - 2)

	labelStyle := m.styles.ListDim.Bold(true)
	valStyle := lipgloss.NewStyle().Foreground(m.theme.Accent).Bold(true)
	errStyle := lipgloss.NewStyle().Foreground(m.theme.Error).Bold(true)

	lastActive := "Never"
	if !m.lastActivity.IsZero() {
		lastActive = m.lastActivity.Format("15:04:05")
	}

	dashboard := statsStyle.Render(
		lipgloss.JoinHorizontal(lipgloss.Center,
			labelStyle.Render("Processed: "), valStyle.Render(fmt.Sprintf("%d", m.statsProcessed)),
			labelStyle.Render("  •  Errors: "), errStyle.Render(fmt.Sprintf("%d", m.statsErrors)),
			labelStyle.Render("  •  Last Active: "), valStyle.Render(lastActive),
		),
	)

	// Live Log (Viewport)
	// We need to reduce the height of the viewport area when dashboard is shown
	logHeight := style.GetHeight() - lipgloss.Height(dashboard) - 1
	m.viewport.Height = logHeight

	return lipgloss.JoinVertical(lipgloss.Left,
		dashboard,
		m.viewport.View(),
	)
}

func (m model) renderKey(key, desc string, active bool) string {
	if active {
		return lipgloss.JoinHorizontal(lipgloss.Left,
			m.styles.Key.Foreground(m.theme.Background).Background(m.theme.Accent).Render(" "+key+" "),
			m.styles.SelectedAccent.Render(" "+desc+" "),
		)
	}
	return fmt.Sprintf("%s %s", m.styles.Key.Render(key), m.styles.Desc.Render(desc))
}

func main() {
	var configPath string
	var socketPath string
	var logFilePath string
	var verbose bool

	rootCmd := &cobra.Command{
		Use:     "straw",
		Short:   "Straw TUI client",
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			level := slog.LevelInfo
			if verbose {
				level = slog.LevelDebug
			}

			if logFilePath == "" {
				stateDir, err := config.DefaultStateDir()
				if err == nil {
					logFilePath = filepath.Join(stateDir, "straw.log")
				}
			}

			return logging.Setup(level, logFilePath, false)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var sock string
			var daemonLogPath string

			cfg, err := config.Load(configPath)
			if err != nil {
				// If config doesn't exist, create a default one
				cfg = &config.Config{
					SocketPath: config.DefaultSocketPath(),
					TUI: config.TUIConfig{
						Theme: "everforest",
					},
				}
				if configPath == "" {
					configPath, _ = config.DefaultConfigPath()
				}
			}

			sock = cfg.SocketPath
			if socketPath != "" {
				sock = socketPath
			}
			if sock == "" {
				sock = config.DefaultSocketPath()
			}

			// For the "Logs" tab, we want to look at the daemon's log file.
			// If not explicitly provided, infer it from state dir.
			daemonLogPath = logFilePath // default to what client is using or what was passed
			if daemonLogPath == "" || strings.Contains(daemonLogPath, "straw.log") {
				stateDir, err := config.DefaultStateDir()
				if err == nil {
					daemonLogPath = filepath.Join(stateDir, "strawd.log")
				}
			}

			p := tea.NewProgram(initialModel(sock, configPath, cfg, daemonLogPath))
			if _, err := p.Run(); err != nil {
				return err
			}
			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to config file")
	rootCmd.PersistentFlags().StringVar(&socketPath, "socket", "", "Override socket path")
	rootCmd.PersistentFlags().StringVar(&logFilePath, "log-file", "", "Path to log file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")

	// Override Cobra's default version flag to avoid conflict with -v (verbose)
	rootCmd.Flags().Bool("version", false, "Print the version")
	rootCmd.SetVersionTemplate(fmt.Sprintf("straw version {{.Version}} (commit: %s, built: %s)\n", commit, date))

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
