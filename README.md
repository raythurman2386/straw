# Straw — File Automation System

[![Go Version](https://img.shields.io/github/go-mod/go-version/raythurman2386/straw)](https://go.dev/)
[![License](https://img.shields.io/github/license/raythurman2386/straw)](LICENSE)

Straw is a modern, terminal-based file automation system built with Go. It features a persistent background daemon (`strawd`) that monitors your filesystem and an interactive TUI client (`straw`) for real-time monitoring and configuration.

## 🚀 Features

- **Recursive Watching:** Automatically tracks changes across entire directory trees.
- **Powerful Rules Engine:** Match files using multiple criteria:
  - **Patterns:** Globbing, Regular Expressions.
  - **Attributes:** Extension, Size (min/max), Age (min/max days), File Type.
  - **Visibility:** Filter by hidden or visible files.
- **Automated Actions:**
  - `move` / `copy`: Organize files into specific directories.
  - `trash`: Safely move files to the system trash (XDG compliant).
  - `shell`: Execute custom scripts with `$FILE` environment variable support.
- **Interactive TUI:**
  - Real-time event logging and status monitoring.
  - Interactive Rule Wizard for easy configuration without manual editing.
- **Live Reload:** Configuration updates instantly via SIGHUP or TUI updates.
- **Lightweight & Efficient:** Written in Go with minimal overhead and JSON-RPC over Unix sockets for communication.

## 🛠 Installation

### Unified Script (Recommended)

Clone the repository and run the installation script:

```bash
git clone https://github.com/raythurman2386/straw.git
cd straw
./install.sh
```

The script will:
1. Build both `straw` (TUI) and `strawd` (Daemon).
2. Install binaries to `/usr/local/bin`.
3. Initialize the default configuration at `~/.config/straw/config.toml`.
4. Set up and start the `strawd` systemd user service.

### Manual Installation

If you prefer to build manually:

```bash
make build
# Binaries are now in bin/
```

To uninstall the service only, run:
```bash
./scripts/uninstall_service.sh
```

## ⚙️ Configuration

Straw uses TOML for configuration. You can find it at `~/.config/straw/config.toml`.

### Example Configuration

```toml
# Folders to watch
[[watch]]
path = "~/Downloads"
recursive = true

# Define rules
[[rules]]
name = "Cleanup Old PDFs"
enabled = true
[rules.match]
extension = ".pdf"
min_age_days = 30
[[rules.actions]]
type = "trash"

[[rules]]
name = "Organize Code Snippets"
enabled = true
[rules.match]
glob = "*.go"
[[rules.actions]]
type = "copy"
target = "~/Projects/snippets"
```

### Available Match Criteria

| Key | Description | Example |
|-----|-------------|---------|
| `glob` | Shell-style pattern matching | `*.txt` |
| `regex` | Regular expression matching | `^temp_.*` |
| `extension` | Specific file extension | `.zip` |
| `min_size` | Minimum size in bytes | `1048576` (1MB) |
| `max_size` | Maximum size in bytes | `1024` |
| `min_age_days` | Files older than X days | `7` |
| `max_age_days` | Files newer than X days | `1` |
| `hidden` | Match hidden files (true/false) | `true` |

## 🕹 Usage

### The Daemon (`strawd`)

The daemon usually runs in the background as a systemd service.

```bash
# Manage the service
systemctl --user status strawd
systemctl --user restart strawd

# View logs via journald
journalctl --user -u strawd -f
```

### The TUI Client (`straw`)

Launch the TUI to see what's happening and manage rules interactively.

```bash
straw
```

- **Live View:** Watch events and rule executions as they happen.
- **Rule Wizard:** Create new rules interactively without touching the config file.

## 🏗 Architecture

- **`strawd`**: The engine. Handles filesystem events and rule evaluation.
- **`straw`**: The interface. A Bubble Tea TUI that communicates via Unix Domain Sockets.
- **IPC**: JSON-RPC over sockets for high-performance, low-latency communication.

## 🔨 Development

See [AGENTS.md](AGENTS.md) for detailed project guidelines and internal structure.

```bash
make build    # Build binaries
make test     # Run all tests
make check    # Run linting and formatting
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
