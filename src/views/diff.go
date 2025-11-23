package views

import (
	"fmt"
	"strings"

	"gg/src/models"
	"gg/src/styles"
	"gg/src/utils"

	"github.com/charmbracelet/lipgloss"
)

// getAutoReloadStatus returns "on" if enabled, "off" otherwise
func getAutoReloadStatus(enabled bool) string {
	if enabled {
		return "on"
	}
	return "off"
}

// getDiffTypeIndicator returns a string indicator for the diff type
func getDiffTypeIndicator(diffType string) string {
	if diffType == "staged" {
		return " (staged)"
	}
	return ""
}

// UpdateContent updates the viewport content with the current file's diff
func UpdateContent(m *models.Model) {
	if len(m.Files) == 0 || m.ActiveTab >= len(m.Files) {
		return
	}

	currentFile := m.Files[m.ActiveTab]
	content := currentFile.Content

	// For untracked files, show file content on right side (like additions)
	if currentFile.Status == "Untracked" {
		leftColWidth := m.LeftViewport.Width
		rightColWidth := m.RightViewport.Width
		rightContentWidth := rightColWidth - 6 // -6 for line numbers

		var rightLines []string
		rightLineNum := 1

		for lineIdx, line := range content {
			// Apply syntax highlighting
			highlighted := line
			highlighted = currentFile.HighlightLine(lineIdx, line)

			// Truncate if needed
			visibleLen := len(utils.StripAnsi(highlighted))
			if visibleLen > rightContentWidth {
				// Truncate the original line safely
				maxLen := min(len(line), rightContentWidth-3)
				if maxLen < 0 {
					maxLen = 0
				}
				line = line[:maxLen] + "..."
				highlighted = line
				highlighted = currentFile.HighlightLine(lineIdx, line)
			}

			// Apply background color for additions
			bgCode := "\x1b[48;2;30;61;30m" // #1e3d1e
			resetBg := "\x1b[49m"

			padding := rightContentWidth - len(utils.StripAnsi(highlighted))
			if padding < 0 {
				padding = 0
			}

			lineNum := fmt.Sprintf("%5d ", rightLineNum)
			left := "      " + styles.NeutralStyle.Render(strings.Repeat(" ", leftColWidth))
			right := styles.LineNumBgRight.Render(lineNum) + bgCode + highlighted + strings.Repeat(" ", padding) + resetBg

			rightLines = append(rightLines, left+styles.DividerStyle.Render("│")+right)
			rightLineNum++
		}

		m.LeftViewport.SetContent("")
		m.RightViewport.SetContent(strings.Join(rightLines, "\n"))
		return
	}

	// Use the actual viewport widths (set in model.go)
	leftColWidth := m.LeftViewport.Width
	rightColWidth := m.RightViewport.Width
	leftContentWidth := leftColWidth - 6   // -6 for line numbers ("12345 ")
	rightContentWidth := rightColWidth - 6 // -6 for line numbers ("12345 ")

	// Calculate full width for headers (full screen width minus center divider)
	fullWidth := m.Width - 1

	var leftLines, rightLines []string
	leftLineNum := 0
	rightLineNum := 0

	for lineIdx, line := range content {
		left, right, isFullWidth, skip := formatLineWithWidths(m, line, leftContentWidth, rightContentWidth, fullWidth, lineIdx, &leftLineNum, &rightLineNum)
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

// formatLineWithWidths formats a single diff line for display with separate left/right widths
func formatLineWithWidths(m *models.Model, line string, leftWidth int, rightWidth int, fullWidth int, lineIdx int, leftLineNum, rightLineNum *int) (string, string, bool, bool) {
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
					fmt.Sscanf(leftNum[0], "%d", leftLineNum)
				}
				if rightNum := strings.Split(strings.TrimSpace(rightPart), ","); len(rightNum) > 0 {
					fmt.Sscanf(rightNum[0], "%d", rightLineNum)
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
		highlighted := text
		if fileRef != nil {
			highlighted = fileRef.HighlightLine(lineIdx, text)
		}

		// Truncate if needed
		visibleLen := len(utils.StripAnsi(highlighted))
		if visibleLen > leftWidth {
			// Truncate the original text and re-highlight
			text = text[:leftWidth-3] + "..."
			highlighted = text
			if fileRef != nil {
				highlighted = fileRef.HighlightLine(lineIdx, text)
			}
		}

		// Apply background color directly with ANSI codes to preserve syntax highlighting
		bgCode := "\x1b[48;2;61;30;30m" // #3d1e1e
		resetBg := "\x1b[49m"

		// Pad to width
		padding := leftWidth - visibleLen
		if padding < 0 {
			padding = 0
		}

		lineNum := fmt.Sprintf("%5d ", *leftLineNum)
		left := styles.LineNumBgLeft.Render(lineNum) + bgCode + highlighted + strings.Repeat(" ", padding) + resetBg

		// Right side empty with neutral background, padded to rightWidth
		emptyStyle := styles.NeutralStyle
		right := "      " + emptyStyle.Render(strings.Repeat(" ", rightWidth))
		*leftLineNum++
		return left, right, false, false
	}

	if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
		// Added line - show on right only with syntax highlighting
		text := line[1:]

		// Apply syntax highlighting
		highlighted := text
		if fileRef != nil {
			highlighted = fileRef.HighlightLine(lineIdx, text)
		}

		// Truncate if needed
		visibleLen := len(utils.StripAnsi(highlighted))
		if visibleLen > rightWidth {
			// Truncate the original text and re-highlight
			text = text[:rightWidth-3] + "..."
			highlighted = text
			if fileRef != nil {
				highlighted = fileRef.HighlightLine(lineIdx, text)
			}
		}

		// Apply background color directly with ANSI codes to preserve syntax highlighting
		bgCode := "\x1b[48;2;30;61;30m" // #1e3d1e
		resetBg := "\x1b[49m"

		// Pad to width
		padding := rightWidth - visibleLen
		if padding < 0 {
			padding = 0
		}

		lineNum := fmt.Sprintf("%5d ", *rightLineNum)

		// Left side empty with neutral background, padded to leftWidth
		emptyStyle := styles.NeutralStyle
		left := "      " + emptyStyle.Render(strings.Repeat(" ", leftWidth))
		right := styles.LineNumBgRight.Render(lineNum) + bgCode + highlighted + strings.Repeat(" ", padding) + resetBg
		*rightLineNum++
		return left, right, false, false
	}

	// Context line - show on both sides with syntax highlighting
	leftHighlighted := line
	if fileRef != nil {
		leftHighlighted = fileRef.HighlightLine(lineIdx, line)
	}
	rightHighlighted := leftHighlighted

	// Truncate if needed for left side
	leftVisibleLen := len(utils.StripAnsi(leftHighlighted))
	if leftVisibleLen > leftWidth {
		maxLen := min(len(line), leftWidth-3)
		if maxLen < 0 {
			maxLen = 0
		}
		leftLine := line[:maxLen] + "..."
		leftHighlighted = leftLine
		if fileRef != nil {
			leftHighlighted = fileRef.HighlightLine(lineIdx, leftLine)
		}
	}

	// Truncate if needed for right side
	rightVisibleLen := len(utils.StripAnsi(rightHighlighted))
	if rightVisibleLen > rightWidth {
		maxLen := min(len(line), rightWidth-3)
		if maxLen < 0 {
			maxLen = 0
		}
		rightLine := line[:maxLen] + "..."
		rightHighlighted = rightLine
		if fileRef != nil {
			rightHighlighted = fileRef.HighlightLine(lineIdx, rightLine)
		}
	}

	leftNum := fmt.Sprintf("%5d ", *leftLineNum)
	rightNum := fmt.Sprintf("%5d ", *rightLineNum)
	left := styles.LineNumStyle.Render(leftNum) + styles.NeutralStyle.Render(utils.PadRight(leftHighlighted, leftWidth))
	right := styles.LineNumStyle.Render(rightNum) + styles.NeutralStyle.Render(utils.PadRight(rightHighlighted, rightWidth))
	*leftLineNum++
	*rightLineNum++
	return left, right, false, false
}

