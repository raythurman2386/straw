# Configuration Examples

Practical configuration examples for common use cases.

## Download Organization

Keep your Downloads folder organized by file type:

```toml
[[watch]]
path = "~/Downloads"
recursive = false

# PDFs go to Documents
[[rules]]
name = "Organize PDFs"
enabled = true
[rules.match]
extension = ".pdf"
[[rules.actions]]
type = "move"
target = "~/Documents/PDFs"

# Images go to Pictures
[[rules]]
name = "Organize Images"
enabled = true
[rules.match]
glob = "*.{jpg,jpeg,png,gif,bmp,webp}"
[[rules.actions]]
type = "move"
target = "~/Pictures/Downloads"

# Archives go to dedicated folder
[[rules]]
name = "Organize Archives"
enabled = true
[rules.match]
glob = "*.{zip,tar,gz,bz2,7z,rar}"
[[rules.actions]]
type = "move"
target = "~/Downloads/Archives"

# Videos go to Videos folder
[[rules]]
name = "Organize Videos"
enabled = true
[rules.match]
glob = "*.{mp4,avi,mkv,mov,wmv}"
[[rules.actions]]
type = "move"
target = "~/Videos/Downloads"

# Audio files go to Music
[[rules]]
name = "Organize Audio"
enabled = true
[rules.match]
glob = "*.{mp3,wav,flac,aac,ogg}"
[[rules.actions]]
type = "move"
target = "~/Music/Downloads"

# Documents go to Documents
[[rules]]
name = "Organize Documents"
enabled = true
[rules.match]
glob = "*.{doc,docx,xls,xlsx,ppt,pptx,odt,ods,odp}"
[[rules.actions]]
type = "move"
target = "~/Documents/Downloads"
```

## Automated Cleanup

Clean up temporary and old files:

```toml
[[watch]]
path = "~/Downloads"
recursive = true

# Delete temp files older than 7 days
[[rules]]
name = "Clean Temp Files"
enabled = true
[rules.match]
glob = "*.tmp"
min_age_days = 7
[[rules.actions]]
type = "trash"

# Delete log files older than 30 days
[[rules]]
name = "Clean Old Logs"
enabled = true
[rules.match]
extension = ".log"
min_age_days = 30
[[rules.actions]]
type = "trash"

# Delete partial downloads older than 3 days
[[rules]]
name = "Clean Partial Downloads"
enabled = true
[rules.match]
glob = "*.part"
min_age_days = 3
[[rules.actions]]
type = "trash"

# Delete empty text files
[[rules]]
name = "Clean Empty Files"
enabled = true
[rules.match]
extension = ".txt"
max_size = 1
[[rules.actions]]
type = "trash"
```

## Backup Strategy

Create automated backups:

```toml
[[watch]]
path = "~/Documents"
recursive = true

[[watch]]
path = "~/Projects"
recursive = true

# Backup important documents
[[rules]]
name = "Backup Documents"
enabled = true
[rules.match]
regex = "^important_|_backup\\."
[[rules.actions]]
type = "copy"
target = "~/Backup/Documents"
[[rules.actions]]
type = "shell"
command = "notify-send 'Document backed up' '$(basename $FILE)'"

# Backup project files
[[rules]]
name = "Backup Project Files"
enabled = true
[rules.match]
extension = ".go"
max_age_days = 1
[[rules.actions]]
type = "copy"
target = "~/Backup/Code/Go"

# Backup large files to external drive
[[rules]]
name = "Archive Large Files"
enabled = true
[rules.match]
min_size = 104857600  # 100MB
min_age_days = 7
[[rules.actions]]
type = "copy"
target = "/mnt/external/Archive"
```

## Screenshot Management

Organize screenshots automatically:

```toml
[[watch]]
path = "~/Desktop"
recursive = false

[[watch]]
path = "~/Downloads"
recursive = false

# macOS screenshots
[[rules]]
name = "Organize macOS Screenshots"
enabled = true
[rules.match]
regex = "^Screen Shot [0-9]{4}-[0-9]{2}-[0-9]{2}"
extension = ".png"
[[rules.actions]]
type = "move"
target = "~/Pictures/Screenshots"

# Linux/GNOME screenshots
[[rules]]
name = "Organize GNOME Screenshots"
enabled = true
[rules.match]
regex = "^Screenshot from [0-9]{4}"
[[rules.actions]]
type = "move"
target = "~/Pictures/Screenshots"

# Windows screenshots (Snipping Tool)
[[rules]]
name = "Organize Windows Screenshots"
enabled = true
[rules.match]
regex = "^Screenshot [0-9]{4}"
[[rules.actions]]
type = "move"
target = "~/Pictures/Screenshots"
```

## Development Workflow

Organize development files:

```toml
[[watch]]
path = "~/Downloads"
recursive = false

# Move source code archives
[[rules]]
name = "Organize Source Archives"
enabled = true
[rules.match]
glob = "*.{tar.gz,tgz,zip}"
regex = "^[a-zA-Z0-9_-]+-[0-9]+\\.[0-9]+"
[[rules.actions]]
type = "move"
target = "~/Downloads/Source"

# Move package files
[[rules]]
name = "Organize Packages"
enabled = true
[rules.match]
glob = "*.{deb,rpm,pkg,dmg,exe,msi}"
[[rules.actions]]
type = "move"
target = "~/Downloads/Packages"

# Extract and organize tool archives
[[rules]]
name = "Extract Tools"
enabled = true
[rules.match]
glob = "*.tar.gz"
regex = "tool-|utility-"
[[rules.actions]]
type = "shell"
command = "tar -xzf '$FILE' -C ~/Tools/ && rm '$FILE'"
```

