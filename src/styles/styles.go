package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

var (
	// Soft, subtle background colors like GitHub
	NeutralStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))             // White for context
	HeaderStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true) // Cyan for headers
	HelpItemStyle1   = lipgloss.NewStyle().Background(lipgloss.Color("#6B5B7C")).Foreground(lipgloss.Color("15")).Padding(0, 1).Bold(true) // Dark pastel purple
	HelpItemStyle2   = lipgloss.NewStyle().Background(lipgloss.Color("#5B7C7C")).Foreground(lipgloss.Color("15")).Padding(0, 1).Bold(true) // Dark pastel teal
	HelpItemStyle3   = lipgloss.NewStyle().Background(lipgloss.Color("#7C6B5B")).Foreground(lipgloss.Color("15")).Padding(0, 1).Bold(true) // Dark pastel brown
	HelpItemStyle4   = lipgloss.NewStyle().Background(lipgloss.Color("#5B7C6B")).Foreground(lipgloss.Color("15")).Padding(0, 1).Bold(true) // Dark pastel green
	HelpItemStyle5   = lipgloss.NewStyle().Background(lipgloss.Color("#7C5B6B")).Foreground(lipgloss.Color("15")).Padding(0, 1).Bold(true) // Dark pastel mauve
	HelpItemStyles   = []lipgloss.Style{HelpItemStyle1, HelpItemStyle2, HelpItemStyle3, HelpItemStyle4, HelpItemStyle5}
	HelpGapStyle     = lipgloss.NewStyle().Background(lipgloss.Color("0"))
	DividerStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))   // Gray divider
	LineNumStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Gray line numbers
	LineNumBgLeft    = lipgloss.NewStyle().Background(lipgloss.Color("#3d1e1e")).Foreground(lipgloss.Color("240"))
	LineNumBgRight   = lipgloss.NewStyle().Background(lipgloss.Color("#1e3d1e")).Foreground(lipgloss.Color("240"))
	ActiveTabStyle   = lipgloss.NewStyle().Background(lipgloss.Color("12")).Foreground(lipgloss.Color("15")).Bold(true).Padding(0, 2)
	InactiveTabStyle = lipgloss.NewStyle().Background(lipgloss.Color("236")).Foreground(lipgloss.Color("7")).Padding(0, 2)
	TabGapStyle      = lipgloss.NewStyle().Background(lipgloss.Color("234"))
	ResetCode        = "\x1b[0m"

	// Table styles shared across views
	TableBorder = table.Border{
		Top:            "─",
		Left:           "│",
		Right:          "│",
		Bottom:         "─",
		TopRight:       "┐",
		TopLeft:        "┌",
		BottomRight:    "┘",
		BottomLeft:     "└",
		TopJunction:    "┬",
		LeftJunction:   "├",
		RightJunction:  "┤",
		BottomJunction: "┴",
		InnerJunction:  "┼",
		InnerDivider:   "│",
	}

	TableHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Background(lipgloss.Color("12")).
				Align(lipgloss.Center).
				Bold(true)

	TableBaseStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Align(lipgloss.Left)
)
