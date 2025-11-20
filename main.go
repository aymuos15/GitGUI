package main

import (
	"fmt"
	"os"

	"diffview/src/diff"
	"diffview/src/io"
	"diffview/src/models"
	"diffview/src/views"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	lines, err := io.ReadDiff()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

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

	m := models.Model{
		Files:         files,
		ActiveTab:     0,
		ViewMode:      viewMode,
		NoDiffMessage: noDiffMessage,
	}

	p := tea.NewProgram(&appWrapper{Model: m}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
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
	return a.Model.Init()
}

func (a *appWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
