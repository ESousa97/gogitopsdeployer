# gogitopsdeployer

Agente GitOps escrito em Go para monitoramento de repositorios e deteccao automatica de novos commits.

## Objetivos
- **Monitoramento Continuo**: Verifica alteracoes no branch principal (`master`/`main`) a cada 30 segundos.
- **Deteccao por Hash**: Compara o `HEAD` local com o remoto via biblioteca `go-git`.
- **Deploy Remoto (SSH)**: Executa comandos em uma VPS via SSH ao detectar novos commits.
- **Arquitetura Modular**: Segue os principios de Domain-Driven Design (DDD) e Inversao de Dependencia.

## Estrutura do Projeto
- `cmd/agent`: Ponto de entrada (Main).
- `internal/config`: Configuracoes tipadas via variaveis de ambiente.
- `internal/gitops`: Servico de abstracao para operacoes Git.
- `internal/ssh`: Servico de execucao de comandos remotos.
- `internal/monitor`: Orquestrador do loop de execucao.

## Como Executar
```powershell
# Instale as dependencias
go mod tidy

# Configure as variaveis (opcional para SSH)
$env:GOGITOPS_SSH_HOST="1.2.3.4"
$env:GOGITOPS_SSH_USER="root"
$env:GOGITOPS_SSH_KEY_PATH="C:\Users\user\.ssh\id_rsa"

# Rodar o agente
go run cmd/agent/main.go
```

## Variaveis de Ambiente
| Variavel | Descricao | Default |
|----------|-----------|---------|
| `GOGITOPS_REPO_URL` | URL do repositorio Git | `https://github.com/ESousa97/gogitopsdeployer` |
| `GOGITOPS_INTERVAL` | Intervalo de checagem | `30s` |
| `GOGITOPS_LOCAL_PATH` | Diretorio de cache local | `./repo-cache` |
| `GOGITOPS_SSH_HOST` | IP/Host da VPS | (Opcional) |
| `GOGITOPS_SSH_USER` | Usuario SSH | (Opcional) |
| `GOGITOPS_SSH_KEY_PATH` | Caminho para Chave Privada | (Opcional) |
| `GOGITOPS_SSH_COMMANDS` | Comandos para deploy | `cd /app && git pull && docker-compose up --build -d` |

---
Projeto de estudo desenvolvido com os padroes **Antigravity**.
