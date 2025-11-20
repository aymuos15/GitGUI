package watcher

import (
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
)

type GitChangeMsg struct {
	ChangeType string
}

func WatchGitChanges() tea.Cmd {
	return func() tea.Msg {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			// If we can't create watcher, just return nil - app continues without auto-reload
			return nil
		}
		defer watcher.Close()

		// Watch git directory files
		gitPaths := []string{
			filepath.Join(".git", "index"),
			filepath.Join(".git", "HEAD"),
			filepath.Join(".git", "refs", "heads"),
			filepath.Join(".git", "refs", "remotes"),
		}

		for _, path := range gitPaths {
			if err := watcher.Add(path); err != nil {
				// If specific path fails, continue - watch what we can
				continue
			}
		}

		// Wait for first change event
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return nil
				}

				// Only trigger on write and create events (ignore chmod, etc.)
				if event.Op&fsnotify.Write == fsnotify.Write ||
				   event.Op&fsnotify.Create == fsnotify.Create ||
				   event.Op&fsnotify.Remove == fsnotify.Remove {

					// Return immediately on first relevant change
					// Watcher will be restarted after refresh, providing natural debouncing
					return GitChangeMsg{ChangeType: "git_change"}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return nil
				}
				// Silently ignore errors and continue watching
				_ = err
			}
		}
	}
}
