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

	"gogitopsdeployer/internal/config"
	"gogitopsdeployer/internal/gitops"
	"gogitopsdeployer/internal/monitor"
	"gogitopsdeployer/internal/notification"
	"gogitopsdeployer/internal/ssh"
	"gogitopsdeployer/internal/storage"
	"gogitopsdeployer/internal/webhook"
)

func main() {
	// 1. Carrega Configuracoes (Config Tipada - P3 Antigravity)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configuracoes: %v", err)
	}

	// 2. Inicializa Storage (SQLite)
	db, err := storage.NewService(cfg.DBPath)
	if err != nil {
		log.Fatalf("Erro ao inicializar banco de dados: %v", err)
	}
	defer db.Close()

	// 3. Verifica Subcomandos
	if len(os.Args) > 1 && os.Args[1] == "history" {
		showHistory(db)
		return
	}

	fmt.Println("=== GOGITOPSDEPLOYER INICIANDO ===")

	// 4. Inicializa Canais de Comunicacao
	triggerChan := make(chan struct{}, 1)

	// 5. Inicializa Servicos (Inversao de Dependencia - P1 Antigravity)
	gitOps := gitops.NewService(cfg)
	sshService := ssh.NewService(cfg)
	notifier := notification.NewService(cfg)
	agent := monitor.NewMonitor(cfg, gitOps, sshService, db, notifier, triggerChan)
	webhookSvc := webhook.NewService(cfg, triggerChan)

	// 6. Roda o Servidor Webhook em Background
	go func() {
		if err := webhookSvc.Start(); err != nil {
			log.Printf("Erro no servidor Webhook: %v\n", err)
		}
	}()

	// 7. Setup de Context para Shutdown Gracioso
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Captura sinais de sistema
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nSinal de interrupcao recebido. Finalizando agente...")
		cancel()
	}()

	// 4. Inicia o Monitor
	fmt.Printf("Monitorando repositorio: %s\n", cfg.RepoURL)
	if err := agent.Start(ctx); err != nil {
		log.Fatalf("Erro fatal no monitor: %v", err)
	}

	fmt.Println("=== GOGITOPSDEPLOYER FINALIZADO ===")
}

// showHistory queries the [storage.Service] for recent deployment events
// and prints them in a formatted table to the standard output.
func showHistory(db *storage.Service) {
	fmt.Println("=== HISTORICO DE DEPLOYS (Ultimos 10) ===")
	history, err := db.GetHistory(10)
	if err != nil {
		log.Fatalf("Erro ao buscar historico: %v", err)
	}

	if len(history) == 0 {
		fmt.Println("Nenhum deploy registrado ainda.")
		return
	}

	fmt.Printf("%-5s | %-8s | %-10s | %-20s\n", "ID", "HASH", "STATUS", "DATA")
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
