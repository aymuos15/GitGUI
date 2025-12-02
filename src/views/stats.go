package views

import (
	"fmt"
	"path/filepath"
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
	case "Untracked":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("13")) // Magenta
	default:
		return lipgloss.NewStyle() // Default
	}
}

// UpdateStatsContent initializes the stats table (should only be called once)
func UpdateStatsContent(m *models.Model) {
	// Calculate column widths to fill full screen width first
	// Reserve space for borders and spacing (approximate: 5 chars for borders/separators)
	availableWidth := m.Width - 5
	statusWidth := 8
	addedWidth := 12
	removedWidth := 12
	fileWidth := availableWidth - statusWidth - addedWidth - removedWidth

	// Build table rows with filtering
	rows := []table.Row{}
	totalAdditions := 0
	totalDeletions := 0
	filteredCount := 0

	statusFilter := strings.ToUpper(m.StatsFilters.Status)
	extFilter := strings.ToLower(m.StatsFilters.Extension)

	for _, file := range m.Files {
		// Apply status filter
		if statusFilter != "" {
			statusLetter := strings.ToUpper(string(file.Status[0]))
			if statusLetter != statusFilter {
				continue
			}
		}

		// Apply extension filter
		if extFilter != "" {
			fileExt := strings.ToLower(filepath.Ext(file.Name))
			// Handle filter with or without leading dot
			if !strings.HasPrefix(extFilter, ".") {
				extFilter = "." + extFilter
			}
			if fileExt != extFilter {
				continue
			}
		}

		// Apply color styling to status - display only first letter
		statusLetter := string(file.Status[0])
		styledStatus := getStatusStyle(file.Status).Render(statusLetter)

		rows = append(rows, table.NewRow(table.RowData{
			"file":    file.Name,
			"status":  styledStatus,
			"added":   file.Additions,
			"removed": file.Deletions,
		}))

		totalAdditions += file.Additions
		totalDeletions += file.Deletions
		filteredCount++
	}

	// Add separator line before Total - use calculated widths to extend end to end
	rows = append(rows, table.NewRow(table.RowData{
		"file":    strings.Repeat("─", fileWidth),
		"status":  strings.Repeat("─", statusWidth),
		"added":   strings.Repeat("─", addedWidth),
		"removed": strings.Repeat("─", removedWidth),
	}))

	// Add Total row at the end
	fileWord := "file"
	if filteredCount != 1 {
		fileWord = "files"
	}
	totalLabel := fmt.Sprintf("Total: %d %s", filteredCount, fileWord)
	if statusFilter != "" || extFilter != "" {
		totalLabel += fmt.Sprintf(" (filtered from %d)", len(m.Files))
	} else {
		totalLabel += " changed"
	}

	rows = append(rows, table.NewRow(table.RowData{
		"file":    totalLabel,
		"status":  "",
		"added":   totalAdditions,
		"removed": totalDeletions,
	}))

	// Define table columns - dynamically sized to fill full width
	columns := []table.Column{
		table.NewColumn("file", "File", fileWidth),
		table.NewColumn("status", "Status", statusWidth).WithStyle(lipgloss.NewStyle().Align(lipgloss.Center)),
		table.NewColumn("added", "Added", addedWidth).WithStyle(lipgloss.NewStyle().Align(lipgloss.Right).Foreground(lipgloss.Color("10"))),
		table.NewColumn("removed", "Removed", removedWidth).WithStyle(lipgloss.NewStyle().Align(lipgloss.Right).Foreground(lipgloss.Color("9"))),
	}

	// Calculate page size based on available height
	// Reserve space for: header row (3 lines with borders), help bar (1 line), and padding
	pageSize := m.Height - 5
	if pageSize < 5 {
		pageSize = 5 // minimum page size
	}

	// Create table with custom styles - dynamic page size based on terminal height
	m.StatsTable = table.New(columns).
		WithRows(rows).
		Focused(true).
		Border(styles.TableBorder).
		HeaderStyle(styles.TableHeaderStyle).
		WithBaseStyle(styles.TableBaseStyle).
		WithPageSize(pageSize).
		WithFooterVisibility(false)
}

// RenderStatsView renders the stats view with a clean modern interface using bubble-table
func RenderStatsView(m *models.Model) string {
	// If in filter mode, show filter input
	if m.FilterMode != "" {
		return RenderFilterInput(m, "stats")
	}

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

	// Render table and help bar
	tableView := m.StatsTable.View()

	// Build active filters indicator
	filterIndicator := buildStatsFilterIndicator(m)

	// Render help bar with left and right sections
	leftHelp := "↑↓:scroll M-s:status M-e:ext ^l:clear"
	diffIndicator := getDiffTypeIndicator(m.DiffType)
	rightHelp := fmt.Sprintf("a:auto-reload[%s] d:diff s:stats l:log%s q:quit", getAutoReloadStatus(m.AutoReloadEnabled), diffIndicator)
	if filterIndicator != "" {
		rightHelp = filterIndicator + " " + rightHelp
	}
	help := RenderHelpBarSplit(leftHelp, rightHelp, m.Width)

	// Calculate heights
	tableHeight := lipgloss.Height(tableView)
	helpHeight := 1 // Help bar is always 1 line

	// Calculate vertical padding to center table, then fill to bottom
	availableHeight := m.Height - helpHeight
	topPadding := (availableHeight - tableHeight) / 2
	bottomPadding := availableHeight - tableHeight - topPadding

	if topPadding < 0 {
		topPadding = 0
	}
	if bottomPadding < 0 {
		bottomPadding = 0
	}

	// Build output: padding + table + padding to fill screen + help at very bottom
	var output strings.Builder
	output.WriteString(strings.Repeat("\n", topPadding))
	output.WriteString(tableView)
	output.WriteString(strings.Repeat("\n", bottomPadding))
	output.WriteString("\n")
	output.WriteString(help)

	return output.String()
}

// buildStatsFilterIndicator builds a string showing active stats filters
func buildStatsFilterIndicator(m *models.Model) string {
	if m.StatsFilters.Status == "" && m.StatsFilters.Extension == "" {
		return ""
	}

	var parts []string
	if m.StatsFilters.Status != "" {
		parts = append(parts, "status:"+m.StatsFilters.Status)
	}
	if m.StatsFilters.Extension != "" {
		parts = append(parts, "ext:"+m.StatsFilters.Extension)
	}

	filterStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
	return filterStyle.Render("[" + strings.Join(parts, " ") + "]")
}
