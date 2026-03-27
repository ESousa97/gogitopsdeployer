package monitor

import (
	"context"
	"fmt"
	"time"

	"gogitopsdeployer/internal/config"
	"gogitopsdeployer/internal/gitops"
)

// Monitor orquestra o loop de checagem.
type Monitor struct {
	cfg    *config.Config
	gitOps *gitops.Service
}

// NewMonitor cria uma nova instancia do orquestrador.
func NewMonitor(cfg *config.Config, gitOps *gitops.Service) *Monitor {
	return &Monitor{
		cfg:    cfg,
		gitOps: gitOps,
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

	fmt.Printf("Monitor iniciado. Verificando a cada %s...\n", m.cfg.Interval)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Monitor finalizando...")
			return nil
		case <-ticker.C:
			changed, hash, err := m.gitOps.CheckForUpdates()
			if err != nil {
				fmt.Printf("Erro ao verificar atualizacoes: %v\n", err)
				continue
			}

			if changed {
				fmt.Printf("Nova versao detectada: [%s]\n", hash)
				
				// Opcional: Baixar as mudancas de fato
				if err := m.gitOps.UpdateLocal(); err != nil {
					fmt.Printf("Erro ao baixar mudancas: %v\n", err)
				}
			} else {
				// Log discreto para estudo
				fmt.Printf("[%s] Nenhuma mudanca detectada (HEAD: %s)\n", 
					time.Now().Format("15:04:05"), hash[:8])
			}
		}
	}
}
