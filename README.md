# Diff Viewer

A side-by-side diff viewer for git diffs in the terminal.

## Features

- **Side-by-side display**: Compare files with aligned columns
- **Color coding**: 
  - Red for removed lines
  - Green for added lines
  - Dim for context lines
- **Line numbers**: Track position in both files
- **Git integration**: Works with `git diff` output
- **Simple and fast**: Lightweight, no complex UI framework

## Installation

```bash
pip install rich
```

## Usage

```bash
git diff | python3 diffview.py
```

### Examples

```bash
# See all unstaged changes
git diff | python3 diffview.py

# See staged changes
git diff --staged | python3 diffview.py

# Compare commits
git diff HEAD~1 | python3 diffview.py

# Compare branches
git diff main develop | python3 diffview.py

# Use with standard diff command
diff file1.txt file2.txt | python3 diffview.py
```

## Keyboard Shortcuts

Currently this is a static display - scroll with your terminal or pipe through `less`:

```bash
git diff | python3 diffview.py | less
```

## How It Works

1. Reads unified diff format from stdin (git diff output)
2. Parses the diff into structured lines
3. Formats as a side-by-side table
4. Applies color coding
5. Displays in your terminal

## Output Format

```
    1 Original line        1 Original line        
    2 Context line        2 Context line         
    3 Removed line ❌                            
                          3 Added line ✓         
```

## Project Structure

```
diffview/
├── diffview.py          - Main application
├── requirements.txt     - Dependencies
├── README.md           - This file
├── QUICKSTART.md       - Quick start guide
└── test_file*.txt      - Test files
```

## Examples

### Review changes before committing

```bash
git add .
git diff --staged | python3 diffview.py
```

### Compare feature branches

```bash
git diff main feature-branch | python3 diffview.py
```

### View commit changes

```bash
git diff HEAD~1 HEAD | python3 diffview.py
```

## Tips

- Pipe through `less` for interactive scrolling: `git diff | python3 diffview.py | less`
- Pipe through `head` to see first N lines: `git diff | python3 diffview.py | head -50`
- Create an alias for quick access:
  ```bash
  alias gd='git diff | python3 ~/diffview/diffview.py'
  ```

## License

MIT
