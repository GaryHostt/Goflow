package connectors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Result represents the result of a connector action
type Result struct {
	Status  string                 `json:"status"`  // "success" or "failed"
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// SlackWebhook sends a message to Slack using an incoming webhook
type SlackWebhook struct {
	WebhookURL string
}

// Execute sends a message to Slack
func (s *SlackWebhook) Execute(message string) Result {
	if s.WebhookURL == "" {
		return Result{
			Status:  "failed",
			Message: "Slack webhook URL is not configured",
		}
	}

	payload := map[string]string{
		"text": message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return Result{
			Status:  "failed",
			Message: fmt.Sprintf("Failed to marshal JSON: %v", err),
		}
	}

	resp, err := http.Post(s.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return Result{
			Status:  "failed",
			Message: fmt.Sprintf("Failed to send request to Slack: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return Result{
			Status:  "failed",
			Message: fmt.Sprintf("Slack returned error status: %d", resp.StatusCode),
		}
	}

	return Result{
		Status:  "success",
		Message: fmt.Sprintf("Message sent to Slack successfully"),
		Data: map[string]interface{}{
			"status_code": resp.StatusCode,
		},
	}
}

