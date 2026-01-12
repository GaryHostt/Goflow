package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// SWAPIConnector fetches Star Wars data from swapi.info
// Reference: https://swapi.info/
type SWAPIConnector struct {
	BaseURL string // Default: https://swapi.info/api
}

// SWAPIConfig represents SWAPI connector configuration
type SWAPIConfig struct {
	Resource string `json:"resource"` // films, people, planets, species, vehicles, starships
	ID       string `json:"id"`       // Resource ID (e.g., "1" for first film)
	Search   string `json:"search"`   // Search query
}

// SWAPIResponse represents a single SWAPI resource
type SWAPIResponse struct {
	Name    string                 `json:"name,omitempty"`
	Title   string                 `json:"title,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
	RawData interface{}            `json:"raw_data,omitempty"`
}

// ExecuteWithContext fetches Star Wars data from SWAPI
func (s *SWAPIConnector) ExecuteWithContext(ctx context.Context, config SWAPIConfig) Result {
	start := time.Now()

	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled before SWAPI request: " + ctx.Err().Error())
	default:
	}

	// Set default base URL if not provided
	if s.BaseURL == "" {
		s.BaseURL = "https://swapi.info/api"
	}

	// Validate resource type
	validResources := map[string]bool{
		"films":     true,
		"people":    true,
		"planets":   true,
		"species":   true,
		"vehicles":  true,
		"starships": true,
	}

	if !validResources[config.Resource] {
		return NewFailureResult(
			fmt.Sprintf("Invalid SWAPI resource: %s. Valid: films, people, planets, species, vehicles, starships", config.Resource),
			start,
		)
	}

	// Build URL
	var url string
	if config.ID != "" {
		// Fetch specific resource by ID
		url = fmt.Sprintf("%s/%s/%s", s.BaseURL, config.Resource, config.ID)
	} else if config.Search != "" {
		// Search resources
		url = fmt.Sprintf("%s/%s?search=%s", s.BaseURL, config.Resource, config.Search)
	} else {
		// List all resources
		url = fmt.Sprintf("%s/%s", s.BaseURL, config.Resource)
	}

	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to create SWAPI request: %v", err), start)
	}

	// Execute request with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)

	// Check if context was cancelled during request
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled during SWAPI request: " + ctx.Err().Error())
	default:
	}

	if err != nil {
		return NewFailureResult(fmt.Sprintf("SWAPI request failed: %v", err), start)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to read SWAPI response: %v", err), start)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		return NewFailureResult(fmt.Sprintf("SWAPI returned HTTP error: %d - %s", resp.StatusCode, string(body)), start)
	}

	// Parse JSON response
	var swapiData interface{}
	if err := json.Unmarshal(body, &swapiData); err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to parse SWAPI response: %v", err), start)
	}

	// Extract name/title for logging
	var resourceName string
	if dataMap, ok := swapiData.(map[string]interface{}); ok {
		if name, exists := dataMap["name"]; exists {
			resourceName = fmt.Sprintf("%v", name)
		} else if title, exists := dataMap["title"]; exists {
			resourceName = fmt.Sprintf("%v", title)
		} else if results, exists := dataMap["results"]; exists {
			// It's a list response
			if resultsList, ok := results.([]interface{}); ok {
				resourceName = fmt.Sprintf("%d results", len(resultsList))
			}
		}
	}

	message := fmt.Sprintf("SWAPI data fetched successfully: %s", resourceName)
	if config.ID != "" {
		message = fmt.Sprintf("SWAPI %s #%s: %s", config.Resource, config.ID, resourceName)
	} else if config.Search != "" {
		message = fmt.Sprintf("SWAPI search for '%s': %s", config.Search, resourceName)
	}

	return NewSuccessResult(message, map[string]interface{}{
		"resource":  config.Resource,
		"id":        config.ID,
		"search":    config.Search,
		"data":      swapiData,
		"url":       url,
		"api_info":  "Star Wars API - https://swapi.info/",
		"cache_hit": resp.Header.Get("X-Cache") == "HIT",
	}, start)
}

// GetFilm fetches a specific Star Wars film by ID
func (s *SWAPIConnector) GetFilm(ctx context.Context, filmID string) Result {
	return s.ExecuteWithContext(ctx, SWAPIConfig{
		Resource: "films",
		ID:       filmID,
	})
}

// GetCharacter fetches a specific Star Wars character by ID
func (s *SWAPIConnector) GetCharacter(ctx context.Context, characterID string) Result {
	return s.ExecuteWithContext(ctx, SWAPIConfig{
		Resource: "people",
		ID:       characterID,
	})
}

// GetPlanet fetches a specific Star Wars planet by ID
func (s *SWAPIConnector) GetPlanet(ctx context.Context, planetID string) Result {
	return s.ExecuteWithContext(ctx, SWAPIConfig{
		Resource: "planets",
		ID:       planetID,
	})
}

// SearchCharacters searches for Star Wars characters by name
func (s *SWAPIConnector) SearchCharacters(ctx context.Context, query string) Result {
	return s.ExecuteWithContext(ctx, SWAPIConfig{
		Resource: "people",
		Search:   query,
	})
}

// DryRunSWAPI simulates a SWAPI call without actually making the request
func (s *SWAPIConnector) DryRunSWAPI(config SWAPIConfig) Result {
	start := time.Now()

	if s.BaseURL == "" {
		s.BaseURL = "https://swapi.info/api"
	}

	var url string
	if config.ID != "" {
		url = fmt.Sprintf("%s/%s/%s", s.BaseURL, config.Resource, config.ID)
	} else if config.Search != "" {
		url = fmt.Sprintf("%s/%s?search=%s", s.BaseURL, config.Resource, config.Search)
	} else {
		url = fmt.Sprintf("%s/%s", s.BaseURL, config.Resource)
	}

	return NewSuccessResult("SWAPI dry run completed", map[string]interface{}{
		"resource": config.Resource,
		"id":       config.ID,
		"search":   config.Search,
		"url":      url,
		"api_info": "Star Wars API - https://swapi.info/",
		"note":     "This is a dry run - no actual SWAPI call was made",
		"example_data": map[string]string{
			"films":     "A New Hope, The Empire Strikes Back, Return of the Jedi",
			"people":    "Luke Skywalker, Darth Vader, Princess Leia",
			"planets":   "Tatooine, Alderaan, Hoth",
			"species":   "Human, Droid, Wookiee",
			"vehicles":  "Sand Crawler, AT-AT, Speeder Bike",
			"starships": "X-wing, TIE Fighter, Millennium Falcon",
		},
	}, start)
}

