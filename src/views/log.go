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
	// Guard against uninitialized dimensions
	if m.Width == 0 {
		return
	}

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

	// Run git log command - fetch all commits with colored graph and branch decorations
	cmd := exec.Command("git", "log", "--graph", "--color=always", "--all", "--decorate=short", "--pretty=format:%Cred%h%Creset - %d %s %Cgreen(%cr)%Creset %C(bold blue)<%an>%Creset", "--abbrev-commit")
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

		// Extract branch info from decorations like "(HEAD -> main, origin/main)"
		// Split into local branch and origin/remote branch
		localBranch := ""
		originBranch := ""
		bracketStart := strings.Index(rest, "(")
		bracketEnd := strings.Index(rest, ")")
		if bracketStart != -1 && bracketEnd != -1 && bracketStart < bracketEnd {
			// Extract decoration content
			decoration := rest[bracketStart+1 : bracketEnd]

			// Only treat as branch decoration if it contains HEAD, tag, or "/" (remote branches)
			// This avoids mistaking time format like "(2 hours ago)" as decoration
			if strings.Contains(decoration, "HEAD") || strings.Contains(decoration, "tag:") || strings.Contains(decoration, "/") {
				rest = strings.TrimSpace(rest[bracketEnd+1:])

				// Parse branch names from decoration
				// Format: HEAD -> branch, origin/branch, tag: ...
				decorParts := strings.Split(decoration, ",")
				for _, part := range decorParts {
					part = strings.TrimSpace(part)
					// Extract local branch name from HEAD pointer
					if strings.HasPrefix(part, "HEAD ->") {
						branch := strings.TrimSpace(strings.TrimPrefix(part, "HEAD ->"))
						localBranch = branch
					} else if !strings.HasPrefix(part, "tag:") {
						// Check if it's a remote branch (starts with known remote prefixes)
						isRemote := strings.HasPrefix(part, "origin/") ||
							strings.HasPrefix(part, "upstream/") ||
							strings.HasPrefix(part, "remote/")

						if isRemote {
							originBranch = part
						} else {
							// It's a local branch (could have slashes like feature/new)
							localBranch = part
						}
					}
				}
			}
		}

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
			"branch":  localBranch,
			"origin":  originBranch,
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

	// Define table columns
	// Order: Hash → Branch → Origin → Graph → Message → Author → Time
	// Calculate widths to fill full screen using ratios
	availableWidth := m.Width - 8 // Reserve space for borders: 2 outer + 6 inner dividers

	// First allocate graph width, then distribute remaining space by ratios
	remainingWidth := availableWidth - graphWidth

	// Assign ratios: Hash(5%) Branch(12%) Origin(12%) Message(45%) Author(13%) Time(13%)
	hashWidth := int(float64(remainingWidth) * 0.05)
	branchWidth := int(float64(remainingWidth) * 0.12)
	originWidth := int(float64(remainingWidth) * 0.12)
	authorWidth := int(float64(remainingWidth) * 0.13)
	timeWidth := int(float64(remainingWidth) * 0.13)

	// Calculate message width as remainder to ensure we use full width
	messageWidth := remainingWidth - hashWidth - branchWidth - originWidth - authorWidth - timeWidth

	// Ensure minimum widths for readability
	if hashWidth < 6 {
		hashWidth = 6
	}
	if branchWidth < 8 {
		branchWidth = 8
	}
	if originWidth < 8 {
		originWidth = 8
	}
	if messageWidth < 15 {
		messageWidth = 15
	}
	if authorWidth < 8 {
		authorWidth = 8
	}
	if timeWidth < 10 {
		timeWidth = 10
	}

	columns := []table.Column{
		table.NewColumn("hash", "Hash", hashWidth),
		table.NewColumn("branch", "Branch", branchWidth),
		table.NewColumn("origin", "Origin", originWidth),
		table.NewColumn("graph", "Graph", graphWidth),
		table.NewColumn("message", "Message", messageWidth),
		table.NewColumn("author", "Author", authorWidth),
		table.NewColumn("time", "Time", timeWidth),
	}

	// Calculate page size based on available height
	// Reserve space for: header row (3 lines with borders), help bar (1 line), and padding
	pageSize := m.Height - 5
	if pageSize < 5 {
		pageSize = 5 // minimum page size
	}

	// Create table with custom styles - dynamic page size based on terminal height
	m.LogTable = table.New(columns).
		WithRows(rows).
		Focused(true).
		Border(styles.TableBorder).
		HeaderStyle(styles.TableHeaderStyle).
		WithBaseStyle(styles.TableBaseStyle).
		WithPageSize(pageSize).
		WithFooterVisibility(false)

	// Just update the table - don't store anything in viewport
	// The table will be rendered fresh each time with its current scroll state
}

// RenderLogView renders the log viewport
func RenderLogView(m *models.Model) string {
	// Guard against uninitialized dimensions
	if m.Width == 0 || m.Height == 0 {
		return "Initializing..."
	}

	// Note: We don't check NoDiffMessage here because logs should always be shown
	// even when there are no current changes to diff

	// Render table and help bar
	tableView := m.LogTable.View()

	// Render help bar with left and right sections
	leftHelp := "↑↓:scroll"
	diffIndicator := getDiffTypeIndicator(m.DiffType)
	rightHelp := fmt.Sprintf("a:auto-reload[%s] d:diff s:stats l:log%s q:quit", getAutoReloadStatus(m.AutoReloadEnabled), diffIndicator)
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
