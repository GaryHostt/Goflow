package connectors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// DiscordWebhook handles Discord webhook integrations
type DiscordWebhook struct {
	WebhookURL string
}

// DiscordMessage represents a Discord message payload
type DiscordMessage struct {
	Content string `json:"content"`
}

// Execute sends a message to Discord
func (d *DiscordWebhook) Execute(message string) Result {
	return d.ExecuteWithContext(context.Background(), message)
}

// ExecuteWithContext sends a message to Discord with context awareness
func (d *DiscordWebhook) ExecuteWithContext(ctx context.Context, message string) Result {
	start := time.Now()

	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled before Discord request: " + ctx.Err().Error())
	default:
	}

	payload := DiscordMessage{Content: message}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to marshal Discord payload: %v", err), start)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", d.WebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to create Discord request: %v", err), start)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)

	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled during Discord request: " + ctx.Err().Error())
	default:
	}

	if err != nil {
		return NewFailureResult(fmt.Sprintf("Discord webhook request failed: %v", err), start)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return NewFailureResult(fmt.Sprintf("Discord returned error status: %d", resp.StatusCode), start)
	}

	return NewSuccessResult("Discord message sent successfully", map[string]interface{}{
		"status_code": resp.StatusCode,
		"message":     message,
	}, start)
}
