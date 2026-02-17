package tui

import "github.com/charmbracelet/lipgloss"

// Theme defines the color palette for the TUI.
type Theme struct {
	Name       string
	Background lipgloss.Color
	Foreground lipgloss.Color
	Secondary  lipgloss.Color // Lighter background
	Tertiary   lipgloss.Color // Darker background/border
	Accent     lipgloss.Color
	Success    lipgloss.Color
	Warning    lipgloss.Color
	Error      lipgloss.Color
	Dim        lipgloss.Color

	// Action specific
	ActionMove  lipgloss.Color
	ActionCopy  lipgloss.Color
	ActionTrash lipgloss.Color
	ActionShell lipgloss.Color
}

// Styles contains lipgloss styles based on a theme.
type Styles struct {
	App          lipgloss.Style
	Header       lipgloss.Style
	HeaderTitle  lipgloss.Style
	HeaderStatus lipgloss.Style
	Tabs         lipgloss.Style
	TabActive    lipgloss.Style
	TabInactive  lipgloss.Style

	// List Styles
	ListTitle    lipgloss.Style
	ListItem     lipgloss.Style
	ListSelected lipgloss.Style
	ListDim      lipgloss.Style

	// Log Styles
	LogContainer lipgloss.Style
	LogTime      lipgloss.Style
	LogAction    lipgloss.Style
	LogFile      lipgloss.Style
	LogMessage   lipgloss.Style
	LogTab       lipgloss.Style

	// Footer
	Footer lipgloss.Style
	Key    lipgloss.Style
	Desc   lipgloss.Style
}

var (
	// Themes is a map of all available themes
	Themes = map[string]Theme{
		"catppuccin": Catppuccin,
		"everforest": Everforest,
		"nord":       Nord,
		"ravenwood":  Ravenwood,
	}

	// Nord theme
	Nord = Theme{
		Name:       "Nord",
		Background: lipgloss.Color("#2e3440"),
		Foreground: lipgloss.Color("#d8dee9"),
		Secondary:  lipgloss.Color("#3b4252"),
		Tertiary:   lipgloss.Color("#4c566a"),
		Accent:     lipgloss.Color("#88c0d0"),
		Success:    lipgloss.Color("#a3be8c"),
		Warning:    lipgloss.Color("#ebcb8b"),
		Error:      lipgloss.Color("#bf616a"),
		Dim:        lipgloss.Color("#616e88"),

		ActionMove:  lipgloss.Color("#81a1c1"),
		ActionCopy:  lipgloss.Color("#a3be8c"),
		ActionTrash: lipgloss.Color("#bf616a"),
		ActionShell: lipgloss.Color("#b48ead"),
	}

	// Catppuccin Macchiato
	Catppuccin = Theme{
		Name:       "Catppuccin",
		Background: lipgloss.Color("#24273a"),
		Foreground: lipgloss.Color("#cad3f5"),
		Secondary:  lipgloss.Color("#363a4f"),
		Tertiary:   lipgloss.Color("#494d64"),
		Accent:     lipgloss.Color("#b7bdf8"), // Lavender
		Success:    lipgloss.Color("#a6da95"), // Green
		Warning:    lipgloss.Color("#eed49f"), // Yellow
		Error:      lipgloss.Color("#ed8796"), // Red
		Dim:        lipgloss.Color("#6e738d"), // Overlay0

		ActionMove:  lipgloss.Color("#8aadf4"), // Blue
		ActionCopy:  lipgloss.Color("#f5a97f"), // Peach
		ActionTrash: lipgloss.Color("#ed8796"), // Red
		ActionShell: lipgloss.Color("#c6a0f6"), // Mauve
	}

	// Everforest Dark
	Everforest = Theme{
		Name:       "Everforest",
		Background: lipgloss.Color("#2d353b"),
		Foreground: lipgloss.Color("#d3c6aa"),
		Secondary:  lipgloss.Color("#374247"), // Surface0
		Tertiary:   lipgloss.Color("#4a555b"), // Surface1
		Accent:     lipgloss.Color("#7fbbb3"), // Blue
		Success:    lipgloss.Color("#a7c080"), // Green
		Warning:    lipgloss.Color("#dbbc7f"), // Yellow
		Error:      lipgloss.Color("#e67e80"), // Red
		Dim:        lipgloss.Color("#859289"), // Grey

		ActionMove:  lipgloss.Color("#7fbbb3"), // Blue
		ActionCopy:  lipgloss.Color("#dbbc7f"), // Yellow
		ActionTrash: lipgloss.Color("#e67e80"), // Red
		ActionShell: lipgloss.Color("#d699b6"), // Purple
	}

	// Ravenwood - A refined forest theme
	Ravenwood = Theme{
		Name:       "Ravenwood",
		Background: lipgloss.Color("#282e33"),
		Foreground: lipgloss.Color("#d3c6aa"),
		Secondary:  lipgloss.Color("#374145"), // Surface
		Tertiary:   lipgloss.Color("#495156"), // Overlay
		Accent:     lipgloss.Color("#7fbbb3"), // Blue/Aqua
		Success:    lipgloss.Color("#7fb38e"), // Green
		Warning:    lipgloss.Color("#dbbc7f"), // Yellow
		Error:      lipgloss.Color("#e67e80"), // Red
		Dim:        lipgloss.Color("#859289"), // Grey

		ActionMove:  lipgloss.Color("#7fbbb3"), // Blue
		ActionCopy:  lipgloss.Color("#dbbc7f"), // Yellow
		ActionTrash: lipgloss.Color("#e67e80"), // Red
		ActionShell: lipgloss.Color("#d699b6"), // Purple
	}
)