## Media Processing

Process media files automatically:

```toml
[[watch]]
path = "~/Downloads"
recursive = true

# Compress images
[[rules]]
name = "Compress Images"
enabled = true
[rules.match]
glob = "*.{jpg,jpeg,png}"
max_size = 10485760  # 10MB
[[rules.actions]]
type = "shell"
command = "convert '$FILE' -resize 1920x1080> -quality 85 '$(dirname $FILE)/$(basename $FILE | cut -d. -f1)_compressed.jpg'"

# Convert videos to MP4
[[rules]]
name = "Convert Videos"
enabled = true
[rules.match]
glob = "*.{avi,mkv,mov}"
[[rules.actions]]
type = "shell"
command = "ffmpeg -i '$FILE' -c:v libx264 -c:a aac '$(dirname $FILE)/$(basename $FILE | cut -d. -f1).mp4'"
[[rules.actions]]
type = "trash"
```

## Security & Monitoring

Monitor and secure your filesystem:

```toml
[[watch]]
path = "~/Downloads"
recursive = true

# Scan executables with antivirus
[[rules]]
name = "Scan Executables"
enabled = true
[rules.match]
glob = "*.{exe,dll,bat,sh,AppImage}"
[[rules.actions]]
type = "shell"
command = "clamscan '$FILE' --stdout | logger -t straw-antivirus"

# Log all downloaded files
[[rules]]
name = "Log Downloads"
enabled = true
[rules.match]
max_age_days = 1
[[rules.actions]]
type = "shell"
command = "echo '{\"timestamp\":\"'$(date -Iseconds)'\",\"file\":\"'$FILE'\",\"size\":'$(stat -c%s '$FILE')'}' >> ~/download-log.json"

# Alert on suspicious files
[[rules]]
name = "Alert Suspicious Files"
enabled = true
[rules.match]
regex = "(password|credential|secret|key).*\\.(txt|log|csv)"
[[rules.actions]]
type = "shell"
command = "notify-send -u critical 'Suspicious file detected' '$FILE'"
```

## Multi-Platform Configuration

A configuration that works across platforms:

```toml
# Platform-agnostic paths using ~
socket_path = "~/.straw/straw.sock"

# Watch common directories
[[watch]]
path = "~/Downloads"
recursive = true

# Platform-specific downloads folder (Windows)
[[watch]]
path = "~/Downloads"
recursive = true

# Organize by type
[[rules]]
name = "Organize Documents"
enabled = true
[rules.match]
glob = "*.{pdf,doc,docx,xls,xlsx,ppt,pptx}"
[[rules.actions]]
type = "move"
target = "~/Documents"

[[rules]]
name = "Organize Images"
enabled = true
[rules.match]
glob = "*.{jpg,jpeg,png,gif,bmp,webp}"
[[rules.actions]]
type = "move"
target = "~/Pictures"

# Cross-platform notification (using shell)
[[rules]]
name = "Notify on Large Download"
enabled = true
[rules.match]
min_size = 104857600
[[rules.actions]]
type = "shell"
# Linux
command = "notify-send 'Large file downloaded' '$(basename $FILE)' 2>/dev/null || echo 'Large file: $FILE'"
```

## Advanced Patterns

Complex real-world examples:

### Date-Based Organization

```toml
[[watch]]
path = "~/Downloads"
recursive = false

# Files with date in name
[[rules]]
name = "Archive Dated Files"
enabled = true
[rules.match]
regex = "[0-9]{4}-[0-9]{2}-[0-9]{2}"
[[rules.actions]]
type = "shell"
command = "YEAR=$(echo '$(basename $FILE)' | grep -oE '[0-9]{4}' | head -1) && mkdir -p ~/Archive/$YEAR && mv '$FILE' ~/Archive/$YEAR/"
```

### Email Attachment Organization

```toml
[[watch]]
path = "~/Downloads"
recursive = false

# Files from email clients often have specific patterns
[[rules]]
name = "Email PDFs"
enabled = true
[rules.match]
extension = ".pdf"
regex = "(attachment|re_|fw_)"
[[rules.actions]]
type = "move"
target = "~/Documents/Email"
```

### Build Artifact Cleanup

```toml
[[watch]]
path = "~/Projects"
recursive = true

# Clean build artifacts
[[rules]]
name = "Clean Node Modules"
enabled = true
[rules.match]
path = "**/node_modules"
min_age_days = 30
[[rules.actions]]
type = "shell"
command = "rm -rf '$FILE'"

# Clean build directories
[[rules]]
name = "Clean Build Directories"
enabled = true
[rules.match]
path = "**/{build,dist,target,out}"
min_age_days = 7
[[rules.actions]]
type = "trash"
```

## Tips for Creating Configurations

1. **Start simple**: Begin with basic rules and add complexity gradually
2. **Test first**: Use `copy` instead of `move` when testing new rules
3. **Use specific patterns**: More specific patterns reduce false matches
4. **Add age limits**: Prevent processing files currently being downloaded
5. **Monitor logs**: Check TUI event log to verify rules are working
6. **Organize rules**: Group related rules together in the config file
7. **Document rules**: Use descriptive names and comments
