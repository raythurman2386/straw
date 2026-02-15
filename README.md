# Straw — File Automation System

[![Go Version](https://img.shields.io/github/go-mod/go-version/raythurman2386/straw)](https://go.dev/)
[![License](https://img.shields.io/github/license/raythurman2386/straw)](LICENSE)

Straw is a modern, terminal-based file automation system built with Go. It runs on **Linux**, **macOS**, and **Windows**. It features a persistent background daemon (`strawd`) that monitors your filesystem and an interactive TUI client (`straw`) for real-time monitoring and configuration.

## 🚀 Features

- **Recursive Watching:** Automatically tracks changes across entire directory trees.
- **Powerful Rules Engine:** Match files using multiple criteria:
  - **Patterns:** Globbing, Regular Expressions.
  - **Attributes:** Extension, Size (min/max), Age (min/max days), File Type.
  - **Visibility:** Filter by hidden or visible files.
- **Automated Actions:**
  - `move` / `copy`: Organize files into specific directories.
  - `trash`: Safely move files to the system trash (Recycle Bin on Windows, XDG trash on Linux, Finder trash on macOS).
  - `shell`: Execute custom scripts with `$FILE` environment variable support.
- **Interactive TUI:**
  - Real-time event logging and status monitoring.
  - Interactive Rule Wizard for easy configuration without manual editing.
- **Live Reload:** Configuration updates instantly via the TUI or daemon reload signal.
- **Cross-Platform:** Runs on Linux, macOS, and Windows (10 1803+) with platform-native behaviors.
- **Lightweight & Efficient:** Written in Go with minimal overhead and JSON-RPC over Unix sockets for communication.

## 🛠 Installation

### Quick Install (Linux/macOS)

One-liner to download and install the latest release:

```bash
curl -fsSL https://raw.githubusercontent.com/raythurman2386/straw/main/install.sh | sh
```

Or with wget:
```bash
wget -qO- https://raw.githubusercontent.com/raythurman2386/straw/main/install.sh | sh
```

This installs both `straw` and `strawd` to `/usr/local/bin`.

### Windows

**Quick Install (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/raythurman2386/straw/main/install.ps1 | iex
```



Or download the `.zip` archive from the [releases page](https://github.com/raythurman2386/straw/releases) and add it to your `PATH`.

### Download Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/raythurman2386/straw/releases).

```bash
# Example for Linux AMD64
VERSION=$(curl -s https://api.github.com/repos/raythurman2386/straw/releases/latest | grep '"tag_name":' | cut -d'"' -f4)
VERSION_NUM="${VERSION#v}"
curl -LO "https://github.com/raythurman2386/straw/releases/download/${VERSION}/straw_${VERSION_NUM}_linux_amd64.tar.gz"
tar -xzf straw_${VERSION_NUM}_linux_amd64.tar.gz
sudo mv straw strawd /usr/local/bin/
```

### Build from Source

Clone the repository and build manually:

```bash
git clone https://github.com/raythurman2386/straw.git
cd straw
make build
# Binaries are now in bin/
```

Or use the full development install script (builds, installs, and sets up the systemd service):

```bash
./scripts/install-from-source.sh
```

To uninstall the service only, run:
```bash
./scripts/uninstall_service.sh
```

## ⚙️ Configuration

Straw uses TOML for configuration. The config file location depends on your OS:

| OS | Default Config Path |
|----|---------------------|
| Linux | `~/.config/straw/config.toml` |
| macOS | `~/Library/Application Support/straw/config.toml` |
| Windows | `%AppData%\straw\config.toml` |

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

The daemon runs in the background and watches your filesystem.

**Linux (systemd):**
```bash
# Manage the service
systemctl --user status strawd
systemctl --user restart strawd

# View logs via journald
journalctl --user -u strawd -f
```

**macOS / Windows:**
```bash
# Run the daemon directly (e.g., from a terminal or startup script)
strawd
```

On all platforms, you can reload the configuration without restarting the daemon by using the TUI or sending the `TRIGGER_RELOAD` IPC command. On Linux/macOS you can also use `pkill -HUP strawd`.

### The TUI Client (`straw`)

Launch the TUI to see what's happening and manage rules interactively.

```bash
straw
```

- **Live View:** Watch events and rule executions as they happen.
- **Rule Wizard:** Create new rules interactively without touching the config file.

## 🏗 Architecture

- **`strawd`**: The engine. Handles filesystem events and rule evaluation.
- **`straw`**: The interface. A Bubble Tea TUI that communicates with the daemon.
- **IPC**: JSON-RPC over Unix Domain Sockets for high-performance, low-latency communication (supported on Linux, macOS, and Windows 10+).

## 🔨 Development

See [AGENTS.md](AGENTS.md) for detailed project guidelines and internal structure.

```bash
make build    # Build binaries
make test     # Run all tests
make check    # Run linting and formatting
```

## 📈 Performance and Benchmarks

Straw is designed for high efficiency and minimal resource footprint, making it ideal for everything from low-power devices like Raspberry Pi to high-performance workstations.

### Resource Footprint (Idle)
Benchmarks taken on **Raspberry Pi (arm64)** running **Go 1.25.7**:

| Component | CPU Usage | Memory (RSS) | Binary Size |
|-----------|-----------|--------------|-------------|
| `strawd` (Daemon) | < 0.1% | ~6.8 MB | 7.2 MB |
| `straw` (TUI) | ~0.0% | ~6.5 MB | 9.4 MB |

### Key Performance Features
- **Minimal Overhead:** Written in Go with a focus on low memory allocation and efficient event handling.
- **IPC Efficiency:** Uses JSON-RPC over Unix Domain Sockets for sub-millisecond communication latency between the TUI and daemon.
- **Compact Binaries:** Statically linked binaries with no external runtime dependencies, optimized for rapid startup.
- **Scalability:** Capable of watching thousands of files with minimal impact on system responsiveness.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
