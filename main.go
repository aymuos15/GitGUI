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

	if len(lines) == 0 {
		fmt.Println("No diff to display")
		os.Exit(0)
	}

	files := diff.ParseDiffIntoFiles(lines)

	if len(files) == 0 {
		fmt.Println("No files in diff")
		os.Exit(0)
	}

	m := models.Model{
		Files:     files,
		ActiveTab: 0,
		ViewMode:  "diff",
	}

	p := tea.NewProgram(&appWrapper{m}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// appWrapper wraps the Model to provide the View method
// This avoids circular imports between models and views packages
type appWrapper struct {
	models.Model
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
		views.UpdateLogContent(&a.Model)
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
