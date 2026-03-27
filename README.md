# gogitopsdeployer

> Agente GitOps em Go para monitoramento resiliente de repositórios e automação de deploys.

![CI](https://github.com/ESousa97/gogitopsdeployer/actions/workflows/ci.yml/badge.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/ESousa97/gogitopsdeployer)
![Go Reference](https://pkg.go.dev/badge/github.com/ESousa97/gogitopsdeployer.svg)
![License](https://img.shields.io/github/license/ESousa97/gogitopsdeployer)
![Go Version](https://img.shields.io/github/go-mod/go-version/ESousa97/gogitopsdeployer)
![Last Commit](https://img.shields.io/github/last-commit/ESousa97/gogitopsdeployer)

---

O **gogitopsdeployer** é um agente leve e modular projetado para automatizar o ciclo de entrega contínua (CD) em ambientes distribuídos. Ele monitora alterações em repositórios Git, dispara execuções remotas via SSH e garante a estabilidade do sistema com mecanismos de rollback automático e notificações em tempo real.

## Demonstração

### CLI History
```bash
gogitopsdeployer history
```

### Discord Notification
O agente envia cartões detalhados para o Discord com o status de cada operação:
- **Sucesso**: Deploy concluído na infraestrutura alvo.
- **Falha**: Erro detectado durante a execução remota.
- **Rollback**: Restauração automática para a última versão estável.

## Stack Tecnológico

| Tecnologia | Papel |
|---|---|
| Go 1.25 | Linguagem principal e runtime |
| go-git | Manipulação nativa de repositórios Git |
| SQLite | Persistência de histórico de deploys |
| SSH (crypto/ssh) | Execução de comandos remotos e gestão de estado |
| Discord Webhooks | Sistema de notificações e alertas |

## Pré-requisitos

- Go >= 1.25.0
- Acesso SSH configurado na máquina alvo (chave privada)
- Banco de dados SQLite (criado automaticamente)

## Instalação e Uso

### Como binário

```bash
go install github.com/ESousa97/gogitopsdeployer/cmd/agent@latest
```

### A partir do source

```bash
git clone https://github.com/ESousa97/gogitopsdeployer.git
cd gogitopsdeployer
cp .env.example .env
# Edite .env com suas configurações
make build
./bin/gogitopsdeployer
```

## Makefile Targets

| Target | Descrição |
|---|---|
| `make build` | Compila o binário do agente em `bin/` |
| `make run` | Executa o agente diretamente via Go |
| `make test` | Executa a suíte de testes unitários |
| `make lint` | Executa análise estática via `go vet` |
| `make tidy` | Limpa e atualiza dependências do `go.mod` |
| `make clean` | Remove arquivos binários e temporários |

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

| Variável | Descrição | Tipo | Padrão |
|---|---|---|---|
| `GOGITOPS_REPO_URL` | URL do repositório Git alvo | String | Repositório atual |
| `GOGITOPS_INTERVAL` | Intervalo entre checagens | Duration | `30s` |
| `GOGITOPS_DB_PATH` | Caminho para o banco SQLite | String | `./deployments.db` |
| `GOGITOPS_SSH_HOST` | Host ou IP da máquina de deploy | String | - |
| `GOGITOPS_SSH_USER` | Usuário para conexão SSH | String | - |
| `GOGITOPS_SSH_KEY_PATH` | Caminho para chave privada SSH | String | - |
| `GOGITOPS_DISCORD_WEBHOOK` | URL do Webhook do Discord | String | - |

## Roadmap

- [x] **Fase 1**: Arquitetura modular e serviços internos.
- [x] **Fase 2**: Monitoramento Git e detecção de commits por Hash.
- [x] **Fase 3**: Execução remota via SSH e Rollback automático.
- [/] **Fase 4**: Notificações ricas e suporte a Webhooks HTTP (HMAC).
- [ ] **Fase 5**: Multi-ambiente e autenticação avançada (OIDC).

## Contribuindo

Veja como contribuir em [CONTRIBUTING.md](CONTRIBUTING.md).

## Licença

Este projeto está licenciado sob a [MIT License](LICENSE).

## Autor

**Enoque Sousa** — [Portfólio](https://enoquesousa.vercel.app) | [GitHub](https://github.com/ESousa97)
