# Git Diff Viewer

A beautiful terminal UI for viewing git diffs, built with Charmbracelet tools.

## Features

- 🎨 Syntax highlighting for git diffs
- 📜 Full scrolling support with vim-style keybindings
- 🚀 Can read from `git diff` or stdin
- ✨ Built with Bubbles, Bubbletea, and Lipgloss

## Installation

```bash
go build -o diffview
```

## Usage

View current git diff:
```bash
./diffview
```

View staged changes:
```bash
git diff --staged | ./diffview
```

View diff between commits:
```bash
git diff HEAD~5..HEAD | ./diffview
```

View diff from a file:
```bash
git show <commit> | ./diffview
```

## Keybindings

- `↑`/`k` - Scroll up
- `↓`/`j` - Scroll down
- `Page Up` - Scroll up one page
- `Page Down` - Scroll down one page
- `g` - Go to top
- `G` - Go to bottom
- `q`/`esc`/`ctrl+c` - Quit

## Color Scheme

- **Green**: Added lines (+)
- **Red**: Removed lines (-)
- **Cyan**: File headers (---, +++)
- **Yellow**: Index information
- **Gray**: Context markers (@@)
- **Blue**: Diff headers (diff --git)
