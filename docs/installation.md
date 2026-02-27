# Installation

Straw can be installed on Linux, macOS, and Windows using multiple methods.

## System Requirements

- **Operating System**: Linux, macOS, or Windows 10 1803+
- **Go Version**: 1.25 or later (only for building from source)

## Quick Install (Recommended)

### Linux/macOS

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

## Manual Installation

### Download Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/raythurman2386/straw/releases).

=== "Linux (AMD64)"

    ```bash
    VERSION=$(curl -s https://api.github.com/repos/raythurman2386/straw/releases/latest | grep '"tag_name":' | cut -d'"' -f4)
    VERSION_NUM="${VERSION#v}"
    curl -LO "https://github.com/raythurman2386/straw/releases/download/${VERSION}/straw_${VERSION_NUM}_linux_amd64.tar.gz"
    tar -xzf straw_${VERSION_NUM}_linux_amd64.tar.gz
    sudo mv straw strawd /usr/local/bin/
    ```

=== "macOS (Apple Silicon)"

    ```bash
    VERSION=$(curl -s https://api.github.com/repos/raythurman2386/straw/releases/latest | grep '"tag_name":' | cut -d'"' -f4)
    VERSION_NUM="${VERSION#v}"
    curl -LO "https://github.com/raythurman2386/straw/releases/download/${VERSION}/straw_${VERSION_NUM}_darwin_arm64.tar.gz"
    tar -xzf straw_${VERSION_NUM}_darwin_arm64.tar.gz
    sudo mv straw strawd /usr/local/bin/
    ```

=== "Windows"

    Download the `.zip` file from the releases page and extract it to a directory in your `PATH`.

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

## Post-Installation

### Linux (systemd)

The install scripts will set up a systemd user service. You can manage it with:

```bash
# Check status
systemctl --user status strawd

# Start the service
systemctl --user start strawd

# Enable auto-start on login
systemctl --user enable strawd

# Restart the service
systemctl --user restart strawd

# View logs
journalctl --user -u strawd -f
```

### macOS

Use `launchd` to run the daemon:

```bash
# Create a plist file
mkdir -p ~/Library/LaunchAgents
cat > ~/Library/LaunchAgents/com.straw.strawd.plist << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.straw.strawd</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/strawd</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/tmp/strawd.out</string>
    <key>StandardErrorPath</key>
    <string>/tmp/strawd.err</string>
</dict>
</plist>
EOF

# Load the service
launchctl load ~/Library/LaunchAgents/com.straw.strawd.plist

# Start the service
launchctl start com.straw.strawd
```

### Windows

Use Task Scheduler or run the daemon directly:

```powershell
# Run directly
strawd

# Or create a scheduled task
$action = New-ScheduledTaskAction -Execute "strawd"
$trigger = New-ScheduledTaskTrigger -AtLogOn
$settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries
Register-ScheduledTask -TaskName "Straw Daemon" -Action $action -Trigger $trigger -Settings $settings -Force
```

## Verify Installation

Check that both binaries are installed correctly:

```bash
straw --version
strawd --version
```

## Uninstallation

=== "Linux/macOS"

    ```bash
    # Remove binaries
    sudo rm /usr/local/bin/straw /usr/local/bin/strawd

    # Remove config (optional)
    rm -rf ~/.config/straw

    # Remove systemd service (Linux only)
    systemctl --user stop strawd
    systemctl --user disable strawd
    rm ~/.config/systemd/user/strawd.service
    systemctl --user daemon-reload
    ```

=== "Windows"

    ```powershell
    # Remove binaries from PATH
    # Remove config (optional)
    Remove-Item -Recurse -Force "$env:AppData\straw"
    ```

## Troubleshooting

### Permission Denied

If you get a "permission denied" error when running the install script, try:

```bash
# Download the script first
curl -fsSL -o install.sh https://raw.githubusercontent.com/raythurman2386/straw/main/install.sh

# Run with bash
bash install.sh
```

### Binary Not Found

Make sure `/usr/local/bin` is in your `PATH`:

```bash
# Check PATH
echo $PATH

# Add to PATH if needed (add to ~/.bashrc or ~/.zshrc)
export PATH="/usr/local/bin:$PATH"
```

### Service Not Starting (Linux)

Check the systemd logs for errors:

```bash
journalctl --user -u strawd -n 50 --no-pager
```

### Socket Connection Failed

If the TUI can't connect to the daemon:

1. Make sure the daemon is running: `systemctl --user status strawd` (Linux) or `ps aux | grep strawd`
2. Check the socket path in your config matches where the daemon creates it
3. Check file permissions on the socket directory
