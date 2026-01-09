package connectors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// DiscordWebhook sends a message to Discord using a webhook
type DiscordWebhook struct {
	WebhookURL string
}

// Execute sends a message to Discord
func (d *DiscordWebhook) Execute(message string) Result {
	if d.WebhookURL == "" {
		return Result{
			Status:  "failed",
			Message: "Discord webhook URL is not configured",
		}
	}

	payload := map[string]string{
		"content": message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return Result{
			Status:  "failed",
			Message: fmt.Sprintf("Failed to marshal JSON: %v", err),
		}
	}

	resp, err := http.Post(d.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return Result{
			Status:  "failed",
			Message: fmt.Sprintf("Failed to send request to Discord: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return Result{
			Status:  "failed",
			Message: fmt.Sprintf("Discord returned error status: %d", resp.StatusCode),
		}
	}

	return Result{
		Status:  "success",
		Message: fmt.Sprintf("Message sent to Discord successfully"),
		Data: map[string]interface{}{
			"status_code": resp.StatusCode,
		},
	}
}

