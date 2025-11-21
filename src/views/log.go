package views

import (
	"fmt"
	"os/exec"
	"strings"

	"gg/src/models"
	"gg/src/styles"
	"gg/src/utils"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

// UpdateLogContent populates the log viewport with git log data
func UpdateLogContent(m *models.Model) {
	// Get HEAD commit hash
	headCmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	headOutput, _ := headCmd.Output()
	headHash := strings.TrimSpace(string(headOutput))

	// Get upstream branch commit hash
	// First, try to get the upstream branch for current branch
	upstreamCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	upstreamOutput, err := upstreamCmd.Output()
	originHash := ""

	if err == nil && len(upstreamOutput) > 0 {
		// Got upstream branch name, now get its commit hash
		upstreamBranch := strings.TrimSpace(string(upstreamOutput))
		originCmd := exec.Command("git", "rev-parse", "--short", upstreamBranch)
		originOutput, err := originCmd.Output()
		if err == nil {
			originHash = strings.TrimSpace(string(originOutput))
		}
	}

	// Fallback: try common remote branch names if no upstream configured
	if originHash == "" {
		for _, remoteBranch := range []string{"origin/master", "origin/main"} {
			originCmd := exec.Command("git", "rev-parse", "--short", remoteBranch)
			originOutput, err := originCmd.Output()
			if err == nil && len(originOutput) > 0 {
				originHash = strings.TrimSpace(string(originOutput))
				break
			}
		}
	}

	// Run git log command - fetch all commits with colored graph
	cmd := exec.Command("git", "log", "--graph", "--color=always", "--all", "--pretty=format:%Cred%h%Creset - %s %Cgreen(%cr)%Creset %C(bold blue)<%an>%Creset", "--abbrev-commit")
	output, err := cmd.Output()

	var logLines []string
	if err != nil {
		logLines = []string{"Error: Unable to fetch git log", "Make sure you're in a git repository"}
	} else {
		logLines = strings.Split(string(output), "\n")
	}

	// Calculate adaptive graph column width by scanning all lines
	graphWidth := 10 // minimum width
	for _, line := range logLines {
		cleanLine := utils.StripAnsi(line)
		graphEnd := 0
		for i, char := range cleanLine {
			if char != '*' && char != '|' && char != '\\' && char != '/' && char != ' ' {
				graphEnd = i
				break
			}
		}
		if graphEnd > graphWidth {
			graphWidth = graphEnd
		}
	}
	// Cap at maximum width of 30
	if graphWidth > 30 {
		graphWidth = 30
	}

	// Build table rows by parsing git log output
	rows := []table.Row{}
	for _, line := range logLines {
		// Strip ANSI codes for parsing, but keep original line for graph
		cleanLine := utils.StripAnsi(line)

		// Skip empty lines
		if strings.TrimSpace(cleanLine) == "" {
			continue
		}

		// Parse format: [graph] hash - message (time) <author>
		// Extract graph characters (* | \ /) before the hash
		graphEnd := 0
		for i, char := range cleanLine {
			if char != '*' && char != '|' && char != '\\' && char != '/' && char != ' ' {
				graphEnd = i
				break
			}
		}

		// Extract graph with ANSI color codes preserved from original line
		// Find the actual position in the original line (accounting for ANSI codes)
		graphPrefix := ""
		originalPos := 0
		cleanPos := 0
		for originalPos < len(line) && cleanPos < graphEnd {
			if line[originalPos] == '\x1b' {
				// Found ANSI escape sequence, include it in graph
				escapeEnd := originalPos
				for escapeEnd < len(line) && line[escapeEnd] != 'm' {
					escapeEnd++
				}
				if escapeEnd < len(line) {
					graphPrefix += line[originalPos : escapeEnd+1]
					originalPos = escapeEnd + 1
				}
			} else {
				graphPrefix += string(line[originalPos])
				originalPos++
				cleanPos++
			}
		}

		cleanLine = strings.TrimSpace(cleanLine[graphEnd:])

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
			"graph":   graphPrefix,
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

	// Define table columns - Hash → Graph → Message → Time → Author
	columns := []table.Column{
		table.NewColumn("hash", "Hash", 8),
		table.NewColumn("graph", "Graph", graphWidth),
		table.NewColumn("message", "Message", 57),
		table.NewColumn("time", "Time", 18),
		table.NewColumn("author", "Author", 20),
	}

	// Create table with custom styles - fixed page size for scrolling
	m.LogTable = table.New(columns).
		WithRows(rows).
		Focused(true).
		Border(styles.TableBorder).
		HeaderStyle(styles.TableHeaderStyle).
		WithBaseStyle(styles.TableBaseStyle).
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
	diffIndicator := getDiffTypeIndicator(m.DiffType)
	rightHelp := fmt.Sprintf("a:auto-reload[%s] d:diff s:stats l:log%s q:quit", getAutoReloadStatus(m.AutoReloadEnabled), diffIndicator)
	help := RenderHelpBarSplit(leftHelp, rightHelp, m.Width)

	return centeredContent + "\n" + help
}
