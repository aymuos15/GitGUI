# DiffView

A beautiful, side-by-side git diff viewer for the terminal. Perfect for code review workflows and integration with OpenCode TUI.

## Features

- 📊 **Side-by-side diffs** - View old and new code directly alongside each other
- 🎨 **Syntax highlighting** - Color-coded additions (green), deletions (red), and context (gray)
- 📈 **File summaries** - Quick overview of changes across all modified files
- 📑 **Tab navigation** - Browse multiple files with interactive tabs (use `n`/`p` keys)
- 🔧 **Git integration** - Works with staged, unstaged, and historical diffs
- 🖥️ **Terminal optimized** - Responsive design that adapts to your terminal width

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

## Output Example

With a single file:
```
📊 Git Diff Viewer

Summary:
  1 file changed, +5 insertions, -2 deletions

Files:
  ✏️  src/app.ts (+5 -2)

✏️  MODIFIED src/app.ts

@@ -10,5 +10,8 @@
   10  function hello()         │   10  function hello()       
   11  −   console.log("hi")     │   11  + console.log("hello")
   12  +   return "greeting"      │   12  + return "greeting"   
   13     }                       │   13     }                  
```

With multiple files (interactive tabs):
```
📊 Git Diff Viewer

Summary:
  2 files changed, +10 insertions, -5 deletions

Files:
  ✏️  src/app.ts (+5 -2)
  ✏️  src/utils.ts (+5 -3)

 app.ts  utils.ts 
────────────────────────────────────────────────────────────────────

✏️  MODIFIED src/app.ts
[diff content...]

File 1/2 — Press 'n' for next, 'p' for prev, 'q' to quit, 'h' for help
>
```

When you have multiple files, use the following keyboard commands:
- `n` or `next` - Go to the next file
- `p` or `prev` - Go to the previous file
- `h` or `help` - Show help message
- `q` or `quit` - Exit the viewer

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
