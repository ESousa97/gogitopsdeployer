package gitops

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ESousa97/gogitopsdeployer/internal/config"
)

func TestEnsureClone(t *testing.T) {
	tempDir := t.TempDir()
	cfg := &config.Config{
		RepoURL:   "https://github.com/go-git/go-git.git",
		LocalPath: filepath.Join(tempDir, "go-git"),
		Interval:  1 * time.Second,
	}

	svc := NewService(cfg)

	// Test case 1: Clone when repository does not exist
	err := svc.EnsureClone()
	if err != nil {
		t.Fatalf("EnsureClone failed on initial clone: %v", err)
	}

	// Verify the directory exists and seems to be a git repository
	if _, err := os.Stat(filepath.Join(cfg.LocalPath, ".git")); os.IsNotExist(err) {
		t.Errorf("Repository was not cloned successfully")
	}

	// Test case 2: Clone when repository already exists
	err = svc.EnsureClone()
	if err != nil {
		t.Fatalf("EnsureClone failed when repository already exists: %v", err)
	}
}

func TestCheckForUpdates(t *testing.T) {
	tempDir := t.TempDir()
	cfg := &config.Config{
		RepoURL:   "https://github.com/go-git/go-git.git",
		LocalPath: filepath.Join(tempDir, "go-git"),
		Interval:  1 * time.Second,
	}

	svc := NewService(cfg)

	// Setup: clone first
	err := svc.EnsureClone()
	if err != nil {
		t.Fatalf("Failed to clone for testing CheckForUpdates: %v", err)
	}

	// Test case 1: Check for updates right after cloning
	changed, currentHash, err := svc.CheckForUpdates()
	if err != nil {
		t.Fatalf("CheckForUpdates failed: %v", err)
	}

	// Usually, right after cloning, HEAD matches remote master/main,
	// but it depends on what the remote's HEAD is vs our local HEAD.
	// If it hasn't changed, changed should be false.
	if currentHash == "" {
		t.Errorf("Expected a non-empty hash")
	}

	// We can't deterministically test `changed == true` against a public repo
	// without mocking the repository, so we just verify it runs without error.
	_ = changed
}
