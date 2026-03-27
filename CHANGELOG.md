# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-03-27

### Added
- **Standard Documentation**: Added comprehensive `README.md`, `CONTRIBUTING.md`, `SECURITY.md`, and `CODE_OF_CONDUCT.md`.
- **Infrastructure Automation**: Integrated `Makefile` for standardized build, test, and linting pipelines.
- **Godoc Documentation**: 100% documentation coverage for all exported types and functions in `internal` and `cmd` packages.
- **GitOps Engine**: High-performance monitoring and change detection for remote Git repositories.
- **SSH Deployment**: Robust implementation of remote command execution with automatic rollback.
- **Persistence**: SQLite-backed storage for deployment history tracking.
- **Notification**: Real-time alerts via Discord webhooks.
- **Immediate Triggers**: HTTP listener for GitHub webhooks with HMAC signature validation.

---
*Initial professional release.*
zed project structure for better scalability.
- Improved error handling and resilience in the monitoring loop.

### Fixed
- Inconsistencies in remote execution state management.
