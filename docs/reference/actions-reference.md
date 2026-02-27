# Actions Reference

Complete reference for all action types available in Straw rules.

## Action Structure

Actions are defined in the `[[rules.actions]]` array:

```toml
[[rules]]
name = "Example Rule"
[rules.match]
extension = ".pdf"
[[rules.actions]]
type = "move"
target = "~/Documents/PDFs"
[[rules.actions]]
type = "shell"
command = "echo 'Moved: $FILE'"
```

## Move Action

Moves a file from its current location to a target directory.

### Syntax

```toml
[[rules.actions]]
type = "move"
target = "path/to/directory"
```

### Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `type` | string | Yes | - | Must be `"move"` |
| `target` | string | Yes | - | Destination directory path |

### Examples

**Basic move:**
```toml
[[rules.actions]]
type = "move"
target = "~/Documents"
```

**Organize by type:**
```toml
[[rules.actions]]
type = "move"
target = "~/Documents/PDFs"
```

**Archive location:**
```toml
[[rules.actions]]
type = "move"
target = "/mnt/external/Archive"
```

### Behavior

- Target directory is created if it doesn't exist
- Original file is removed from source location
- File permissions are preserved
- On conflict with existing file, behavior is platform-dependent

### Platform Notes

| Platform | Notes |
|----------|-------|
| Linux | Uses atomic move when possible |
| macOS | Preserves extended attributes |
| Windows | Handles path length limits automatically |

## Copy Action

Creates a copy of a file in a target directory while keeping the original.

### Syntax

```toml
[[rules.actions]]
type = "copy"
target = "path/to/directory"
```

### Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `type` | string | Yes | - | Must be `"copy"` |
| `target` | string | Yes | - | Destination directory path |

### Examples

**Create backup:**
```toml
[[rules.actions]]
type = "copy"
target = "~/Backup"
```

**Mirror to external drive:**
```toml
[[rules.actions]]
type = "copy"
target = "/Volumes/Backup/Downloads"
```

**Multi-location backup:**
```toml
[[rules.actions]]
type = "copy"
target = "~/Backup"
[[rules.actions]]
type = "copy"
target = "/mnt/nas/backup"
```

### Behavior

- Target directory is created if it doesn't exist
- Original file remains in place
- File permissions and timestamps are copied
- On conflict with existing file, may overwrite

## Trash Action

Moves a file to the system trash/recycle bin for safe deletion.

### Syntax

```toml
[[rules.actions]]
type = "trash"
```

### Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `type` | string | Yes | - | Must be `"trash"` |

### Examples

**Basic trash:**
```toml
[[rules.actions]]
type = "trash"
```

**With notification:**
```toml
[[rules.actions]]
type = "trash"
[[rules.actions]]
type = "shell"
command = "notify-send 'File moved to trash' '$(basename $FILE)'"
```

### Platform-Specific Locations

| Platform | Trash Location |
|----------|----------------|
| Linux | `~/.local/share/Trash/` |
| macOS | `~/.Trash/` |
| Windows | Recycle Bin |

### Behavior

- File is recoverable from trash
- Original file is removed from source
- File metadata is preserved in trash info files (Linux)

## Shell Action

Executes a custom shell command with the matched file path available as an environment variable.

### Syntax

```toml
[[rules.actions]]
type = "shell"
command = "your-command-here"
```

### Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `type` | string | Yes | - | Must be `"shell"` |
| `command` | string | Yes | - | Shell command to execute |

### Environment Variables

| Variable | Description |
|----------|-------------|
| `$FILE` | Full path to the matched file |

### Examples

**Log to file:**
```toml
[[rules.actions]]
type = "shell"
command = "echo '$(date -Iseconds) $FILE' >> ~/file-events.log"
```

**Send notification:**
```toml
[[rules.actions]]
type = "shell"
command = "notify-send 'New file' '$FILE'"
```

**Compress file:**
```toml
[[rules.actions]]
type = "shell"
command = "gzip -k '$FILE'"
```

**Upload to cloud:**
```toml
[[rules.actions]]
type = "shell"
command = "rclone copy '$FILE' gdrive:uploads/"
```

**Process with custom script:**
```toml
[[rules.actions]]
type = "shell"
command = "/usr/local/bin/process-file.sh '$FILE'"
```

