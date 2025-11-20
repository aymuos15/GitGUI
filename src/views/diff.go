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

	// Calculate column width based on new layout
	sidebarWidth := int(float64(m.Width) * 0.4)
	diffWidth := m.Width - sidebarWidth - 1 // 1 for divider
	colWidth := diffWidth/2 - 8             // Account for line numbers (6 chars + space)
	fullWidth := diffWidth - 1              // Full width minus divider

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

// RenderDiffView renders the side-by-side diff view with static sidebar
func RenderDiffView(m *models.Model) string {
	if !m.Ready {
		return "Loading..."
	}

	// Calculate sidebar width
	sidebarWidth := int(float64(m.Width) * 0.4)
	diffWidth := m.Width - sidebarWidth - 1 // 1 for divider

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
			Width(diffWidth)

		// Center vertically (account for tab bar taking 1 line)
		verticalPadding := (m.Height - 3) / 2
		content := strings.Repeat("\n", verticalPadding) + messageStyle.Render(m.NoDiffMessage)

		// Render help bar
		rightHelp := "d:diff l:log q:quit"
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
	diffContent := strings.Join(combined, "\n")

	// Render sidebar
	sidebar := RenderSidebar(m, sidebarWidth)

	// Combine diff and sidebar horizontally (line by line)
	diffLines := strings.Split(diffContent, "\n")
	sidebarLines := strings.Split(sidebar, "\n")

	var bodyContent []string
	maxBodyLines := len(diffLines)
	if len(sidebarLines) > maxBodyLines {
		maxBodyLines = len(sidebarLines)
	}

	for i := 0; i < maxBodyLines; i++ {
		left := ""
		right := ""
		if i < len(diffLines) {
			left = diffLines[i]
		}
		if i < len(sidebarLines) {
			right = sidebarLines[i]
		}
		bodyContent = append(bodyContent, left+right)
	}

	// Assemble final output: tabs + body + help
	// Tab bar already spans full width, just add it on its own line
	body := strings.Join(bodyContent, "\n")

	// Render help bar with left and right sections
	leftHelp := "↑↓:scroll h/←→:file 1-9:jump"
	rightHelp := "d:diff s:stats l:log q:quit"
	help := RenderHelpBarSplit(leftHelp, rightHelp, m.Width)

	return fmt.Sprintf("%s%s\n%s", tabBar, body, help)
}

// RenderHelpBar renders help items as styled tabs
func RenderHelpBar(helpText string, width int) string {
	// Split help text by spaces to get individual items
	items := strings.Fields(helpText)

	var styledItems []string

	for i, item := range items {
		// Cycle through different styles
		styleIndex := i % len(styles.HelpItemStyles)
		styledItems = append(styledItems, styles.HelpItemStyles[styleIndex].Render(item))
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

// RenderSidebar renders a static sidebar with stats and helper info
func RenderSidebar(m *models.Model, sidebarWidth int) string {
	sidebarStyle := lipgloss.NewStyle().
		Width(sidebarWidth).
		Height(m.Height-2).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color("240")).
		PaddingLeft(1).
		PaddingRight(1)

	// Calculate stats
	totalAdditions := 0
	totalDeletions := 0
	for _, file := range m.Files {
		totalAdditions += file.Additions
		totalDeletions += file.Deletions
	}

	// Build sidebar content
	sections := []string{}

	// File info section
	fileInfo := fmt.Sprintf("Files: %d", len(m.Files))
	sections = append(sections, lipgloss.NewStyle().Bold(true).Render(fileInfo))

	// Stats section
	statsText := fmt.Sprintf("Added: +%d\nRemoved: -%d", totalAdditions, totalDeletions)
	sections = append(sections, statsText)

	// Current file info
	if m.ActiveTab < len(m.Files) {
		file := m.Files[m.ActiveTab]
		currentFile := fmt.Sprintf("\nFile: %s\n+%d -%d", utils.Truncate(file.Name, sidebarWidth-4), file.Additions, file.Deletions)
		sections = append(sections, currentFile)
	}

	content := strings.Join(sections, "\n")
	content = utils.PadRight(content, sidebarWidth-2)

	return sidebarStyle.Render(content)
}
