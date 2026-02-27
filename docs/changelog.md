# Changelog

All notable changes to Straw will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial documentation site with MkDocs
- Comprehensive user guides and reference documentation

## [1.0.0] - 2024-XX-XX

### Added
- Initial release of Straw file automation system
- Background daemon (`strawd`) for filesystem monitoring
- Interactive TUI client (`straw`) for real-time monitoring
- Rule-based file processing engine
- Support for glob and regex pattern matching
- File attribute matching (extension, size, age, hidden)
- Four action types: move, copy, trash, shell
- Cross-platform support (Linux, macOS, Windows)
- JSON-RPC IPC over Unix sockets
- Live configuration reloading
- Systemd service support on Linux
- Interactive Rule Wizard in TUI
- Event logging and monitoring
- Multi-directory watching with recursive support

### Features

#### Daemon (strawd)
- Filesystem watching using fsnotify
- Rule evaluation and action execution
- Configuration management
- IPC server for client communication
- Signal-based configuration reloading
- Comprehensive logging

#### TUI Client (straw)
- Real-time event monitoring
- Rule management interface
- Interactive Rule Wizard
- Configuration editing
- Status monitoring
- Keyboard-driven interface

#### Configuration
- TOML-based configuration
- Platform-specific default paths
- Watch directory configuration
- Rule definition with multiple criteria
- Multiple actions per rule

#### Actions
- `move`: Move files to target directory
- `copy`: Copy files to target directory
- `trash`: Move files to system trash
- `shell`: Execute custom shell commands

#### IPC API
- JSON-RPC 2.0 protocol
- Methods: GetStatus, GetEvents, GetRules, AddRule, RemoveRule, EnableRule, DisableRule, ReloadConfig, GetWatches, TriggerAction
- Unix Domain Sockets (named pipes on Windows)

### Platform Support
- Linux with systemd integration
- macOS with launchd support
- Windows 10 1803+ with Task Scheduler support

### Performance
- Minimal resource footprint (~7MB RAM, <0.1% CPU idle)
- Efficient event-driven architecture
- Concurrent action execution
- Scalable to hundreds of watched directories

[Unreleased]: https://github.com/raythurman2386/straw/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/raythurman2386/straw/releases/tag/v1.0.0
