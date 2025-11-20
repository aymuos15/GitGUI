package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	leftBgStyle      = lipgloss.NewStyle().Background(lipgloss.Color("52")).Foreground(lipgloss.Color("9"))  // Red bg for removed
	rightBgStyle     = lipgloss.NewStyle().Background(lipgloss.Color("22")).Foreground(lipgloss.Color("10")) // Green bg for added
	neutralStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))                                   // White for context
	headerStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true)                       // Cyan for headers
	titleStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))                       // Blue for file names
	helpStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	dividerStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))   // Gray divider
	lineNumStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Gray line numbers
	lineNumBgLeft    = lipgloss.NewStyle().Background(lipgloss.Color("52")).Foreground(lipgloss.Color("240"))
	lineNumBgRight   = lipgloss.NewStyle().Background(lipgloss.Color("22")).Foreground(lipgloss.Color("240"))
	activeTabStyle   = lipgloss.NewStyle().Background(lipgloss.Color("12")).Foreground(lipgloss.Color("15")).Bold(true).Padding(0, 2)
	inactiveTabStyle = lipgloss.NewStyle().Background(lipgloss.Color("236")).Foreground(lipgloss.Color("7")).Padding(0, 2)
	tabGapStyle      = lipgloss.NewStyle().Background(lipgloss.Color("234"))
)

type diffLine struct {
	leftNum   string
	rightNum  string
	leftText  string
	rightText string
	lineType  string // "add", "remove", "context", "header"
}

type fileDiff struct {
	name    string
	content []string
}

type model struct {
	leftViewport  viewport.Model
	rightViewport viewport.Model
	files         []fileDiff
	activeTab     int
	ready         bool
	width         int
	height        int
	leftLineNum   int
	rightLineNum  int
	indexInfo     string // Store index information
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "tab", "right", "l":
			if m.activeTab < len(m.files)-1 {
				m.activeTab++
				m.updateContent()
			}
		case "shift+tab", "left", "h":
			if m.activeTab > 0 {
				m.activeTab--
				m.updateContent()
			}
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			tabNum := int(msg.String()[0] - '1')
			if tabNum < len(m.files) {
				m.activeTab = tabNum
				m.updateContent()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Calculate viewport height: total - tabs (if multiple files) - help line
		viewportHeight := msg.Height - 1 // 1 for help
		if len(m.files) > 1 {
			viewportHeight -= 1 // 1 for tabs
		}

		if !m.ready {
			m.leftViewport = viewport.New(msg.Width/2-1, viewportHeight)
			m.rightViewport = viewport.New(msg.Width/2-1, viewportHeight)
			m.ready = true
		} else {
			m.leftViewport.Width = msg.Width/2 - 1
			m.rightViewport.Width = msg.Width/2 - 1
			m.leftViewport.Height = viewportHeight
			m.rightViewport.Height = viewportHeight
		}

		m.updateContent()
	}

	// Sync both viewports
	m.leftViewport, cmd = m.leftViewport.Update(msg)
	m.rightViewport.YOffset = m.leftViewport.YOffset
	m.rightViewport.YPosition = m.leftViewport.YPosition

	return m, cmd
}

func (m *model) updateContent() {
	if len(m.files) == 0 || m.activeTab >= len(m.files) {
		return
	}

	content := m.files[m.activeTab].content
	colWidth := m.width/2 - 8 // Account for line numbers (6 chars + space)
	fullWidth := m.width - 1  // Full width minus divider
	var leftLines, rightLines []string

	m.leftLineNum = 0
	m.rightLineNum = 0
	m.indexInfo = ""

	for _, line := range content {
		left, right, isFullWidth, skip := m.formatLine(line, colWidth, fullWidth)
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

	m.leftViewport.SetContent(strings.Join(leftLines, "\n"))
	m.rightViewport.SetContent(strings.Join(rightLines, "\n"))
}

func (m *model) formatLine(line string, width int, fullWidth int) (string, string, bool, bool) {
	if len(line) == 0 {
		return "", "", false, false
	}

	// Skip diff --git, ---, +++ lines (filename info redundant with tabs)
	if strings.HasPrefix(line, "diff --git") ||
		strings.HasPrefix(line, "---") ||
		strings.HasPrefix(line, "+++") {
		return "", "", false, true // skip = true
	}

	// Capture and skip index line, we'll show it differently
	if strings.HasPrefix(line, "index ") {
		m.indexInfo = strings.TrimPrefix(line, "index ")
		return "", "", false, true // skip = true
	}

	if strings.HasPrefix(line, "@@") {
		// Extract line numbers from @@ header
		parts := strings.Split(line, "@@")
		hunkInfo := ""
		if len(parts) >= 2 {
			nums := strings.TrimSpace(parts[1])
			hunkInfo = nums
			// Parse "-leftStart,leftCount +rightStart,rightCount"
			if strings.Contains(nums, "-") && strings.Contains(nums, "+") {
				leftPart := strings.Split(strings.Split(nums, "+")[0], "-")[1]
				rightPart := strings.Split(nums, "+")[1]

				if leftNum := strings.Split(strings.TrimSpace(leftPart), ","); len(leftNum) > 0 {
					fmt.Sscanf(leftNum[0], "%d", &m.leftLineNum)
				}
				if rightNum := strings.Split(strings.TrimSpace(rightPart), ","); len(rightNum) > 0 {
					fmt.Sscanf(rightNum[0], "%d", &m.rightLineNum)
				}
			}
		}

		// Build compact header: "index • @@ hunk @@"
		compactHeader := ""
		if m.indexInfo != "" {
			compactHeader = fmt.Sprintf("index %s • @@ %s @@", m.indexInfo, hunkInfo)
		} else {
			compactHeader = fmt.Sprintf("@@ %s @@", hunkInfo)
		}

		formatted := headerStyle.Render(padRight(truncate(compactHeader, fullWidth), fullWidth))
		return formatted, "", true, false
	}

	// Handle diff lines
	if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
		// Removed line - show on left only with full background
		text := line[1:]
		if len(text) > width {
			text = text[:width-3] + "..."
		}
		// Use lipgloss Width to set fixed width for background
		contentStyle := leftBgStyle.Copy().Width(width)
		lineNum := fmt.Sprintf("%5d ", m.leftLineNum)
		left := lineNumBgLeft.Render(lineNum) + contentStyle.Render(text)

		// Right side empty with neutral background
		emptyStyle := neutralStyle.Copy().Width(width)
		right := "      " + emptyStyle.Render("")
		m.leftLineNum++
		return left, right, false, false
	}

	if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
		// Added line - show on right only with full background
		text := line[1:]
		if len(text) > width {
			text = text[:width-3] + "..."
		}
		// Use lipgloss Width to set fixed width for background
		contentStyle := rightBgStyle.Copy().Width(width)
		lineNum := fmt.Sprintf("%5d ", m.rightLineNum)

		// Left side empty with neutral background
		emptyStyle := neutralStyle.Copy().Width(width)
		left := "      " + emptyStyle.Render("")
		right := lineNumBgRight.Render(lineNum) + contentStyle.Render(text)
		m.rightLineNum++
		return left, right, false, false
	}

	// Context line - show on both sides
	text := truncate(line, width)
	leftNum := fmt.Sprintf("%5d ", m.leftLineNum)
	rightNum := fmt.Sprintf("%5d ", m.rightLineNum)
	left := lineNumStyle.Render(leftNum) + neutralStyle.Render(padRight(text, width))
	right := lineNumStyle.Render(rightNum) + neutralStyle.Render(padRight(text, width))
	m.leftLineNum++
	m.rightLineNum++
	return left, right, false, false
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func padRight(s string, length int) string {
	// Strip ANSI codes for length calculation
	visibleLen := len(stripAnsi(s))
	if visibleLen >= length {
		return s
	}
	return s + strings.Repeat(" ", length-visibleLen)
}

