# Contributing to gogitopsdeployer

First of all, thank you for considering contributing to **gogitopsdeployer**! It is people like you who make the open-source community an amazing place to learn, inspire, and create.

## Principles & Guidelines

This project follows the **Antigravity** engineering pillars. Before proposing any changes, please ensure your code aligns with these principles:

1. **Extreme Modularization**: Single responsibility per file. Logic vs. Presentation separation.
2. **Stateless by Design**: No local memory state. Use distributed stores or persistence for stateful data.
3. **Typed Configuration**: Absolute zero hardcoded credentials or URLs. All configuration must be validated at boot time.
4. **Composition over Inheritance**: Small, reusable building blocks combined into complex systems.
5. **Contract First**: Typed interfaces between modules and shared types for communication.

## How to Contribute

### Reporting Bugs
- Use the **GitHub Issues** tab.
- Describe the unexpected behavior, steps to reproduce, and your environment.

### Proposing Features
- Open an issue to discuss the proposal before starting implementation.
- This ensures alignment with the project's roadmap and architectural pillars.

### Pull Requests
1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/amazing-feature`).
3. Ensure all tests pass (`make test`).
4. Run static analysis (`make lint`).
5. Commit your changes following [Conventional Commits](https://www.conventionalcommits.org/).
6. Push to the branch and open a PR.

## Coding Standards

- **Language**: Variable names, functions, and files in **English**.
- **Documentation**: All exported items must have **Godoc** comments.
- **Formatting**: Use `gofmt` (handled by `make tidy` or your IDE).

---
Let's build a more resilient GitOps future together!

## Issues and Discussions
Check the [Issues](https://github.com/ESousa97/gogitopsdeployer/issues) tab before starting work to avoid duplication.
