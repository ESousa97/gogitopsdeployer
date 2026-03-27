package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"

	"gogitopsdeployer/internal/config"
)

// Service gerencia o recebimento de webhooks do GitHub.
type Service struct {
	cfg         *config.Config
	triggerChan chan struct{}
}

// NewService cria uma nova instancia do servico de Webhook.
func NewService(cfg *config.Config, triggerChan chan struct{}) *Service {
	return &Service{
		cfg:         cfg,
		triggerChan: triggerChan,
	}
}

// Start inicia o servidor HTTP para escutar webhooks.
func (s *Service) Start() error {
	http.HandleFunc("/webhook", s.handleWebhook)
	fmt.Printf("[Webhook] Servidor escutando na porta %s...\n", s.cfg.WebhookPort)
	return http.ListenAndServe(":"+s.cfg.WebhookPort, nil)
}

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

	// Validacao de Assinatura HMAC (se configurado)
	if s.cfg.WebhookSecret != "" {
		signature := r.Header.Get("X-Hub-Signature-256")
		if !s.validateSignature(payload, signature) {
			fmt.Println("[Webhook] Assinatura invalida recebida.")
			http.Error(w, "Invalid signature", http.StatusForbidden)
			return
		}
	}

	// Verifica se e um evento de push
	event := r.Header.Get("X-GitHub-Event")
	if event != "push" && event != "ping" {
		fmt.Printf("[Webhook] Evento ignorado: %s\n", event)
		w.WriteHeader(http.StatusOK)
		return
	}

	if event == "ping" {
		fmt.Println("[Webhook] Ping recebido do GitHub.")
		w.WriteHeader(http.StatusOK)
		return
	}

	fmt.Println("[Webhook] Push detectado! Disparando deploy imediato...")
	
	// Envia sinal para o monitor (sem bloquear)
	select {
	case s.triggerChan <- struct{}{}:
	default:
		fmt.Println("[Webhook] Deploy ja em fila ou em processamento.")
	}

	w.WriteHeader(http.StatusAccepted)
}

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
