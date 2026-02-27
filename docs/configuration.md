# Configuration

Straw uses TOML for configuration. This page covers the configuration file format, locations, and all available options.

## Configuration File Location

Straw looks for its configuration file in platform-specific locations:

| OS | Default Config Path |
|----|---------------------|
| Linux | `~/.config/straw/config.toml` |
| macOS | `~/Library/Application Support/straw/config.toml` |
| Windows | `%AppData%\straw\config.toml` |

You can also specify a custom config path:

```bash
strawd --config /path/to/config.toml
straw --config /path/to/config.toml
```

## Configuration Structure

A Straw configuration file consists of:

1. **Global Settings** - Top-level configuration options
2. **Watch Directories** - Folders to monitor for changes
3. **Rules** - Conditions and actions for file processing

### Basic Example

```toml
# Global settings
socket_path = "/tmp/straw.sock"

# Watch directories
[[watch]]
path = "~/Downloads"
recursive = true

# Rules
[[rules]]
name = "Example Rule"
enabled = true
[rules.match]
extension = ".tmp"
[[rules.actions]]
type = "trash"
```

## Global Settings

### `socket_path`

- **Type**: String
- **Default**: Platform-specific (e.g., `/tmp/straw.sock` on Linux)
- **Description**: Path to the Unix socket used for IPC between the daemon and TUI

```toml
socket_path = "/var/run/straw.sock"
```

## Watch Configuration

Define directories to monitor using the `[[watch]]` array:

```toml
[[watch]]
path = "~/Downloads"
recursive = true
```

### Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `path` | String | (required) | Directory path to watch. Supports `~` for home directory. |
| `recursive` | Boolean | `false` | Watch subdirectories recursively |

### Multiple Watch Directories

```toml
[[watch]]
path = "~/Downloads"
recursive = true

[[watch]]
path = "~/Desktop"
recursive = false

[[watch]]
path = "~/Documents/Incoming"
recursive = true
```

## Rules Configuration

Rules define what to do with files that match certain criteria:

```toml
[[rules]]
name = "Rule Name"
enabled = true
[rules.match]
# Match criteria here
[[rules.actions]]
# Actions here
```

### Rule Structure

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `name` | String | (required) | Human-readable name for the rule |
| `enabled` | Boolean | `true` | Whether the rule is active |
| `match` | Object | (required) | Criteria for matching files |
| `actions` | Array | (required) | List of actions to execute on matched files |

### Match Criteria

The `[rules.match]` section supports the following criteria:

| Key | Type | Description | Example |
|-----|------|-------------|---------|
| `glob` | String | Shell-style pattern matching | `"*.txt"` |
| `regex` | String | Regular expression matching | `"^temp_.*"` |
| `extension` | String | Specific file extension | `".zip"` |
| `min_size` | Integer | Minimum size in bytes | `1048576` (1MB) |
| `max_size` | Integer | Maximum size in bytes | `1024` |
| `min_age_days` | Integer | Files older than X days | `7` |
| `max_age_days` | Integer | Files newer than X days | `1` |
| `hidden` | Boolean | Match hidden files | `true` |

Multiple criteria can be combined - all must match for the rule to trigger (AND logic).

#### Examples

**Match by extension:**
```toml
[rules.match]
extension = ".pdf"
```

**Match by glob pattern:**
```toml
[rules.match]
glob = "*.log"
```

**Match by regex:**
```toml
[rules.match]
regex = "^backup_.*\\.zip$"
```

**Match by size:**
```toml
[rules.match]
min_size = 104857600  # 100MB
max_size = 1073741824  # 1GB
```

**Match by age:**
```toml
[rules.match]
min_age_days = 30  # Files older than 30 days
```

**Match hidden files:**
```toml
[rules.match]
hidden = true
```

**Combined criteria:**
```toml
[rules.match]
extension = ".log"
min_age_days = 7
max_size = 10485760  # 10MB
```

### Actions

Each rule can have one or more actions executed when a file matches:

```toml
[[rules.actions]]
type = "move"
target = "~/Archive"
```

Available action types:

| Type | Description | Required Parameters |
|------|-------------|---------------------|
| `move` | Move file to target directory | `target` |
| `copy` | Copy file to target directory | `target` |
| `trash` | Move file to system trash | None |
| `shell` | Execute shell command | `command` |

#### Move Action

```toml
[[rules.actions]]
type = "move"
target = "~/Documents/Organized"
```

#### Copy Action

```toml
[[rules.actions]]
type = "copy"
target = "~/Backup/Downloads"
```

#### Trash Action

```toml
[[rules.actions]]
type = "trash"
```

#### Shell Action

Execute custom shell commands with the `$FILE` environment variable:

```toml
[[rules.actions]]
type = "shell"
command = "echo 'New file: $FILE' >> ~/file-log.txt"
```

## Complete Example

```toml
# Global settings
socket_path = "/tmp/straw.sock"

# Watch Downloads folder recursively
[[watch]]
path = "~/Downloads"
recursive = true

# Rule 1: Organize PDFs
[[rules]]
name = "Organize PDFs"
enabled = true
[rules.match]
extension = ".pdf"
[[rules.actions]]
type = "move"
target = "~/Documents/PDFs"

# Rule 2: Clean up old temp files
[[rules]]
name = "Clean Old Temp Files"
enabled = true
[rules.match]
glob = "*.tmp"
min_age_days = 7
[[rules.actions]]
type = "trash"

# Rule 3: Backup important files
[[rules]]
name = "Backup Downloads"
enabled = true
[rules.match]
regex = "^important_.*"
[[rules.actions]]
type = "copy"
target = "~/Backup/Important"
[[rules.actions]]
type = "shell"
command = "notify-send 'File backed up' '$FILE'"

# Rule 4: Log large downloads
[[rules]]
name = "Log Large Downloads"
enabled = true
[rules.match]
min_size = 104857600  # 100MB
[[rules.actions]]
type = "shell"
command = "echo '$(date): Large file downloaded - $FILE ($(stat -c%s "$FILE"))' >> ~/large-files.log"
```

## Configuration Validation

Straw validates the configuration on startup. Common errors:

| Error | Solution |
|-------|----------|
| `watch path does not exist` | Ensure the directory exists before starting the daemon |
| `invalid regex pattern` | Check your regex syntax |
| `action type not supported` | Use one of: move, copy, trash, shell |
| `target directory does not exist` | Create the target directory or use a full path |

## Live Reload

Straw supports live configuration reloading:

- **Via TUI**: Use the reload function in the TUI interface
- **Via Signal**: Send `HUP` signal to the daemon: `pkill -HUP strawd`
- **Via IPC**: Send the `TRIGGER_RELOAD` command through the IPC interface

The daemon will reload the configuration without restarting, preserving active watches.
