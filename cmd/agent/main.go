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
	"gogitopsdeployer/internal/ssh"
)

func main() {
	fmt.Println("=== GOGITOPSDEPLOYER INICIANDO ===")

	// 1. Carrega Configuracoes (Config Tipada - P3 Antigravity)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configuracoes: %v", err)
	}

	// 2. Inicializa Servicos (Inversao de Dependencia - P1 Antigravity)
	gitOps := gitops.NewService(cfg)
	sshService := ssh.NewService(cfg)
	agent := monitor.NewMonitor(cfg, gitOps, sshService)

	// 3. Setup de Context para Shutdown Gracioso
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
