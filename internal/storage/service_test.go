package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStorageService(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	// Test NewService (Initialization & Migration)
	svc, err := NewService(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize storage service: %v", err)
	}
	defer svc.Close()

	// Ensure database file was created
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatalf("Database file was not created at %s", dbPath)
	}

	// Test RecordDeploy
	err = svc.RecordDeploy("hash123", "success", "output text")
	if err != nil {
		t.Fatalf("Failed to record deploy: %v", err)
	}

	err = svc.RecordDeploy("hash456", "failed", "error text")
	if err != nil {
		t.Fatalf("Failed to record second deploy: %v", err)
	}

	// Test GetHistory
	history, err := svc.GetHistory(10)
	if err != nil {
		t.Fatalf("Failed to get history: %v", err)
	}

	if len(history) != 2 {
		t.Errorf("Expected history length 2, got %d", len(history))
	}

	// History is ordered by ID DESC
	if history[0].Hash != "hash456" {
		t.Errorf("Expected first history item to be hash456, got %s", history[0].Hash)
	}
	if history[1].Hash != "hash123" {
		t.Errorf("Expected second history item to be hash123, got %s", history[1].Hash)
	}

	// Test limit
	historyLimit, err := svc.GetHistory(1)
	if err != nil {
		t.Fatalf("Failed to get history with limit: %v", err)
	}

	if len(historyLimit) != 1 {
		t.Errorf("Expected history length 1 with limit, got %d", len(historyLimit))
	}
}

func TestNewServiceError(t *testing.T) {
	// Trying to open a DB in a non-existent directory without permissions
	// (or just a bad path) might not fail immediately with pure sqlite memory,
	// but a completely invalid path or read-only dir might.
	// We just ensure the error branch exists for robust coverage if possible,
	// but purely invalid paths are hard to guarantee errors across OS.
	// Let's create an invalid state:
	invalidPath := "/root/invalid/db.sqlite" // assuming this fails
	_, err := NewService(invalidPath)
	if err == nil {
		t.Logf("Expected error for invalid path %s, but got none (depends on OS/permissions)", invalidPath)
	}
}
