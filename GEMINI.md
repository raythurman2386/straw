# Gemini Context: Straw

Straw is a terminal-based, daemon-backed file automation system built with Go. It watches specific directories for filesystem events, evaluates them against user-defined rules, and executes automated actions. It runs on Linux, macOS, and Windows (10 1803+).

## Project Overview

- **Architecture**: Client-Server model.
    - `strawd`: The background daemon (engine) that watches folders and executes rules.
    - `straw`: The TUI client (interface) built with Bubble Tea for real-time monitoring and management.
        - **Dashboard**: High-level statistics (processed count, errors, last active).
        - **Wizards**: Interactive forms built with `huh`.
        - **Navigation**: Centrally managed keybindings and dynamic footer help via `bubbles/key` and `bubbles/help`.
- **Communication (IPC)**: JSON-RPC over Unix Domain Sockets (works on all supported platforms including Windows 10+). Default socket path is OS-dependent (e.g., `/tmp/straw.sock` on Linux).
- **Core Technologies**:
    - **Language**: Go 1.25+
    - **TUI**: [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Lip Gloss](https://github.com/charmbracelet/lipgloss), and [huh](https://github.com/charmbracelet/huh) for forms.
    - **File Watching**: [fsnotify](https://github.com/fsnotify/fsnotify)
    - **CLI**: [Cobra](https://github.com/spf13/cobra)
    - **Configuration**: TOML ([go-toml/v2](https://github.com/pelletier/go-toml/v2))
    - **Logging**: [log](https://github.com/charmbracelet/log) as an `slog` handler.

## Project Structure

- `cmd/`: Entry points for `straw` (client) and `strawd` (daemon).
- `internal/`: Core logic packages.
    - `actions/`: Implements file operations (`move`, `copy`, `trash`, `shell`). Platform-specific code uses build tags (`trash_unix.go`, `trash_windows.go`).
    - `config/`: Configuration loading and validation.
    - `ipc/`: Socket server/client and JSON-RPC types.
    - `rules/`: The rules evaluation engine. Platform-specific hidden file detection (`hidden_unix.go`, `hidden_windows.go`).
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
- **Reload Config**: Via TUI, or `pkill -HUP strawd` on Linux/macOS, or the `TRIGGER_RELOAD` IPC command on any platform.

## Development Conventions

- **Formatting**: Always use `go fmt ./...` before committing.
- **Style**: Adhere to standard Go idioms. Use `slog` for structured logging.
- **Testing**: Add unit tests in `*_test.go` files alongside the implementation.
- **Rules Engine**: Rules match files based on `glob`, `regex`, `extension`, `size`, `age`, and `file_type`.
- **Actions**:
    - `move`: Moves file to target directory (handles cross-device fallback).
    - `copy`: Copies file to target directory.
    - `trash`: Moves file to native OS trash (XDG on Linux, Finder trash on macOS, Recycle Bin on Windows).
    - `shell`: Executes a shell command with `$FILE` substitution.

## Configuration Details

Default paths are OS-dependent (resolved via `os.UserConfigDir()` and `os.UserCacheDir()`):

| OS | Config Path | Log Path |
|----|-------------|----------|
| Linux | `~/.config/straw/config.toml` | `~/.cache/straw/strawd.log` |
| macOS | `~/Library/Application Support/straw/config.toml` | `~/Library/Caches/straw/strawd.log` |
| Windows | `%AppData%\straw\config.toml` | `%LocalAppData%\straw\strawd.log` |
