# Quick Start

## Installation

```bash
pip install rich
```

## Basic Usage

View git diff in side-by-side format:

```bash
git diff | python3 diffview.py
```

## Common Commands

```bash
# Staged changes
git diff --staged | python3 diffview.py

# Latest commit
git diff HEAD~1 | python3 diffview.py

# Compare branches
git diff main develop | python3 diffview.py

# Standard diff
diff file1.txt file2.txt | python3 diffview.py
```

## Output

The viewer displays two columns:
- **Left**: Original lines
- **Right**: Modified lines

Colors:
- 🔴 Red = Removed
- 🟢 Green = Added
- ⚪ Dim = Unchanged

## Scrolling

```bash
# Scroll with less
git diff | python3 diffview.py | less

# View first 50 lines
git diff | python3 diffview.py | head -50
```

## Create an Alias

Add to your `.bashrc` or `.zshrc`:

```bash
alias gd='git diff | python3 ~/diffview/diffview.py'
```

Then use:
```bash
gd
gd | less
```

That's it! Happy diffing! 🎉
