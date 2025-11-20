package views

import (
	"os/exec"
	"strings"

	"diffview/src/models"
	"diffview/src/utils"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

// UpdateLogContent populates the log viewport with git log data
func UpdateLogContent(m *models.Model) {
	// Get HEAD commit hash
	headCmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	headOutput, _ := headCmd.Output()
	headHash := strings.TrimSpace(string(headOutput))

	// Get origin commit hash (try origin/master, then origin/main)
	originCmd := exec.Command("git", "rev-parse", "--short", "origin/master")
	originOutput, err := originCmd.Output()
	originHash := strings.TrimSpace(string(originOutput))

	if err != nil || originHash == "" {
		// Try origin/main instead
		originCmd = exec.Command("git", "rev-parse", "--short", "origin/main")
		originOutput, _ = originCmd.Output()
		originHash = strings.TrimSpace(string(originOutput))
	}

	// Run git log command - fetch all commits
	cmd := exec.Command("git", "log", "--graph", "--pretty=format:%Cred%h%Creset - %s %Cgreen(%cr)%Creset %C(bold blue)<%an>%Creset", "--abbrev-commit")
	output, err := cmd.Output()

	var logLines []string
	if err != nil {
		logLines = []string{"Error: Unable to fetch git log", "Make sure you're in a git repository"}
	} else {
		logLines = strings.Split(string(output), "\n")
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

		// Create row with optional styling based on HEAD or origin
		row := table.NewRow(table.RowData{
			"hash":    hash,
			"message": message,
			"time":    time,
			"author":  author,
		})

		// Apply color styling for special commits
		if hash == headHash {
			// HEAD commit - pink/magenta
			row = row.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("212")))
		} else if hash == originHash {
			// Origin commit - orange
			row = row.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("208")))
		}

		rows = append(rows, row)
	}

	// Define table columns - message gets good width
	columns := []table.Column{
		table.NewColumn("hash", "Hash", 10),
		table.NewColumn("message", "Message", 65),
		table.NewColumn("time", "Time", 18),
		table.NewColumn("author", "Author", 20),
	}

	// Create table with custom styles - fixed page size for scrolling
	m.LogTable = table.New(columns).
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
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("12")).
			Align(lipgloss.Center).
			Bold(true)).
		WithBaseStyle(lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Align(lipgloss.Left)).
		WithPageSize(20).
		WithFooterVisibility(false)

	// Just update the table - don't store anything in viewport
	// The table will be rendered fresh each time with its current scroll state
}

// RenderLogView renders the log viewport
func RenderLogView(m *models.Model) string {
	// Center the table vertically and horizontally
	centeredContent := lipgloss.Place(
		m.Width,
		m.Height-1, // Leave space for help at bottom
		lipgloss.Center,
		lipgloss.Center,
		m.LogTable.View(),
	)

	// Render help bar with left and right sections
	leftHelp := "↑↓:scroll"
	rightHelp := "d:diff s:stats l:log q:quit"
	help := RenderHelpBarSplit(leftHelp, rightHelp, m.Width)

	return centeredContent + "\n" + help
}
