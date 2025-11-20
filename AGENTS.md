# Agent Guidelines for diffview

## Build & Test Commands
- Build: `make build` or `go build -o gg`
- Install: `make install` (builds and copies to ~/bin/gg)
- Run: `./gg` or `git diff | ./gg`
- Test single package: `go test ./src/diff` (no tests currently exist)
- Format: `go fmt ./...`
- Vet: `go vet ./...`

## Code Style
- **Imports**: Standard library first, then external packages, then local packages (diffview/src/...)
- **Formatting**: Use `go fmt` - tabs for indentation, standard Go conventions
- **Naming**: Exported types/functions start with capital letter (e.g., `FileDiff`, `ParseDiffIntoFiles`)
- **Comments**: Exported functions have comment starting with function name (e.g., `// CalculateStats computes additions...`)
- **Error Handling**: Return errors explicitly, check with `if err != nil`, use `fmt.Fprintf(os.Stderr, ...)` for error output
- **Types**: Struct types defined in src/models/types.go, prefer explicit types over interface{}
- **Receivers**: Use single-letter receiver names (e.g., `func (f *FileDiff)`)
- **Packages**: Organize by feature (diff, models, views, utils, io, highlighting, styles)
- **Architecture**: Main package in root, avoid circular imports (e.g., appWrapper pattern to bridge models/views)
- **Dependencies**: Charmbracelet tools (bubbletea, bubbles, lipgloss) + Chroma for syntax highlighting
