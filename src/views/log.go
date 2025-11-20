package views

import (
	"os/exec"
	"strings"

	"diffview/src/models"
	"diffview/src/styles"
	"diffview/src/utils"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// UpdateLogContent populates the log viewport with git log data
func UpdateLogContent(m *models.Model) {
	// Run git log command - fetch 100 commits for scrolling
	cmd := exec.Command("git", "log", "--graph", "--pretty=format:%Cred%h%Creset - %s %Cgreen(%cr)%Creset %C(bold blue)<%an>%Creset", "--abbrev-commit", "-100")
	output, err := cmd.Output()

	var logLines []string
	if err != nil {
		logLines = []string{"Error: Unable to fetch git log", "Make sure you're in a git repository"}
	} else {
		logLines = strings.Split(string(output), "\n")
	}

	// Define table columns
	columns := []table.Column{
		{Title: "Hash", Width: 10},
		{Title: "Message", Width: 65},
		{Title: "Time", Width: 18},
		{Title: "Author", Width: 20},
	}

	// Build table rows by parsing git log output
	rows := []table.Row{}
	for _, line := range logLines {
		// Strip ANSI codes
		cleanLine := utils.StripAnsi(line)

		// Skip empty lines
		if strings.TrimSpace(cleanLine) == "" {
			continue
		}

		// Parse format: hash - message (time) <author>
		// Remove graph characters (* | \ /)
		cleanLine = strings.TrimLeft(cleanLine, "*|\\ /")
		cleanLine = strings.TrimSpace(cleanLine)

		// Split by " - " for hash and rest
		parts := strings.SplitN(cleanLine, " - ", 2)
		if len(parts) < 2 {
			continue
		}

		hash := strings.TrimSpace(parts[0])
		rest := parts[1]

		// Find message, time, and author
		// Format: message (time) <author>
		timeStart := strings.LastIndex(rest, "(")
		authorStart := strings.LastIndex(rest, "<")

		if timeStart == -1 || authorStart == -1 {
			continue
		}

		message := strings.TrimSpace(rest[:timeStart])
		time := strings.TrimSpace(rest[timeStart+1 : strings.Index(rest[timeStart:], ")")+timeStart])
		author := strings.TrimSpace(rest[authorStart+1 : len(rest)-1])

		// Truncate if too long
		if len(message) > 65 {
			message = message[:62] + "..."
		}
		if len(author) > 20 {
			author = author[:17] + "..."
		}
		if len(time) > 18 {
			time = time[:15] + "..."
		}

		rows = append(rows, table.Row{hash, message, time, author})
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

	// Remove highlight from selected row
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("0")).
		Bold(false)

	// Cell style
	s.Cell = s.Cell.Foreground(lipgloss.Color("15"))

	t.SetStyles(s)

	// Title
	title := styles.HeaderStyle.Render("Git Log (↑↓ to scroll)")

	// Build the content with table
	var content strings.Builder
	content.WriteString(title)
	content.WriteString("\n\n")
	content.WriteString(t.View())

	// Set viewport content
	m.LogViewport.SetContent(content.String())
}

// RenderLogView renders the log viewport
func RenderLogView(m *models.Model) string {
	logContent := m.LogViewport.View()

	// Calculate padding to push help to the very bottom
	lines := strings.Split(logContent, "\n")
	currentHeight := len(lines)
	totalHeight := m.Height - 1 // Reserve 1 line for help

	// Add empty lines if needed to fill the screen
	if currentHeight < totalHeight {
		padding := strings.Repeat("\n", totalHeight-currentHeight)
		logContent += padding
	}

	// Render help bar with tab-styled items
	helpText := "↑↓:scroll d:diff s:stats l:log q:quit"
	help := RenderHelpBar(helpText, m.Width)

	return logContent + "\n" + help
}
