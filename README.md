# Git Diff Viewer

A beautiful terminal UI for viewing git diffs, built with Charmbracelet tools.

## Features

- 🎨 **Syntax highlighting** - Full syntax highlighting for Python, JavaScript, Go, and many more languages
- 📊 **Side-by-side view** - Compare old and new code directly alongside each other
- 📑 **Tabbed interface** - Easy navigation between multiple changed files
- 📈 **Statistics view** - Press `s` to see git diff --stat style summary
- 🎯 **Color-coded help bar** - Each keyboard shortcut displayed in distinct dark pastel colors for quick reference
- 🔢 **Line numbers** - See exact line numbers for both versions
- 📜 **Full scrolling** - Navigate with vim-style keybindings (j/k, arrows, page up/down)
- 🚀 **Flexible input** - Works with `git diff`, piped input, or any diff format
- ⚡ **High performance** - Cached syntax highlighting for smooth scrolling on large files
- ✨ **Beautiful UI** - Built with Charmbracelet tools (Bubbles, Bubbletea, Lipgloss, Chroma)

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

### Navigation
- `↑`/`↓` or `j`/`k` - Scroll up/down
- `Page Up`/`Page Down` - Scroll one page
- `g` - Jump to top
- `G` - Jump to bottom

### File Switching
- `tab`/`h`/`l` or `←`/`→` - Switch between files
- `1`-`9` - Jump directly to file 1-9

### Views
- `s` - Toggle statistics view (git diff --stat style)

### General
- `q`/`esc`/`ctrl+c` - Quit

## Syntax Highlighting

The viewer automatically detects the file type and applies appropriate syntax highlighting:

- **Python** - Keywords (def, class, if, etc.), strings, numbers, comments
- **JavaScript/TypeScript** - Functions, variables, strings, JSX
- **Go** - Keywords, types, functions
- **And many more** - Supports 200+ languages via Chroma

Combined with diff colors:
- 🟥 **Red background** - Removed lines
- 🟩 **Green background** - Added lines
- ⚪ **No background** - Unchanged context lines

## UI Design

The interface features a clean, color-coded help bar at the bottom where each keyboard shortcut is displayed in a distinct dark pastel color (purple, teal, brown, green, mauve) with white text, making it easy to quickly scan available commands.
