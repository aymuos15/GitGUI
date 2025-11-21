package main

import (
	"fmt"
	"os"

	"gg/src/diff"
	"gg/src/io"
	"gg/src/models"
	"gg/src/views"
	"gg/src/watcher"

	tea "github.com/charmbracelet/bubbletea"
)

// processDiffLines processes diff lines and returns files, message, and view mode
func processDiffLines(lines []string) ([]models.FileDiff, string, string) {
	var files []models.FileDiff
	var noDiffMessage string
	var viewMode string

	if len(lines) == 0 {
		noDiffMessage = "No diff to display"
		files = []models.FileDiff{} // Empty slice
		viewMode = "log"            // Default to log view when no diff
	} else {
		files = diff.ParseDiffIntoFiles(lines)
		if len(files) == 0 {
			noDiffMessage = "No files in diff"
			viewMode = "log" // Default to log view when no files
		} else {
			viewMode = "diff" // Show diff view when there are files
		}
	}

	return files, noDiffMessage, viewMode
}

func main() {
	lines, err := io.ReadDiff()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	files, noDiffMessage, viewMode := processDiffLines(lines)

	m := models.Model{
		Files:             files,
		ActiveTab:         0,
		ViewMode:          viewMode,
		NoDiffMessage:     noDiffMessage,
		AutoReloadEnabled: true, // Enable auto-reload by default
	}

	p := tea.NewProgram(&appWrapper{Model: m}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// RefreshDataMsg contains refreshed git diff data
type RefreshDataMsg struct {
	Files         []models.FileDiff
	NoDiffMessage string
	ViewMode      string
}

// refreshDiffData reads git diff and returns RefreshDataMsg
func refreshDiffData() tea.Msg {
	lines, err := io.ReadDiff()
	if err != nil {
		// On error, return empty data
		return RefreshDataMsg{
			Files:         []models.FileDiff{},
			NoDiffMessage: "Error reading diff",
			ViewMode:      "log",
		}
	}

	files, noDiffMessage, viewMode := processDiffLines(lines)

	return RefreshDataMsg{
		Files:         files,
		NoDiffMessage: noDiffMessage,
		ViewMode:      viewMode,
	}
}

// appWrapper wraps the Model to provide the View method
// This avoids circular imports between models and views packages
type appWrapper struct {
	models.Model
	logTableInit   bool
	statsTableInit bool
}

func (a *appWrapper) Init() tea.Cmd {
	// Start both model init and watcher
	return tea.Batch(
		a.Model.Init(),
		watcher.WatchGitChanges(),
	)
}

func (a *appWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle git change messages before passing to model
	switch msg := msg.(type) {
	case watcher.GitChangeMsg:
		// Only trigger refresh if auto-reload is enabled
		if a.AutoReloadEnabled {
			// Trigger refresh and restart watcher
			// Use Sequence to ensure refresh completes before watcher restarts
			// This forces Bubble Tea to render immediately
			return a, tea.Sequence(
				tea.Cmd(refreshDiffData),
				watcher.WatchGitChanges(),
			)
		} else {
			// Just restart watcher without refreshing
			return a, watcher.WatchGitChanges()
		}

	case RefreshDataMsg:
		// Update model with refreshed data
		a.Files = msg.Files
		a.NoDiffMessage = msg.NoDiffMessage
		// Don't change ViewMode - keep user in their current view
		a.ActiveTab = 0 // Reset to first tab

		// Reinitialize all views with new data
		if a.ViewMode == "diff" {
			views.UpdateContent(&a.Model)
		}

		// Always reinitialize stats and log tables with new data
		views.UpdateStatsContent(&a.Model)
		views.UpdateLogContent(&a.Model)
		a.statsTableInit = true
		a.logTableInit = true

		return a, nil
	}

	updatedModel, cmd := a.Model.Update(msg)
	a.Model = updatedModel.(models.Model)

	// Update content after model changes
	if a.ViewMode == "diff" {
		views.UpdateContent(&a.Model)
	} else if a.ViewMode == "log" {
		// Only initialize log table once
		if !a.logTableInit {
			views.UpdateLogContent(&a.Model)
			a.logTableInit = true
		}
	} else if a.ViewMode == "stats" {
		// Only initialize stats table once
		if !a.statsTableInit {
			views.UpdateStatsContent(&a.Model)
			a.statsTableInit = true
		}
	}

	return a, cmd
}

func (a *appWrapper) View() string {
	switch a.ViewMode {
	case "stats":
		return views.RenderStatsView(&a.Model)
	case "log":
		return views.RenderLogView(&a.Model)
	default:
		return views.RenderDiffView(&a.Model)
	}
}
