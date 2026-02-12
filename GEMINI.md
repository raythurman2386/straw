# Gemini Context: Straw

Straw is a terminal-based, daemon-backed file automation system built with Go. It watches specific directories for filesystem events, evaluates them against user-defined rules, and executes automated actions.

## Project Overview

- **Architecture**: Client-Server model.
    - `strawd`: The background daemon (engine) that watches folders and executes rules.
    - `straw`: The TUI client (interface) built with Bubble Tea for real-time monitoring and management.
- **Communication (IPC)**: JSON-RPC over Unix Domain Sockets (default: `/tmp/straw.sock`).
- **Core Technologies**:
    - **Language**: Go 1.25+
    - **TUI**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) & [Lip Gloss](https://github.com/charmbracelet/lipgloss)
    - **File Watching**: [fsnotify](https://github.com/fsnotify/fsnotify)
    - **CLI**: [Cobra](https://github.com/spf13/cobra)
    - **Configuration**: TOML ([go-toml/v2](https://github.com/pelletier/go-toml/v2))

## Project Structure

- `cmd/`: Entry points for `straw` (client) and `strawd` (daemon).
- `internal/`: Core logic packages.
    - `actions/`: Implements file operations (`move`, `copy`, `trash`, `shell`).
    - `config/`: Configuration loading and validation.
    - `ipc/`: Unix socket server/client and JSON-RPC types.
    - `rules/`: The rules evaluation engine.
    - `watcher/`: Filesystem event watching wrapper.
    - `tui/`: Shared TUI components and themes.
    - `logging/`: Structured logging setup using `slog`.
- `bin/`: Built binaries (generated).
- `scripts/`: Installation and service management scripts.

## Building and Running

### Commands
- **Build**: `make build` (outputs to `bin/`)
- **Install**: `make install` (installs to `/usr/local/bin` by default)
- **Test**: `make test`
- **Lint/Check**: `make check` (runs `fmt`, `vet`, `test`, and `golangci-lint`)
- **Clean**: `make clean`

### Execution
- **Start Daemon**: `strawd`
- **Start TUI**: `straw`
- **Reload Config**: `pkill -HUP strawd` or via TUI.

## Development Conventions

- **Formatting**: Always use `go fmt ./...` before committing.
- **Style**: Adhere to standard Go idioms. Use `slog` for structured logging.
- **Testing**: Add unit tests in `*_test.go` files alongside the implementation.
- **Rules Engine**: Rules match files based on `glob`, `regex`, `extension`, `size`, `age`, and `file_type`.
- **Actions**:
    - `move`: Moves file to target directory (handles cross-device fallback).
    - `copy`: Copies file to target directory.
    - `trash`: Moves file to XDG-compliant trash.
    - `shell`: Executes a shell command with `$FILE` substitution.

## Configuration Details
Default config path: `~/.config/straw/config.toml`
Default log path: `~/.local/state/straw/strawd.log`
