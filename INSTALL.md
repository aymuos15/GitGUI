# Installation

The `dif` command has been installed to `~/bin/dif` and is ready to use!

## Usage

From any git repository:

```bash
# View unstaged changes
dif

# View staged changes
git diff --staged | dif

# View specific commits
git diff HEAD~3..HEAD | dif

# View a specific commit
git show abc123 | dif

# Compare branches
git diff main...feature-branch | dif
```

## Keybindings

- **↑/↓** or **j/k** - Scroll up/down
- **h/l** or **←/→** or **tab** - Switch between files
- **1-9** - Jump directly to file 1-9
- **g** - Go to top
- **G** - Go to bottom
- **q** or **esc** - Quit

## Reinstallation

If you need to reinstall or update:

```bash
cd ~/diffview
go build -o diffview
cp diffview ~/bin/dif
```

## Global Installation (requires sudo)

For system-wide installation:

```bash
cd ~/diffview
sudo cp diffview /usr/local/bin/dif
```
