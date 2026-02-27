# Rules

Rules are the core of Straw's automation capabilities. They define what files to match and what actions to take.

## Rule Structure

A rule consists of:
- **Name**: Human-readable identifier
- **Enabled**: Whether the rule is active
- **Match**: Criteria for matching files
- **Actions**: What to do with matched files

```toml
[[rules]]
name = "My Rule"
enabled = true
[rules.match]
# Match criteria
[[rules.actions]]
# Actions
```

## Match Criteria

Rules can match files based on multiple criteria. All criteria must match (AND logic).

### Pattern Matching

#### Glob Patterns

Shell-style pattern matching for filenames:

```toml
[rules.match]
glob = "*.pdf"           # Match all PDF files
glob = "backup_*.zip"    # Match files starting with "backup_" and ending with ".zip"
glob = "**/*.log"        # Match all .log files in any subdirectory
```

**Supported glob syntax:**
- `*` - Match any sequence of characters (except `/`)
- `?` - Match any single character
- `**` - Match any sequence of characters (including `/`)
- `[abc]` - Match any character in brackets
- `[!abc]` - Match any character NOT in brackets

#### Regular Expressions

Full regex support for complex matching:

```toml
[rules.match]
regex = "^temp_.*\\.txt$"   # Match files starting with "temp_" and ending with ".txt"
regex = "\\.(jpg|jpeg|png|gif)$"  # Match image files
regex = "^[0-9]{4}-[0-9]{2}-[0-9]{2}_"  # Match files with date prefix
```

**Note:** Double backslashes are required in TOML strings.

### File Attributes

#### Extension

Match files by their extension:

```toml
[rules.match]
extension = ".pdf"
extension = ".tar.gz"  # Multi-part extensions supported
```

#### Size

Match files by size in bytes:

```toml
[rules.match]
min_size = 1048576      # At least 1MB
max_size = 1073741824   # At most 1GB
```

Common size values:
- 1 KB = 1024 bytes
- 1 MB = 1048576 bytes
- 1 GB = 1073741824 bytes

#### Age

Match files by age in days:

```toml
[rules.match]
min_age_days = 7        # Older than 1 week
max_age_days = 1        # Newer than 1 day
```

#### Hidden Files

Match hidden files (dotfiles on Unix, system attribute on Windows):

```toml
[rules.match]
hidden = true   # Match hidden files only
hidden = false  # Match visible files only
```

### Combining Criteria

Multiple criteria are combined with AND logic:

```toml
[rules.match]
extension = ".log"
min_age_days = 30
max_size = 10485760
```

This matches files that are:
- `.log` files AND
- Older than 30 days AND
- Smaller than 10MB

## Rule Evaluation

Rules are evaluated in order. The first matching rule's actions are executed on a file.

### Evaluation Order

```toml
# Rule 1 - Evaluated first
[[rules]]
name = "Important Files"
[rules.match]
regex = "^important_.*"

# Rule 2 - Evaluated second
[[rules]]
name = "All PDFs"
[rules.match]
extension = ".pdf"
```

A file named `important_document.pdf` will only trigger "Important Files", not "All PDFs".

### Disabling Rules

Disable a rule without removing it:

```toml
[[rules]]
name = "Temporary Rule"
enabled = false
```

Disabled rules are skipped during evaluation.

## Rule Examples

### Organize Downloads

```toml
[[rules]]
name = "Organize PDF Downloads"
enabled = true
[rules.match]
glob = "*.pdf"
[[rules.actions]]
type = "move"
target = "~/Documents/PDFs"
```

### Clean Up Old Files

```toml
[[rules]]
name = "Clean Old Temp Files"
enabled = true
[rules.match]
glob = "*.tmp"
min_age_days = 7
[[rules.actions]]
type = "trash"
```

### Archive Large Files

```toml
[[rules]]
name = "Archive Large Downloads"
enabled = true
[rules.match]
min_size = 104857600  # 100MB
[[rules.actions]]
type = "move"
target = "~/Archive/Large"
```

### Backup and Notify

```toml
[[rules]]
name = "Backup Important Files"
enabled = true
[rules.match]
regex = "^important_.*"
[[rules.actions]]
type = "copy"
target = "~/Backup"
[[rules.actions]]
type = "shell"
command = "notify-send 'Important file backed up' '$FILE'"
```

### Multiple Match Criteria

```toml
[[rules]]
name = "Old Log Files"
enabled = true
[rules.match]
extension = ".log"
min_age_days = 30
max_size = 10485760
[[rules.actions]]
type = "trash"
```

### Complex Pattern

```toml
[[rules]]
name = "Screenshot Organization"
enabled = true
[rules.match]
regex = "^(Screenshot|Screen Shot)\\s.*\\.(png|jpg)$"
[[rules.actions]]
type = "move"
target = "~/Pictures/Screenshots"
```

## Best Practices

1. **Order matters**: Put specific rules before general ones
2. **Use descriptive names**: Make rules self-documenting
3. **Test patterns**: Verify glob/regex patterns before deploying
4. **Start with copy**: Use `copy` instead of `move` when testing new rules
5. **Use age/size limits**: Prevent accidental processing of files currently in use
6. **Disable, don't delete**: Use `enabled = false` to temporarily disable rules

## Debugging Rules

### Enable Verbose Logging

```bash
strawd --verbose
```

### Check Rule Evaluation

Watch the TUI event log to see:
- Which files are being matched
- Which rules are triggering
- Action execution status

### Test Patterns

Test your patterns before adding them:

```bash
# Test glob pattern
ls ~/Downloads/*.pdf

# Test regex pattern (with grep)
ls ~/Downloads | grep -E "^important_.*\.pdf$"
```

### Common Issues

| Issue | Solution |
|-------|----------|
| Rule never matches | Check pattern syntax and file paths |
| Wrong rule matches | Reorder rules (specific first) |
| Actions not executed | Check target directory exists and has write permissions |
| Rule matches too much | Add more specific criteria |
