package views

import (
	"fmt"
	"strings"

	"diffview/src/highlighting"
	"diffview/src/models"
	"diffview/src/styles"
	"diffview/src/utils"

	"github.com/charmbracelet/lipgloss"
)

// UpdateContent updates the viewport content with the current file's diff
func UpdateContent(m *models.Model) {
	if len(m.Files) == 0 || m.ActiveTab >= len(m.Files) {
		return
	}

	content := m.Files[m.ActiveTab].Content
	colWidth := m.Width/2 - 8 // Account for line numbers (6 chars + space)
	fullWidth := m.Width - 1  // Full width minus divider
	var leftLines, rightLines []string

	m.LeftLineNum = 0
	m.RightLineNum = 0

	for lineIdx, line := range content {
		left, right, isFullWidth, skip := formatLine(m, line, colWidth, fullWidth, lineIdx)
		if skip {
			// Skip this line entirely
			continue
		}
		if isFullWidth {
			// Header lines that span full width
			leftLines = append(leftLines, left)
			rightLines = append(rightLines, "")
		} else {
			leftLines = append(leftLines, left)
			rightLines = append(rightLines, right)
		}
	}

	m.LeftViewport.SetContent(strings.Join(leftLines, "\n"))
	m.RightViewport.SetContent(strings.Join(rightLines, "\n"))
}

// formatLine formats a single diff line for display
func formatLine(m *models.Model, line string, width int, fullWidth int, lineIdx int) (string, string, bool, bool) {
	if len(line) == 0 {
		return "", "", false, false
	}

	// Skip diff --git, ---, +++ lines (filename info redundant with tabs)
	if strings.HasPrefix(line, "diff --git") ||
		strings.HasPrefix(line, "---") ||
		strings.HasPrefix(line, "+++") {
		return "", "", false, true // skip = true
	}

	// Skip index line - not useful to users
	if strings.HasPrefix(line, "index ") {
		return "", "", false, true // skip = true
	}

	if strings.HasPrefix(line, "@@") {
		// Extract line numbers from @@ header
		parts := strings.Split(line, "@@")
		if len(parts) >= 2 {
			nums := strings.TrimSpace(parts[1])
			// Parse "-leftStart,leftCount +rightStart,rightCount"
			if strings.Contains(nums, "-") && strings.Contains(nums, "+") {
				leftPart := strings.Split(strings.Split(nums, "+")[0], "-")[1]
				rightPart := strings.Split(nums, "+")[1]

				if leftNum := strings.Split(strings.TrimSpace(leftPart), ","); len(leftNum) > 0 {
					fmt.Sscanf(leftNum[0], "%d", &m.LeftLineNum)
				}
				if rightNum := strings.Split(strings.TrimSpace(rightPart), ","); len(rightNum) > 0 {
					fmt.Sscanf(rightNum[0], "%d", &m.RightLineNum)
				}
			}
		}

		// Just show the hunk header
		formatted := styles.HeaderStyle.Render(utils.PadRight(utils.Truncate(line, fullWidth), fullWidth))
		return formatted, "", true, false
	}

	// Get current file reference for syntax highlighting
	var fileRef *models.FileDiff
	if m.ActiveTab < len(m.Files) {
		fileRef = &m.Files[m.ActiveTab]
	}

	// Handle diff lines
	if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
		// Removed line - show on left only with syntax highlighting
		text := line[1:]

		// Apply syntax highlighting
		highlighted := highlighting.HighlightCode(text, fileRef, lineIdx)

		// Truncate if needed
		visibleLen := len(utils.StripAnsi(highlighted))
		if visibleLen > width {
			// Truncate the original text and re-highlight
			text = text[:width-3] + "..."
			highlighted = highlighting.HighlightCode(text, fileRef, lineIdx)
		}

		// Apply background color directly with ANSI codes to preserve syntax highlighting
		bgCode := "\x1b[48;2;61;30;30m" // #3d1e1e
		resetBg := "\x1b[49m"

		// Pad to width
		padding := width - visibleLen
		if padding < 0 {
			padding = 0
		}

		lineNum := fmt.Sprintf("%5d ", m.LeftLineNum)
		left := styles.LineNumBgLeft.Render(lineNum) + bgCode + highlighted + strings.Repeat(" ", padding) + resetBg

		// Right side empty with neutral background
		emptyStyle := styles.NeutralStyle
		right := "      " + emptyStyle.Render("")
		m.LeftLineNum++
		return left, right, false, false
	}

	if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
		// Added line - show on right only with syntax highlighting
		text := line[1:]

		// Apply syntax highlighting
		highlighted := highlighting.HighlightCode(text, fileRef, lineIdx)

		// Truncate if needed
		visibleLen := len(utils.StripAnsi(highlighted))
		if visibleLen > width {
			// Truncate the original text and re-highlight
			text = text[:width-3] + "..."
			highlighted = highlighting.HighlightCode(text, fileRef, lineIdx)
		}

		// Apply background color directly with ANSI codes to preserve syntax highlighting
		bgCode := "\x1b[48;2;30;61;30m" // #1e3d1e
		resetBg := "\x1b[49m"

		// Pad to width
		padding := width - visibleLen
		if padding < 0 {
			padding = 0
		}

		lineNum := fmt.Sprintf("%5d ", m.RightLineNum)

		// Left side empty with neutral background
		emptyStyle := styles.NeutralStyle
		left := "      " + emptyStyle.Render("")
		right := styles.LineNumBgRight.Render(lineNum) + bgCode + highlighted + strings.Repeat(" ", padding) + resetBg
		m.RightLineNum++
		return left, right, false, false
	}

	// Context line - show on both sides with syntax highlighting
	highlighted := highlighting.HighlightCode(line, fileRef, lineIdx)

	// Truncate if needed
	visibleLen := len(utils.StripAnsi(highlighted))
	if visibleLen > width {
		line = line[:width-3] + "..."
		highlighted = highlighting.HighlightCode(line, fileRef, lineIdx)
	}

	leftNum := fmt.Sprintf("%5d ", m.LeftLineNum)
	rightNum := fmt.Sprintf("%5d ", m.RightLineNum)
	left := styles.LineNumStyle.Render(leftNum) + styles.NeutralStyle.Render(utils.PadRight(highlighted, width))
	right := styles.LineNumStyle.Render(rightNum) + styles.NeutralStyle.Render(utils.PadRight(highlighted, width))
	m.LeftLineNum++
	m.RightLineNum++
	return left, right, false, false
}

