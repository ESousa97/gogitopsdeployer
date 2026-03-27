// Package monitor acts as the core orchestrator for the GitOps agent.
// It coordinates polling intervals, webhook triggers, repository updates,
// SSH deployments, and notifications.
package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/esousa97/gogitopsdeployer/internal/config"
	"github.com/esousa97/gogitopsdeployer/internal/gitops"
	"github.com/esousa97/gogitopsdeployer/internal/notification"
	"github.com/esousa97/gogitopsdeployer/internal/ssh"
	"github.com/esousa97/gogitopsdeployer/internal/storage"
)

// Monitor coordinates the reconciliation loop between the Git repository
// and the target deployment environment.
type Monitor struct {
	cfg          *config.Config
	gitOps       *gitops.Service
	sshService   *ssh.Service
	storage      *storage.Service
	notification *notification.Service
	triggerChan  chan struct{}
}

// NewMonitor initializes the [Monitor] with all necessary dependencies
// and the communication channel for webhook triggers.
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

// Start begins the main monitoring loop. It performs an initial clone
// and then listens for timer ticks or webhook triggers to execute reconciliation.
// It respects the provided [context.Context] for graceful shutdown.
func (m *Monitor) Start(ctx context.Context) error {
	// Ensure initial clone
	if err := m.gitOps.EnsureClone(); err != nil {
		return fmt.Errorf("initial clone failure: %v", err)
	}

	ticker := time.NewTicker(m.cfg.Interval)
	defer ticker.Stop()

	fmt.Printf("Monitor started. Polling every %s. Webhook active on port %s.\n", m.cfg.Interval, m.cfg.WebhookPort)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Monitor shutting down...")
			return nil
		case <-m.triggerChan:
			fmt.Println("[Monitor] Webhook trigger received!")
			m.performCheck()
		case <-ticker.C:
			m.performCheck()
		}
	}
}

// performCheck executes the core logic: check for updates, pull changes,
// run SSH commands, and handle rollbacks if necessary.
func (m *Monitor) performCheck() {
	changed, hash, err := m.gitOps.CheckForUpdates()
	if err != nil {
		fmt.Printf("Error checking for updates: %v\n", err)
		return
	}

	if changed {
		fmt.Printf("New version detected: [%s]\n", hash)

		// 1. Pull changes locally
		if err := m.gitOps.UpdateLocal(); err != nil {
			fmt.Printf("Error pulling changes: %v\n", err)
			return
		}

		// 2. Trigger SSH commands on VPS (if configured)
		if m.cfg.SSHHost != "" {
			output, err := m.sshService.RunCommands()
			if err != nil {
				fmt.Printf("[Monitor] Deploy FAILED: %v\n", err)

				// Record failure
				m.storage.RecordDeploy(hash, config.StatusFailed, output)

				// Notify Discord (Failure)
				m.notification.Notify(config.StatusFailed, "Deploy failed. Initiating auto-rollback...", hash)

				// 3. AUTO-ROLLBACK
				fmt.Println("[Monitor] Executing emergency rollback...")
				rbOutput, rbErr := m.sshService.RunRollback()
				if rbErr != nil {
					fmt.Printf("[Monitor] ROLLBACK FAILED: %v\n", rbErr)
					m.notification.Notify(config.StatusFailed, fmt.Sprintf("CRITICAL: Rollback also failed!\n%s", rbOutput), hash)
				} else {
					fmt.Println("[Monitor] Rollback executed successfully.")
					m.storage.RecordDeploy(hash, config.StatusRollback, rbOutput)
					m.notification.Notify(config.StatusRollback, "System restored to previous stable version.", hash)
				}
				return
			}

			// 4. Success
			fmt.Println("[Monitor] Deployment finished successfully.")
			m.storage.RecordDeploy(hash, config.StatusSuccess, output)
			m.notification.Notify(config.StatusSuccess, "Deploy successful.", hash)
		}
	} else {
		// Discreet log for analysis
		fmt.Printf("[%s] No changes detected (HEAD: %s)\n",
			time.Now().Format("15:04:05"), hash[:8])
	}
}
