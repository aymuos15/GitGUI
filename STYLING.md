# Diff Viewer - Styling Guide

## Current Theme

The diff viewer comes with a professional dark theme inspired by modern code editors.

### Color Scheme

**Background & UI:**
- Panel Background: Dark gray (`$panel`)
- Surface: Darker shade (`$surface`)
- Primary: Accent color for borders (`$primary`)
- Text: Light gray (`$text`)

**Diff Colors:**
- 🔴 **Removed Lines**: Red with dark background
- 🟢 **Added Lines**: Green with dark background
- ⚪ **Unchanged Lines**: Normal text (dim)

### Layout

```
┌────────────────────────────────────────────────┐
│  Diff Viewer - git file  ↔  git file           │  Header (compact)
├─────────────────────┬─────────────────────────┤
│                     │                         │
│  Before (removed)   │  After (added)          │
│  (left panel)       │  (right panel)          │
│                     │                         │
├─────────────────────┴─────────────────────────┤
│  q: Quit  j/k: Scroll  f/b: Page  g/G: Top/End │  Footer (compact)
└────────────────────────────────────────────────┘
```

## Customizing the Theme

### Change Background Color

Edit `diffview.py` and modify the CSS section:

```python
CSS = """
Screen {
    background: #1e1e1e;  /* Custom dark color */
}
"""
```

### Change Remove/Add Colors

Modify the color codes in the `update_display()` method:

```python
if line.line_type == "add":
    # Change from green to another color
    formatted_line = f"[#00ff00]{line_num_str}[/#00ff00] [#00ff00]{display_content}[/#00ff00]"
elif line.line_type == "remove":
    # Change from red to another color
    formatted_line = f"[#ff0000]{line_num_str}[/#ff0000] [#ff0000]{display_content}[/#ff0000]"
```

### Color Palette Options

**Popular Color Combinations:**

Dark Terminal Themes:
- Solarized Dark
- Dracula
- One Dark
- Nord
- Gruvbox Dark

Light Terminal Themes:
- Solarized Light
- GitHub Light
- Nord Light

## Textual CSS Variables

The viewer uses Textual's built-in color variables:

| Variable | Usage |
|----------|-------|
| `$primary` | Borders, accents |
| `$secondary` | Alternative accents |
| `$panel` | Panel backgrounds |
| `$boost` | Highlighted panels |
| `$text` | Primary text |
| `$text-muted` | Secondary text |
| `$success` | Positive/added (green) |
| `$error` | Negative/removed (red) |
| `$warning` | Warnings (yellow) |

## Creating Custom Themes

### Theme 1: Minimal (No Borders)

```python
CSS = """
DiffPanel {
    border: none;
    padding: 0;
}
"""
```

### Theme 2: Bordered (Like Git)

```python
CSS = """
DiffPanel {
    border: solid $primary;
}
"""
```

### Theme 3: High Contrast

```python
CSS = """
DiffPanel > Static {
    color: $text;
    text-style: bold;
}
"""
```

## Text Markup

The viewer uses Textual's markup syntax for colors:

```
[red]This text is red[/red]
[green]This text is green[/green]
[yellow]This text is yellow[/yellow]
[blue]This text is blue[/blue]
[bold]This text is bold[/bold]
[dim]This text is dimmed[/dim]
```

## Performance Considerations

- **Color rendering**: Minimal performance impact
- **Theme switching**: Not yet implemented but could be added
- **Terminal support**: Works with 256-color and true-color terminals

## Future Styling Features

Potential enhancements:

- [ ] Dark/Light mode toggle
- [ ] Multiple built-in themes
- [ ] Custom theme configuration file
- [ ] User-defined color palettes
- [ ] Syntax highlighting by file type
- [ ] Line highlighting on hover
- [ ] Selected section highlighting

## Troubleshooting

### Colors not showing?

1. Ensure your terminal supports colors
2. Check terminal theme settings
3. Try with a different terminal emulator

### Text hard to read?

1. Adjust terminal theme brightness
2. Try a different color combination
3. Enable bold text for better contrast

### Borders misaligned?

1. Clear terminal cache
2. Resize the window
3. Restart the application

## Tips for Best Results

1. **Use a dark terminal theme** - Works best with modern dark themes
2. **Enable bold fonts** - Improves readability
3. **Use a monospace font** - Ensures proper alignment
4. **Ensure sufficient contrast** - Check terminal settings
5. **Use a modern terminal** - WezTerm, Alacritty, Kitty, etc.

---

The current styling provides professional appearance while maintaining readability. Feel free to customize it to match your preferences!
