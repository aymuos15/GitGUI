package models

import (
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/evertras/bubble-table/table"
)

type DiffLine struct {
	LeftNum   string
	RightNum  string
	LeftText  string
	RightText string
	LineType  string // "add", "remove", "context", "header"
}

type FileDiff struct {
	Name           string
	Content        []string
	HighlightCache map[int]string   // Cache: line index -> highlighted line
	Lexer          chroma.Lexer     // Cached lexer for this file type
	Style          *chroma.Style    // Cached style
	Formatter      chroma.Formatter // Cached formatter
	Additions      int              // Number of added lines
	Deletions      int              // Number of deleted lines
}

// CalculateStats computes additions and deletions for a file
func (f *FileDiff) CalculateStats() {
	f.Additions = 0
	f.Deletions = 0

	for _, line := range f.Content {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			f.Additions++
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			f.Deletions++
		}
	}
}

// InitSyntaxHighlighting initializes the lexer, style, and formatter for a file
func (f *FileDiff) InitSyntaxHighlighting() {
	if f.Lexer != nil {
		return // Already initialized
	}

	// Get lexer based on file extension
	ext := filepath.Ext(f.Name)
	f.Lexer = lexers.Get(ext)

	if f.Lexer == nil {
		f.Lexer = lexers.Fallback
	}

	f.Lexer = chroma.Coalesce(f.Lexer)

	// Cache style and formatter
	f.Style = styles.Get("monokai")
	if f.Style == nil {
		f.Style = styles.Fallback
	}

	f.Formatter = formatters.TTY16m
	f.HighlightCache = make(map[int]string)
}

// HighlightLine highlights a single line using cached lexer/style/formatter
func (f *FileDiff) HighlightLine(lineIdx int, code string) string {
	// Check cache first
	if cached, exists := f.HighlightCache[lineIdx]; exists {
		return cached
	}

	// Ensure lexer is initialized
	f.InitSyntaxHighlighting()

	// Tokenize and format
	iterator, err := f.Lexer.Tokenise(nil, code)
	if err != nil {
		return code
	}

	var buf strings.Builder
	err = f.Formatter.Format(&buf, f.Style, iterator)
	if err != nil {
		return code
	}

	result := strings.TrimRight(buf.String(), "\n")

	// Cache the result
	f.HighlightCache[lineIdx] = result

	return result
}

type Model struct {
	LeftViewport      viewport.Model
	RightViewport     viewport.Model
	LogViewport       viewport.Model
	Files             []FileDiff
	ActiveTab         int
	Ready             bool
	Width             int
	Height            int
	LeftLineNum       int
	RightLineNum      int
	PendingRemoved    []string    // Track consecutive removed lines for word-level diff
	PendingAdded      []string    // Track consecutive added lines for word-level diff
	ViewMode          string      // "diff", "stats", or "log"
	NoDiffMessage     string      // Message to display when there's no diff
	StatsTable        table.Model // Scrollable stats table
	LogTable          table.Model // Scrollable log table
	SidebarWidth      int         // Width of the right sidebar
	AutoReloadEnabled bool        // Toggle for automatic reload on git changes
}
