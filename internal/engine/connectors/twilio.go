package connectors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// TwilioSMS handles Twilio SMS integrations
type TwilioSMS struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}

// TwilioConfig represents Twilio configuration
type TwilioConfig struct {
	To      string `json:"to"`       // Recipient phone number (e.g., "+15551234567")
	Message string `json:"message"`  // SMS message body
}

// ExecuteWithContext sends an SMS via Twilio
func (t *TwilioSMS) ExecuteWithContext(ctx context.Context, config TwilioConfig) Result {
	start := time.Now()

	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled before Twilio request: " + ctx.Err().Error())
	default:
	}

	// Validate phone number format
	if config.To == "" || config.Message == "" {
		return NewFailureResult("Twilio requires 'to' and 'message' fields", start)
	}

	// Prepare Twilio API request
	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", t.AccountSID)

	// Create form data
	formData := url.Values{}
	formData.Set("To", config.To)
	formData.Set("From", t.FromNumber)
	formData.Set("Body", config.Message)

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to create Twilio request: %v", err), start)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(t.AccountSID, t.AuthToken)

	// Execute request with timeout
	client := &http.Client{
		Timeout: 15 * time.Second, // Twilio can be slow
	}
	resp, err := client.Do(req)

	// Check if context was cancelled during request
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled during Twilio request: " + ctx.Err().Error())
	default:
	}

	if err != nil {
		return NewFailureResult(fmt.Sprintf("Twilio API request failed: %v", err), start)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode >= 400 {
		return NewFailureResult(fmt.Sprintf("Twilio returned error status: %d", resp.StatusCode), start)
	}

	// Parse response
	var twilioResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&twilioResp); err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to parse Twilio response: %v", err), start)
	}

	return NewSuccessResult("SMS sent successfully via Twilio", map[string]interface{}{
		"status_code": resp.StatusCode,
		"to":          config.To,
		"sid":         twilioResp["sid"],
		"status":      twilioResp["status"],
	}, start)
}

