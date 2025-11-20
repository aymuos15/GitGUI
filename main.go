package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Soft, subtle background colors like GitHub
	leftBgStyle      = lipgloss.NewStyle().Background(lipgloss.Color("#3d1e1e"))       // Very subtle red tint
	rightBgStyle     = lipgloss.NewStyle().Background(lipgloss.Color("#1e3d1e"))       // Very subtle green tint
	neutralStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))             // White for context
	headerStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true) // Cyan for headers
	titleStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")) // Blue for file names
	helpStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	dividerStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))   // Gray divider
	lineNumStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Gray line numbers
	lineNumBgLeft    = lipgloss.NewStyle().Background(lipgloss.Color("#3d1e1e")).Foreground(lipgloss.Color("240"))
	lineNumBgRight   = lipgloss.NewStyle().Background(lipgloss.Color("#1e3d1e")).Foreground(lipgloss.Color("240"))
	activeTabStyle   = lipgloss.NewStyle().Background(lipgloss.Color("12")).Foreground(lipgloss.Color("15")).Bold(true).Padding(0, 2)
	inactiveTabStyle = lipgloss.NewStyle().Background(lipgloss.Color("236")).Foreground(lipgloss.Color("7")).Padding(0, 2)
	tabGapStyle      = lipgloss.NewStyle().Background(lipgloss.Color("234"))

	// ANSI codes for highlighting changed text within lines
	redHighlight   = "\x1b[48;2;90;40;40m" // Brighter red background for changed text
	greenHighlight = "\x1b[48;2;40;90;40m" // Brighter green background for changed text
	resetCode      = "\x1b[0m"
)

type diffLine struct {
	leftNum   string
	rightNum  string
	leftText  string
	rightText string
	lineType  string // "add", "remove", "context", "header"
}

type fileDiff struct {
	name           string
	content        []string
	highlightCache map[int]string   // Cache: line index -> highlighted line
	lexer          chroma.Lexer     // Cached lexer for this file type
	style          *chroma.Style    // Cached style
	formatter      chroma.Formatter // Cached formatter
	additions      int              // Number of added lines
	deletions      int              // Number of deleted lines
}

// calculateStats computes additions and deletions for a file
func (f *fileDiff) calculateStats() {
	f.additions = 0
	f.deletions = 0

	for _, line := range f.content {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			f.additions++
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			f.deletions++
		}
	}
}

// initSyntaxHighlighting initializes the lexer, style, and formatter for a file
func (f *fileDiff) initSyntaxHighlighting() {
	if f.lexer != nil {
		return // Already initialized
	}

	// Get lexer based on file extension
	ext := filepath.Ext(f.name)
	f.lexer = lexers.Get(ext)

	if f.lexer == nil {
		f.lexer = lexers.Fallback
	}

	f.lexer = chroma.Coalesce(f.lexer)

	// Cache style and formatter
	f.style = styles.Get("monokai")
	if f.style == nil {
		f.style = styles.Fallback
	}

	f.formatter = formatters.TTY16m
	f.highlightCache = make(map[int]string)
}

// highlightLine highlights a single line using cached lexer/style/formatter
func (f *fileDiff) highlightLine(lineIdx int, code string) string {
	// Check cache first
	if cached, exists := f.highlightCache[lineIdx]; exists {
		return cached
	}

	// Ensure lexer is initialized
	f.initSyntaxHighlighting()

	// Tokenize and format
	iterator, err := f.lexer.Tokenise(nil, code)
	if err != nil {
		return code
	}

	var buf strings.Builder
	err = f.formatter.Format(&buf, f.style, iterator)
	if err != nil {
		return code
	}

	result := strings.TrimRight(buf.String(), "\n")

	// Cache the result
	f.highlightCache[lineIdx] = result

	return result
}

// highlightCode applies syntax highlighting to a line of code (now uses file cache)
func highlightCode(code string, filename string, fileRef *fileDiff, lineIdx int) string {
	if fileRef == nil {
		// Fallback: no caching available
		return code
	}

	return fileRef.highlightLine(lineIdx, code)
}

// findDiffChars finds character-level differences between two strings
func findDiffChars(old, new string) ([]int, []int) {
	// Simple character-by-character comparison
	// Returns indices of changed characters in old and new strings
	oldChars := []rune(old)
	newChars := []rune(new)

	// Find common prefix
	prefixLen := 0
	for prefixLen < len(oldChars) && prefixLen < len(newChars) && oldChars[prefixLen] == newChars[prefixLen] {
		prefixLen++
	}

	// Find common suffix
	suffixLen := 0
	for suffixLen < len(oldChars)-prefixLen && suffixLen < len(newChars)-prefixLen &&
		oldChars[len(oldChars)-1-suffixLen] == newChars[len(newChars)-1-suffixLen] {
		suffixLen++
	}

	// Mark changed regions
	var oldChanges, newChanges []int
	for i := prefixLen; i < len(oldChars)-suffixLen; i++ {
		oldChanges = append(oldChanges, i)
	}
	for i := prefixLen; i < len(newChars)-suffixLen; i++ {
		newChanges = append(newChanges, i)
	}

	return oldChanges, newChanges
}

