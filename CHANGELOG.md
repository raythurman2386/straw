# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.4-beta] - 2026-02-25

### Fixed
- **Installation**: Improved script robustness with dependency checks for `curl`, `wget`, and `tar`.
- **Portability**: Added OS and `systemctl` checks to installation scripts to prevent failures on non-systemd/non-Linux systems.
- **Build**: Set `CGO_ENABLED=0` in Makefile to ensure portability and fix build errors in environments without C headers.
- **Documentation**: Corrected documentation URL in systemd service file.

### Added
- **Themes**: New "Ravenwood" and "Nord" color schemes.
- **Cross-platform support**: Linux, macOS, and Windows (10 1803+) are now all supported.
- **Windows trash**: Native Recycle Bin integration via `SHFileOperationW` (shell32.dll).
- **Windows hidden files**: Detects `FILE_ATTRIBUTE_HIDDEN` in addition to dot-prefix convention.
- **Platform-aware config paths**: Uses `os.UserConfigDir()` / `os.UserCacheDir()` for OS-appropriate defaults.
- **Windows signal handling**: Graceful shutdown via `os.Interrupt` (Ctrl+C); config reload via IPC `TRIGGER_RELOAD`.
- **CI matrix**: GitHub Actions now tests on Ubuntu, macOS, and Windows.
- **GoReleaser**: Windows builds (amd64, 386, arm64) with zip archives and Scoop bucket.

### Planned
- Dry-run mode: Evaluate rules without modifying disk.
- Desktop notifications.

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
- **Release Workflow**: GoReleaser with GitHub Actions for cross-platform builds (Linux, macOS, Windows, Raspberry Pi).
- **Installer**: One-liner install script for downloading pre-built binaries from GitHub releases.
- **Packaging**: DEB, RPM, APK, and Arch Linux packages via GoReleaser.

[Unreleased]: https://github.com/raythurman2386/straw/compare/v0.1.4-beta...HEAD
[0.1.4-beta]: https://github.com/raythurman2386/straw/compare/v0.1.3-beta...v0.1.4-beta
[0.1.0-beta]: https://github.com/raythurman2386/straw/releases/tag/v0.1.0-beta
