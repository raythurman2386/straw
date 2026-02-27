# Contributing

Thank you for your interest in contributing to Straw! This guide will help you get started.

## Getting Started

### Prerequisites

- Go 1.25 or later
- Git
- Make (optional, for using Makefile commands)

### Setting Up Development Environment

1. **Fork and clone the repository:**
   ```bash
   git clone https://github.com/YOUR_USERNAME/straw.git
   cd straw
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Build the project:**
   ```bash
   make build
   # or
   go build -o bin/straw ./cmd/straw
   go build -o bin/strawd ./cmd/strawd
   ```

4. **Run tests:**
   ```bash
   make test
   # or
   go test ./...
   ```

## Development Workflow

### Branching Strategy

- `main` - Stable release branch
- `develop` - Development branch
- `feature/*` - Feature branches
- `bugfix/*` - Bug fix branches

### Making Changes

1. **Create a new branch:**
   ```bash
   git checkout -b feature/my-feature
   ```

2. **Make your changes** following the code style guidelines

3. **Test your changes:**
   ```bash
   make test
   make check
   ```

4. **Commit with a clear message:**
   ```bash
   git commit -m "feat: add new feature"
   ```

5. **Push and create a pull request:**
   ```bash
   git push origin feature/my-feature
   ```

## Code Style

### Go Conventions

We follow standard Go conventions:

- Use `gofmt` for formatting
- Use `goimports` for import organization
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Project-Specific Guidelines

See the project's [AGENTS.md](https://github.com/raythurman2386/straw/blob/main/AGENTS.md) for detailed project guidelines.

Key points:

1. **Imports ordering:**
   ```go
   import (
       "encoding/json"
       "fmt"
       "os"

       "github.com/spf13/cobra"

       "straw/internal/config"
   )
   ```

2. **Error handling:**
   ```go
   result, err := doSomething()
   if err != nil {
       return fmt.Errorf("failed to do something: %w", err)
   }
   ```

3. **Testing:**
   - Use table-driven tests
   - Name tests descriptively
   - Clean up test resources

## Testing

### Running Tests

```bash
# All tests
make test

# Specific package
go test ./internal/config

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

### Writing Tests

Example test structure:

```go
func TestFeature_Description(t *testing.T) {
    t.Run("Test case 1", func(t *testing.T) {
        // Test code
    })
    
    t.Run("Test case 2", func(t *testing.T) {
        // Test code
    })
}
```

## Linting and Formatting

```bash
# Run all checks
make check

# Format code
make fmt

# Run linter
golangci-lint run ./...

# Vet code
go vet ./...
```

## Documentation

- Update documentation for any user-facing changes
- Update [CHANGELOG.md](../changelog.md) for significant changes
- Add examples for new features
- Update configuration examples if needed

## Commit Message Format

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting)
- `refactor`: Code refactoring
- `test`: Test changes
- `chore`: Build/tooling changes

Examples:
```
feat(rules): add support for regex matching
fix(actions): handle spaces in file paths
docs: update installation instructions
```

## Pull Request Process

1. **Update documentation** for any changed functionality
2. **Add tests** for new features
3. **Ensure all tests pass**
4. **Update CHANGELOG.md** if applicable
5. **Request review** from maintainers
6. **Address review feedback**
7. **Squash commits** if requested

## Reporting Issues

### Bug Reports

Include:
- Operating system and version
- Straw version (`straw --version`)
- Steps to reproduce
- Expected behavior
- Actual behavior
- Configuration file (sanitized)
- Relevant log output

### Feature Requests

Include:
- Use case description
- Proposed solution
- Alternatives considered

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback
- Assume good intentions

## Getting Help

- **GitHub Discussions**: For questions and ideas
- **GitHub Issues**: For bugs and feature requests
- **Documentation**: Check the docs first

## Release Process

Releases are managed by maintainers:

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create git tag
4. Build release binaries
5. Create GitHub release
6. Update documentation

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