// highlightChangedWords applies brighter background to changed portions using simple word diff
func highlightChangedWords(oldLine, newLine string, isOldLine bool) string {
	// Split into words
	oldWords := strings.Fields(oldLine)
	newWords := strings.Fields(newLine)

	// Find which words changed
	maxLen := len(oldWords)
	if len(newWords) > maxLen {
		maxLen = len(newWords)
	}

	changedIndices := make(map[int]bool)
	for i := 0; i < maxLen; i++ {
		oldWord := ""
		newWord := ""
		if i < len(oldWords) {
			oldWord = oldWords[i]
		}
		if i < len(newWords) {
			newWord = newWords[i]
		}

		if oldWord != newWord {
			changedIndices[i] = true
		}
	}

	// If nothing changed or everything changed, don't highlight
	if len(changedIndices) == 0 || len(changedIndices) == maxLen {
		return ""
	}

	// Build result with highlighting on changed words
	words := oldWords
	if !isOldLine {
		words = newWords
	}

	var result strings.Builder
	highlightCode := redHighlight
	if !isOldLine {
		highlightCode = greenHighlight
	}

	for i, word := range words {
		if i > 0 {
			result.WriteString(" ")
		}

		if changedIndices[i] {
			result.WriteString(highlightCode)
			result.WriteString(word)
			result.WriteString(resetCode)
		} else {
			result.WriteString(word)
		}
	}

	return result.String()
}

type model struct {
	leftViewport   viewport.Model
	rightViewport  viewport.Model
	files          []fileDiff
	activeTab      int
	ready          bool
	width          int
	height         int
	leftLineNum    int
	rightLineNum   int
	pendingRemoved []string // Track consecutive removed lines for word-level diff
	pendingAdded   []string // Track consecutive added lines for word-level diff
	showStats      bool     // Toggle stats view with 's' key
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
		case "s":
			// Toggle stats view
			m.showStats = !m.showStats
		case "tab", "right", "l":
			if !m.showStats && m.activeTab < len(m.files)-1 {
				m.activeTab++
				m.updateContent()
			}
		case "shift+tab", "left", "h":
			if !m.showStats && m.activeTab > 0 {
				m.activeTab--
				m.updateContent()
			}
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			if !m.showStats {
				tabNum := int(msg.String()[0] - '1')
				if tabNum < len(m.files) {
					m.activeTab = tabNum
					m.updateContent()
				}
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

	for lineIdx, line := range content {
		left, right, isFullWidth, skip := m.formatLine(line, colWidth, fullWidth, lineIdx)
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

func (m *model) formatLine(line string, width int, fullWidth int, lineIdx int) (string, string, bool, bool) {
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
					fmt.Sscanf(leftNum[0], "%d", &m.leftLineNum)
				}
				if rightNum := strings.Split(strings.TrimSpace(rightPart), ","); len(rightNum) > 0 {
					fmt.Sscanf(rightNum[0], "%d", &m.rightLineNum)
				}
			}
		}

		// Just show the hunk header
		formatted := headerStyle.Render(padRight(truncate(line, fullWidth), fullWidth))
		return formatted, "", true, false
	}

	// Get current file reference for syntax highlighting
	var fileRef *fileDiff
	filename := ""
	if m.activeTab < len(m.files) {
		fileRef = &m.files[m.activeTab]
		filename = fileRef.name
	}

	// Handle diff lines
	if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
		// Removed line - show on left only with syntax highlighting
		text := line[1:]

		// Apply syntax highlighting
		highlighted := highlightCode(text, filename, fileRef, lineIdx)

		// Truncate if needed
		visibleLen := len(stripAnsi(highlighted))
		if visibleLen > width {
			// Truncate the original text and re-highlight
			text = text[:width-3] + "..."
			highlighted = highlightCode(text, filename, fileRef, lineIdx)
		}

		// Apply background color directly with ANSI codes to preserve syntax highlighting
		bgCode := "\x1b[48;2;61;30;30m" // #3d1e1e
		resetBg := "\x1b[49m"

		// Pad to width
		padding := width - visibleLen
		if padding < 0 {
			padding = 0
		}

		lineNum := fmt.Sprintf("%5d ", m.leftLineNum)
		left := lineNumBgLeft.Render(lineNum) + bgCode + highlighted + strings.Repeat(" ", padding) + resetBg

		// Right side empty with neutral background
		emptyStyle := neutralStyle.Copy().Width(width)
		right := "      " + emptyStyle.Render("")
		m.leftLineNum++
		return left, right, false, false
	}

	if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
		// Added line - show on right only with syntax highlighting
		text := line[1:]

		// Apply syntax highlighting
		highlighted := highlightCode(text, filename, fileRef, lineIdx)

		// Truncate if needed
		visibleLen := len(stripAnsi(highlighted))
		if visibleLen > width {
			// Truncate the original text and re-highlight
			text = text[:width-3] + "..."
			highlighted = highlightCode(text, filename, fileRef, lineIdx)
		}

		// Apply background color directly with ANSI codes to preserve syntax highlighting
		bgCode := "\x1b[48;2;30;61;30m" // #1e3d1e
		resetBg := "\x1b[49m"

		// Pad to width
		padding := width - visibleLen
		if padding < 0 {
			padding = 0
		}

		lineNum := fmt.Sprintf("%5d ", m.rightLineNum)

		// Left side empty with neutral background
		emptyStyle := neutralStyle.Copy().Width(width)
		left := "      " + emptyStyle.Render("")
		right := lineNumBgRight.Render(lineNum) + bgCode + highlighted + strings.Repeat(" ", padding) + resetBg
		m.rightLineNum++
		return left, right, false, false
	}

	// Context line - show on both sides with syntax highlighting
	highlighted := highlightCode(line, filename, fileRef, lineIdx)

	// Truncate if needed
	visibleLen := len(stripAnsi(highlighted))
	if visibleLen > width {
		line = line[:width-3] + "..."
		highlighted = highlightCode(line, filename, fileRef, lineIdx)
	}

	leftNum := fmt.Sprintf("%5d ", m.leftLineNum)
	rightNum := fmt.Sprintf("%5d ", m.rightLineNum)
	left := lineNumStyle.Render(leftNum) + neutralStyle.Render(padRight(highlighted, width))
	right := lineNumStyle.Render(rightNum) + neutralStyle.Render(padRight(highlighted, width))
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

	// If stats view is active, show stats instead
	if m.showStats {
		return m.renderStatsView()
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
	help := helpStyle.Render("↑↓:scroll h/l:file 1-9:jump s:stats q:quit")

	if tabBar != "" {
		return fmt.Sprintf("%s%s\n%s", tabBar, content, help)
	}
	return fmt.Sprintf("%s\n%s", content, help)
}

