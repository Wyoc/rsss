package tui

import "github.com/charmbracelet/lipgloss"

// Theme represents a color theme
type Theme struct {
	Primary    string
	Secondary  string
	Accent     string
	Background string
	Error      string
	Success    string
}

// Themes available in the application
var Themes = map[string]Theme{
	"default": {
		Primary:    "62",
		Secondary:  "240",
		Accent:     "205",
		Background: "230",
		Error:      "196",
		Success:    "46",
	},
	"dark": {
		Primary:    "39",
		Secondary:  "245",
		Accent:     "212",
		Background: "235",
		Error:      "196",
		Success:    "46",
	},
	"ocean": {
		Primary:    "33",
		Secondary:  "39",
		Accent:     "45",
		Background: "195",
		Error:      "196",
		Success:    "46",
	},
}

// Styles contains all styled components
type Styles struct {
	Title    lipgloss.Style
	Selected lipgloss.Style
	Normal   lipgloss.Style
	Error    lipgloss.Style
	Success  lipgloss.Style
	Accent   lipgloss.Style
	Menu     lipgloss.Style
	Header   lipgloss.Style
}

// NewStyles creates styles based on the given theme name
func NewStyles(themeName string) *Styles {
	theme, exists := Themes[themeName]
	if !exists {
		theme = Themes["default"]
	}

	return &Styles{
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Primary)).
			Bold(true).
			Padding(0, 1),

		Selected: lipgloss.NewStyle().
			Background(lipgloss.Color(theme.Primary)).
			Foreground(lipgloss.Color(theme.Background)).
			Padding(0, 1),

		Normal: lipgloss.NewStyle().
			Padding(0, 1),

		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Error)).
			Bold(true),

		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Success)).
			Bold(true),

		Accent: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Accent)).
			Bold(true),

		Menu: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.Secondary)).
			Padding(1, 2),

		Header: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Primary)).
			Bold(true).
			Align(lipgloss.Center).
			Width(80),
	}
}