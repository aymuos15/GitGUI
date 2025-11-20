package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Soft, subtle background colors like GitHub
	LeftBgStyle      = lipgloss.NewStyle().Background(lipgloss.Color("#3d1e1e"))       // Very subtle red tint
	RightBgStyle     = lipgloss.NewStyle().Background(lipgloss.Color("#1e3d1e"))       // Very subtle green tint
	NeutralStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))             // White for context
	HeaderStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true) // Cyan for headers
	TitleStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")) // Blue for file names
	HelpStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Background(lipgloss.Color("236")).Padding(0, 1)
	HelpItemStyle    = lipgloss.NewStyle().Background(lipgloss.Color("205")).Foreground(lipgloss.Color("0")).Padding(0, 2).Bold(true)
	HelpGapStyle     = lipgloss.NewStyle().Background(lipgloss.Color("0"))
	DividerStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))   // Gray divider
	LineNumStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Gray line numbers
	LineNumBgLeft    = lipgloss.NewStyle().Background(lipgloss.Color("#3d1e1e")).Foreground(lipgloss.Color("240"))
	LineNumBgRight   = lipgloss.NewStyle().Background(lipgloss.Color("#1e3d1e")).Foreground(lipgloss.Color("240"))
	ActiveTabStyle   = lipgloss.NewStyle().Background(lipgloss.Color("12")).Foreground(lipgloss.Color("15")).Bold(true).Padding(0, 2)
	InactiveTabStyle = lipgloss.NewStyle().Background(lipgloss.Color("236")).Foreground(lipgloss.Color("7")).Padding(0, 2)
	TabGapStyle      = lipgloss.NewStyle().Background(lipgloss.Color("234"))

	// ANSI codes for highlighting changed text within lines
	RedHighlight   = "\x1b[48;2;90;40;40m" // Brighter red background for changed text
	GreenHighlight = "\x1b[48;2;40;90;40m" // Brighter green background for changed text
	ResetCode      = "\x1b[0m"
)
