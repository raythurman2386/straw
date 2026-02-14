# AGENTS.md — Straw Project Guidelines

## Project Overview

Straw is a Go-based file automation system with a TUI client (`straw`) and daemon (`strawd`). Runs on Linux, macOS, and Windows. Uses Bubble Tea for TUI, Cobra for CLI, and JSON-RPC over Unix sockets for IPC.

## Build Commands

```bash
# Build binaries
make build              # Builds bin/straw and bin/strawd
go build -o bin/straw ./cmd/straw
go build -o bin/strawd ./cmd/strawd

# Install
make install            # Builds and installs to /usr/local/bin

# Clean
make clean              # Removes bin/ and straw.log
```

## Test Commands

```bash
# Run all tests
make test
go test ./...

# Run tests for a specific package
go test ./internal/config
go test ./internal/rules
go test ./internal/actions

# Run a single test
go test -run TestConfig_Validate ./internal/config
go test -run TestConfig_Validate/Valid_config ./internal/config

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

## Lint Commands

```bash
# Run all checks (fmt, vet, test, lint)
make check

# Format code
make fmt
go fmt ./...

# Vet code
make vet
go vet ./...

# Run golangci-lint
golangci-lint run ./...
```

## Code Style Guidelines

### Imports
- Standard library imports first
- Third-party imports second (grouped by blank line)
- Internal package imports last (e.g., `straw/internal/config`)
- Use `goimports` or `gofmt` to organize imports

```go
import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/charmbracelet/bubbletea"

    "straw/internal/config"
    "straw/internal/ipc"
)
```

### Naming Conventions
- **Exported**: PascalCase (e.g., `Config`, `Load`, `NewEngine`)
- **Unexported**: camelCase (e.g., `loadConfig`, `validatePath`)
- **Interfaces**: End with `-er` suffix (e.g., `Executor`, `Watcher`)
- **Constants**: PascalCase for exported, camelCase for unexported
- **Packages**: Short, lowercase, single word (e.g., `config`, `rules`, `ipc`)

### Types and Structs
- Use struct tags for serialization (toml, json)
- Pointer receivers for methods that modify state
- Value receivers for read-only methods
- Document exported types and functions with comments

```go
// Config represents the application configuration.
type Config struct {
    SocketPath string        `toml:"socket_path"`
    Watch      []WatchFolder `toml:"watch"`
}
```

### Error Handling
- Return errors explicitly, never ignore
- Wrap errors with context using `fmt.Errorf("...: %w", err)`
- Use `errors.New()` for static error messages
- Check errors immediately after function calls

```go
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
```

### Testing
- Use table-driven tests with `t.Run()` for subtests
- Create temporary directories/files and clean up with `defer os.RemoveAll()`
- Test function naming: `Test<Type>_<Method>` or `Test<Method>`
- Use descriptive subtest names (e.g., "Valid config", "Missing watch")

```go
func TestConfig_Validate(t *testing.T) {
    t.Run("Valid config", func(t *testing.T) {
        // test code
    })
}
```

### Comments
- Start with the name of the thing being documented
- Use complete sentences
- Document all exported functions, types, and constants

```go
// Load reads and parses the configuration from the specified path.
func Load(path string) (*Config, error) { ... }
```

### Formatting
- Use `gofmt` for automatic formatting
- Line length: aim for 100 chars max (Go doesn't enforce)
- Indentation: tabs (Go standard)
- Opening braces: same line as statement

### Platform-Specific Code
- Use Go build tags (`//go:build`) to split platform code into separate files
- File naming convention: `<name>_unix.go` (for `//go:build !windows`) and `<name>_windows.go`
- Keep main application code platform-agnostic where possible
- Use stdlib functions for OS paths: `os.UserConfigDir()`, `os.UserCacheDir()`, `os.TempDir()`
- Use `filepath.Join()` and `filepath.Separator` instead of hardcoded path separators
- Guard Unix-only syscalls (e.g., `os.Chmod` for socket permissions) with `runtime.GOOS` checks
- On Windows, only `os.Interrupt` is meaningful for signal handling; `SIGHUP`/`SIGTERM` have no effect
- Existing platform-split files:
  - `cmd/strawd/signal_unix.go` / `signal_windows.go` — signal handling
  - `internal/actions/trash_unix.go` / `trash_windows.go` — trash implementation
  - `internal/rules/hidden_unix.go` / `hidden_windows.go` — hidden file detection

## Project Structure

```
cmd/              # Main applications
├── straw/      # TUI client
└── strawd/     # Daemon
internal/         # Private application code
├── actions/      # Action execution (move, copy, trash, shell)
├── config/       # Configuration loading and validation
├── ipc/          # Inter-process communication (JSON-RPC)
├── logging/      # Logging setup
├── rules/        # Rules engine
├── tui/          # UI components and themes
└── watcher/      # File system watching
```

## Go Version

Requires Go 1.25 or later (specified in go.mod).

## Dependencies

Key dependencies:
- `github.com/spf13/cobra` — CLI framework
- `github.com/charmbracelet/bubbletea` — TUI framework
- `github.com/pelletier/go-toml/v2` — TOML parsing
- `github.com/bmatcuk/doublestar/v4` — Glob matching

Add new dependencies with:
```bash
go get github.com/user/repo
go mod tidy
```
