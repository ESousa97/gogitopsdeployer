// Package notification provides integration with external communication
// platforms like Discord to broadcast deployment status.
package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gogitopsdeployer/internal/config"
)

// Service handles the dispatching of notifications to external webhooks.
type Service struct {
	cfg *config.Config
}

// NewService creates a new [Service] with the provided [config.Config].
func NewService(cfg *config.Config) *Service {
	return &Service{cfg: cfg}
}

// DiscordPayload defines the top-level structure for a Discord webhook request,
// which essentially consists of a list of embeds.
type DiscordPayload struct {
	// Embeds is the list of rich cards to display in the Discord message.
	Embeds []DiscordEmbed `json:"embeds"`
}

// DiscordEmbed represents a single rich-content card in a Discord message.
type DiscordEmbed struct {
	// Title is the main heading of the embed card.
	Title string `json:"title"`
	// Description is the body text of the embed, supporting basic markdown.
	Description string `json:"description"`
	// Color is the hexadecimal integer color code for the left side of the embed.
	Color int `json:"color"`
	// Footer contains metadata like the timestamp or application name.
	Footer struct {
		Text string `json:"text"`
	} `json:"footer"`
}

// Notify sends a formatted alert to Discord based on the deployment status.
// It maps [config.StatusSuccess], [config.StatusFailed], and [config.StatusRollback]
// to specific colors and creates an informative embed message.
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
