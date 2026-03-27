package notification

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/esousa97/gogitopsdeployer/internal/config"
)

func TestNotify(t *testing.T) {
	// Mock HTTP Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %v", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %v", r.Header.Get("Content-Type"))
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read body: %v", err)
		}

		var payload DiscordPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("Failed to unmarshal payload: %v", err)
		}

		if len(payload.Embeds) != 1 {
			t.Errorf("Expected 1 embed, got %d", len(payload.Embeds))
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.Config{
		DiscordWebhookURL: server.URL,
	}

	svc := NewService(cfg)

	// Test 1: Success Status
	err := svc.Notify(config.StatusSuccess, "Deploy completed", "hash123")
	if err != nil {
		t.Fatalf("Failed to notify success: %v", err)
	}

	// Test 2: Failed Status
	err = svc.Notify(config.StatusFailed, "Deploy failed", "hash456")
	if err != nil {
		t.Fatalf("Failed to notify failure: %v", err)
	}

	// Test 3: Rollback Status
	err = svc.Notify(config.StatusRollback, "Rollback executed", "hash789")
	if err != nil {
		t.Fatalf("Failed to notify rollback: %v", err)
	}

	// Test 4: Empty URL (should do nothing)
	emptySvc := NewService(&config.Config{DiscordWebhookURL: ""})
	err = emptySvc.Notify(config.StatusSuccess, "msg", "hash")
	if err != nil {
		t.Fatalf("Expected no error for empty webhook URL, got %v", err)
	}

	// Test 5: Server error response
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer errorServer.Close()

	errorSvc := NewService(&config.Config{DiscordWebhookURL: errorServer.URL})
	err = errorSvc.Notify(config.StatusSuccess, "msg", "hash")
	if err == nil {
		t.Fatalf("Expected error for server error response, got none")
	}
}
