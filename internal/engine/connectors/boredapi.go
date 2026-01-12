package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// BoredAPIConnector fetches random activity suggestions from Bored API
// Reference: http://www.boredapi.com/
type BoredAPIConnector struct {
	BaseURL string // Default: https://www.boredapi.com/api/activity
}

// BoredAPIConfig represents Bored API connector configuration
type BoredAPIConfig struct {
	Type         string  `json:"type"`         // education, recreational, social, diy, charity, cooking, relaxation, music, busywork
	Participants int     `json:"participants"` // Number of participants
	MinPrice     float64 `json:"min_price"`    // Minimum price (0.0 to 1.0)
	MaxPrice     float64 `json:"max_price"`    // Maximum price (0.0 to 1.0)
}

// ExecuteWithContext fetches a random activity suggestion from Bored API
func (b *BoredAPIConnector) ExecuteWithContext(ctx context.Context, config BoredAPIConfig) Result {
	start := time.Now()

	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled before Bored API request: " + ctx.Err().Error())
	default:
	}

	// Set default base URL if not provided
	if b.BaseURL == "" {
		b.BaseURL = "https://www.boredapi.com/api/activity"
	}

	// Build URL with query parameters
	url := b.BaseURL
	queryParams := ""

	if config.Type != "" {
		if queryParams == "" {
			queryParams = "?"
		}
		queryParams += fmt.Sprintf("type=%s", config.Type)
	}

	if config.Participants > 0 {
		if queryParams == "" {
			queryParams = "?"
		} else {
			queryParams += "&"
		}
		queryParams += fmt.Sprintf("participants=%d", config.Participants)
	}

	if config.MinPrice > 0 {
		if queryParams == "" {
			queryParams = "?"
		} else {
			queryParams += "&"
		}
		queryParams += fmt.Sprintf("minprice=%.1f", config.MinPrice)
	}

	if config.MaxPrice > 0 {
		if queryParams == "" {
			queryParams = "?"
		} else {
			queryParams += "&"
		}
		queryParams += fmt.Sprintf("maxprice=%.1f", config.MaxPrice)
	}

	url += queryParams

	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to create Bored API request: %v", err), start)
	}

	// Execute request with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)

	// Check if context was cancelled during request
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled during Bored API request: " + ctx.Err().Error())
	default:
	}

	if err != nil {
		return NewFailureResult(fmt.Sprintf("Bored API request failed: %v", err), start)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to read Bored API response: %v", err), start)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		return NewFailureResult(fmt.Sprintf("Bored API returned HTTP error: %d - %s", resp.StatusCode, string(body)), start)
	}

	// Parse JSON response
	var activityData map[string]interface{}
	if err := json.Unmarshal(body, &activityData); err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to parse Bored API response: %v", err), start)
	}

	// Extract activity for logging
	activity := "Random activity"
	if activityStr, ok := activityData["activity"].(string); ok {
		activity = activityStr
	}

	message := fmt.Sprintf("Bored API activity: %s", activity)

	return NewSuccessResult(message, map[string]interface{}{
		"activity":  activityData,
		"url":       url,
		"api_info":  "Bored API - Find something to do!",
	}, start)
}

// GetRandomActivity fetches a completely random activity
func (b *BoredAPIConnector) GetRandomActivity(ctx context.Context) Result {
	return b.ExecuteWithContext(ctx, BoredAPIConfig{})
}

// GetActivityByType fetches an activity of a specific type
func (b *BoredAPIConnector) GetActivityByType(ctx context.Context, activityType string) Result {
	return b.ExecuteWithContext(ctx, BoredAPIConfig{
		Type: activityType,
	})
}

// DryRunBoredAPI simulates a Bored API call without actually making the request
func (b *BoredAPIConnector) DryRunBoredAPI(config BoredAPIConfig) Result {
	start := time.Now()

	return NewSuccessResult("Bored API dry run completed", map[string]interface{}{
		"type":         config.Type,
		"participants": config.Participants,
		"api_info":     "Bored API - http://www.boredapi.com/",
		"note":         "This is a dry run - no actual Bored API call was made",
		"example_activity": map[string]interface{}{
			"activity":     "Learn Express.js",
			"type":         "education",
			"participants": 1,
			"price":        0.1,
			"link":         "https://expressjs.com/",
			"key":          "3943506",
			"accessibility": 0.1,
		},
	}, start)
}

