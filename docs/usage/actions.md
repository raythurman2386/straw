# Actions

Actions define what happens to files that match a rule's criteria.

## Action Types

Straw supports four action types:

| Type | Description | Use Case |
|------|-------------|----------|
| `move` | Move file to target directory | Organizing files |
| `copy` | Copy file to target directory | Backing up files |
| `trash` | Move file to system trash | Safe deletion |
| `shell` | Execute custom shell command | Custom processing |

## Move Action

Moves a file to a target directory.

```toml
[[rules.actions]]
type = "move"
target = "~/Documents/PDFs"
```

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `type` | String | Yes | Must be `"move"` |
| `target` | String | Yes | Destination directory path |

### Examples

**Organize by type:**
```toml
[[rules.actions]]
type = "move"
target = "~/Documents/PDFs"
```

**Archive old files:**
```toml
[[rules.actions]]
type = "move"
target = "~/Archive/2024"
```

### Behavior

- The target directory is created if it doesn't exist
- Files are moved (not copied), so they no longer exist at the original location
- If a file with the same name exists in the target, it may be overwritten (platform-dependent)
- Original file permissions are preserved

## Copy Action

Copies a file to a target directory, keeping the original.

```toml
[[rules.actions]]
type = "copy"
target = "~/Backup/Downloads"
```

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `type` | String | Yes | Must be `"copy"` |
| `target` | String | Yes | Destination directory path |

### Examples

**Backup important files:**
```toml
[[rules.actions]]
type = "copy"
target = "~/Backup/Important"
```

**Mirror to external drive:**
```toml
[[rules.actions]]
type = "copy"
target = "/mnt/backup/Downloads"
```

### Behavior

- The target directory is created if it doesn't exist
- The original file remains in place
- File permissions and timestamps are preserved
- If a file with the same name exists, it may be overwritten

## Trash Action

Moves a file to the system trash/recycle bin.

```toml
[[rules.actions]]
type = "trash"
```

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `type` | String | Yes | Must be `"trash"` |

### Platform Behavior

| Platform | Trash Location |
|----------|----------------|
| Linux | XDG trash (`~/.local/share/Trash/`) |
| macOS | Finder trash (`~/.Trash/`) |
| Windows | Recycle Bin |

### Examples

**Clean up temp files:**
```toml
[[rules.actions]]
type = "trash"
```

**Safe deletion with notification:**
```toml
[[rules.actions]]
type = "trash"
[[rules.actions]]
type = "shell"
command = "notify-send 'File moved to trash' '$(basename $FILE)'"
```

### Behavior

- Files can be recovered from trash
- Original file permissions are preserved
- Files are recoverable until trash is emptied

## Shell Action

Executes a custom shell command with the file path available as `$FILE`.

```toml
[[rules.actions]]
type = "shell"
command = "echo 'New file: $FILE' >> ~/file-log.txt"
```

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `type` | String | Yes | Must be `"shell"` |
| `command` | String | Yes | Shell command to execute |

### Environment Variables

| Variable | Description |
|----------|-------------|
| `$FILE` | Full path to the matched file |

### Examples

**Log file creation:**
```toml
[[rules.actions]]
type = "shell"
command = "echo '$(date): $FILE' >> ~/file-events.log"
```

**Send notification:**
```toml
[[rules.actions]]
type = "shell"
command = "notify-send 'File processed' '$FILE'"
```

**Compress and move:**
```toml
[[rules.actions]]
type = "shell"
command = "gzip -c '$FILE' > '/archive/$(basename $FILE).gz' && rm '$FILE'"
```

**Upload to cloud:**
```toml
[[rules.actions]]
type = "shell"
command = "rclone copy '$FILE' remote:backup/"
```

**Process with custom script:**
```toml
[[rules.actions]]
type = "shell"
command = "/usr/local/bin/process-file.sh '$FILE'"
```

### Behavior

- Commands run in a shell (`/bin/sh` on Unix, `cmd.exe` on Windows)
- The working directory is the directory containing the file
- Command output is logged but not displayed in the TUI
- If the command fails (non-zero exit), it's logged as an error
- Use single quotes around `$FILE` to handle paths with spaces

### Security Considerations

⚠️ **Warning:** Shell actions can execute arbitrary code. Only use commands you trust.

- Avoid unsanitized user input in commands
- Be careful with file paths containing special characters
- Test commands manually before adding them to rules

## Multiple Actions

A rule can have multiple actions that execute in sequence:

```toml
[[rules]]
name = "Backup and Log"
enabled = true
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

Actions execute in the order defined. If one action fails, subsequent actions may still execute (depends on implementation).

## Action Selection Guide

| Goal | Recommended Action |
|------|-------------------|
| Organize files | `move` |
| Create backups | `copy` |
| Delete safely | `trash` |
| Custom processing | `shell` |
| Rename files | `shell` with `mv` command |
| Compress files | `shell` with compression tool |
| Upload to cloud | `shell` with cloud CLI |
| Send notifications | `shell` with notification tool |

## Troubleshooting

### Move/Copy: Permission Denied

Ensure the target directory exists and is writable:
```bash
mkdir -p ~/target/directory
chmod 755 ~/target/directory
```

### Trash: File Not in Trash

Check the appropriate trash location for your platform:
```bash
# Linux
ls ~/.local/share/Trash/files/

# macOS
ls ~/.Trash/
```

### Shell: Command Not Found

Use absolute paths for commands:
```toml
# Instead of:
command = "notify-send 'Done'"

# Use:
command = "/usr/bin/notify-send 'Done'"
```

### Shell: File Path Issues

Always quote `$FILE` to handle spaces:
```toml
# Good:
command = "echo '$FILE' >> log.txt"

# Bad (will break on paths with spaces):
command = "echo $FILE >> log.txt"
```

### Shell: Command Hangs

Avoid interactive commands:
```toml
# Don't use commands that prompt for input
command = "rm -i '$FILE'"  # This will hang!

# Use non-interactive alternatives
command = "rm -f '$FILE'"
```
