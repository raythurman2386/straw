# Contributing to Straw

First off, thank you for considering contributing to Straw! It's people like you that make Straw such a great tool.

## How Can I Contribute?

### Reporting Bugs

If you find a bug, please create an issue on GitHub. Before creating an issue, please check if the issue has already been reported.

When creating an issue, please include:
- A clear and descriptive title
- Steps to reproduce the issue
- Expected behavior and what actually happened
- Your environment (OS, Go version, Straw version)

### Suggesting Enhancements

If you have an idea for a new feature or an improvement to an existing one, please create an issue to discuss it.

### Pull Requests

We welcome pull requests! To ensure a smooth process:
1. Fork the repository
2. Create a new branch for your changes
3. If you're adding a new feature or fixing a bug, please add tests
4. Ensure all tests pass (`make test`)
5. Submit a pull request with a clear description of your changes

## Development Setup

Straw is written in Go and runs on Linux, macOS, and Windows. You'll need Go 1.25 or later installed.

1. Clone your fork: `git clone https://github.com/youruser/straw.git`
2. Install dependencies: `go mod download`
3. Build the project: `make build`
4. Run tests: `make test`

### Cross-Platform Notes

When contributing, please keep cross-platform compatibility in mind:

- Use `os.UserConfigDir()`, `os.UserCacheDir()`, and `os.TempDir()` instead of hardcoded paths.
- Use `filepath.Join()` and `filepath.Separator` instead of hardcoded `/` in file paths.
- Place platform-specific code in files with build tags (e.g., `foo_unix.go`, `foo_windows.go`).
- If your change involves shell commands in tests, branch on `runtime.GOOS` to provide Windows-compatible equivalents.
- CI runs tests on Linux, macOS, and Windows -- all three must pass.

## Code of Conduct

Please be respectful and kind to others in the community. We follow the standard [Contributor Covenant Code of Conduct](https://www.contributorcovenant.org/version/2/1/code_of_conduct/).

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
