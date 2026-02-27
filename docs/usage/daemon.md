# Daemon (strawd)

The Straw daemon (`strawd`) is the core component that monitors filesystem events and executes rules.

## Overview

The daemon:
- Watches configured directories for file changes
- Evaluates rules against new/modified files
- Executes actions on matching files
- Provides IPC endpoint for the TUI client
- Supports live configuration reloading

## Starting the Daemon

### Linux (systemd)

If installed via the install script, the daemon runs as a systemd user service:

```bash
# Check status
systemctl --user status strawd

# Start the service
systemctl --user start strawd

# Enable auto-start on login
systemctl --user enable strawd

# Restart
systemctl --user restart strawd

# Stop
systemctl --user stop strawd
```

### Manual Start

Run the daemon directly:

```bash
# Default config location
strawd

# Custom config
strawd --config /path/to/config.toml

# With verbose logging
strawd --verbose
```

## Command-Line Options

```
strawd [flags]

Flags:
  -c, --config string   Path to config file (default is platform-specific)
  -v, --verbose         Enable verbose logging
  -h, --help            Help for strawd
  --version             Version for strawd
```

## Configuration Reloading

The daemon supports live configuration reloading without restarting:

### Method 1: Signal

```bash
# Send HUP signal (Linux/macOS)
pkill -HUP strawd
```

### Method 2: TUI

Use the reload function in the TUI interface.

### Method 3: IPC

Send the `TRIGGER_RELOAD` command via the IPC interface.

## Logging

### Linux (journald)

```bash
# View all logs
journalctl --user -u strawd

# Follow logs in real-time
journalctl --user -u strawd -f

# View last 50 lines
journalctl --user -u strawd -n 50

# View since last boot
journalctl --user -u strawd --since today
```

### Direct Output

When running directly (not as a service):

```bash
# Log to file
strawd 2>&1 | tee ~/strawd.log

# Log rotation (Linux)
strawd 2>&1 | rotatelogs ~/strawd.log 10M
```

## Troubleshooting

### Daemon Won't Start

1. **Check config file exists:**
   ```bash
   ls -la ~/.config/straw/config.toml
   ```

2. **Validate config syntax:**
   ```bash
   # Use a TOML validator or just start the daemon and check error output
   strawd 2>&1 | head -20
   ```

3. **Check permissions:**
   ```bash
   # Ensure the daemon can access watched directories
   ls -la ~/Downloads  # or your watched directories
   ```

### High CPU Usage

- Reduce the number of watched directories
- Use non-recursive watching where possible
- Add more specific match criteria to reduce rule evaluation

### Rules Not Executing

1. **Check rule is enabled:**
   ```toml
   enabled = true
   ```

2. **Verify match criteria:**
   - Test your glob/regex patterns
   - Ensure file attributes match (size, age, etc.)

3. **Check logs for errors:**
   ```bash
   journalctl --user -u strawd -f
   ```

### Socket Connection Issues

1. **Check socket path:**
   ```bash
   # Default location
   ls -la /tmp/straw.sock
   
   # Or check your config
   grep socket_path ~/.config/straw/config.toml
   ```

2. **Verify daemon is running:**
   ```bash
   pgrep -a strawd
   ```

3. **Check file permissions:**
   ```bash
   ls -la $(dirname /tmp/straw.sock)
   ```

## Performance

### Resource Usage

Typical resource usage on Raspberry Pi (idle):

| Metric | Value |
|--------|-------|
| CPU | < 0.1% |
| Memory (RSS) | ~6.8 MB |
| Binary Size | 7.2 MB |

### Optimization Tips

1. **Limit recursive watching** to only necessary directories
2. **Use specific match criteria** to reduce rule evaluation overhead
3. **Combine multiple criteria** in a single rule instead of multiple rules
4. **Use `trash` action** over shell commands for better performance

## Advanced Topics

### Running as System Service

To run the daemon as a system-wide service (instead of user service):

```bash
# Create system service file
sudo tee /etc/systemd/system/strawd.service << 'EOF'
[Unit]
Description=Straw File Automation Daemon
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/strawd
Restart=always
RestartSec=5
User=straw
Group=straw

[Install]
WantedBy=multi-user.target
EOF

# Create dedicated user
sudo useradd -r -s /bin/false straw

# Set permissions
sudo mkdir -p /etc/straw
sudo chown straw:straw /etc/straw

# Reload and start
sudo systemctl daemon-reload
sudo systemctl enable strawd
sudo systemctl start strawd
```

### Debugging

Enable debug logging:

```bash
# Set environment variable
STRAW_DEBUG=1 strawd

# Or use verbose flag
strawd --verbose
```

### Monitoring

Monitor daemon health:

```bash
# Check if process is running
pgrep strawd > /dev/null && echo "Running" || echo "Not running"

# Monitor resource usage
watch -n 1 'ps aux | grep strawd | grep -v grep'

# Check open files
lsof -p $(pgrep strawd)
```
