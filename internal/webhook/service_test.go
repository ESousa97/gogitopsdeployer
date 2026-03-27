package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ESousa97/gogitopsdeployer/internal/config"
)

func TestHandleWebhook(t *testing.T) {
	cfg := &config.Config{
		WebhookSecret: "secret",
		WebhookPort:   "8080",
	}

	triggerChan := make(chan struct{}, 1)
	svc := NewService(cfg, triggerChan)

	// Test 1: GET method not allowed
	req := httptest.NewRequest(http.MethodGet, "/webhook", nil)
	w := httptest.NewRecorder()
	svc.handleWebhook(w, req)

	if w.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected StatusMethodNotAllowed, got %v", w.Result().StatusCode)
	}

	// Test 2: Invalid signature
	payload := []byte(`{"ref":"refs/heads/master"}`)
	req = httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(payload))
	req.Header.Set("X-Hub-Signature-256", "sha256=invalid")
	w = httptest.NewRecorder()
	svc.handleWebhook(w, req)

	if w.Result().StatusCode != http.StatusForbidden {
		t.Errorf("Expected StatusForbidden, got %v", w.Result().StatusCode)
	}

	// Test 3: Valid signature, ping event
	mac := hmac.New(sha256.New, []byte(cfg.WebhookSecret))
	mac.Write(payload)
	signature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	req = httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(payload))
	req.Header.Set("X-Hub-Signature-256", signature)
	req.Header.Set("X-GitHub-Event", "ping")
	w = httptest.NewRecorder()
	svc.handleWebhook(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("Expected StatusOK for ping, got %v", w.Result().StatusCode)
	}

	// Test 4: Valid signature, push event
	req = httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(payload))
	req.Header.Set("X-Hub-Signature-256", signature)
	req.Header.Set("X-GitHub-Event", "push")
	w = httptest.NewRecorder()
	svc.handleWebhook(w, req)

	if w.Result().StatusCode != http.StatusAccepted {
		t.Errorf("Expected StatusAccepted for push, got %v", w.Result().StatusCode)
	}

	// Check if trigger was sent
	select {
	case <-triggerChan:
		// Success
	default:
		t.Error("Trigger channel should have received a message")
	}

	// Test 5: Ignored event
	req = httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(payload))
	req.Header.Set("X-Hub-Signature-256", signature)
	req.Header.Set("X-GitHub-Event", "pull_request")
	w = httptest.NewRecorder()
	svc.handleWebhook(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("Expected StatusOK for ignored event, got %v", w.Result().StatusCode)
	}
}

func TestValidateSignature(t *testing.T) {
	svc := &Service{
		cfg: &config.Config{
			WebhookSecret: "mysecret",
		},
	}

	payload := []byte("test payload")
	mac := hmac.New(sha256.New, []byte("mysecret"))
	mac.Write(payload)
	validSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	if !svc.validateSignature(payload, validSignature) {
		t.Errorf("Expected valid signature to pass validation")
	}

	if svc.validateSignature(payload, "invalid") {
		t.Errorf("Expected invalid signature to fail validation")
	}

	if svc.validateSignature(payload, "sha256=invalid") {
		t.Errorf("Expected invalid hex signature to fail validation")
	}
}
