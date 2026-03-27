package monitor

import (
	"context"
	"fmt"
	"time"

	"gogitopsdeployer/internal/config"
	"gogitopsdeployer/internal/gitops"
	"gogitopsdeployer/internal/ssh"
	"gogitopsdeployer/internal/storage"
)

// Monitor orquestra o loop de checagem.
type Monitor struct {
	cfg         *config.Config
	gitOps      *gitops.Service
	sshService  *ssh.Service
	storage     *storage.Service
	triggerChan chan struct{}
}

// NewMonitor cria uma nova instancia do orquestrador.
func NewMonitor(cfg *config.Config, gitOps *gitops.Service, sshService *ssh.Service, storage *storage.Service, triggerChan chan struct{}) *Monitor {
	return &Monitor{
		cfg:         cfg,
		gitOps:      gitOps,
		sshService:  sshService,
		storage:     storage,
		triggerChan: triggerChan,
	}
}

// Start inicia o loop de monitoramento.
func (m *Monitor) Start(ctx context.Context) error {
	// Garante o clone inicial
	if err := m.gitOps.EnsureClone(); err != nil {
		return fmt.Errorf("initial clone failure: %v", err)
	}

	ticker := time.NewTicker(m.cfg.Interval)
	defer ticker.Stop()

	fmt.Printf("Monitor iniciado. Polling a cada %s. Webhook ativo na porta %s.\n", m.cfg.Interval, m.cfg.WebhookPort)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Monitor finalizando...")
			return nil
		case <-m.triggerChan:
			fmt.Println("[Monitor] Trigger de Webhook recebido!")
			m.performCheck()
		case <-ticker.C:
			m.performCheck()
		}
	}
}

// performCheck executa a verificacao de mudancas e deploy.
func (m *Monitor) performCheck() {
	changed, hash, err := m.gitOps.CheckForUpdates()
	if err != nil {
		fmt.Printf("Erro ao verificar atualizacoes: %v\n", err)
		return
	}

	if changed {
		fmt.Printf("Nova versao detectada: [%s]\n", hash)
		
		// 1. Baixar as mudancas localmente
		if err := m.gitOps.UpdateLocal(); err != nil {
			fmt.Printf("Erro ao baixar mudancas: %v\n", err)
			return
		}

		// 2. Disparar comandos SSH na VPS (se configurado)
		if m.cfg.SSHHost != "" {
			output, err := m.sshService.RunCommands()
			status := "success"
			if err != nil {
				fmt.Printf("Erro ao executar comandos SSH: %v\n", err)
				status = "failed"
				if output == "" {
					output = err.Error()
				}
			}
			
			// 3. Persistir o resultado do deploy
			if err := m.storage.RecordDeploy(hash, status, output); err != nil {
				fmt.Printf("Erro ao salvar no banco: %v\n", err)
			}
		}
	} else {
		// Log discreto para estudo
		fmt.Printf("[%s] Nenhuma mudanca detectada (HEAD: %s)\n", 
			time.Now().Format("15:04:05"), hash[:8])
	}
}
