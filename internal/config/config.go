package config

import (
	"errors"
	"os"
	"time"
)

// Config define a estrutura de configuracao do agente GitOps.
type Config struct {
	RepoURL   string
	Interval  time.Duration
	LocalPath string

	// Database Settings
	DBPath string

	// Notification Settings
	DiscordWebhookURL string

	// Webhook Settings
	WebhookPort   string
	WebhookSecret string

	// SSH Settings
	SSHHost         string
	SSHUser         string
	SSHKeyPath      string
	SSHCommands     []string
	RollbackCommand string
}

const (
	StatusSuccess  = "success"
	StatusFailed   = "failed"
	StatusRollback = "rollback"
)

// LoadConfig carrega as configuracoes do ambiente ou usa valores default.
func LoadConfig() (*Config, error) {
	repoURL := getEnv("GOGITOPS_REPO_URL", "https://github.com/ESousa97/gogitopsdeployer")

	intervalStr := os.Getenv("GOGITOPS_INTERVAL")
	interval := 30 * time.Second
	if intervalStr != "" {
		if parsed, err := time.ParseDuration(intervalStr); err == nil {
			interval = parsed
		}
	}

	localPath := getEnv("GOGITOPS_LOCAL_PATH", "./repo-cache")
	dbPath := getEnv("GOGITOPS_DB_PATH", "./deployments.db")

	webhookPort := getEnv("GOGITOPS_WEBHOOK_PORT", "8080")
	webhookSecret := os.Getenv("GOGITOPS_WEBHOOK_SECRET")

	discordWebhookURL := os.Getenv("GOGITOPS_DISCORD_WEBHOOK")

	sshHost := os.Getenv("GOGITOPS_SSH_HOST")
	sshUser := os.Getenv("GOGITOPS_SSH_USER")
	sshKeyPath := os.Getenv("GOGITOPS_SSH_KEY_PATH")
	sshCommandsStr := getEnv("GOGITOPS_SSH_COMMANDS", "cd /app && git pull && docker-compose up --build -d")
	rollbackCommand := getEnv("GOGITOPS_ROLLBACK_COMMAND", "cd /app && git checkout HEAD^ && docker-compose up -d")

	cfg := &Config{
		RepoURL:           repoURL,
		Interval:          interval,
		LocalPath:         localPath,
		DBPath:            dbPath,
		DiscordWebhookURL: discordWebhookURL,
		WebhookPort:       webhookPort,
		WebhookSecret:     webhookSecret,
		SSHHost:           sshHost,
		SSHUser:           sshUser,
		SSHKeyPath:        sshKeyPath,
		SSHCommands:       []string{sshCommandsStr}, // Por enquanto um comando composto
		RollbackCommand:   rollbackCommand,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Validate garante que a configuracao e valida antes da inicializacao.
func (c *Config) Validate() error {
	if c.RepoURL == "" {
		return errors.New("REPO_URL cannot be empty")
	}
	if c.Interval < 1*time.Second {
		return errors.New("INTERVAL must be at least 1 second")
	}
	if c.LocalPath == "" {
		return errors.New("LOCAL_PATH cannot be empty")
	}

	// Se SSHHost estiver definido, exige outras credenciais
	if c.SSHHost != "" {
		if c.SSHUser == "" {
			return errors.New("SSH_USER is required when SSH_HOST is provided")
		}
		if c.SSHKeyPath == "" {
			return errors.New("SSH_KEY_PATH is required when SSH_HOST is provided")
		}
	}

	return nil
}
