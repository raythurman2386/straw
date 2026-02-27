# Architecture

Overview of Straw's architecture and design decisions.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Straw System                           │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐        IPC (JSON-RPC)     ┌───────────┐  │
│  │    straw     │  ←───────────────────────→ │  strawd   │  │
│  │   (TUI)      │       Unix Socket          │  (Daemon) │  │
│  └──────────────┘                            └─────┬─────┘  │
│        ↑                                           │        │
│        │                                           │        │
│   User Input                                  File System   │
│                                                  Events     │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## Components

### 1. Daemon (strawd)

The daemon is the core component that runs continuously:

**Responsibilities:**
- Filesystem monitoring
- Rule evaluation
- Action execution
- Configuration management
- IPC server

**Key Modules:**

```
internal/
├── config/     # Configuration loading and validation
├── watcher/    # Filesystem watching (fsnotify)
├── rules/      # Rule evaluation engine
├── actions/    # Action execution
├── ipc/        # IPC server (JSON-RPC)
└── logging/    # Logging infrastructure
```

**Architecture Pattern:**
- Event-driven architecture
- Goroutine per watched directory
- Channel-based communication
- Graceful shutdown handling

### 2. TUI Client (straw)

The TUI provides the user interface:

**Responsibilities:**
- Display real-time events
- Rule management interface
- Configuration editing
- IPC client

**Key Modules:**

```
cmd/straw/
├── main.go      # Entry point
├── wizard.go    # Rule wizard
├── settings.go  # Settings management
└── delegate.go  # UI delegation

internal/tui/
└── theme.go     # UI theming
```

**Framework:**
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Huh](https://github.com/charmbracelet/huh) - Forms

### 3. Inter-Process Communication (IPC)

**Protocol:** JSON-RPC 2.0 over Unix Domain Sockets

**Why JSON-RPC?**
- Simple and lightweight
- Language agnostic
- Easy debugging
- Well-documented standard

**Message Format:**
```json
{
  "jsonrpc": "2.0",
  "method": "methodName",
  "params": { ... },
  "id": 1
}
```

**Available Methods:**
- `GetStatus` - Get daemon status
- `GetEvents` - Get recent events
- `GetRules` - Get configured rules
- `AddRule` - Add a new rule
- `RemoveRule` - Remove a rule
- `ReloadConfig` - Reload configuration
- `GetWatches` - Get watched directories

### 4. Configuration System

**Format:** TOML

**Design Goals:**
- Human-readable
- Comments support
- Clear structure
- Easy to validate

**Structure:**
```toml
# Global settings
socket_path = "/tmp/straw.sock"

# Watch directories
[[watch]]
path = "~/Downloads"
recursive = true

# Rules
[[rules]]
name = "Example"
[rules.match]
extension = ".pdf"
[[rules.actions]]
type = "move"
target = "~/Documents"
```

### 5. Rules Engine

**Evaluation Flow:**

```
File Event
    ↓
Match Criteria Evaluation
    ↓
Action Execution
    ↓
Log Result
```

**Criteria Types:**
- Pattern matching (glob, regex)
- File attributes (extension, size, age)
- Visibility (hidden files)

**Action Types:**
- `move` - Move file
- `copy` - Copy file
- `trash` - Move to trash
- `shell` - Execute command

### 6. Filesystem Watcher

**Implementation:** [fsnotify](https://github.com/fsnotify/fsnotify)

**Features:**
- Cross-platform support
- Recursive watching
- Event debouncing
- Resource efficient

**Event Types:**
- Create
- Write
- Remove
- Rename
- Chmod

## Data Flow

### File Processing Flow

```
1. File Created
        ↓
2. Watcher Detects Event
        ↓
3. Event Sent to Rules Engine
        ↓
4. Rules Evaluated (in order)
        ↓
5. First Matching Rule Selected
        ↓
6. Actions Executed
        ↓
7. Result Logged
        ↓
8. Event Sent to TUI (via IPC)
```

### Configuration Flow

```
1. Config File Changed
        ↓
2. TUI Detects Change
        ↓
3. TUI Sends Reload Request
        ↓
4. Daemon Validates Config
        ↓
5. Daemon Applies Changes
        ↓
6. New Watches Started
        ↓
7. Old Watches Stopped
        ↓
8. Acknowledgment Sent to TUI
```

## Platform Support

### Cross-Platform Strategy

**Platform-Specific Code:**
- Build tags (`//go:build`)
- Separate files per platform
- Common interfaces

**Current Split:**
- `signal_unix.go` / `signal_windows.go`
- `trash_unix.go` / `trash_windows.go`
- `hidden_unix.go` / `hidden_windows.go`

### Platform Considerations

**Linux:**
- systemd service support
- XDG directories
- inotify for file watching

**macOS:**
- launchd service support
- FSEvents for file watching
- Finder trash integration

**Windows:**
- Task Scheduler support
- ReadDirectoryChangesW for watching
- Recycle Bin integration

## Performance Characteristics

### Resource Usage

| Component | CPU (Idle) | Memory (RSS) | Binary Size |
|-----------|------------|--------------|-------------|
| strawd | < 0.1% | ~6.8 MB | 7.2 MB |
| straw | ~0.0% | ~6.5 MB | 9.4 MB |

### Scalability

- **Watched directories:** Hundreds (limited by OS)
- **Rules:** Thousands (limited by evaluation time)
- **Events:** Platform-dependent (inotify/FSEvents limits)
- **Concurrent actions:** Configurable

### Optimization Strategies

1. **Event batching** - Process multiple events together
2. **Rule ordering** - Most specific rules first
3. **Lazy evaluation** - Stop evaluating once match fails
4. **Goroutine pooling** - Limit concurrent actions

## Security Considerations

### IPC Security

- Unix socket with restricted permissions
- User-level access only
- No network exposure

### File Operations

- Respects filesystem permissions
- No privilege escalation
- Safe path handling

### Shell Actions

- User-defined only
- No default shell commands
- Environment variable isolation

## Error Handling

### Strategy

- Fail fast on configuration errors
- Log and continue on runtime errors
- Graceful degradation

### Error Types

1. **Configuration Errors** - Fatal, prevent startup
2. **Runtime Errors** - Logged, operation continues
3. **Action Errors** - Logged, next action executes
4. **IPC Errors** - Logged, retry with backoff

## Future Architecture Plans

### Potential Enhancements

1. **Plugin System** - Allow custom actions
2. **Database Backend** - Persistent event storage
3. **Web UI** - Browser-based interface
4. **Metrics** - Prometheus/StatsD integration
5. **Remote Management** - Secure remote control

### Scalability Improvements

1. **Distributed Watching** - Multiple daemon instances
2. **Rule Compilation** - Compile rules for faster evaluation
3. **Async Actions** - Non-blocking action execution
4. **Event Streaming** - Kafka/RabbitMQ integration