// renderStatsView renders the stats view with a clean modern interface using bubbles table
func (m model) renderStatsView() string {
	// Calculate totals
	totalAdditions := 0
	totalDeletions := 0

	for _, file := range m.files {
		totalAdditions += file.additions
		totalDeletions += file.deletions
	}

	// Define table columns
	columns := []table.Column{
		{Title: "File", Width: 50},
		{Title: "Added", Width: 10},
		{Title: "Removed", Width: 10},
	}

	// Build table rows
	rows := []table.Row{}
	for _, file := range m.files {
		fileName := file.name
		if len(fileName) > 50 {
			fileName = "..." + fileName[len(fileName)-47:]
		}

		rows = append(rows, table.Row{
			fileName,
			fmt.Sprintf("%d", file.additions),
			fmt.Sprintf("%d", file.deletions),
		})
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

	// Remove highlight from selected row - make it look the same as normal cells
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("0")).
		Bold(false)

	// Cell style
	s.Cell = s.Cell.Foreground(lipgloss.Color("15"))

	t.SetStyles(s)

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Render("Change Summary")

	// Build the content with table
	var boxContent strings.Builder
	boxContent.WriteString(title)
	boxContent.WriteString("\n\n")

	tableView := t.View()
	boxContent.WriteString(tableView)
	boxContent.WriteString("\n\n")

	// Calculate separator width from actual table width
	// The table adds internal spacing, so we measure the first line
	tableLines := strings.Split(tableView, "\n")
	var separatorWidth int
	if len(tableLines) > 0 {
		// Remove ANSI codes to get actual width
		separatorWidth = lipgloss.Width(tableLines[0])
	} else {
		separatorWidth = 74 // fallback
	}

	boxContent.WriteString(strings.Repeat("─", separatorWidth))
	boxContent.WriteString("\n")

	// Summary section - align with table columns
	fileCount := len(m.files)
	fileWord := "file"
	if fileCount != 1 {
		fileWord = "files"
	}

	summaryLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(fmt.Sprintf("Total: %d %s changed", fileCount, fileWord))

	// Format totals to align with the Added and Removed columns
	// Table columns: File(50) + padding + Added(10) + padding + Removed(10)
	totalAddStr := fmt.Sprintf("%10d", totalAdditions)
	totalDelStr := fmt.Sprintf("%10d", totalDeletions)

	// Apply color to the formatted strings
	totalAddStyled := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Bold(true).
		Render(totalAddStr)

	totalDelStyled := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Bold(true).
		Render(totalDelStr)

	// Build summary row with proper spacing to match table layout
	summary := fmt.Sprintf("%-50s  %s  %s", summaryLabel, totalAddStyled, totalDelStyled)
	boxContent.WriteString(summary)

	// Create a box around the content
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(2, 3).
		Align(lipgloss.Center)

	box := boxStyle.Render(boxContent.String())

	// Center the box vertically and horizontally
	centeredBox := lipgloss.Place(
		m.width,
		m.height-1, // Leave space for help at bottom
		lipgloss.Center,
		lipgloss.Center,
		box,
	)

	// Help text at the bottom
	help := helpStyle.Render("↑↓:scroll h/l:file 1-9:jump s:stats q:quit")

	return centeredBox + "\n" + help
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

	// Initialize syntax highlighting and calculate stats for all files
	for i := range files {
		files[i].initSyntaxHighlighting()
		files[i].calculateStats()
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
