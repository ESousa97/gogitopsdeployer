<div align="center">
  <h1>Go GitOps Deployer</h1>
  <p>A lightweight, modular GitOps agent for automated repository monitoring and remote deployment via SSH.</p>

  <img src="assets/github-go.png" alt="Go GitOps Deployer Banner" width="600px">

  <br>

  ![CI](https://github.com/esousa97/gogitopsdeployer/actions/workflows/ci.yml/badge.svg)
  [![Go Report Card](https://goreportcard.com/badge/github.com/esousa97/gogitopsdeployer?v=2)](https://goreportcard.com/report/github.com/esousa97/gogitopsdeployer)
  [![CodeFactor](https://www.codefactor.io/repository/github/esousa97/gogitopsdeployer/badge)](https://www.codefactor.io/repository/github/esousa97/gogitopsdeployer)
  ![Go Reference](https://pkg.go.dev/badge/github.com/esousa97/gogitopsdeployer.svg)
  ![License](https://img.shields.io/github/license/esousa97/gogitopsdeployer)
  ![Go Version](https://img.shields.io/github/go-mod/go-version/esousa97/gogitopsdeployer)
  ![Last Commit](https://img.shields.io/github/last-commit/esousa97/gogitopsdeployer)
</div>

---

**gogitopsdeployer** is a GitOps agent designed for simplicity and reliability. It monitors a Git repository for changes and automatically triggers a deployment process to remote servers via SSH. It features built-in rollback mechanisms, Discord notifications, and a GitHub webhook listener for instantaneous updates.

## Roadmap

- [x] **Phase 1: Foundation** — Project structure, Configuration (Env Vars), and standard Logging.
- [x] **Phase 2: GitOps Core** — Repository monitoring (`go-git`), commit hash comparison, and change detection.
- [x] **Phase 3: Deployment & Infrastructure** — SSH client implementation, remote command execution, and rollback logic.
- [x] **Phase 4: Persistence & Observability** — SQLite storage for deployment history and Discord webhook notifications.
- [x] **Phase 5: Immediate Triggers** — GitHub Webhook listener with HMAC signature validation.

## Quick Start

### Installation

```bash
# Via go install
go install github.com/esousa97/gogitopsdeployer/cmd/agent@latest
```

### From source

```bash
git clone https://github.com/esousa97/gogitopsdeployer.git
cd gogitopsdeployer
make build
```

## Makefile Targets

| Target       | Description                               |
| ------------ | ----------------------------------------- |
| `make build` | Compiles the agent binary in `bin/`       |
| `make run`   | Executes the agent directly via Go        |
| `make test`  | Runs the unit test suite                  |
| `make lint`  | Performs static analysis via `go vet`     |
| `make tidy`  | Cleans and updates `go.mod` dependencies  |
| `make clean` | Removes binary and temporary files        |

## Architecture

The project follows **Dependency Inversion** and **Modular Architecture** principles, ensuring that the business core remains infrastructure-agnostic.

- `cmd/agent`: Entry point and CLI commands.
- `internal/gitops`: Abstraction for Git operations and commit detection.
- `internal/ssh`: Remote execution engine and rollback logic.
- `internal/monitor`: Resilient reconciliation loop orchestrator.
- `internal/storage`: SQLite persistence layer.
- `internal/notification`: Discord notification integration.
- `internal/webhook`: GitHub Push event receiver.

## API Reference

Detailed package and function documentation is available via Godoc:
"Check the full documentation at [pkg.go.dev](https://pkg.go.dev/github.com/esousa97/gogitopsdeployer)."

## Configuration

| Variable                   | Description                       | Type     | Default            |
| -------------------------- | ------------------------------- | -------- | ------------------ |
| `GOGITOPS_REPO_URL`        | Target Git repository URL       | String   | Current repo       |
| `GOGITOPS_INTERVAL`        | Check interval duration         | Duration | `30s`              |
| `GOGITOPS_DB_PATH`         | Path to SQLite database         | String   | `./deployments.db` |
| `GOGITOPS_SSH_HOST`        | Deployment machine Host or IP   | String   | -                  |
| `GOGITOPS_SSH_USER`        | SSH connection username         | String   | -                  |
| `GOGITOPS_SSH_KEY_PATH`    | Path to private SSH key         | String   | -                  |
| `GOGITOPS_DISCORD_WEBHOOK` | Discord Webhook URL             | String   | -                  |

## Contributing

See how to contribute in [CONTRIBUTING.md](CONTRIBUTING.md).

## License

This project is licensed under the [MIT License](LICENSE).

<div align="center">

## Author

**Enoque Sousa**

[![LinkedIn](https://img.shields.io/badge/LinkedIn-0077B5?style=flat&logo=linkedin&logoColor=white)](https://www.linkedin.com/in/enoque-sousa-bb89aa168/)
[![GitHub](https://img.shields.io/badge/GitHub-100000?style=flat&logo=github&logoColor=white)](https://github.com/ESousa97)
[![Portfolio](https://img.shields.io/badge/Portfolio-FF5722?style=flat&logo=target&logoColor=white)](https://enoquesousa.vercel.app)

**[⬆ Back to top](#go-gitops-deployer)**

Made with ❤️ by [Enoque Sousa](https://github.com/ESousa97)

**Project Status:** Active Development

</div>
