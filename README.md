# gogitopsdeployer

Agente GitOps escrito em Go para monitoramento de repositorios e deteccao automatica de novos commits.

## Objetivos
- **Monitoramento Continuo**: Verifica alteracoes no branch principal (`master`/`main`) a cada 30 segundos.
- **Deteccao por Hash**: Compara o `HEAD` local com o remoto via biblioteca `go-git`.
- **Deploy Remoto (SSH)**: Executa comandos em uma VPS via SSH ao detectar novos commits.
- **Auto-Rollback**: Se o deploy falhar na VPS, o agente executa automaticamente `git checkout HEAD^` para restaurar a versao estavel.
- **Notificacoes (Discord)**: Envia cards ricos (Embeds) para um canal do Discord informando Sucesso, Falha ou Rollback.
- **Persistencia (SQLite)**: Registra cada tentativa de deploy no banco de dados local.
- **CLI History**: Subcomando para visualizar o historico de deploys.
- **Webhook Support**: Disparo imediato via GitHub Webhooks com validacao HMAC.
- **Arquitetura Modular**: Segue os principios de Domain-Driven Design (DDD) e Inversao de Dependencia.

## Estrutura do Projeto
- `cmd/agent`: Ponto de entrada (Main) e subcomandos CLI.
- `internal/config`: Configuracoes tipadas via variaveis de ambiente.
- `internal/gitops`: Servico de abstracao para operacoes Git.
- `internal/ssh`: Servico de execucao de comandos remotos e rollback.
- `internal/notification`: Integracao com Webhooks do Discord.
- `internal/storage`: Persistencia de dados usando SQLite.
- `internal/monitor`: Orquestrador resiliente do loop de execucao.

## Como Executar
```powershell
# Instale as dependencias
go mod tidy

# Configure as variaveis (opcional)
$env:GOGITOPS_DISCORD_WEBHOOK="https://discord.com/api/webhooks/..."
$env:GOGITOPS_SSH_HOST="1.2.3.4"

# Rodar o agente
go run cmd/agent/main.go
```

## Variaveis de Ambiente
| Variavel | Descricao | Default |
|----------|-----------|---------|
| `GOGITOPS_REPO_URL` | URL do repositorio Git | `https://github.com/ESousa97/gogitopsdeployer` |
| `GOGITOPS_INTERVAL` | Intervalo de checagem | `30s` |
| `GOGITOPS_DB_PATH` | Caminho do banco SQLite | `./deployments.db` |
| `GOGITOPS_DISCORD_WEBHOOK` | URL do Webhook do Discord | (Opcional) |
| `GOGITOPS_ROLLBACK_COMMAND`| Comando de emergencia | `git checkout HEAD^ && docker-compose up -d` |
| `GOGITOPS_WEBHOOK_PORT` | Porta do servidor HTTP | `8080` |
| `GOGITOPS_SSH_USER` | Usuario SSH | (Opcional) |
| `GOGITOPS_SSH_KEY_PATH` | Caminho para Chave Privada | (Opcional) |
| `GOGITOPS_SSH_COMMANDS` | Comandos para deploy | `cd /app && git pull && docker-compose up --build -d` |

---
Projeto de estudo desenvolvido com os padroes **Antigravity**.
