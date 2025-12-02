package models

import (
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/evertras/bubble-table/table"
)

type FileDiff struct {
	Name           string
	Content        []string
	HighlightCache map[int]string   // Cache: line index -> highlighted line
	Lexer          chroma.Lexer     // Cached lexer for this file type
	Style          *chroma.Style    // Cached style
	Formatter      chroma.Formatter // Cached formatter
	Additions      int              // Number of added lines
	Deletions      int              // Number of deleted lines
	Status         string           // File status: "Modified", "New", "Deleted", "Renamed"
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

// SearchMatch represents a match position in diff view search
type SearchMatch struct {
	LineIdx int // Index in the content array
	Col     int // Column position in the line
}

type Model struct {
	LeftViewport      viewport.Model
	RightViewport     viewport.Model
	Files             []FileDiff
	ActiveTab         int
	Ready             bool
	Width             int
	Height            int
	ViewMode          string      // "diff", "stats", or "log"
	NoDiffMessage     string      // Message to display when there's no diff
	DiffType          string      // "working", "staged", or "none"
	StatsTable        table.Model // Scrollable stats table
	LogTable          table.Model // Scrollable log table
	AutoReloadEnabled bool        // Toggle for automatic reload on git changes
	ViewChanged       bool        // Flag to indicate view has changed

	// Filter/Search state
	FilterMode    string            // "", "author", "path", "date_from", "date_to", "search", "status", "extension"
	FilterInput   textinput.Model   // Text input for entering filter values
	LogFilters    LogFilterState    // Active filters for log view
	DiffSearch    DiffSearchState   // Search state for diff view
	StatsFilters  StatsFilterState  // Active filters for stats view
}

// LogFilterState holds active filters for the log view
type LogFilterState struct {
	Author   string // Filter by author name
	Path     string // Filter by file path
	DateFrom string // Filter from date (YYYY-MM-DD)
	DateTo   string // Filter to date (YYYY-MM-DD)
	Search   string // Search in commit messages
}

// HasActiveFilters returns true if any log filter is active
func (f LogFilterState) HasActiveFilters() bool {
	return f.Author != "" || f.Path != "" || f.DateFrom != "" || f.DateTo != "" || f.Search != ""
}

// DiffSearchState holds search state for the diff view
type DiffSearchState struct {
	Query        string        // Current search query
	Matches      []SearchMatch // All match positions
	CurrentMatch int           // Index of currently highlighted match
}

// StatsFilterState holds active filters for the stats view
type StatsFilterState struct {
	Status    string // Filter by status: "N", "M", "D", "R", "U" or ""
	Extension string // Filter by file extension
}