// RenderDiffView renders the side-by-side diff view
func RenderDiffView(m *models.Model) string {
	if !m.Ready {
		return "Loading..."
	}

	// Render tabs (only if multiple files)
	var tabBar string
	if len(m.Files) > 1 {
		var tabs []string
		for i, file := range m.Files {
			style := styles.InactiveTabStyle
			if i == m.ActiveTab {
				style = styles.ActiveTabStyle
			}
			tabLabel := file.Name
			if len(tabLabel) > 20 {
				tabLabel = tabLabel[:17] + "..."
			}
			tabs = append(tabs, style.Render(tabLabel))
		}
		tabBar = lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

		// Add gap to fill the rest of the width
		tabBarWidth := len(utils.StripAnsi(tabBar))
		if tabBarWidth < m.Width {
			gap := styles.TabGapStyle.Render(strings.Repeat(" ", m.Width-tabBarWidth))
			tabBar = tabBar + gap
		}
		tabBar = tabBar + "\n"
	}

	divider := styles.DividerStyle.Render("│")

	leftView := m.LeftViewport.View()
	rightView := m.RightViewport.View()

	// Split into lines and join with divider
	leftLines := strings.Split(leftView, "\n")
	rightLines := strings.Split(rightView, "\n")

	var combined []string
	maxLines := len(leftLines)
	if len(rightLines) > maxLines {
		maxLines = len(rightLines)
	}

	for i := 0; i < maxLines; i++ {
		left := ""
		right := ""
		if i < len(leftLines) {
			left = leftLines[i]
		}
		if i < len(rightLines) {
			right = rightLines[i]
		}

		// If right is empty, this is a full-width header line
		if right == "" && left != "" {
			combined = append(combined, left)
		} else {
			combined = append(combined, left+divider+right)
		}
	}

	// Build output with tabs at top (if multiple files) and minimal help at bottom
	content := strings.Join(combined, "\n")

	// Render help bar with tab-styled items
	helpText := "↑↓:scroll h/←→:file 1-9:jump s:stats l:log q:quit"
	help := RenderHelpBar(helpText, m.Width)

	if tabBar != "" {
		return fmt.Sprintf("%s%s\n%s", tabBar, content, help)
	}
	return fmt.Sprintf("%s\n%s", content, help)
}

// RenderHelpBar renders help items as styled tabs
func RenderHelpBar(helpText string, width int) string {
	// Split help text by spaces to get individual items
	items := strings.Fields(helpText)

	var styledItems []string
	gapBetween := styles.HelpGapStyle.Render("  ") // Two spaces for visual separation

	for i, item := range items {
		// Add visual brackets to make it look like a tab
		tabItem := " " + item + " "
		styledItems = append(styledItems, styles.HelpItemStyle.Render(tabItem))
		// Add gap between items (but not after the last one)
		if i < len(items)-1 {
			styledItems = append(styledItems, gapBetween)
		}
	}

	// Join items with gaps
	helpBar := lipgloss.JoinHorizontal(lipgloss.Top, styledItems...)

	// Calculate remaining space and fill with background
	helpBarWidth := lipgloss.Width(helpBar)
	if helpBarWidth < width {
		gap := styles.HelpGapStyle.Render(strings.Repeat(" ", width-helpBarWidth))
		helpBar = helpBar + gap
	}

	return helpBar
}
