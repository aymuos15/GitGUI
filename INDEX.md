# Diff Viewer - Project Index

Welcome to the side-by-side diff viewer project! Here's a guide to navigate the project files.

## 📖 Documentation Guide

Start here based on what you want to do:

### Just Want to Use It?
→ **[QUICKSTART.md](QUICKSTART.md)** - Get up and running in 2 minutes

### Want Full Details?
→ **[README.md](README.md)** - Complete project documentation

### Curious About Features?
→ **[FEATURES.md](FEATURES.md)** - Detailed feature breakdown

### Need Technical Details?
→ **[PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)** - Architecture and implementation

## 🔧 Project Structure

```
diffview/
├── INDEX.md                     ← You are here
├── 
├── 📚 DOCUMENTATION
│   ├── README.md               - Full documentation
│   ├── QUICKSTART.md           - Quick start guide  
│   ├── FEATURES.md             - Feature showcase
│   └── PROJECT_SUMMARY.md      - Technical deep-dive
│
├── 💻 CODE
│   └── diffview.py             - Main application (341 lines)
│
├── ⚙️  CONFIGURATION
│   ├── requirements.txt        - Python dependencies
│   └── .gitignore              - Git configuration
│
└── 🧪 TESTING
    ├── test_file1.txt          - Original test file
    └── test_file2.txt          - Modified test file
```

## 🚀 Quick Commands

```bash
# Install dependencies
pip install -r requirements.txt

# Compare two files
python3 diffview.py file1.txt file2.txt
python3 diffview.py test_file1.txt test_file2.txt
python3 diffview.py old.py new.py

# Check Python syntax
python3 -m py_compile diffview.py
```

## 📚 Reading Order

1. **[QUICKSTART.md](QUICKSTART.md)** (5 min)
   - Installation
   - Basic usage
   - Key shortcuts

2. **[USAGE_EXAMPLES.md](USAGE_EXAMPLES.md)** (10 min) ⭐ **NEW**
   - Git integration examples
   - Real-world workflows
   - Practical tips and tricks

3. **[diffview.py](diffview.py)** (10 min)
   - Well-commented source code
   - Three main classes
   - Easy to follow structure

4. **[FEATURES.md](FEATURES.md)** (10 min)
   - Visual display capabilities
   - Navigation features
   - Performance characteristics

5. **[PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)** (15 min)
   - Architecture overview
   - Implementation details
   - Future enhancements

## 🎮 Keyboard Shortcuts Quick Reference

| Key | Action |
|-----|--------|
| `q` | Quit |
| `j`/`k` | Scroll down/up |
| `f`/`b` | Page down/up |
| `g`/`G` | Top/bottom |

## 💡 Key Concepts

### What is a Diff?
A diff compares two files and shows what's different. This viewer shows both files side-by-side.

### Color Coding
- 🔴 **Red** = Lines removed from the first file
- 🟢 **Green** = Lines added to the second file  
- ⚪ **Dim** = Lines unchanged in both files

### Line Numbers
- Left panel: Line numbers from file 1
- Right panel: Line numbers from file 2
- Proper alignment maintained even with different line counts

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────┐
│         DiffViewer (App)            │
│  - Loads files                      │
│  - Computes diff                    │
│  - Handles keyboard                 │
│  - Manages panels                   │
├──────────────────┬──────────────────┤
│  DiffPanel Left  │  DiffPanel Right │
│  - Displays      │  - Displays      │
│  - Scrolls       │  - Scrolls       │
│  - Shows colors  │  - Shows colors  │
├──────────────────┴──────────────────┤
│        Powered by Textual            │
└─────────────────────────────────────┘
```

## 🧪 Testing

Test files are included to verify functionality:

```bash
python3 diffview.py test_file1.txt test_file2.txt
```

This demonstrates:
- ✓ Line modifications
- ✓ Line deletions
- ✓ Line insertions
- ✓ Unchanged sections

## 📊 Project Statistics

- **Main code**: 341 lines of Python
- **Classes**: 3 (DiffLine, DiffPanel, DiffViewer)
- **Documentation**: 4 markdown files
- **Dependencies**: Just Textual
- **Python version**: 3.9+

## 🤔 FAQ

**Q: Does it work on Windows?**
A: Yes! Textual works on Windows, macOS, and Linux.

**Q: Can I compare binary files?**
A: No, it's designed for text files.

**Q: Is it fast?**
A: Yes! Efficient rendering and diff algorithm make it suitable for files up to ~10,000 lines.

**Q: Can I modify the code?**
A: Absolutely! The code is well-structured and easy to extend.

**Q: What terminal do I need?**
A: Any modern terminal works. WezTerm, Alacritty, Kitty, Ghostty, etc.

## 🔗 Related Resources

- [Textual Documentation](https://textual.textualize.io/)
- [Python difflib](https://docs.python.org/3/library/difflib.html)
- [Vim Cheatsheet](https://vim.rtorr.com/) (for keybindings)

## 📝 License

This project is provided as-is for educational and personal use.

## ✨ What's Next?

1. Try it out with your own files
2. Read the code to understand how it works
3. Consider extending it with features from the enhancement ideas
4. Share your improvements!

---

**Happy diffing!** 🎉
