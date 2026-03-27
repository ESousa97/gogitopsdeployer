# Contributing to gogitopsdeployer

Thank you for your interest in contributing to **gogitopsdeployer**! This project follows the **Antigravity** engineering principles: modularity, statelessness, and robust error handling.

## Development Setup

### Prerequisites
- **Go**: Version 1.25.0 or higher.
- **Git**: Latest version.
- **Make**: For running build automation.

### Initializing Environment
```bash
git clone https://github.com/ESousa97/gogitopsdeployer.git
cd gogitopsdeployer
make tidy
```

## Engineering Standards

### Code Style
- Follow [Effective Go](https://golang.org/doc/effective_go.html).
- Run `make lint` before every commit.
- Every exported function, type, or constant **must** have a Godoc comment.
- No "magic values". Use configuration or constants.

### Project Structure
- `cmd/`: CLI entry points. Keep `main.go` lean.
- `internal/`: Private business logic. Organized by bounded contexts.
- `pkg/`: (If created) Libraries safe for external use.

## Workflow

1. **Fork** the repository.
2. Create a **feature branch** (`feat/your-feature` or `fix/your-fix`).
3. Implement changes with **unit tests**.
4. Run `make test` and `make lint`.
5. Submit a **Pull Request** with a clear description of changes.

## Commit Convention
We prefer [Conventional Commits](https://www.conventionalcommits.org/):
- `feat:` for new features.
- `fix:` for bug fixes.
- `docs:` for documentation updates.
- `refactor:` for code restructuring.

## Issues and Discussions
Check the [Issues](https://github.com/ESousa97/gogitopsdeployer/issues) tab before starting work to avoid duplication.
