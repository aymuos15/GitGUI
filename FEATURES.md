# Diff Viewer - Features & Capabilities

## Visual Display

### Color Coding
```
🔴 RED    - Lines removed from the original file
🟢 GREEN  - Lines added to the modified file
⚪ DIM    - Lines unchanged in both files
```

### Line Numbers
- Left panel shows line numbers from file 1
- Right panel shows line numbers from file 2
- Empty spaces (∅) indicate alignment gaps

### Layout
```
┌─────────────────────────┬─────────────────────────┐
│  File 1 (Original)      │  File 2 (Modified)      │
├─────────────────────────┼─────────────────────────┤
│ 1 Hello world           │ 1 Hello world           │
│ 2 This is a test file   │ 2 This is a test file   │
│ 3 Line 3 remains...     │ 3 Line 3 remains...     │
│ 4 Line 4 will change ❌ │ 4 Line 4 has modified ✅│
│ 5 Line 5 stays...       │ 5 Line 5 stays...       │
│ 6 Line 6 is removed ❌  │   (empty)               │
│   (empty)               │ 6 A new line inserted ✅│
│ 7 Line 7 is also...     │ 7 Line 7 is also...     │
│ 8 Line 8 ends file      │ 8 Line 8 ends file      │
│   (empty)               │ 9 Line 9 is new ✅      │
└─────────────────────────┴─────────────────────────┘
```

## Navigation Features

### Single Line Scrolling
- **j** - Move down one line
- **k** - Move up one line
- Smooth line-by-line navigation with synchronized panels

### Page Navigation
- **f** (forward) - Scroll down ~10 lines (one page)
- **b** (backward) - Scroll up ~10 lines (one page)
- Useful for quickly scanning large diffs

### Jump Navigation
- **g** - Jump to the very beginning of the diff
- **G** (Shift+G) - Jump to the very end of the diff
- Instant navigation for large files

### Exit
- **q** - Quit the application gracefully

## Synchronization Features

### Linked Scrolling
- Both panels scroll in perfect synchronization
- When you scroll left panel, right panel follows automatically
- Maintains alignment between corresponding sections

### Line Alignment
- Removed lines on left paired with added lines on right
- Empty spaces maintain perfect alignment
- Easy to track what changed where

## Diff Analysis Capabilities

### Diff Operations Detected
1. **Equal** - Identical lines (appear on both sides)
2. **Replace** - Lines that changed (removals left, additions right)
3. **Delete** - Lines only in original file (left only)
4. **Insert** - Lines only in modified file (right only)

### Statistics Tracking
- Counts added lines automatically
- Counts removed lines automatically
- Could be extended to show in UI

## File Format Support

### Text File Types
- Works with any text-based file
- Handles various line endings (Unix, Windows, Mac)
- UTF-8 and ASCII encodings supported

### Examples
- Source code (.py, .js, .java, .c, etc.)
- Configuration files (.json, .yaml, .toml, etc.)
- Documentation (.md, .txt, .rst, etc.)
- Any plain text file

## Error Handling

### Robust Error Management
- **File Not Found**: Clear error message displayed
- **Permission Denied**: Caught and reported
- **Encoding Issues**: Handled gracefully
- **Invalid Paths**: Validated before processing

## Performance Optimizations

### Memory Efficient
- Uses ScrollView for viewport-based rendering
- Only renders visible portion of content
- Can handle moderately large files efficiently

### Time Efficient
- Uses Python's optimized `difflib.SequenceMatcher`
- O(n*m) complexity suitable for files up to 10,000 lines
- Fast computation even for complex diffs

## Accessibility Features

### Keyboard-Only Navigation
- Complete keyboard navigation (no mouse required)
- Vim-like keybindings for familiar UX
- Standard terminal key combinations

### Clear Visual Feedback
- Bordered panels for clear separation
- Color coding for quick understanding
- Line numbers for precise reference
- Title bars showing file names

### Responsive Interface
- Adjusts to terminal size
- Maintains layout at any resolution
- Smooth scrolling behavior

## Advanced Diff Matching

### Smart Alignment
- Handles line insertions/deletions gracefully
- Pairs replaced lines intelligently
- Maintains readability even with significant changes

### Handles Edge Cases
- Empty files
- Files with only additions
- Files with only deletions
- Files with only modifications

## Future Enhancement Hooks

The codebase is structured for easy additions:
- Stats panel integration point
- Search functionality skeleton ready
- Multiple view modes easily implementable
- Configurable keybindings support
- Plugin architecture feasible

---

All features are production-ready and fully tested!
