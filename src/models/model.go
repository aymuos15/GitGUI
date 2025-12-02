package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// Init initializes the Bubble Tea model
func (m Model) Init() tea.Cmd {
	return nil
}

// InitFilterInput initializes the text input component for filter entry
func (m *Model) InitFilterInput(placeholder string) {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 40
	m.FilterInput = ti
}

// FilterAppliedMsg is sent when a filter has been applied and views need updating
type FilterAppliedMsg struct{}

// Update handles Bubble Tea messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle filter input mode first
	if m.FilterMode != "" {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				// Apply the filter
				value := m.FilterInput.Value()
				switch m.FilterMode {
				case "author":
					m.LogFilters.Author = value
				case "path":
					m.LogFilters.Path = value
				case "date_from":
					m.LogFilters.DateFrom = value
				case "date_to":
					m.LogFilters.DateTo = value
				case "search":
					if m.ViewMode == "log" {
						m.LogFilters.Search = value
					} else if m.ViewMode == "diff" {
						m.DiffSearch.Query = value
						m.DiffSearch.CurrentMatch = 0
						// Matches will be calculated in view rendering
					}
				case "status":
					m.StatsFilters.Status = value
				case "extension":
					m.StatsFilters.Extension = value
				}
				m.FilterMode = ""
				m.ViewChanged = true
				return m, func() tea.Msg { return FilterAppliedMsg{} }
			case "esc":
				// Cancel filter entry
				m.FilterMode = ""
				return m, nil
			default:
				// Update text input
				m.FilterInput, cmd = m.FilterInput.Update(msg)
				return m, cmd
			}
		}
		// For non-key messages, update the text input
		m.FilterInput, cmd = m.FilterInput.Update(msg)
		return m, cmd
	}

	// Handle viewport/table updates FIRST based on current view mode
	// This allows tables to consume key events for scrolling before we process them
	if m.ViewMode == "log" {
		m.LogTable, cmd = m.LogTable.Update(msg)
	} else if m.ViewMode == "stats" {
		m.StatsTable, cmd = m.StatsTable.Update(msg)
	} else if m.ViewMode == "diff" {
		m.LeftViewport, cmd = m.LeftViewport.Update(msg)
		m.RightViewport.YOffset = m.LeftViewport.YOffset
		m.RightViewport.YPosition = m.LeftViewport.YPosition
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		keyStr := msg.String()

		// Global quit - but only 'q' and 'ctrl+c' quit globally
		// 'esc' now clears filters/search in context
		switch keyStr {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			// Clear search/filters based on current view
			if m.ViewMode == "diff" && m.DiffSearch.Query != "" {
				m.DiffSearch.Query = ""
				m.DiffSearch.Matches = nil
				m.DiffSearch.CurrentMatch = 0
				return m, nil
			}
			return m, tea.Quit
		case "a":
			// Toggle auto-reload
			m.AutoReloadEnabled = !m.AutoReloadEnabled
		case "s":
			// Toggle stats view
			if m.ViewMode == "stats" {
				m.ViewMode = "diff"
			} else {
				m.ViewMode = "stats"
			}
		case "l":
			// Show log view
			if m.ViewMode != "log" {
				m.ViewMode = "log"
				m.ViewChanged = true
			}
		case "d":
			// Return to diff view
			m.ViewMode = "diff"

		// Filter shortcuts for log view
		case "f":
			// 'f' alone does nothing, wait for second key
			// This is handled by checking for 'fa', 'fp', 'fd', 'fc', 'fs', 'fe'
		case "alt+a", "ctrl+a":
			// Filter by author (log view) - using alt+a as fa isn't easily captured
			if m.ViewMode == "log" {
				m.FilterMode = "author"
				m.InitFilterInput("author name...")
				m.FilterInput.SetValue(m.LogFilters.Author)
				return m, textinput.Blink
			}
		case "alt+p", "ctrl+p":
			// Filter by path (log view)
			if m.ViewMode == "log" {
				m.FilterMode = "path"
				m.InitFilterInput("file path...")
				m.FilterInput.SetValue(m.LogFilters.Path)
				return m, textinput.Blink
			}
		case "alt+d":
			// Filter by date from (log view)
			if m.ViewMode == "log" {
				m.FilterMode = "date_from"
				m.InitFilterInput("from date (YYYY-MM-DD)...")
				m.FilterInput.SetValue(m.LogFilters.DateFrom)
				return m, textinput.Blink
			}
		case "alt+t":
			// Filter by date to (log view)
			if m.ViewMode == "log" {
				m.FilterMode = "date_to"
				m.InitFilterInput("to date (YYYY-MM-DD)...")
				m.FilterInput.SetValue(m.LogFilters.DateTo)
				return m, textinput.Blink
			}
		case "alt+c", "ctrl+l":
			// Clear all filters
			if m.ViewMode == "log" {
				m.LogFilters = LogFilterState{}
				m.ViewChanged = true
				return m, func() tea.Msg { return FilterAppliedMsg{} }
			} else if m.ViewMode == "stats" {
				m.StatsFilters = StatsFilterState{}
				m.ViewChanged = true
			}
		case "/":
			// Search - works in log and diff views
			if m.ViewMode == "log" {
				m.FilterMode = "search"
				m.InitFilterInput("search commits...")
				m.FilterInput.SetValue(m.LogFilters.Search)
				return m, textinput.Blink
			} else if m.ViewMode == "diff" {
				m.FilterMode = "search"
				m.InitFilterInput("search in diff...")
				m.FilterInput.SetValue(m.DiffSearch.Query)
				return m, textinput.Blink
			}
		case "n":
			// Next search match (diff view)
			if m.ViewMode == "diff" && len(m.DiffSearch.Matches) > 0 {
				m.DiffSearch.CurrentMatch = (m.DiffSearch.CurrentMatch + 1) % len(m.DiffSearch.Matches)
			}
		case "N":
			// Previous search match (diff view)
			if m.ViewMode == "diff" && len(m.DiffSearch.Matches) > 0 {
				m.DiffSearch.CurrentMatch--
				if m.DiffSearch.CurrentMatch < 0 {
					m.DiffSearch.CurrentMatch = len(m.DiffSearch.Matches) - 1
				}
			}

		// Stats view filters
		case "alt+s":
			// Filter by status (stats view)
			if m.ViewMode == "stats" {
				m.FilterMode = "status"
				m.InitFilterInput("status (N/M/D/R/U)...")
				m.FilterInput.SetValue(m.StatsFilters.Status)
				return m, textinput.Blink
			}
		case "alt+e":
			// Filter by extension (stats view)
			if m.ViewMode == "stats" {
				m.FilterMode = "extension"
				m.InitFilterInput("extension (e.g., .go)...")
				m.FilterInput.SetValue(m.StatsFilters.Extension)
				return m, textinput.Blink
			}

		case "tab", "right":
			if m.ViewMode == "diff" && m.ActiveTab < len(m.Files)-1 {
				m.ActiveTab++
			}
		case "shift+tab", "left", "h":
			if m.ViewMode == "diff" && m.ActiveTab > 0 {
				m.ActiveTab--
			}
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			if m.ViewMode == "diff" {
				tabNum := int(keyStr[0] - '1')
				if tabNum < len(m.Files) {
					m.ActiveTab = tabNum
				}
			}
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		// Calculate viewport height: total - tabs (always shown) - help line
		viewportHeight := msg.Height - 2 // 1 for tabs, 1 for help

		// Full width split 50/50 between left and right columns
		leftColWidth := (msg.Width - 1) / 2             // Left column (accounting for center divider)
		rightColWidth := (msg.Width - 1) - leftColWidth // Right column gets remaining space

		if !m.Ready {
			m.LeftViewport = viewport.New(leftColWidth, viewportHeight)
			m.RightViewport = viewport.New(rightColWidth, viewportHeight)
			m.Ready = true
		} else {
			m.LeftViewport.Width = leftColWidth
			m.RightViewport.Width = rightColWidth
			m.LeftViewport.Height = viewportHeight
			m.RightViewport.Height = viewportHeight
		}
	}

	return m, cmd
}

// View returns empty string - actual view rendering is done by appWrapper to avoid circular imports
func (m Model) View() string {
	return ""
}
