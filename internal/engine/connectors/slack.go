package connectors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SlackWebhook handles Slack webhook integrations
type SlackWebhook struct {
	WebhookURL string
}

// SlackMessage represents a Slack message payload
type SlackMessage struct {
	Text string `json:"text"`
}

// Execute sends a message to Slack (legacy method - no context)
func (s *SlackWebhook) Execute(message string) Result {
	return s.ExecuteWithContext(context.Background(), message)
}

// ExecuteWithContext sends a message to Slack with context awareness
// PRODUCTION: Respects cancellation and timeouts
func (s *SlackWebhook) ExecuteWithContext(ctx context.Context, message string) Result {
	start := time.Now()

	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled before Slack request: " + ctx.Err().Error())
	default:
	}

	payload := SlackMessage{Text: message}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to marshal Slack payload: %v", err), start)
	}

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "POST", s.WebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to create Slack request: %v", err), start)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute request with context awareness
	client := &http.Client{
		Timeout: 10 * time.Second, // Maximum 10 seconds per request
	}
	resp, err := client.Do(req)

	// Check if context was cancelled during request
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled during Slack request: " + ctx.Err().Error())
	default:
	}

	if err != nil {
		return NewFailureResult(fmt.Sprintf("Slack webhook request failed: %v", err), start)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return NewFailureResult(fmt.Sprintf("Slack returned error status: %d", resp.StatusCode), start)
	}

	return NewSuccessResult("Slack message sent successfully", map[string]interface{}{
		"status_code": resp.StatusCode,
		"message":     message,
	}, start)
}