func GetTheme(name string) Theme {
	if t, ok := Themes[name]; ok {
		return t
	}
	return Ravenwood
}

// AllThemes returns all available themes
func AllThemes() []Theme {
	return []Theme{
		Everforest,
		Catppuccin,
		Nord,
		Ravenwood,
	}
}

func GetStyles(t Theme) Styles {
	return Styles{
		App: lipgloss.NewStyle().
			Padding(1, 2),

		Header: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(t.Tertiary).
			Padding(0, 0, 1, 0).
			MarginBottom(1),

		HeaderTitle: lipgloss.NewStyle().
			Foreground(t.Background).
			Background(t.Accent).
			Bold(true).
			Padding(0, 1).
			MarginRight(1),

		HeaderStatus: lipgloss.NewStyle().
			Foreground(t.Dim).
			Italic(true),

		Tabs: lipgloss.NewStyle().
			MarginLeft(2),

		TabActive: lipgloss.NewStyle().
			Foreground(t.Accent).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(t.Accent).
			Padding(0, 1).
			Bold(true),

		TabInactive: lipgloss.NewStyle().
			Foreground(t.Dim).
			Padding(0, 1),

		// List
		ListTitle: lipgloss.NewStyle().
			Foreground(t.Accent).
			Bold(true).
			MarginLeft(2),

		ListItem: lipgloss.NewStyle().
			PaddingLeft(2),

		ListSelected: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(t.Accent).
			Foreground(t.Accent).
			PaddingLeft(1).
			Bold(true),

		ListDim: lipgloss.NewStyle().
			Foreground(t.Dim),

		// Logs
		LogContainer: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Tertiary),

		LogTime: lipgloss.NewStyle().
			Foreground(t.Dim).
			Width(10),

		LogAction: lipgloss.NewStyle().
			Bold(true).
			Width(12).
			Align(lipgloss.Center).
			MarginRight(2),

		LogFile: lipgloss.NewStyle().
			Foreground(t.Foreground).
			MarginRight(2),

		LogMessage: lipgloss.NewStyle().
			Foreground(t.Dim).
			Italic(true),

		LogTab: lipgloss.NewStyle().
			Foreground(t.Foreground).
			Faint(true),

		// Footer
		Footer: lipgloss.NewStyle().
			MarginTop(1).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(t.Tertiary).
			PaddingTop(1),

		Key: lipgloss.NewStyle().
			Foreground(t.Accent).
			Bold(true),

		Desc: lipgloss.NewStyle().
			Foreground(t.Dim).
			MarginRight(2),
	}
}

// ActionColor returns the color for a specific action type
func (t Theme) ActionColor(action string) lipgloss.Color {
	switch action {
	case "move":
		return t.ActionMove
	case "copy":
		return t.ActionCopy
	case "trash":
		return t.ActionTrash
	case "shell":
		return t.ActionShell
	default:
		return t.Foreground
	}
}
