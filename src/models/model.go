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
