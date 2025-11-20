package models

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// Init initializes the Bubble Tea model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles Bubble Tea messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

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
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "s":
			// Toggle stats view
			if m.ViewMode == "stats" {
				m.ViewMode = "diff"
			} else {
				m.ViewMode = "stats"
			}
		case "l":
			// Show log view
			m.ViewMode = "log"
		case "d":
			// Return to diff view
			m.ViewMode = "diff"
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
				tabNum := int(msg.String()[0] - '1')
				if tabNum < len(m.Files) {
					m.ActiveTab = tabNum
				}
			}
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		// Calculate viewport height: total - tabs (if multiple files) - help line
		viewportHeight := msg.Height - 1 // 1 for help
		if len(m.Files) > 1 {
			viewportHeight -= 1 // 1 for tabs
		}

		if !m.Ready {
			m.LeftViewport = viewport.New(msg.Width/2-1, viewportHeight)
			m.RightViewport = viewport.New(msg.Width/2-1, viewportHeight)
			m.LogViewport = viewport.New(msg.Width, viewportHeight)
			m.Ready = true
		} else {
			m.LeftViewport.Width = msg.Width/2 - 1
			m.RightViewport.Width = msg.Width/2 - 1
			m.LeftViewport.Height = viewportHeight
			m.RightViewport.Height = viewportHeight
			m.LogViewport.Width = msg.Width
			m.LogViewport.Height = viewportHeight
		}
	}

	return m, cmd
}

// View renders the appropriate view based on the current view mode
// This method is required by the tea.Model interface
func (m Model) View() string {
	// Implementation is provided by the views package to avoid circular imports
	return ""
}
