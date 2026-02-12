# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Dry-run mode: Evaluate rules without modifying disk (In Progress).
- Desktop Notifications (In Progress).
- Windows Support (In Progress).

## [0.4.0] - 2026-02-11

### Added
- **Hot Reload**: Support for `SIGHUP` and IPC `TRIGGER_RELOAD`.
- **Rule Wizard**: Guided multi-step creation of rules in TUI.
- **Persistence**: Automatic writing of rules to `config.toml` from TUI.
- **Structured Logging**: Full `slog` implementation.
- **Service Management**: Systemd User Service integration.
- **Installation Scripts**: Unified `install.sh` and `setup.sh`.
- MIT License, CONTRIBUTING.md, and Issue/PR templates.

## [0.3.0] - 2026-02-10

### Added
- **TUI Core**: Root Bubble Tea model and update loop.
- **Activity Log**: Real-time streaming viewport.
- **Navigation**: Tabbed system (Activity, Rules, Creation).
- **Themes**: Everforest and Catppuccin support with `lipgloss`.
- **UI Visuals**: Consistent bordered containers and status indicators.

## [0.2.0] - 2026-02-09

### Added
- **Watcher**: Recursive directory watching with `fsnotify`.
- **Rules Engine**: Predicates for Glob, Regex, Extension, Size, Age, and Type.
- **Actions**: `move`, `copy`, `trash` (XDG), and `shell` executions.
- **Robustness**: Cross-filesystem move fallback and collision handling.

## [0.1.0] - 2026-02-09

### Added
- Initial project scaffold and directory layout.
- TOML configuration schema and validation.
- JSON-RPC over Unix Domain Sockets protocol.
- IPC Server/Client implementation.
- `straw` and `strawd` CLI core with `cobra`.
