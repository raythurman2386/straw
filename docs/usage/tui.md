# TUI Client (straw)

The Straw TUI client (`straw`) provides an interactive interface for monitoring and managing the Straw daemon.

## Overview

The TUI client:
- Connects to the daemon via IPC
- Displays real-time filesystem events
- Shows rule execution status
- Provides interactive rule creation wizard
- Allows configuration reloading

## Launching the TUI

```bash
# Default config location
straw

# Custom config
straw --config /path/to/config.toml
```

## Interface

The TUI consists of several views:

### Main View

The main view displays:
- **Event Log**: Real-time filesystem events
- **Rule Status**: Active rules and match counts
- **Watch Status**: Monitored directories
- **System Status**: Daemon connection status

### Navigation

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate through lists |
| `Tab` | Switch between panels |
| `Enter` | Select/confirm |
| `r` | Reload configuration |
| `w` | Open Rule Wizard |
| `q` or `Ctrl+c` | Quit |

## Rule Wizard

The Rule Wizard helps you create rules interactively without editing the config file manually.

### Accessing the Wizard

Press `w` in the main view to open the Rule Wizard.

### Wizard Steps

1. **Rule Name**
   - Enter a descriptive name for your rule
   - Example: "Organize PDF Downloads"

2. **Enable Rule**
   - Toggle whether the rule is active
   - Default: Enabled

3. **Match Criteria**
   - Select match type:
     - **Glob**: Pattern matching (e.g., `*.pdf`)
     - **Regex**: Regular expression
     - **Extension**: File extension
     - **Size**: File size range
     - **Age**: File age range
     - **Hidden**: Hidden file status
   - Enter the criteria value
   - Add additional criteria (optional)

4. **Actions**
   - Select action type:
     - **Move**: Move file to directory
     - **Copy**: Copy file to directory
     - **Trash**: Move to system trash
     - **Shell**: Execute command
   - Configure action parameters
   - Add additional actions

5. **Review & Save**
   - Review the rule configuration
   - Save to config file
   - Or go back to edit

### Example Wizard Session

```
Rule Name: Clean Old Logs
Enabled: Yes

Match Criteria:
- Extension: .log
- Min Age: 7 days

Actions:
1. Trash

Save this rule? [Yes/No]
```

## Views and Panels

### Event Log Panel

Shows filesystem events as they happen:

```
[14:32:15] Created: /home/user/Downloads/report.pdf
[14:32:16] Rule matched: "Organize PDFs"
[14:32:16] Action: Move to /home/user/Documents/PDFs
[14:32:16] Completed: /home/user/Documents/PDFs/report.pdf
```

### Rules Panel

Lists all configured rules with status:

```
Rules (3)
---------
[✓] Organize PDFs        (12 matches)
[✓] Clean Temp Files     (45 matches)
[✗] Backup Important     (disabled)
```

### Watches Panel

Shows watched directories:

```
Watches (2)
-----------
[✓] ~/Downloads (recursive)
[✓] ~/Desktop
```

## Commands

### Reload Configuration

Press `r` to reload the configuration without restarting the daemon.

### Create New Rule

Press `w` to open the Rule Wizard and create a new rule.

### Filter Events

Use `/` to filter the event log by:
- Filename pattern
- Rule name
- Action type

### Export Log

Press `e` to export the current event log to a file.

## Keyboard Shortcuts

### Global

| Key | Action |
|-----|--------|
| `q` | Quit |
| `?` | Show help |
| `Tab` | Next panel |
| `Shift+Tab` | Previous panel |

### Event Log

| Key | Action |
|-----|--------|
| `↑/↓` | Scroll |
| `Page Up/Down` | Scroll faster |
| `Home/End` | Jump to start/end |
| `/` | Filter |
| `e` | Export |
| `c` | Clear |

### Rules Panel

| Key | Action |
|-----|--------|
| `↑/↓` | Navigate rules |
| `Enter` | Toggle enable/disable |
| `d` | Delete rule |
| `e` | Edit rule |

### Wizard

| Key | Action |
|-----|--------|
| `↑/↓` | Navigate fields |
| `Tab` | Next field |
| `Enter` | Confirm |
| `Esc` | Cancel/Go back |

## Configuration

### TUI Settings

The TUI uses the same configuration file as the daemon. Additional TUI-specific settings can be added:

```toml
[tui]
refresh_rate = 100  # milliseconds
show_timestamps = true
max_events = 1000
```

### Themes

The TUI supports theming via the configuration:

```toml
[tui.theme]
primary = "#6c5ce7"
secondary = "#a29bfe"
success = "#00b894"
warning = "#fdcb6e"
error = "#d63031"
```

## Troubleshooting

### TUI Won't Start

1. **Check daemon is running:**
   ```bash
   systemctl --user status strawd
   ```

2. **Verify socket path:**
   ```bash
   ls -la /tmp/straw.sock
   ```

3. **Check config file:**
   ```bash
   straw --config ~/.config/straw/config.toml
   ```

### Connection Lost

If the TUI shows "Connection Lost":

1. Check if daemon crashed: `systemctl --user status strawd`
2. Restart the daemon: `systemctl --user restart strawd`
3. Press `r` in the TUI to reconnect

### Display Issues

If the TUI doesn't render correctly:

1. **Check terminal support:**
   - Ensure your terminal supports Unicode
   - Try a different terminal emulator

2. **Resize terminal:**
   - Minimum recommended size: 80x24
   - The TUI adapts to smaller sizes but may truncate content

3. **Disable mouse support:**
   ```toml
   [tui]
   mouse = false
   ```

## Tips

1. **Use the Rule Wizard** for quick rule creation without editing TOML
2. **Filter events** to focus on specific file types or rules
3. **Export logs** periodically for debugging or auditing
4. **Reload config** after manual edits instead of restarting the daemon
5. **Use keyboard shortcuts** for faster navigation
