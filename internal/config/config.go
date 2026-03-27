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
}

// LoadConfig carrega as configuracoes do ambiente ou usa valores default.
func LoadConfig() (*Config, error) {
	repoURL := os.Getenv("GOGITOPS_REPO_URL")
	if repoURL == "" {
		// Repositorio default para estudo (go-git proprio repositorio)
		repoURL = "https://github.com/go-git/go-git"
	}

	intervalStr := os.Getenv("GOGITOPS_INTERVAL")
	interval := 30 * time.Second
	if intervalStr != "" {
		parsed, err := time.ParseDuration(intervalStr)
		if err == nil {
			interval = parsed
		}
	}

	localPath := os.Getenv("GOGITOPS_LOCAL_PATH")
	if localPath == "" {
		localPath = "./repo-cache"
	}

	cfg := &Config{
		RepoURL:   repoURL,
		Interval:  interval,
		LocalPath: localPath,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
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
	return nil
}