// RenderDiffView renders the side-by-side diff view with static sidebar
func RenderDiffView(m *models.Model) string {
	if !m.Ready {
		return "Loading..."
	}

	// Render tabs (always show, spanning full width)
	var tabBar string
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

	// Add gap to fill the full screen width
	tabBarWidth := len(utils.StripAnsi(tabBar))
	if tabBarWidth < m.Width {
		gap := styles.TabGapStyle.Render(strings.Repeat(" ", m.Width-tabBarWidth))
		tabBar = tabBar + gap
	}
	tabBar = tabBar + "\n"

	// If there's no diff to display, show a centered message
	if m.NoDiffMessage != "" {
		messageStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Bold(true).
			Align(lipgloss.Center).
			Width(m.Width)

		// Center vertically (account for tab bar taking 1 line)
		verticalPadding := (m.Height - 3) / 2
		content := strings.Repeat("\n", verticalPadding) + messageStyle.Render(m.NoDiffMessage)

		// Render help bar
		diffIndicator := getDiffTypeIndicator(m.DiffType)
		rightHelp := fmt.Sprintf("a:auto-reload[%s] d:diff l:log%s q:quit", getAutoReloadStatus(m.AutoReloadEnabled), diffIndicator)
		help := RenderHelpBarSplit("", rightHelp, m.Width)

		return tabBar + content + "\n" + help
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

	// Build diff content
	body := strings.Join(combined, "\n")

	// Render help bar with left and right sections
	leftHelp := "↑↓:scroll h/←→:file 1-9:jump"
	diffIndicator := getDiffTypeIndicator(m.DiffType)
	rightHelp := fmt.Sprintf("a:auto-reload[%s] d:diff s:stats l:log%s q:quit", getAutoReloadStatus(m.AutoReloadEnabled), diffIndicator)
	help := RenderHelpBarSplit(leftHelp, rightHelp, m.Width)

	return fmt.Sprintf("%s%s\n%s", tabBar, body, help)
}

// RenderHelpBarSplit renders help items with left and right sections
func RenderHelpBarSplit(leftText string, rightText string, width int) string {
	// Render left section
	var leftBar string
	var leftWidth int

	if leftText != "" {
		leftItems := strings.Fields(leftText)
		var styledLeftItems []string
		for i, item := range leftItems {
			styleIndex := i % len(styles.HelpItemStyles)
			styledLeftItems = append(styledLeftItems, styles.HelpItemStyles[styleIndex].Render(item))
		}
		leftBar = lipgloss.JoinHorizontal(lipgloss.Top, styledLeftItems...)
		leftWidth = lipgloss.Width(leftBar)
	}

	// Render right section
	var rightBar string
	var rightWidth int

	if rightText != "" {
		rightItems := strings.Fields(rightText)
		var styledRightItems []string
		for i, item := range rightItems {
			styleIndex := i % len(styles.HelpItemStyles)
			styledRightItems = append(styledRightItems, styles.HelpItemStyles[styleIndex].Render(item))
		}
		rightBar = lipgloss.JoinHorizontal(lipgloss.Top, styledRightItems...)
		rightWidth = lipgloss.Width(rightBar)
	}

	// Calculate gap size
	totalUsed := leftWidth + rightWidth
	gapSize := width - totalUsed
	if gapSize < 0 {
		gapSize = 0
	}

	// Create gap with background style
	gap := styles.HelpGapStyle.Render(strings.Repeat(" ", gapSize))

	// Combine left, gap, and right
	result := leftBar + gap + rightBar

	// If still shorter than width, pad the end
	resultWidth := lipgloss.Width(result)
	if resultWidth < width {
		padding := styles.HelpGapStyle.Render(strings.Repeat(" ", width-resultWidth))
		result = result + padding
	}

	return result
}
