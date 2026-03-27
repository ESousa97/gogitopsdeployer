package monitor

import (
	"context"
	"fmt"
	"time"

	"gogitopsdeployer/internal/config"
	"gogitopsdeployer/internal/gitops"
	"gogitopsdeployer/internal/notification"
	"gogitopsdeployer/internal/ssh"
	"gogitopsdeployer/internal/storage"
)

// Monitor orquestra o loop de checagem.
type Monitor struct {
	cfg          *config.Config
	gitOps       *gitops.Service
	sshService   *ssh.Service
	storage      *storage.Service
	notification *notification.Service
	triggerChan  chan struct{}
}

// NewMonitor cria uma nova instancia do orquestrador.
func NewMonitor(cfg *config.Config, gitOps *gitops.Service, sshService *ssh.Service, storage *storage.Service, notification *notification.Service, triggerChan chan struct{}) *Monitor {
	return &Monitor{
		cfg:          cfg,
		gitOps:       gitOps,
		sshService:   sshService,
		storage:      storage,
		notification: notification,
		triggerChan:  triggerChan,
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
			if err != nil {
				fmt.Printf("[Monitor] Deploy FAILED: %v\n", err)
				
				// Persiste Falha
				m.storage.RecordDeploy(hash, config.StatusFailed, output)
				
				// Notifica Discord (Falha)
				m.notification.Notify(config.StatusFailed, "Deploy failed. Initiating auto-rollback...", hash)

				// 3. AUTO-ROLLBACK
				fmt.Println("[Monitor] Executando Rollback de emergencia...")
				rbOutput, rbErr := m.sshService.RunRollback()
				if rbErr != nil {
					fmt.Printf("[Monitor] ROLLBACK FAILED: %v\n", rbErr)
					m.notification.Notify(config.StatusFailed, fmt.Sprintf("CRITICAL: Rollback also failed!\n%s", rbOutput), hash)
				} else {
					fmt.Println("[Monitor] Rollback executado com sucesso.")
					m.storage.RecordDeploy(hash, config.StatusRollback, rbOutput)
					m.notification.Notify(config.StatusRollback, "System restored to previous stable version.", hash)
				}
				return
			}
			
			// 4. Sucesso
			fmt.Println("[Monitor] Deploy finalizado com sucesso.")
			m.storage.RecordDeploy(hash, config.StatusSuccess, output)
			m.notification.Notify(config.StatusSuccess, "Deploy successful.", hash)
		}
	} else {
		// Log discreto para estudo
		fmt.Printf("[%s] Nenhuma mudanca detectada (HEAD: %s)\n", 
			time.Now().Format("15:04:05"), hash[:8])
	}
}
