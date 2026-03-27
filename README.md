# gogitopsdeployer

Agente GitOps escrito em Go para monitoramento de repositorios e deteccao automatica de novos commits.

## Objetivos
- **Monitoramento Continuo**: Verifica alteracoes no branch principal (`master`/`main`) a cada 30 segundos.
- **Deteccao por Hash**: Compara o `HEAD` local com o remoto via biblioteca `go-git`.
- **Arquitetura Modular**: Segue os principios de Domain-Driven Design (DDD) e Inversao de Dependencia.

## Estrutura do Projeto
- `cmd/agent`: Ponto de entrada (Main).
- `internal/config`: Configuracoes tipadas via variaveis de ambiente.
- `internal/gitops`: Servico de abstracao para operacoes de rede e sistema de arquivos Git.
- `internal/monitor`: Orquestrador do loop de execucao.

## Como Executar
```powershell
# Instale as dependencias
go mod tidy

# Rodar o agente
go run cmd/agent/main.go
```

## Variaveis de Ambiente
| Variavel | Descricao | Default |
|----------|-----------|---------|
| `GOGITOPS_REPO_URL` | URL do repositorio Git | `https://github.com/go-git/go-git` |
| `GOGITOPS_INTERVAL` | Intervalo de checagem | `30s` |
| `GOGITOPS_LOCAL_PATH` | Diretorio de cache local | `./repo-cache` |

---
Projeto de estudo desenvolvido com os padroes **Antigravity**.
