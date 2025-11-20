# Diff Viewer - Usage Examples

## Piping from Git (Most Common!)

The most powerful feature is piping `git diff` output directly into the viewer:

### See all uncommitted changes
```bash
git diff | python3 diffview.py
```

### See changes between commits
```bash
git diff HEAD~1 | python3 diffview.py
```

### See changes between branches
```bash
git diff main develop | python3 diffview.py
```

### See staged changes
```bash
git diff --staged | python3 diffview.py
```

### See changes in a specific commit
```bash
git diff HEAD~1 HEAD | python3 diffview.py
```

### See changes in a specific file
```bash
git diff HEAD -- src/app.py | python3 diffview.py
```

## Piping from Standard Diff

### Compare two files
```bash
diff file1.txt file2.txt | python3 diffview.py
```

### With unified diff format (default)
```bash
diff -u old.py new.py | python3 diffview.py
```

### Recursive directory diff
```bash
diff -u -r dir1/ dir2/ | python3 diffview.py
```

## Direct File Comparison

Sometimes you just want to compare files without using git:

```bash
# Compare two specific files
python3 diffview.py version1.txt version2.txt

# Compare code files
python3 diffview.py old_script.py new_script.py

# Compare configuration files
python3 diffview.py config.old.json config.new.json
```

## Advanced Git Usage

### Create aliases for quick access

Add to your `.bashrc` or `.zshrc`:
```bash
# View all changes
alias gd='git diff | python3 ~/diffview/diffview.py'

# View staged changes
alias gds='git diff --staged | python3 ~/diffview/diffview.py'

# View changes from last commit
alias gdh='git diff HEAD~1 | python3 ~/diffview/diffview.py'
```

Then use:
```bash
gd           # See all changes
gds          # See staged changes
gdh          # See changes from last commit
```

### View changes from pull request

```bash
# If you have a PR branch checked out
git diff main | python3 diffview.py

# Or compare remote branches
git diff origin/main origin/feature-branch | python3 diffview.py
```

### Review before committing

```bash
# See what you're about to commit
git diff --staged | python3 diffview.py

# Stage, review, then commit
git add .
git diff --staged | python3 diffview.py
# If happy:
git commit -m "My changes"
```

## Practical Workflows

### Code Review Workflow

```bash
# 1. Fetch latest changes
git fetch origin

# 2. Review what changed
git diff origin/main | python3 diffview.py

# 3. View specific feature branch
git diff origin/feature-branch | python3 diffview.py

# 4. See the actual commits
git log -p | less  # or use your preferred method
```

### Before Rebasing

```bash
# See what commits you're about to rebase
git diff main | python3 diffview.py

# Then rebase safely
git rebase main
```

### Comparing Different Versions

```bash
# Compare current code with a tag
git diff v1.0.0 | python3 diffview.py

# Compare two tags
git diff v1.0.0..v2.0.0 | python3 diffview.py
```

## Tips & Tricks

### Combining with other tools

```bash
# View diff, then go straight to git add
git diff | python3 diffview.py
git add .  # After reviewing in diffview

# Check diff before push
git diff origin/main | python3 diffview.py
git push origin main
```

### Large diffs

For very large diffs with many files:
```bash
# View just one file
git diff -- path/to/file.py | python3 diffview.py

# Or a directory
git diff -- src/components/ | python3 diffview.py
```

### Creating patches

```bash
# Save diff to file
git diff > changes.patch

# View it later
cat changes.patch | python3 diffview.py

# Or directly
python3 diffview.py file1 file2 > changes.patch
```

## Integration with CI/CD

### Pre-commit hook

Create `.git/hooks/pre-commit`:
```bash
#!/bin/bash
git diff --staged | python3 ~/diffview/diffview.py
# Review changes, then exit the viewer
# If you quit (q), the commit will proceed
```

Make it executable:
```bash
chmod +x .git/hooks/pre-commit
```

## Keyboard Shortcuts While Reviewing

Once in the diff viewer:

| Key | Action |
|-----|--------|
| `j` | Scroll down to see more changes |
| `k` | Scroll up to previous changes |
| `f` | Page down through large diffs |
| `b` | Page up |
| `g` | Jump to the beginning |
| `G` | Jump to the end |
| `q` | Exit and continue with git |

## Real-World Examples

### Review PR before merging

```bash
# Checkout the PR branch
git checkout feature/new-feature

# See what you're about to merge
git diff main | python3 diffview.py

# If happy, merge it
git checkout main
git merge feature/new-feature
```

### Compare development vs production

```bash
# What code is different from production?
git diff main production | python3 diffview.py

# What will be released?
git diff v1.2.3 main | python3 diffview.py
```

### Identify what changed in a bug fix

```bash
# Find the commit that fixed something
git log --oneline | grep -i fix

# See the actual changes
git diff commit-hash~1 commit-hash | python3 diffview.py
```

---

The key insight: **Any place you'd use `git diff` to see changes, you can pipe it to diffview for a beautiful side-by-side comparison!**
