# Quick Start

Get up and running with Straw in 5 minutes.

## 1. Install Straw

Install both the daemon (`strawd`) and TUI client (`straw`):

```bash
# Linux/macOS
curl -fsSL https://raw.githubusercontent.com/raythurman2386/straw/main/install.sh | sh

# Windows PowerShell
irm https://raw.githubusercontent.com/raythurman2386/straw/main/install.ps1 | iex
```

## 2. Start the Daemon

The daemon watches your filesystem and executes rules:

=== "Linux (systemd)"

    ```bash
    # Start and enable the service
    systemctl --user start strawd
    systemctl --user enable strawd
    ```

=== "macOS"

    ```bash
    # Run directly or use launchd
    strawd &
    ```

=== "Windows"

    ```powershell
    # Run directly
    strawd
    ```

## 3. Configure Your First Rule

Create a basic configuration file:

```bash
# Create the config directory
mkdir -p ~/.config/straw

# Create the config file
cat > ~/.config/straw/config.toml << 'EOF'
# Watch your Downloads folder
[[watch]]
path = "~/Downloads"
recursive = true

# Rule: Organize PDFs
[[rules]]
name = "Organize PDFs"
enabled = true
[rules.match]
extension = ".pdf"
[[rules.actions]]
type = "move"
target = "~/Documents/PDFs"
EOF
```

## 4. Launch the TUI

Open a new terminal and launch the TUI:

```bash
straw
```

You should see:
- Live event feed showing filesystem changes
- Status of watched directories
- Active rules and their match counts

## 5. Test Your Rule

1. Download a PDF file to your Downloads folder
2. Watch the TUI - you should see the rule match and execute
3. Check `~/Documents/PDFs` - your PDF should be there!

## What's Next?

- Learn more about [Configuration](configuration.md)
- Explore [Rules](usage/rules.md) in detail
- See more [Action Types](usage/actions.md)
- Check out [Example Configurations](reference/examples.md)

## Common Commands

```bash
# Check daemon status (Linux)
systemctl --user status strawd

# View daemon logs (Linux)
journalctl --user -u strawd -f

# Restart the daemon
systemctl --user restart strawd  # Linux
pkill -HUP strawd                # Signal reload

# Open the TUI
straw
```
