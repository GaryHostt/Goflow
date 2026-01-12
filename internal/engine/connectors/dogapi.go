package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DogAPIConnector fetches random dog images from Dog CEO API
// Reference: https://dog.ceo/dog-api/
type DogAPIConnector struct {
	BaseURL string // Default: https://dog.ceo/api
}

// DogAPIConfig represents Dog CEO API connector configuration
type DogAPIConfig struct {
	Endpoint string `json:"endpoint"` // breed, breeds/list, breeds/image/random
	Breed    string `json:"breed"`    // Specific breed (e.g., "husky", "corgi")
	SubBreed string `json:"sub_breed"` // Sub-breed (e.g., "australian" for shepherd/australian)
	Count    int    `json:"count"`    // Number of images (default: 1)
}

// ExecuteWithContext fetches dog images from Dog CEO API
func (d *DogAPIConnector) ExecuteWithContext(ctx context.Context, config DogAPIConfig) Result {
	start := time.Now()

	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled before Dog API request: " + ctx.Err().Error())
	default:
	}

	// Set default base URL if not provided
	if d.BaseURL == "" {
		d.BaseURL = "https://dog.ceo/api"
	}

	// Build URL based on configuration
	var url string
	
	if config.Endpoint == "breeds/list" || config.Endpoint == "breeds/list/all" {
		// List all breeds
		url = fmt.Sprintf("%s/breeds/list/all", d.BaseURL)
	} else if config.Breed != "" {
		if config.SubBreed != "" {
			// Specific sub-breed image
			if config.Count > 1 {
				url = fmt.Sprintf("%s/breed/%s/%s/images/random/%d", d.BaseURL, config.Breed, config.SubBreed, config.Count)
			} else {
				url = fmt.Sprintf("%s/breed/%s/%s/images/random", d.BaseURL, config.Breed, config.SubBreed)
			}
		} else {
			// Specific breed image
			if config.Count > 1 {
				url = fmt.Sprintf("%s/breed/%s/images/random/%d", d.BaseURL, config.Breed, config.Count)
			} else {
				url = fmt.Sprintf("%s/breed/%s/images/random", d.BaseURL, config.Breed)
			}
		}
	} else {
		// Random dog image
		if config.Count > 1 {
			url = fmt.Sprintf("%s/breeds/image/random/%d", d.BaseURL, config.Count)
		} else {
			url = fmt.Sprintf("%s/breeds/image/random", d.BaseURL)
		}
	}

	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to create Dog API request: %v", err), start)
	}

	// Execute request with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)

	// Check if context was cancelled during request
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled during Dog API request: " + ctx.Err().Error())
	default:
	}

	if err != nil {
		return NewFailureResult(fmt.Sprintf("Dog API request failed: %v", err), start)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to read Dog API response: %v", err), start)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		return NewFailureResult(fmt.Sprintf("Dog API returned HTTP error: %d - %s", resp.StatusCode, string(body)), start)
	}

	// Parse JSON response
	var dogData map[string]interface{}
	if err := json.Unmarshal(body, &dogData); err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to parse Dog API response: %v", err), start)
	}

	// Extract message for logging
	imageCount := 0
	if message, ok := dogData["message"].(string); ok {
		imageCount = 1
		_ = message
	} else if messages, ok := dogData["message"].([]interface{}); ok {
		imageCount = len(messages)
	}

	breedInfo := "random breed"
	if config.Breed != "" {
		breedInfo = config.Breed
		if config.SubBreed != "" {
			breedInfo = fmt.Sprintf("%s %s", config.SubBreed, config.Breed)
		}
	}

	message := fmt.Sprintf("Dog API: %d image(s) of %s", imageCount, breedInfo)

	return NewSuccessResult(message, map[string]interface{}{
		"breed":      config.Breed,
		"sub_breed":  config.SubBreed,
		"count":      imageCount,
		"data":       dogData,
		"url":        url,
		"api_info":   "Dog CEO API - The internet's biggest collection of open source dog pictures",
	}, start)
}

// GetRandomDogImage fetches a random dog image
func (d *DogAPIConnector) GetRandomDogImage(ctx context.Context) Result {
	return d.ExecuteWithContext(ctx, DogAPIConfig{})
}

// GetRandomDogImages fetches multiple random dog images
func (d *DogAPIConnector) GetRandomDogImages(ctx context.Context, count int) Result {
	return d.ExecuteWithContext(ctx, DogAPIConfig{
		Count: count,
	})
}

// GetBreedImage fetches a random image of a specific breed
func (d *DogAPIConnector) GetBreedImage(ctx context.Context, breed string) Result {
	return d.ExecuteWithContext(ctx, DogAPIConfig{
		Breed: breed,
	})
}

// GetAllBreeds fetches list of all dog breeds
func (d *DogAPIConnector) GetAllBreeds(ctx context.Context) Result {
	return d.ExecuteWithContext(ctx, DogAPIConfig{
		Endpoint: "breeds/list/all",
	})
}

// DryRunDogAPI simulates a Dog API call without actually making the request
func (d *DogAPIConnector) DryRunDogAPI(config DogAPIConfig) Result {
	start := time.Now()

	if d.BaseURL == "" {
		d.BaseURL = "https://dog.ceo/api"
	}

	url := fmt.Sprintf("%s/breeds/image/random", d.BaseURL)
	if config.Breed != "" {
		url = fmt.Sprintf("%s/breed/%s/images/random", d.BaseURL, config.Breed)
	}

	return NewSuccessResult("Dog API dry run completed", map[string]interface{}{
		"breed":    config.Breed,
		"count":    config.Count,
		"url":      url,
		"api_info": "Dog CEO API - https://dog.ceo/dog-api/",
		"note":     "This is a dry run - no actual Dog API call was made",
		"example_response": map[string]interface{}{
			"message": "https://images.dog.ceo/breeds/husky/n02110185_10175.jpg",
			"status":  "success",
		},
	}, start)
}

