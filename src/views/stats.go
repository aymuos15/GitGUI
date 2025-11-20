package views

import (
	"fmt"
	"strings"

	"diffview/src/models"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

// UpdateStatsContent initializes the stats table (should only be called once)
func UpdateStatsContent(m *models.Model) {
	// Build table rows
	rows := []table.Row{}
	for _, file := range m.Files {
		rows = append(rows, table.NewRow(table.RowData{
			"file":    file.Name,
			"added":   file.Additions,
			"removed": file.Deletions,
		}))
	}

	// Define table columns - file gets fixed good width
	columns := []table.Column{
		table.NewColumn("file", "File", 50),
		table.NewColumn("added", "Added", 10).WithStyle(lipgloss.NewStyle().Align(lipgloss.Right)),
		table.NewColumn("removed", "Removed", 10).WithStyle(lipgloss.NewStyle().Align(lipgloss.Right)),
	}

	// Create table with custom styles - fixed page size for scrolling
	m.StatsTable = table.New(columns).
		WithRows(rows).
		Focused(true).
		Border(table.Border{
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
		}).
		HeaderStyle(lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Bold(true)).
		WithBaseStyle(lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Align(lipgloss.Left)).
		WithPageSize(15).
		WithFooterVisibility(false)
}

// RenderStatsView renders the stats view with a clean modern interface using bubble-table
func RenderStatsView(m *models.Model) string {
	// If there's no diff to display, show a centered message
	if m.NoDiffMessage != "" {
		messageStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Bold(true).
			Align(lipgloss.Center).
			Width(m.Width)

		// Center vertically
		verticalPadding := (m.Height - 2) / 2
		content := strings.Repeat("\n", verticalPadding) + messageStyle.Render(m.NoDiffMessage)

		// Render help bar
		helpText := "l:log q:quit"
		help := RenderHelpBar(helpText, m.Width)

		return content + "\n" + help
	}

	// Calculate totals
	totalAdditions := 0
	totalDeletions := 0

	for _, file := range m.Files {
		totalAdditions += file.Additions
		totalDeletions += file.Deletions
	}

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Render("Change Summary")

	// Build the content with table
	var boxContent strings.Builder
	boxContent.WriteString(title)
	boxContent.WriteString("\n\n")

	tableView := m.StatsTable.View()
	boxContent.WriteString(tableView)
	boxContent.WriteString("\n\n")

	// Calculate separator width from actual table width
	tableLines := strings.Split(tableView, "\n")
	var separatorWidth int
	if len(tableLines) > 0 {
		separatorWidth = lipgloss.Width(tableLines[0])
	} else {
		separatorWidth = 74 // fallback
	}

	boxContent.WriteString(strings.Repeat("─", separatorWidth))
	boxContent.WriteString("\n")

	// Summary section
	fileCount := len(m.Files)
	fileWord := "file"
	if fileCount != 1 {
		fileWord = "files"
	}

	summaryLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(fmt.Sprintf("Total: %d %s changed", fileCount, fileWord))

	// Format totals
	totalAddStr := fmt.Sprintf("%10d", totalAdditions)
	totalDelStr := fmt.Sprintf("%10d", totalDeletions)

	// Apply color to the formatted strings
	totalAddStyled := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Bold(true).
		Render(totalAddStr)

	totalDelStyled := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Bold(true).
		Render(totalDelStr)

	// Build summary row with proper spacing to match table layout
	summary := fmt.Sprintf("%-50s  %s  %s", summaryLabel, totalAddStyled, totalDelStyled)
	boxContent.WriteString(summary)

	// Create a box around the content
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(2, 3).
		Align(lipgloss.Center)

	box := boxStyle.Render(boxContent.String())

	// Center the box vertically and horizontally
	centeredBox := lipgloss.Place(
		m.Width,
		m.Height-1, // Leave space for help at bottom
		lipgloss.Center,
		lipgloss.Center,
		box,
	)

	// Render help bar with tab-styled items
	helpText := "d:diff s:stats l:log q:quit"
	help := RenderHelpBar(helpText, m.Width)

	return centeredBox + "\n" + help
}