func stripAnsi(s string) string {
	// Simple ANSI code stripper for length calculation
	result := ""
	inEscape := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		result += string(r)
	}
	return result
}

func (m model) View() string {
	if !m.ready {
		return "Loading..."
	}

	// Render tabs (only if multiple files)
	var tabBar string
	if len(m.files) > 1 {
		var tabs []string
		for i, file := range m.files {
			style := inactiveTabStyle
			if i == m.activeTab {
				style = activeTabStyle
			}
			tabLabel := file.name
			if len(tabLabel) > 20 {
				tabLabel = tabLabel[:17] + "..."
			}
			tabs = append(tabs, style.Render(tabLabel))
		}
		tabBar = lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

		// Add gap to fill the rest of the width
		tabBarWidth := len(stripAnsi(tabBar))
		if tabBarWidth < m.width {
			gap := tabGapStyle.Render(strings.Repeat(" ", m.width-tabBarWidth))
			tabBar = tabBar + gap
		}
		tabBar = tabBar + "\n"
	}

	divider := dividerStyle.Render("│")

	leftView := m.leftViewport.View()
	rightView := m.rightViewport.View()

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

	// Minimal help text
	help := helpStyle.Render("↑↓:scroll h/l:file 1-9:jump q:quit")

	if tabBar != "" {
		return fmt.Sprintf("%s%s\n%s", tabBar, content, help)
	}
	return fmt.Sprintf("%s\n%s", content, help)
}

func parseDiffIntoFiles(lines []string) []fileDiff {
	var files []fileDiff
	var currentFile *fileDiff

	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git") {
			// Extract file name
			parts := strings.Fields(line)
			fileName := "unknown"
			if len(parts) >= 4 {
				fileName = parts[3]
				if strings.HasPrefix(fileName, "b/") {
					fileName = fileName[2:]
				}
			}

			// Save previous file if exists
			if currentFile != nil {
				files = append(files, *currentFile)
			}

			// Start new file
			currentFile = &fileDiff{
				name:    fileName,
				content: []string{line},
			}
		} else if currentFile != nil {
			currentFile.content = append(currentFile.content, line)
		}
	}

	// Don't forget the last file
	if currentFile != nil {
		files = append(files, *currentFile)
	}

	return files
}

func readDiff() ([]string, error) {
	var input io.Reader

	// Check if stdin has data
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Data is being piped to stdin
		input = os.Stdin
	} else {
		// No piped data, run git diff
		cmd := exec.Command("git", "diff")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return nil, fmt.Errorf("failed to create pipe: %w", err)
		}

		if err := cmd.Start(); err != nil {
			return nil, fmt.Errorf("failed to run git diff: %w", err)
		}

		input = stdout
		defer cmd.Wait()
	}

	var lines []string
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading diff: %w", err)
	}

	return lines, nil
}

func main() {
	lines, err := readDiff()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(lines) == 0 {
		fmt.Println("No diff to display")
		os.Exit(0)
	}

	files := parseDiffIntoFiles(lines)

	if len(files) == 0 {
		fmt.Println("No files in diff")
		os.Exit(0)
	}

	m := model{
		files:     files,
		activeTab: 0,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
