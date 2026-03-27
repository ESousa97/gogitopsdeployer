// Package main is the entry point for the gogitopsdeployer agent.
// It initializes all internal services, sets up the orchestrator,
// and manages the application's lifecycle, including graceful shutdowns.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ESousa97/gogitopsdeployer/internal/config"
	"github.com/ESousa97/gogitopsdeployer/internal/gitops"
	"github.com/ESousa97/gogitopsdeployer/internal/monitor"
	"github.com/ESousa97/gogitopsdeployer/internal/notification"
	"github.com/ESousa97/gogitopsdeployer/internal/ssh"
	"github.com/ESousa97/gogitopsdeployer/internal/storage"
	"github.com/ESousa97/gogitopsdeployer/internal/webhook"
)

func main() {
	// 1. Load Configurations (Typed Config - Antigravity P3)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configurations: %v", err)
	}

	// 2. Initialize Storage (SQLite)
	db, err := storage.NewService(cfg.DBPath)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	// 3. Check Subcommands
	if len(os.Args) > 1 && os.Args[1] == "history" {
		showHistory(db)
		return
	}

	fmt.Println("=== GOGITOPSDEPLOYER STARTING ===")

	// 4. Initialize Communication Channels
	triggerChan := make(chan struct{}, 1)

	// 5. Initialize Services (Dependency Inversion - Antigravity P1)
	gitOps := gitops.NewService(cfg)
	sshService := ssh.NewService(cfg)
	notifier := notification.NewService(cfg)
	agent := monitor.NewMonitor(cfg, gitOps, sshService, db, notifier, triggerChan)
	webhookSvc := webhook.NewService(cfg, triggerChan)

	// 6. Run Webhook Server in Background
	go func() {
		if err := webhookSvc.Start(); err != nil {
			log.Printf("Webhook server error: %v\n", err)
		}
	}()

	// 7. Context Setup for Graceful Shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Capture system signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nInterrupt signal received. Shutting down agent...")
		cancel()
	}()

	// 8. Start the Monitor
	fmt.Printf("Monitoring repository: %s\n", cfg.RepoURL)
	if err := agent.Start(ctx); err != nil {
		log.Fatalf("Fatal error in monitor: %v", err)
	}

	fmt.Println("=== GOGITOPSDEPLOYER FINISHED ===")
}

// showHistory queries the [storage.Service] for recent deployment events
// and prints them in a formatted table to the standard output.
func showHistory(db *storage.Service) {
	fmt.Println("=== DEPLOYMENT HISTORY (Latest 10) ===")
	history, err := db.GetHistory(10)
	if err != nil {
		log.Fatalf("Error fetching history: %v", err)
	}

	if len(history) == 0 {
		fmt.Println("No deployments recorded yet.")
		return
	}

	fmt.Printf("%-5s | %-8s | %-10s | %-20s\n", "ID", "HASH", "STATUS", "DATE")
	fmt.Println("------------------------------------------------------------")
	for _, d := range history {
		hashShort := d.Hash
		if len(hashShort) > 8 {
			hashShort = hashShort[:8]
		}
		fmt.Printf("%-5d | %-8s | %-10s | %-20s\n",
			d.ID, hashShort, d.Status, d.CreatedAt.Format("2006-01-02 15:04:05"))
	}
}
