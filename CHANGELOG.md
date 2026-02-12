# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- Dry-run mode: Evaluate rules without modifying disk.
- Desktop notifications.
- Windows support.

## [0.1.0-beta] - 2026-02-11

### Added
- **CLI**: `straw` TUI client and `strawd` daemon with Cobra, `--version` flag.
- **Configuration**: TOML schema with validation, live reload via `SIGHUP` and IPC.
- **Watcher**: Recursive directory watching with `fsnotify`.
- **Rules Engine**: Match predicates for glob, regex, extension, size, age, and file type.
- **Actions**: `move`, `copy`, `trash` (XDG compliant), and `shell` execution.
- **TUI**: Bubble Tea interface with activity log, rules list, and rule creation wizard.
- **Themes**: Everforest and Catppuccin color schemes with `lipgloss`.
- **IPC**: JSON-RPC over Unix Domain Sockets for client-daemon communication.
- **Service Management**: Systemd user service integration with install/uninstall scripts.
- **Release Workflow**: GoReleaser with GitHub Actions for cross-platform builds (Linux, macOS, Raspberry Pi).
- **Installer**: One-liner install script for downloading pre-built binaries from GitHub releases.
- **Packaging**: DEB, RPM, APK, and Arch Linux packages via GoReleaser.

[Unreleased]: https://github.com/raythurman2386/straw/compare/v0.1.0-beta...HEAD
[0.1.0-beta]: https://github.com/raythurman2386/straw/releases/tag/v0.1.0-beta
