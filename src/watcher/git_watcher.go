package watcher

import (
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
)

type GitChangeMsg struct {}

func WatchGitChanges() tea.Cmd {
	return func() tea.Msg {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			// If we can't create watcher, just return nil - app continues without auto-reload
			return nil
		}
		defer watcher.Close()

		// Watch git directory files for staging/commit changes
		gitPaths := []string{
			filepath.Join(".git", "index"),
			filepath.Join(".git", "HEAD"),
			filepath.Join(".git", "refs", "heads"),
			filepath.Join(".git", "refs", "remotes"),
			filepath.Join(".git", "refs", "remotes", "origin"), // Watch origin remotes specifically
		}

		for _, path := range gitPaths {
			if err := watcher.Add(path); err != nil {
				// If specific path fails, continue - watch what we can
				continue
			}
		}

		// Get all tracked files and watch their parent directories
		// This handles editors that use atomic saves (create temp file, rename)
		cmd := exec.Command("git", "ls-files")
		output, err := cmd.Output()
		if err == nil {
			// Track unique directories to avoid adding the same directory multiple times
			dirsToWatch := make(map[string]bool)

			files := strings.Split(strings.TrimSpace(string(output)), "\n")
			for _, file := range files {
				if file != "" {
					// Get the directory containing this file
					dir := filepath.Dir(file)
					if dir == "" || dir == "." {
						// File is in root directory, watch current directory
						dir = "."
					}
					dirsToWatch[dir] = true
				}
			}

			// Add all unique directories to watcher
			for dir := range dirsToWatch {
				if err := watcher.Add(dir); err != nil {
					// If specific directory fails, continue - watch what we can
					continue
				}
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
					event.Op&fsnotify.Remove == fsnotify.Remove ||
					event.Op&fsnotify.Rename == fsnotify.Rename {

					// Return immediately on first relevant change
					// Watcher will be restarted after refresh, providing natural debouncing
					return GitChangeMsg{}
				}

			case _, ok := <-watcher.Errors:
				if !ok {
					return nil
				}
			}
		}
	}
}
