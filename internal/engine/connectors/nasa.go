package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// NASAAPIConnector fetches data from NASA's APIs
// Reference: https://api.nasa.gov/
type NASAAPIConnector struct {
	BaseURL string // Default: https://api.nasa.gov
	APIKey  string // NASA API key (use "DEMO_KEY" for testing)
}

// NASAAPIConfig represents NASA API connector configuration
type NASAAPIConfig struct {
	Endpoint string `json:"endpoint"` // apod (Astronomy Picture of the Day), mars-photos, neo (Near Earth Objects), etc.
	Date     string `json:"date"`     // Optional: YYYY-MM-DD format
	Count    int    `json:"count"`    // Optional: number of results
}

// ExecuteWithContext fetches data from NASA API
func (n *NASAAPIConnector) ExecuteWithContext(ctx context.Context, config NASAAPIConfig) Result {
	start := time.Now()

	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled before NASA API request: " + ctx.Err().Error())
	default:
	}

	// Set defaults
	if n.BaseURL == "" {
		n.BaseURL = "https://api.nasa.gov"
	}
	if n.APIKey == "" {
		n.APIKey = "DEMO_KEY" // NASA provides a demo key for testing
	}
	if config.Endpoint == "" {
		config.Endpoint = "planetary/apod" // Astronomy Picture of the Day
	}

	// Build URL with query parameters
	url := fmt.Sprintf("%s/%s?api_key=%s", n.BaseURL, config.Endpoint, n.APIKey)

	if config.Date != "" {
		url += fmt.Sprintf("&date=%s", config.Date)
	}

	if config.Count > 0 {
		url += fmt.Sprintf("&count=%d", config.Count)
	}

	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to create NASA API request: %v", err), start)
	}

	// Execute request with timeout
	client := &http.Client{
		Timeout: 15 * time.Second, // NASA API can be slower
	}
	resp, err := client.Do(req)

	// Check if context was cancelled during request
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled during NASA API request: " + ctx.Err().Error())
	default:
	}

	if err != nil {
		return NewFailureResult(fmt.Sprintf("NASA API request failed: %v", err), start)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to read NASA API response: %v", err), start)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		return NewFailureResult(fmt.Sprintf("NASA API returned HTTP error: %d - %s", resp.StatusCode, string(body)), start)
	}

	// Parse JSON response
	var nasaData interface{}
	if err := json.Unmarshal(body, &nasaData); err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to parse NASA API response: %v", err), start)
	}

	// Extract title for APOD if available
	resourceName := config.Endpoint
	if dataMap, ok := nasaData.(map[string]interface{}); ok {
		if title, exists := dataMap["title"]; exists {
			resourceName = fmt.Sprintf("%v", title)
		}
	}

	message := fmt.Sprintf("NASA API data fetched: %s", resourceName)

	return NewSuccessResult(message, map[string]interface{}{
		"endpoint": config.Endpoint,
		"data":     nasaData,
		"url":      url,
		"api_info": "NASA API - https://api.nasa.gov/",
	}, start)
}

// GetAPOD fetches Astronomy Picture of the Day
func (n *NASAAPIConnector) GetAPOD(ctx context.Context, date string) Result {
	return n.ExecuteWithContext(ctx, NASAAPIConfig{
		Endpoint: "planetary/apod",
		Date:     date,
	})
}

// GetRandomAPOD fetches random APODs
func (n *NASAAPIConnector) GetRandomAPOD(ctx context.Context, count int) Result {
	return n.ExecuteWithContext(ctx, NASAAPIConfig{
		Endpoint: "planetary/apod",
		Count:    count,
	})
}

// DryRunNASAAPI simulates a NASA API call without actually making the request
func (n *NASAAPIConnector) DryRunNASAAPI(config NASAAPIConfig) Result {
	start := time.Now()

	if n.BaseURL == "" {
		n.BaseURL = "https://api.nasa.gov"
	}
	if config.Endpoint == "" {
		config.Endpoint = "planetary/apod"
	}

	url := fmt.Sprintf("%s/%s?api_key=DEMO_KEY", n.BaseURL, config.Endpoint)

	return NewSuccessResult("NASA API dry run completed", map[string]interface{}{
		"endpoint": config.Endpoint,
		"url":      url,
		"api_info": "NASA API - https://api.nasa.gov/",
		"note":     "This is a dry run - no actual NASA API call was made",
		"example_apod": map[string]interface{}{
			"title":       "The Eagle Nebula from Kitt Peak",
			"explanation": "The Eagle Nebula is a star-forming region...",
			"url":         "https://apod.nasa.gov/apod/image/...",
			"media_type":  "image",
			"date":        "2026-01-12",
		},
	}, start)
}

