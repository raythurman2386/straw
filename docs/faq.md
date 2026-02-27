# Frequently Asked Questions

## General Questions

### What is Straw?

Straw is a file automation system that monitors your filesystem for changes and automatically processes files based on rules you define. It consists of a background daemon (`strawd`) and an interactive TUI client (`straw`).

### Is Straw free?

Yes, Straw is open-source software released under the MIT License. It's free to use, modify, and distribute.

### What platforms does Straw support?

Straw supports Linux, macOS, and Windows (10 1803+).

### Do I need to know how to code to use Straw?

No! While Straw uses configuration files, the TUI provides an interactive Rule Wizard that guides you through creating rules without writing any code.

## Installation

### How do I install Straw?

The easiest way is using the install script:

**Linux/macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/raythurman2386/straw/main/install.sh | sh
```

**Windows:**
```powershell
irm https://raw.githubusercontent.com/raythurman2386/straw/main/install.ps1 | iex
```

See the [Installation Guide](installation.md) for more options.

### Can I install Straw without root/admin privileges?

Yes, you can install Straw to any directory in your user space. Just add that directory to your PATH.

### How do I uninstall Straw?

See the [Installation Guide](installation.md#uninstallation) for platform-specific instructions.

## Configuration

### Where is the configuration file?

The default locations are:
- **Linux:** `~/.config/straw/config.toml`
- **macOS:** `~/Library/Application Support/straw/config.toml`
- **Windows:** `%AppData%\straw\config.toml`

### What format is the configuration file?

Straw uses TOML (Tom's Obvious, Minimal Language). It's designed to be easy to read and write.

### Can I use YAML or JSON instead?

Currently, only TOML is supported. TOML was chosen for its readability and comment support.

### How do I reload the configuration?

You can reload the configuration without restarting the daemon:
- Press `r` in the TUI
- Send a HUP signal: `pkill -HUP strawd`
- Use the IPC API

## Rules

### How do rules work?

Rules consist of:
1. **Match criteria** - Conditions a file must meet
2. **Actions** - What to do with matching files

When a file is created or modified, Straw checks it against all enabled rules in order. The first matching rule's actions are executed.

### What match criteria are available?

- `glob` - Shell-style pattern matching
- `regex` - Regular expressions
- `extension` - File extension
- `min_size` / `max_size` - File size in bytes
- `min_age_days` / `max_age_days` - File age
- `hidden` - Hidden file status

See [Match Criteria Reference](reference/match-criteria.md) for details.

### Can I use multiple criteria?

Yes! Multiple criteria are combined with AND logic. All criteria must match for the rule to trigger.

### Can I have multiple actions per rule?

Yes, a rule can have multiple actions that execute in sequence.

### Why isn't my rule matching?

Common causes:
1. Rule is disabled (`enabled = false`)
2. Pattern syntax error
3. File doesn't meet all criteria
4. Another rule matched first

Check the TUI event log to see what's happening.

### How do I test a rule?

1. Use the TUI event log to monitor rule execution
2. Create a test file that should match
3. Watch for the rule to trigger
4. Use `copy` action instead of `move` while testing

## Usage

### Do I need to run both straw and strawd?

Yes, they serve different purposes:
- `strawd` (daemon) - Monitors files and executes rules (must be running)
- `straw` (TUI) - Interactive interface (optional, for monitoring)

You can run the daemon without the TUI, but you need the daemon for the TUI to work.

### How do I start the daemon on boot?

**Linux (systemd):**
```bash
systemctl --user enable strawd
```

**macOS:**
Use `launchd` with a plist file. See the [Installation Guide](installation.md#macos).

**Windows:**
Use Task Scheduler. See the [Installation Guide](installation.md#windows).

### Can I use Straw without the TUI?

Yes! The daemon works independently. The TUI is only for monitoring and interactive management.

### How do I check if the daemon is running?

**Linux:**
```bash
systemctl --user status strawd
```

**All platforms:**
```bash
pgrep strawd
```

## Troubleshooting

### The TUI shows "Connection Lost"

This means the TUI can't connect to the daemon:
1. Check if the daemon is running
2. Verify the socket path in your config
3. Check file permissions on the socket

### Rules aren't executing

1. Check the daemon is running
2. Verify the config file is valid TOML
3. Check the TUI event log for errors
4. Ensure watched directories exist

### High CPU usage

- Reduce the number of watched directories
- Use non-recursive watching where possible
- Add more specific match criteria
- Check for rapidly changing files (e.g., temp files)

### Permission denied errors

1. Ensure the daemon has read access to watched directories
2. Ensure the daemon has write access to target directories
3. Check that target directories exist

### "Socket already in use" error

Another instance of the daemon is already running:
```bash
# Find and kill existing daemon
pkill strawd

# Or on Linux
systemctl --user stop strawd
```

## Advanced

### Can I run multiple daemons?

Not recommended. Each daemon should have its own socket path and configuration.

### Can Straw watch network drives?

Yes, but performance may vary depending on the filesystem and network latency.

### Does Straw work with Docker?

Yes, but you'll need to mount the directories you want to watch as volumes.

### Can I extend Straw with plugins?

Currently, no plugin system exists. However, you can use the `shell` action to execute custom scripts.

### Is there a web interface?

Not currently, but you could build one using the IPC API. See the [API Reference](development/api.md).

## Getting Help

### Where can I get help?

- Check this FAQ
- Read the [documentation](https://raythurman2386.github.io/straw)
- Search [GitHub Issues](https://github.com/raythurman2386/straw/issues)
- Open a new issue if you found a bug

### How do I report a bug?

Include:
- Operating system and version
- Straw version (`straw --version`)
- Steps to reproduce
- Expected vs actual behavior
- Configuration file (sanitized)
- Relevant logs

### How can I contribute?

See the [Contributing Guide](development/contributing.md) for details.
