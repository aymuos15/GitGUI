package views

import (
	"fmt"
	"strings"

	"gg/src/models"
	"gg/src/styles"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

// getStatusStyle returns the styled status text with appropriate color
func getStatusStyle(status string) lipgloss.Style {
	switch status {
	case "New":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // Green
	case "Deleted":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("9")) // Red
	case "Renamed":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("11")) // Yellow
	case "Modified":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("12")) // Blue
	default:
		return lipgloss.NewStyle() // Default
	}
}

// UpdateStatsContent initializes the stats table (should only be called once)
func UpdateStatsContent(m *models.Model) {
	// Calculate totals
	totalAdditions := 0
	totalDeletions := 0

	for _, file := range m.Files {
		totalAdditions += file.Additions
		totalDeletions += file.Deletions
	}

	// Build table rows
	rows := []table.Row{}
	for _, file := range m.Files {
		// Apply color styling to status - display only first letter
		statusLetter := string(file.Status[0])
		styledStatus := getStatusStyle(file.Status).Render(statusLetter)

		rows = append(rows, table.NewRow(table.RowData{
			"file":    file.Name,
			"status":  styledStatus,
			"added":   file.Additions,
			"removed": file.Deletions,
		}))
	}

	// Add separator line before Total
	rows = append(rows, table.NewRow(table.RowData{
		"file":    strings.Repeat("─", 50),
		"status":  strings.Repeat("─", 6),
		"added":   strings.Repeat("─", 10),
		"removed": strings.Repeat("─", 10),
	}))

	// Add Total row at the end
	fileCount := len(m.Files)
	fileWord := "file"
	if fileCount != 1 {
		fileWord = "files"
	}
	totalLabel := fmt.Sprintf("Total: %d %s changed", fileCount, fileWord)

	rows = append(rows, table.NewRow(table.RowData{
		"file":    totalLabel,
		"status":  "",
		"added":   totalAdditions,
		"removed": totalDeletions,
	}))

	// Define table columns - file gets fixed good width
	columns := []table.Column{
		table.NewColumn("file", "File", 50),
		table.NewColumn("status", "Status", 6).WithStyle(lipgloss.NewStyle().Align(lipgloss.Center)),
		table.NewColumn("added", "Added", 10).WithStyle(lipgloss.NewStyle().Align(lipgloss.Right).Foreground(lipgloss.Color("10"))),
		table.NewColumn("removed", "Removed", 10).WithStyle(lipgloss.NewStyle().Align(lipgloss.Right).Foreground(lipgloss.Color("9"))),
	}

	// Create table with custom styles - fixed page size for scrolling
	m.StatsTable = table.New(columns).
		WithRows(rows).
		Focused(true).
		Border(styles.TableBorder).
		HeaderStyle(styles.TableHeaderStyle).
		WithBaseStyle(styles.TableBaseStyle).
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
		diffIndicator := getDiffTypeIndicator(m.DiffType)
		rightHelp := fmt.Sprintf("a:auto-reload[%s] l:log%s q:quit", getAutoReloadStatus(m.AutoReloadEnabled), diffIndicator)
		help := RenderHelpBarSplit("", rightHelp, m.Width)

		return content + "\n" + help
	}

	// Center the table vertically and horizontally
	centeredContent := lipgloss.Place(
		m.Width,
		m.Height-1, // Leave space for help at bottom
		lipgloss.Center,
		lipgloss.Center,
		m.StatsTable.View(),
	)

	// Render help bar with left and right sections
	diffIndicator := getDiffTypeIndicator(m.DiffType)
	rightHelp := fmt.Sprintf("a:auto-reload[%s] d:diff s:stats l:log%s q:quit", getAutoReloadStatus(m.AutoReloadEnabled), diffIndicator)
	help := RenderHelpBarSplit("↑↓:scroll", rightHelp, m.Width)

	return centeredContent + "\n" + help
}
