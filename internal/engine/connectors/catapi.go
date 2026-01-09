package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// CatAPI handles The Cat API integrations
// API Documentation: https://thecatapi.com/
type CatAPI struct {
	APIKey string // Optional for basic usage
}

// CatConfig represents Cat API query configuration
type CatConfig struct {
	Limit      int    `json:"limit"`       // Number of cats (default: 1)
	HasBreeds  bool   `json:"has_breeds"`  // Filter to only cats with breed info
	BreedID    string `json:"breed_id"`    // Specific breed (e.g., "beng" for Bengal)
	Category   string `json:"category"`    // Category ID (e.g., "boxes", "hats")
}

// CatImage represents a cat image from The Cat API
type CatImage struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Breeds []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Temperament string `json:"temperament"`
		Origin      string `json:"origin"`
		Description string `json:"description"`
	} `json:"breeds"`
}

// ExecuteWithContext fetches cat images from The Cat API
func (c *CatAPI) ExecuteWithContext(ctx context.Context, config CatConfig) Result {
	start := time.Now()

	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled before Cat API request: " + ctx.Err().Error())
	default:
	}

	// Default values
	if config.Limit == 0 {
		config.Limit = 1
	}
	if config.Limit > 10 {
		config.Limit = 10 // Reasonable limit
	}

	// Build API URL
	apiURL := fmt.Sprintf("https://api.thecatapi.com/v1/images/search?limit=%d", config.Limit)
	
	if config.HasBreeds {
		apiURL += "&has_breeds=1"
	}
	if config.BreedID != "" {
		apiURL += "&breed_ids=" + config.BreedID
	}
	if config.Category != "" {
		apiURL += "&category_ids=" + config.Category
	}

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to create Cat API request: %v", err), start)
	}

	// Add API key if provided
	if c.APIKey != "" {
		req.Header.Set("x-api-key", c.APIKey)
	}

	// Execute request with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)

	// Check if context was cancelled during request
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled during Cat API request: " + ctx.Err().Error())
	default:
	}

	if err != nil {
		return NewFailureResult(fmt.Sprintf("Cat API request failed: %v", err), start)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode >= 400 {
		return NewFailureResult(fmt.Sprintf("Cat API returned error status: %d", resp.StatusCode), start)
	}

	// Parse response
	var cats []CatImage
	if err := json.NewDecoder(resp.Body).Decode(&cats); err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to parse Cat API response: %v", err), start)
	}

	return NewSuccessResult("Cat images fetched successfully", map[string]interface{}{
		"cats":  cats,
		"count": len(cats),
	}, start)
}

