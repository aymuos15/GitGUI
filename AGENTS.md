# Agents Guidelines

## Build & Test Commands

**Language**: Go 1.25+

- **Build**: `go build -o gg` - Compiles the binary
- **Run**: `./gg` - Execute from git repository root
- **Run tests**: `go test ./...` - Run all package tests (use `-v` for verbose)
- **Single test**: `go test -run TestName ./package` - Run specific test
- **Format**: `go fmt ./...` - Auto-format all files
- **Lint**: `golangci-lint run` (if installed) - Run static analysis
- **Check**: `go vet ./...` - Check for potential issues

## Code Style Guidelines

### Imports
- Organize in 3 groups: stdlib, external packages, local imports (separated by blank lines)
- Use absolute import paths: `"gg/src/models"`
- No underscore imports or aliases unless necessary

### Formatting
- Use `go fmt` standard formatting (2-space indentation, enforced by Go)
- Max line length: no strict limit, but keep under 120 characters when practical
- Use camelCase for function/variable names, PascalCase for exported identifiers

### Types & Structs
- Define types near their usage location
- Document exported types with comment lines: `// TypeName does something`
- Use struct embedding for composition (e.g., `struct { Model }`)
- Use descriptive field names; avoid abbreviations

### Error Handling
- Use `%w` for error wrapping: `fmt.Errorf("context: %w", err)`
- Return errors explicitly; never silently fail
- Check errors immediately after calls: `if err != nil { return nil, err }`
- Wrap command execution errors with context describing the operation

### Functions & Methods
- Keep functions focused and under 50 lines when possible
- Document exported functions with receiver name: `// Update handles Bubble Tea messages`
- Use receiver methods on types for Bubble Tea interfaces (Init, Update, View)
- Pass pointers to mutable structs; use values for immutable data

### Naming Conventions
- Use descriptive package names (diff, models, views, io, watcher)
- Constants in UPPER_SNAKE_CASE if exported, lowerCamelCase if private
- Boolean variables/functions use "is", "has", "enabled" prefixes
- Avoid generic names (data, result, temp)

### Package Organization
- `main.go`: Entry point and app wrapper orchestration
- `src/models/`: Core data structures and Bubble Tea model logic
- `src/diff/`: Git diff parsing and file diff operations
- `src/views/`: UI rendering for diff, log, and stats views
- `src/io/`: Git command execution and file I/O
- `src/watcher/`: Git change detection and filesystem monitoring
- `src/highlighting/`: Syntax highlighting utilities
- `src/styles/`: Lipgloss styling and theming
- `src/utils/`: Text processing and formatting helpers

### Testing
- Place tests in `*_test.go` files in the same package
- Use table-driven tests for multiple scenarios
- Mock external dependencies (git commands via exec.Command)

### Git Integration
- Use `fmt.Errorf("message: %w", err)` to wrap git command errors
- Always check `cmd.Wait()` after reading pipe output
- Handle both working tree and staged changes gracefully

No Cursor, Copilot, or MCP rules found for this repository.
