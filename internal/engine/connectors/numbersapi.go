package connectors

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// NumbersAPIConnector fetches interesting facts about numbers
// Reference: http://numbersapi.com/
type NumbersAPIConnector struct {
	BaseURL string // Default: http://numbersapi.com
}

// NumbersAPIConfig represents Numbers API connector configuration
type NumbersAPIConfig struct {
	Number string `json:"number"` // The number to get facts about (e.g., "42", "random")
	Type   string `json:"type"`   // trivia, math, date, year
}

// ExecuteWithContext fetches a number fact from Numbers API
func (n *NumbersAPIConnector) ExecuteWithContext(ctx context.Context, config NumbersAPIConfig) Result {
	start := time.Now()

	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled before Numbers API request: " + ctx.Err().Error())
	default:
	}

	// Set default base URL if not provided
	if n.BaseURL == "" {
		n.BaseURL = "http://numbersapi.com"
	}

	// Default values
	if config.Number == "" {
		config.Number = "random"
	}
	if config.Type == "" {
		config.Type = "trivia"
	}

	// Build URL
	url := fmt.Sprintf("%s/%s/%s?json", n.BaseURL, config.Number, config.Type)

	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to create Numbers API request: %v", err), start)
	}

	// Execute request with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)

	// Check if context was cancelled during request
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled during Numbers API request: " + ctx.Err().Error())
	default:
	}

	if err != nil {
		return NewFailureResult(fmt.Sprintf("Numbers API request failed: %v", err), start)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to read Numbers API response: %v", err), start)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		return NewFailureResult(fmt.Sprintf("Numbers API returned HTTP error: %d - %s", resp.StatusCode, string(body)), start)
	}

	// Numbers API returns plain text or JSON
	// For JSON format, we requested ?json parameter
	message := string(body)
	if len(message) > 100 {
		message = message[:100] + "..."
	}

	return NewSuccessResult(fmt.Sprintf("Numbers API fact: %s", message), map[string]interface{}{
		"number":   config.Number,
		"type":     config.Type,
		"fact":     string(body),
		"url":      url,
		"api_info": "Numbers API - An API for interesting facts about numbers",
	}, start)
}

// GetTriviaFact fetches a trivia fact about a number
func (n *NumbersAPIConnector) GetTriviaFact(ctx context.Context, number string) Result {
	return n.ExecuteWithContext(ctx, NumbersAPIConfig{
		Number: number,
		Type:   "trivia",
	})
}

// GetMathFact fetches a math fact about a number
func (n *NumbersAPIConnector) GetMathFact(ctx context.Context, number string) Result {
	return n.ExecuteWithContext(ctx, NumbersAPIConfig{
		Number: number,
		Type:   "math",
	})
}

// GetDateFact fetches a fact about a date (format: month/day)
func (n *NumbersAPIConnector) GetDateFact(ctx context.Context, monthDay string) Result {
	return n.ExecuteWithContext(ctx, NumbersAPIConfig{
		Number: monthDay,
		Type:   "date",
	})
}

// GetYearFact fetches a fact about a year
func (n *NumbersAPIConnector) GetYearFact(ctx context.Context, year string) Result {
	return n.ExecuteWithContext(ctx, NumbersAPIConfig{
		Number: year,
		Type:   "year",
	})
}

// DryRunNumbersAPI simulates a Numbers API call without actually making the request
func (n *NumbersAPIConnector) DryRunNumbersAPI(config NumbersAPIConfig) Result {
	start := time.Now()

	if n.BaseURL == "" {
		n.BaseURL = "http://numbersapi.com"
	}

	if config.Number == "" {
		config.Number = "random"
	}
	if config.Type == "" {
		config.Type = "trivia"
	}

	url := fmt.Sprintf("%s/%s/%s", n.BaseURL, config.Number, config.Type)

	return NewSuccessResult("Numbers API dry run completed", map[string]interface{}{
		"number":   config.Number,
		"type":     config.Type,
		"url":      url,
		"api_info": "Numbers API - http://numbersapi.com/",
		"note":     "This is a dry run - no actual Numbers API call was made",
		"example_fact": "42 is the answer to the Ultimate Question of Life, the Universe, and Everything.",
	}, start)
}

