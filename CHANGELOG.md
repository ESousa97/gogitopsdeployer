# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-03-26

### Added
- Professional documentation suite (README, CONTRIBUTING, LICENSE, etc.).
- Robust Makefile for build automation.
- Full Godoc coverage for internal packages and CLI.
- Modular architecture implementation based on Domain-Driven Design.
- GitOps monitoring loop with commit hash detection.
- SSH execution engine with automatic rollback capability.
- Discord notification system for deployment lifecycle events.
- SQLite persistence layer for deployment history.
- Webhook support with HMAC validation.

### Changed
- Standardized project structure for better scalability.
- Improved error handling and resilience in the monitoring loop.

### Fixed
- Inconsistencies in remote execution state management.
