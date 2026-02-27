# Straw

**Modern File Automation System**

Straw is a modern, terminal-based file automation system built with Go. It runs on **Linux**, **macOS**, and **Windows** and features a persistent background daemon (`strawd`) that monitors your filesystem and an interactive TUI client (`straw`) for real-time monitoring and configuration.

## Features

- **Recursive Watching**: Automatically tracks changes across entire directory trees
- **Powerful Rules Engine**: Match files using multiple criteria
  - **Patterns**: Globbing, Regular Expressions
  - **Attributes**: Extension, Size (min/max), Age (min/max days), File Type
  - **Visibility**: Filter by hidden or visible files
- **Automated Actions**:
  - `move` / `copy`: Organize files into specific directories
  - `trash`: Safely move files to the system trash
  - `shell`: Execute custom scripts with `$FILE` environment variable support
- **Interactive TUI**:
  - Real-time event logging and status monitoring
  - Interactive Rule Wizard for easy configuration without manual editing
- **Live Reload**: Configuration updates instantly via the TUI or daemon reload signal
- **Cross-Platform**: Runs on Linux, macOS, and Windows (10 1803+)
- **Lightweight & Efficient**: Written in Go with minimal overhead

## Quick Start

Install Straw with one command:

```bash
# Linux/macOS
curl -fsSL https://raw.githubusercontent.com/raythurman2386/straw/main/install.sh | sh

# Windows PowerShell
irm https://raw.githubusercontent.com/raythurman2386/straw/main/install.ps1 | iex
```

Start the daemon and TUI:

```bash
# Start the daemon
strawd

# Launch the TUI (in another terminal)
straw
```

## Architecture

- **`strawd`**: The engine. Handles filesystem events and rule evaluation.
- **`straw`**: The interface. A Bubble Tea TUI that communicates with the daemon.
- **IPC**: JSON-RPC over Unix Domain Sockets for high-performance, low-latency communication.

## Next Steps

- [Installation Guide](installation.md) - Detailed installation instructions
- [Configuration](configuration.md) - Learn how to configure Straw
- [Usage Guide](usage/daemon.md) - Learn how to use Straw effectively
- [Development](development/contributing.md) - Contribute to the project

## License

This project is licensed under the MIT License.
