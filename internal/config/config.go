// Package config provides typed configuration management for the GitOps agent,
// using environment variables as the source of truth.
package config

import (
	"errors"
	"os"
	"time"
)

// Config defines the configuration structure for the GitOps agent.
// It includes settings for Git, SSH, notifications, and storage.
type Config struct {
	// RepoURL is the URL of the remote Git repository to monitor.
	RepoURL string
	// Interval defines how often the agent checks for new commits.
	Interval time.Duration
	// LocalPath is the directory where the repository is cloned for analysis.
	LocalPath string

	// Database Settings

	// DBPath is the file path to the SQLite deployment history database.
	DBPath string

	// Notification Settings

	// DiscordWebhookURL is the endpoint for sending deployment alerts.
	DiscordWebhookURL string

	// Webhook Settings

	// WebhookPort is the HTTP port for the inbound webhook server.
	WebhookPort string
	// WebhookSecret is the HMAC key for validating GitHub webhook payloads.
	WebhookSecret string

	// SSH Settings

	// SSHHost is the target machine address for deployments.
	SSHHost string
	// SSHUser is the username for the SSH connection.
	SSHUser string
	// SSHKeyPath is the absolute path to the private SSH key.
	SSHKeyPath string
	// SSHCommands is a list of commands to execute on a successful commit detection.
	SSHCommands []string
	// RollbackCommand is the command executed if the primary deployment fails.
	RollbackCommand string
}

const (
	// StatusSuccess represents a successfully completed deployment.
	StatusSuccess = "success"
	// StatusFailed represents a deployment that encountered an error.
	StatusFailed = "failed"
	// StatusRollback represents a state where the system was reverted to a previous version.
	StatusRollback = "rollback"
)

// LoadConfig retrieves configuration from environment variables or applies default values.
// It returns a pointer to [Config] or an error if validation fails.
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

// Validate ensures that the [Config] structure contains all required fields
// and valid values before starting the services.
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
