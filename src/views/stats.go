package views

import (
	"fmt"
	"strings"

	"diffview/src/models"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// RenderStatsView renders the stats view with a clean modern interface using bubbles table
func RenderStatsView(m *models.Model) string {
	// Calculate totals
	totalAdditions := 0
	totalDeletions := 0

	for _, file := range m.Files {
		totalAdditions += file.Additions
		totalDeletions += file.Deletions
	}

	// Define table columns
	columns := []table.Column{
		{Title: "File", Width: 50},
		{Title: "Added", Width: 10},
		{Title: "Removed", Width: 10},
	}

	// Build table rows
	rows := []table.Row{}
	for _, file := range m.Files {
		fileName := file.Name
		if len(fileName) > 50 {
			fileName = "..." + fileName[len(fileName)-47:]
		}

		rows = append(rows, table.Row{
			fileName,
			fmt.Sprintf("%d", file.Additions),
			fmt.Sprintf("%d", file.Deletions),
		})
	}

	// Create table with custom styles
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(len(rows)),
	)

	// Custom table styles
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("240"))

	// Remove highlight from selected row - make it look the same as normal cells
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("0")).
		Bold(false)

	// Cell style
	s.Cell = s.Cell.Foreground(lipgloss.Color("15"))

	t.SetStyles(s)

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Render("Change Summary")

	// Build the content with table
	var boxContent strings.Builder
	boxContent.WriteString(title)
	boxContent.WriteString("\n\n")

	tableView := t.View()
	boxContent.WriteString(tableView)
	boxContent.WriteString("\n\n")

	// Calculate separator width from actual table width
	// The table adds internal spacing, so we measure the first line
	tableLines := strings.Split(tableView, "\n")
	var separatorWidth int
	if len(tableLines) > 0 {
		// Remove ANSI codes to get actual width
		separatorWidth = lipgloss.Width(tableLines[0])
	} else {
		separatorWidth = 74 // fallback
	}

	boxContent.WriteString(strings.Repeat("─", separatorWidth))
	boxContent.WriteString("\n")

	// Summary section - align with table columns
	fileCount := len(m.Files)
	fileWord := "file"
	if fileCount != 1 {
		fileWord = "files"
	}

	summaryLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(fmt.Sprintf("Total: %d %s changed", fileCount, fileWord))

	// Format totals to align with the Added and Removed columns
	// Table columns: File(50) + padding + Added(10) + padding + Removed(10)
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
