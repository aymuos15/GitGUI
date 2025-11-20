# DiffView

A beautiful, interactive git diff viewer for the terminal. Like OpenCode, but for viewing diffs. Full TUI with tabs, scrolling, and keyboard navigation.

## Features

- 📊 **Side-by-side diffs** - View old and new code directly alongside each other
- 🎨 **Syntax highlighting** - Color-coded additions (green), deletions (red), and context (gray)
- 📈 **File summaries** - Quick overview of changes across all modified files
- 📑 **Interactive tabs** - Browse multiple files with highlighted tabs
- ⌨️ **Smooth scrolling** - Navigate with j/k, arrows, page up/down
- 🔧 **Git integration** - Works with staged, unstaged, and historical diffs
- 🖥️ **Full TUI** - Clean interface with no terminal background, just like OpenCode

## Installation

```bash
npm install -g diffview
```

Or build from source:

```bash
git clone <repository>
cd diffview
npm install
npm run build
npm install -g .
```

## Usage

### View unstaged changes

```bash
diffview
```

### View staged changes

```bash
diffview --staged
# or
diffview -s
```

### View changes against a specific commit/branch

```bash
diffview HEAD
diffview main
diffview v1.0.0
```

### Disable colored output

```bash
diffview --no-color
```

### Show help

```bash
diffview --help
```

## Integration with OpenCode

To use DiffView within OpenCode, you can run it as a custom command:

```bash
diffview | opencode
```

Or integrate it as a custom tool in your `opencode.json`:

```json
{
  "tools": {
    "show-diff": {
      "description": "Show side-by-side git diff",
      "command": "diffview"
    }
  }
}
```

## Interface

DiffView launches a full terminal UI with tabs and scrolling:

```
┏─────────────────────────────────────────────┬─────────────────────────────────────────────┓
┃ 📊 Git Diff Viewer                          │                                             ┃
┃ 2 files changed  •  +10  -5                 │                                             ┃
┃                                             │                                             ┃
┃ Current: src/app.ts                        │                                             ┃
┣─────────────────────────────────────────────┼─────────────────────────────────────────────┫
┃ Files                                       │ ✏️  MODIFIED src/app.ts                      ┃
┃ ──────────────────────────               │                                             ┃
┃  ◆ app.ts (+5,-2)                          │ @@ -10,5 +10,8 @@                          ┃
┃  ◆ utils.ts (+5,-3)                        │    10 − function hello()  │   10 + function │
┃                                             │    11 −   console.log...  │   11 +  console │
┃                                             │    12 +   return "gree... │   12 +  return  │
┃                                             │    13    }               │   13    }       │
┃                                             │                                             ┃
┣─────────────────────────────────────────────┴─────────────────────────────────────────────┫
┃ [n]ext  [p]rev  [j/k] scroll  [?] help  [q] quit                  File 1/2            ┃
┗─────────────────────────────────────────────────────────────────────────────────────────┛
```

**Color scheme:**
- **Red** for deleted lines
- **Green** for added lines
- **Blue** for file paths and hunk headers
- **Cyan** for highlighted active file in sidebar
- **White background** with dark text for clean appearance

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `n` | Next file |
| `p` | Previous file |
| `j` / `↓` | Scroll down |
| `k` / `↑` | Scroll up |
| `Page Down` | Scroll down (full page) |
| `Page Up` | Scroll up (full page) |
| `?` / `h` | Show help dialog |
| `q` / `Ctrl+C` | Quit |

## Architecture

### Parser (`src/parser.ts`)
Parses raw git diff output into structured data:
- Extracts file changes and status
- Identifies hunks (sections of changes)
- Tracks line numbers and change types

### Viewer (`src/viewer.ts`)
Renders the parsed diff data:
- Formats output for terminal display
- Handles side-by-side layout
- Applies color coding and truncation

### CLI (`src/index.ts`)
Command-line interface:
- Argument parsing
- Summary calculation
- Tool invocation

## Development

```bash
# Install dependencies
npm install

# Build TypeScript
npm run build

# Run development mode
npm run dev

# Run the CLI
node dist/index.js
```

## License

MIT
