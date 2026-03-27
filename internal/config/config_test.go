package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Limpa variáveis de ambiente para garantir um ambiente limpo
	os.Clearenv()

	// Testa valores padrão
	os.Setenv("GOGITOPS_REPO_URL", "https://github.com/example/repo")
	os.Setenv("GOGITOPS_LOCAL_PATH", "/tmp/repo")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.RepoURL != "https://github.com/example/repo" {
		t.Errorf("Expected RepoURL 'https://github.com/example/repo', got '%s'", cfg.RepoURL)
	}

	if cfg.Interval != 30*time.Second {
		t.Errorf("Expected Interval 30s, got %v", cfg.Interval)
	}

	if cfg.LocalPath != "/tmp/repo" {
		t.Errorf("Expected LocalPath '/tmp/repo', got '%s'", cfg.LocalPath)
	}

	if cfg.DBPath != "./deployments.db" {
		t.Errorf("Expected DBPath './deployments.db', got '%s'", cfg.DBPath)
	}

	if cfg.WebhookPort != "8080" {
		t.Errorf("Expected WebhookPort '8080', got '%s'", cfg.WebhookPort)
	}

	// Testa validação com erro (REPO_URL vazio)
	os.Setenv("GOGITOPS_REPO_URL", "")
	_, err = LoadConfig()
	if err == nil {
		t.Error("Expected error for empty REPO_URL, got none")
	}

	// Testa validação SSH
	os.Setenv("GOGITOPS_REPO_URL", "https://github.com/example/repo")
	os.Setenv("GOGITOPS_SSH_HOST", "192.168.1.10")
	_, err = LoadConfig()
	if err == nil {
		t.Error("Expected error for missing SSH_USER, got none")
	}

	os.Setenv("GOGITOPS_SSH_USER", "user")
	_, err = LoadConfig()
	if err == nil {
		t.Error("Expected error for missing SSH_KEY_PATH, got none")
	}

	os.Setenv("GOGITOPS_SSH_KEY_PATH", "/path/to/key")
	_, err = LoadConfig()
	if err != nil {
		t.Errorf("Expected no error with full SSH config, got %v", err)
	}
}

func TestValidate(t *testing.T) {
	cfg := &Config{
		RepoURL:   "https://github.com/test/repo",
		Interval:  30 * time.Second,
		LocalPath: "/tmp",
	}

	if err := cfg.Validate(); err != nil {
		t.Errorf("Valid config failed validation: %v", err)
	}

	cfg.Interval = 0
	if err := cfg.Validate(); err == nil {
		t.Error("Expected error for interval < 1s")
	}

	cfg.Interval = 30 * time.Second
	cfg.LocalPath = ""
	if err := cfg.Validate(); err == nil {
		t.Error("Expected error for empty LocalPath")
	}
}
