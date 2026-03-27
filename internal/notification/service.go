package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gogitopsdeployer/internal/config"
)

// Service gerencia o envio de notificacoes.
type Service struct {
	cfg *config.Config
}

// NewService cria uma nova instancia do servico de notificacao.
func NewService(cfg *config.Config) *Service {
	return &Service{cfg: cfg}
}

// DiscordPayload representa a estrutura da mensagem para o Discord.
type DiscordPayload struct {
	Embeds []DiscordEmbed `json:"embeds"`
}

// DiscordEmbed representa um card rico no Discord.
type DiscordEmbed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int    `json:"color"`
	Footer      struct {
		Text string `json:"text"`
	} `json:"footer"`
}

// Notify envia uma notificacao para o Discord informando o status do deploy.
func (s *Service) Notify(status, message, hash string) error {
	if s.cfg.DiscordWebhookURL == "" {
		return nil // Discord nao configurado
	}

	color := 0x2ecc71 // Verde (Success)
	if status == config.StatusFailed {
		color = 0xe74c3c // Vermelho (Failed)
	} else if status == config.StatusRollback {
		color = 0xf1c40f // Amarelo (Rollback)
	}

	payload := DiscordPayload{
		Embeds: []DiscordEmbed{
			{
				Title:       fmt.Sprintf("Deploy Status: %s", status),
				Description: fmt.Sprintf("**Hash:** %s\n**Log:**\n```\n%s\n```", hash, message),
				Color:       color,
				Footer: struct {
					Text string `json:"text"`
				}{
					Text: fmt.Sprintf("Agente: gogitopsdeployer | %s", time.Now().Format("15:04:05")),
				},
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(s.cfg.DiscordWebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord returned status: %d", resp.StatusCode)
	}

	return nil
}
