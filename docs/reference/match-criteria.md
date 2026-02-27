# Match Criteria Reference

Complete reference for all file matching criteria available in Straw rules.

## Overview

Match criteria are defined in the `[rules.match]` section of a rule:

```toml
[[rules]]
name = "Example Rule"
[rules.match]
# Criteria go here
```

All criteria are combined with AND logic - a file must match ALL specified criteria for the rule to trigger.

## Criteria Types

### Pattern Matching

#### `glob`

Shell-style glob pattern matching for filenames.

```toml
[rules.match]
glob = "*.pdf"
```

**Syntax:**

| Pattern | Matches |
|---------|---------|
| `*` | Any sequence of characters (except `/`) |
| `?` | Any single character |
| `**` | Any sequence of characters (including `/`) |
| `[abc]` | Any character in brackets |
| `[!abc]` | Any character NOT in brackets |
| `{a,b,c}` | Any of the comma-separated patterns |

**Examples:**

```toml
# Match all PDFs
glob = "*.pdf"

# Match files starting with "temp_"
glob = "temp_*"

# Match files in any subdirectory
glob = "**/*.log"

# Match specific extensions
glob = "*.{jpg,jpeg,png,gif}"

# Match files with date prefix
glob = "2024-??-??_*.txt"
```

#### `regex`

Regular expression matching for complex patterns.

```toml
[rules.match]
regex = "^backup_.*\\.zip$"
```

**Note:** Double backslashes are required in TOML strings. The regex is matched against the filename (not full path) by default.

**Examples:**

```toml
# Match files starting with "temp_"
regex = "^temp_"

# Match image files
regex = "\\.(jpg|jpeg|png|gif|bmp|webp)$"

# Match files with date format (YYYY-MM-DD)
regex = "^\\d{4}-\\d{2}-\\d{2}_"

# Match files with numbers
regex = "[0-9]+"

# Match specific prefix and suffix
regex = "^report_.*\\.pdf$"
```

### File Attributes

#### `extension`

Match files by their extension (case-insensitive).

```toml
[rules.match]
extension = ".pdf"
```

**Format:**
- Include the leading dot
- Case-insensitive matching
- Multi-part extensions supported

**Examples:**

```toml
# Single extension
extension = ".pdf"

# Multi-part extension
extension = ".tar.gz"

# Case insensitive (matches .PDF, .pdf, .Pdf, etc.)
extension = ".txt"
```

#### `min_size`

Match files with size greater than or equal to the specified value (in bytes).

```toml
[rules.match]
min_size = 1048576  # 1 MB
```

**Common values:**

| Size | Bytes |
|------|-------|
| 1 KB | 1024 |
| 1 MB | 1048576 |
| 10 MB | 10485760 |
| 100 MB | 104857600 |
| 1 GB | 1073741824 |

**Examples:**

```toml
# Files larger than 100 MB
min_size = 104857600

# Files larger than 1 GB
min_size = 1073741824

# Combine with max_size for a range
min_size = 10485760   # 10 MB
max_size = 104857600  # 100 MB
```

#### `max_size`

Match files with size less than or equal to the specified value (in bytes).

```toml
[rules.match]
max_size = 10485760  # 10 MB
```

**Examples:**

```toml
# Files smaller than 1 MB
max_size = 1048576

# Files smaller than 100 KB
max_size = 102400

# Combine with min_size for a range
min_size = 1024       # 1 KB
max_size = 1048576    # 1 MB
```

#### `min_age_days`

Match files that are at least N days old.

```toml
[rules.match]
min_age_days = 7
```

**Calculation:** File modification time is compared to current time.

**Examples:**

```toml
# Files older than 1 week
min_age_days = 7

# Files older than 30 days
min_age_days = 30

# Files older than 1 year
min_age_days = 365

# Combine with max_age_days for a range
min_age_days = 7
max_age_days = 30
```

#### `max_age_days`

Match files that are at most N days old (recent files).

```toml
[rules.match]
max_age_days = 1
```

**Examples:**

```toml
# Files created/modified today
max_age_days = 1

# Files created/modified in the last week
max_age_days = 7

# Files created/modified in the last hour (use fractions)
max_age_days = 0.041  # 1 hour / 24 hours
```

#### `hidden`

Match files based on their hidden status.

```toml
[rules.match]
hidden = true
```

**Platform behavior:**

| Platform | Hidden File Detection |
|----------|----------------------|
| Linux/macOS | Files starting with `.` |
| Windows | Files with hidden attribute |

**Examples:**

```toml
# Match only hidden files
hidden = true

# Match only visible files
hidden = false
```

## Combining Criteria

Multiple criteria are combined with AND logic:

```toml
[rules.match]
extension = ".log"
min_age_days = 30
max_size = 10485760
```

This matches files that meet ALL criteria:
1. Extension is `.log` AND
2. File is at least 30 days old AND
3. File is at most 10 MB

## OR Logic

To achieve OR logic (match any of multiple criteria), create multiple rules:

```toml
# Rule 1: Match PDFs
[[rules]]
name = "PDF Files"
[rules.match]
extension = ".pdf"
[[rules.actions]]
type = "move"
target = "~/Documents"

# Rule 2: Match DOC files
[[rules]]
name = "DOC Files"
[rules.match]
extension = ".doc"
[[rules.actions]]
type = "move"
target = "~/Documents"
```

## Complete Examples

### Download Organization

```toml
[[rules]]
name = "Organize Downloads"
[rules.match]
glob = "*.pdf"
max_age_days = 1
[[rules.actions]]
type = "move"
target = "~/Downloads/Recent/PDFs"
```

### Cleanup Old Logs

```toml
[[rules]]
name = "Clean Old Logs"
[rules.match]
extension = ".log"
min_age_days = 30
max_size = 104857600
[[rules.actions]]
type = "trash"
```

### Archive Large Files

```toml
[[rules]]
name = "Archive Large Downloads"
[rules.match]
min_size = 104857600
min_age_days = 7
[[rules.actions]]
type = "move"
target = "~/Archive"
```

### Process Screenshots

```toml
[[rules]]
name = "Process Screenshots"
[rules.match]
regex = "^(Screenshot|Screen Shot)"
extension = ".png"
max_age_days = 1
[[rules.actions]]
type = "move"
target = "~/Pictures/Screenshots"
```

### Backup Important Files

```toml
[[rules]]
name = "Backup Important"
[rules.match]
regex = "^important_|_backup\\."
hidden = false
[[rules.actions]]
type = "copy"
target = "~/Backup"
```

## Pattern Testing

Before adding patterns to your config, test them:

### Testing Glob Patterns

```bash
# List matching files
ls ~/Downloads/*.pdf

# Test with find
find ~/Downloads -name "*.pdf" -type f

# Test complex patterns
find ~/Downloads -name "temp_*" -o -name "*.tmp"
```

### Testing Regex Patterns

```bash
# Test with grep
grep -E "^temp_.*\\.txt$" <<< "temp_file.txt"

# Test with files
ls ~/Downloads | grep -E "^important_.*"

# Test with find and regex
find ~/Downloads -regextype posix-extended -regex ".*important.*"
```

## Performance Considerations

1. **Simple criteria first**: Put simpler checks (extension) before complex ones (regex)
2. **Specific patterns**: More specific patterns are faster to evaluate
3. **Combine criteria**: Use multiple criteria to reduce false positives
4. **Avoid over-broad patterns**: `**/*` can be slow on large directory trees
