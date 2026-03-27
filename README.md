# gogitopsdeployer

> A lightweight, modular GitOps agent for automated repository monitoring and remote deployment via SSH.

![CI](https://github.com/ESousa97/gogitopsdeployer/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/ESousa97/gogitopsdeployer)](https://goreportcard.com/report/github.com/ESousa97/gogitopsdeployer)
[![Go Reference](https://pkg.go.dev/badge/github.com/ESousa97/gogitopsdeployer.svg)](https://pkg.go.dev/github.com/ESousa97/gogitopsdeployer)
![License](https://img.shields.io/github/license/ESousa97/gogitopsdeployer)
![Go Version](https://img.shields.io/github/go-mod/go-version/ESousa97/gogitopsdeployer)
![Last Commit](https://img.shields.io/github/last-commit/ESousa97/gogitopsdeployer)

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
go install github.com/ESousa97/gogitopsdeployer/cmd/agent@latest
```

# From source

```bash
git clone https://github.com/ESousa97/gogitopsdeployer.git
cd gogitopsdeployer
make build
```

## Makefile Targets

| Target       | Descrição                                 |
| ------------ | ----------------------------------------- |
| `make build` | Compila o binário do agente em `bin/`     |
| `make run`   | Executa o agente diretamente via Go       |
| `make test`  | Executa a suíte de testes unitários       |
| `make lint`  | Executa análise estática via `go vet`     |
| `make tidy`  | Limpa e atualiza dependências do `go.mod` |
| `make clean` | Remove arquivos binários e temporários    |

## Arquitetura

O projeto segue princípios de **Inversão de Dependência** e **Arquitetura Modular**, garantindo que o núcleo de negócio seja agnóstico a infraestrutura.

- `cmd/agent`: Ponto de entrada e comandos CLI.
- `internal/gitops`: Abstração de operações Git e detecção de commits.
- `internal/ssh`: Motor de execução remota e lógica de rollback.
- `internal/monitor`: Orquestrador resiliente do loop de reconciliação.
- `internal/storage`: Camada de persistência SQLite.

## API Reference

A documentação detalhada de pacotes e funções está disponível via Godoc:
"Veja a documentação completa em [pkg.go.dev](https://pkg.go.dev/github.com/ESousa97/gogitopsdeployer)."

## Configuração

| Variável                   | Descrição                       | Tipo     | Padrão             |
| -------------------------- | ------------------------------- | -------- | ------------------ |
| `GOGITOPS_REPO_URL`        | URL do repositório Git alvo     | String   | Repositório atual  |
| `GOGITOPS_INTERVAL`        | Intervalo entre checagens       | Duration | `30s`              |
| `GOGITOPS_DB_PATH`         | Caminho para o banco SQLite     | String   | `./deployments.db` |
| `GOGITOPS_SSH_HOST`        | Host ou IP da máquina de deploy | String   | -                  |
| `GOGITOPS_SSH_USER`        | Usuário para conexão SSH        | String   | -                  |
| `GOGITOPS_SSH_KEY_PATH`    | Caminho para chave privada SSH  | String   | -                  |
| `GOGITOPS_DISCORD_WEBHOOK` | URL do Webhook do Discord       | String   | -                  |

## Contribuindo

Veja como contribuir em [CONTRIBUTING.md](CONTRIBUTING.md).

## Licença

Este projeto está licenciado sob a [MIT License](LICENSE).

## Autor

**Enoque Sousa** — [Portfólio](https://enoquesousa.vercel.app) | [GitHub](https://github.com/ESousa97)
