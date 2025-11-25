# Changelog

## [0.1.3] - 2025-11-25

### Added
- **Release Tags in Log View**: Release tags (e.g., v0.1.2) are now displayed inline in the message column for tagged commits
- **Dynamic Page Size**: Stats and log views now automatically adjust page size based on terminal height for better usability

### Fixed
- **Log View Always Shows History**: Log view now displays commit history even when there are no current diff changes
- **Log View Updates on Switch**: Fixed issue where log content wasn't updating when switching to log view
- **Safe String Truncation**: Fixed index out of range panics in diff view by implementing safe truncation in all code paths
- **Robust Log View Dimensions**: Added author column with adaptive sizing and guards against uninitialized dimensions

## [0.1.2] - 2025-11-22

### Fixed
- **Fullscreen Layout for Stats and Log Views**: Tables now properly expand to fill the entire terminal width, matching the diff view behavior
- Stats and log tables now dynamically calculate column widths to use available screen space
- Help bar positioning is now consistent across all views (diff, stats, log)
- Separator line in stats view now extends end-to-end across all columns
- Tables are vertically centered with help bar fixed at the bottom

## [0.1.1] - 2025-11-22

### Added
- **Separate Branch and Origin Columns in Log View**: Two new dedicated columns display local branch names and remote tracking branches for better clarity
- Branch detection now intelligently distinguishes between:
  - Local branches (e.g., `master`, `feature/new`, `bugfix`)
  - Remote branches (e.g., `origin/master`, `upstream/main`)
- Support for multiple remote prefixes: `origin/`, `upstream/`, `remote/`
- Automatic tag filtering while preserving branch information

### Changed
- Updated git log command with `--decorate=short` flag to fetch branch decorations
- Reorganized log view table columns: Hash → Branch → Origin → Graph → Message → Time
- Removed Author column to make room for branch information (prioritized for better UX)
- Applied code formatting with `go fmt` for consistency

### Technical Details
- Implemented smart git decoration parsing to extract local and remote branch refs
- Only treats git decorations as branch info, avoiding confusion with time format
- Handles edge cases like branches with slashes (e.g., `feature/test/branch`)

## [0.1.0] - Pilot Release

### Initial Release
- Side-by-side diff view with syntax highlighting
- Tabbed interface for multiple changed files
- Git log viewer with commit history
- Statistics view with file changes summary
- Auto-reload on git changes with fsnotify
- Color-coded help bar with keyboard shortcuts
- Line numbers for both old and new code versions
- Full scrolling with vim-style keybindings
- Untracked file detection and display
