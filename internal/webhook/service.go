// Package webhook provides an HTTP server to receive and validate
// push event notifications from GitHub.
package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/esousa97/gogitopsdeployer/internal/config"
)

// Service implements an HTTP handler that listens for inbound repository
// notifications, validates their signatures, and triggers the monitor.
type Service struct {
	cfg         *config.Config
	triggerChan chan struct{}
}

// NewService initializes a new [Service] with the required [config.Config]
// and a shared channel to signal updates to the orchestrator.
func NewService(cfg *config.Config, triggerChan chan struct{}) *Service {
	return &Service{
		cfg:         cfg,
		triggerChan: triggerChan,
	}
}

// Start launches the HTTP server and blocks until it is shut down.
// It exposes the /webhook endpoint on the configured [config.WebhookPort].
func (s *Service) Start() error {
	http.HandleFunc("/webhook", s.handleWebhook)
	fmt.Printf("[Webhook] Server listening on port %s...\n", s.cfg.WebhookPort)
	return http.ListenAndServe(":"+s.cfg.WebhookPort, nil)
}

// handleWebhook processes the HTTP POST request, validates the HMAC-SHA256
// signature (if a secret is configured), checks for push/ping events,
// and sends a non-blocking signal to the monitor channel.
func (s *Service) handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// HMAC Signature Validation (if configured)
	if s.cfg.WebhookSecret != "" {
		signature := r.Header.Get("X-Hub-Signature-256")
		if !s.validateSignature(payload, signature) {
			fmt.Println("[Webhook] Invalid signature received.")
			http.Error(w, "Invalid signature", http.StatusForbidden)
			return
		}
	}

	// Check if it is a push event
	event := r.Header.Get("X-GitHub-Event")
	if event != "push" && event != "ping" {
		fmt.Printf("[Webhook] Event ignored: %s\n", event)
		w.WriteHeader(http.StatusOK)
		return
	}

	if event == "ping" {
		fmt.Println("[Webhook] Ping received from GitHub.")
		w.WriteHeader(http.StatusOK)
		return
	}

	fmt.Println("[Webhook] Push detected! Triggering immediate deployment...")

	// Send signal to monitor (non-blocking)
	select {
	case s.triggerChan <- struct{}{}:
	default:
		fmt.Println("[Webhook] Deployment already queued or in progress.")
	}

	w.WriteHeader(http.StatusAccepted)
}

// validateSignature performs a cryptographic comparison between the
// payload's HMAC and the signature provided in the header.
func (s *Service) validateSignature(payload []byte, signature string) bool {
	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}

	actualSig, err := hex.DecodeString(signature[7:])
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, []byte(s.cfg.WebhookSecret))
	mac.Write(payload)
	expectedSig := mac.Sum(nil)

	return hmac.Equal(actualSig, expectedSig)
}