**Move with timestamp:**
```toml
[[rules.actions]]
type = "shell"
command = "mv '$FILE' '/archive/$(date +%Y%m%d)_$(basename $FILE)'"
```

**Scan with antivirus:**
```toml
[[rules.actions]]
type = "shell"
command = "clamscan '$FILE' --move=/quarantine"
```

### Platform-Specific Shells

| Platform | Shell Used |
|----------|------------|
| Linux/macOS | `/bin/sh` |
| Windows | `cmd.exe` |

### Best Practices

1. **Quote `$FILE`**: Always use single quotes around `$FILE` to handle paths with spaces
2. **Use absolute paths**: For commands, use full paths to avoid PATH issues
3. **Avoid interactive commands**: Commands that prompt for input will hang
4. **Test first**: Test commands manually before adding to rules
5. **Check exit codes**: Failed commands are logged as errors

### Common Commands

**File operations:**
```toml
# Rename with timestamp
command = "mv '$FILE' '$(dirname $FILE)/$(date +%Y%m%d)_$(basename $FILE)'"

# Change permissions
command = "chmod 644 '$FILE'"

# Change ownership
command = "chown user:group '$FILE'"
```

**Notifications:**
```toml
# Linux
cOMMAND = "notify-send 'Straw' 'Processed: $(basename $FILE)'"

# macOS
cOMMAND = "osascript -e 'display notification \"$(basename $FILE)\" with title \"Straw\"'"

# Windows PowerShell
cOMMAND = 'powershell -Command "Add-Type -AssemblyName System.Windows.Forms; [System.Windows.Forms.MessageBox]::Show(\'$FILE processed\')"'
```

**Logging:**
```toml
# Simple log
command = "echo '$(date): $FILE' >> ~/straw.log"

# JSON log
command = "echo '{\"time\":\"'$(date -Iseconds)'\",\"file\":\"'$FILE'\"}' >> ~/straw.json.log"
```

## Multiple Actions

Rules can have multiple actions that execute sequentially:

```toml
[[rules]]
name = "Backup and Log"
[rules.match]
extension = ".pdf"
[[rules.actions]]
type = "copy"
target = "~/Backup"
[[rules.actions]]
type = "shell"
command = "echo 'Backed up: $FILE' >> ~/backup.log"
[[rules.actions]]
type = "move"
target = "~/Documents/PDFs"
```

### Execution Order

Actions execute in the order defined in the config file:
1. Action 1
2. Action 2
3. Action 3
4. etc.

### Error Handling

- If an action fails, subsequent actions may still execute
- Errors are logged but don't stop the rule execution chain
- Use shell commands for more complex error handling

## Action Comparison

| Feature | Move | Copy | Trash | Shell |
|---------|------|------|-------|-------|
| Preserves original | No | Yes | No | Depends |
| Reversible | Manual | Yes | Yes (trash) | Depends |
| Cross-device | Yes | Yes | Yes | Yes |
| Performance | Fast | Slower | Fast | Depends |
| Custom logic | No | No | No | Yes |

## Common Patterns

### Safe Deletion

```toml
[[rules]]
name = "Clean Temp Files"
[rules.match]
glob = "*.tmp"
min_age_days = 7
[[rules.actions]]
type = "trash"
```

### Backup Strategy

```toml
[[rules]]
name = "Backup and Archive"
[rules.match]
extension = ".pdf"
[[rules.actions]]
type = "copy"
target = "~/Backup"
[[rules.actions]]
type = "move"
target = "~/Documents/Archive"
```

### Multi-Step Processing

```toml
[[rules]]
name = "Process Downloads"
[rules.match]
glob = "*.zip"
[[rules.actions]]
type = "copy"
target = "~/Backup/Originals"
[[rules.actions]]
type = "shell"
command = "unzip -q '$FILE' -d ~/Extracted/"
[[rules.actions]]
type = "shell"
command = "rm '$FILE'"
```

### Conditional Logging

```toml
[[rules]]
name = "Log Large Files"
[rules.match]
min_size = 104857600
[[rules.actions]]
type = "shell"
command = "echo '$(date): Large file detected: $FILE ($(stat -c%s '$FILE') bytes)' >> ~/large-files.log"
```
