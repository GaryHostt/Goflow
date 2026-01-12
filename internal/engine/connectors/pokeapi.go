package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// PokeAPIConnector fetches Pokemon data from PokeAPI
// Reference: https://pokeapi.co/
type PokeAPIConnector struct {
	BaseURL string // Default: https://pokeapi.co/api/v2
}

// PokeAPIConfig represents PokeAPI connector configuration
type PokeAPIConfig struct {
	Resource string `json:"resource"` // pokemon, berry, item, move, ability, type, etc.
	ID       string `json:"id"`       // Pokemon ID or name (e.g., "1", "bulbasaur")
}

// ExecuteWithContext fetches Pokemon data from PokeAPI
func (p *PokeAPIConnector) ExecuteWithContext(ctx context.Context, config PokeAPIConfig) Result {
	start := time.Now()

	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled before PokeAPI request: " + ctx.Err().Error())
	default:
	}

	// Set default base URL if not provided
	if p.BaseURL == "" {
		p.BaseURL = "https://pokeapi.co/api/v2"
	}

	// Default to pokemon resource if not specified
	if config.Resource == "" {
		config.Resource = "pokemon"
	}

	// Validate ID
	if config.ID == "" {
		return NewFailureResult("Pokemon ID or name is required", start)
	}

	// Build URL
	url := fmt.Sprintf("%s/%s/%s", p.BaseURL, config.Resource, config.ID)

	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to create PokeAPI request: %v", err), start)
	}

	// Execute request with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)

	// Check if context was cancelled during request
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled during PokeAPI request: " + ctx.Err().Error())
	default:
	}

	if err != nil {
		return NewFailureResult(fmt.Sprintf("PokeAPI request failed: %v", err), start)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to read PokeAPI response: %v", err), start)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		return NewFailureResult(fmt.Sprintf("PokeAPI returned HTTP error: %d - %s", resp.StatusCode, string(body)), start)
	}

	// Parse JSON response
	var pokeData map[string]interface{}
	if err := json.Unmarshal(body, &pokeData); err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to parse PokeAPI response: %v", err), start)
	}

	// Extract name for logging
	resourceName := config.ID
	if name, ok := pokeData["name"].(string); ok {
		resourceName = name
	}

	message := fmt.Sprintf("PokeAPI %s fetched: %s", config.Resource, resourceName)

	return NewSuccessResult(message, map[string]interface{}{
		"resource":  config.Resource,
		"id":        config.ID,
		"data":      pokeData,
		"url":       url,
		"api_info":  "PokeAPI - The RESTful Pokemon API",
	}, start)
}

// GetPokemon fetches a specific Pokemon by ID or name
func (p *PokeAPIConnector) GetPokemon(ctx context.Context, idOrName string) Result {
	return p.ExecuteWithContext(ctx, PokeAPIConfig{
		Resource: "pokemon",
		ID:       idOrName,
	})
}

// GetBerry fetches a specific berry by ID or name
func (p *PokeAPIConnector) GetBerry(ctx context.Context, idOrName string) Result {
	return p.ExecuteWithContext(ctx, PokeAPIConfig{
		Resource: "berry",
		ID:       idOrName,
	})
}

// GetMove fetches a specific move by ID or name
func (p *PokeAPIConnector) GetMove(ctx context.Context, idOrName string) Result {
	return p.ExecuteWithContext(ctx, PokeAPIConfig{
		Resource: "move",
		ID:       idOrName,
	})
}

// DryRunPokeAPI simulates a PokeAPI call without actually making the request
func (p *PokeAPIConnector) DryRunPokeAPI(config PokeAPIConfig) Result {
	start := time.Now()

	if p.BaseURL == "" {
		p.BaseURL = "https://pokeapi.co/api/v2"
	}

	if config.Resource == "" {
		config.Resource = "pokemon"
	}

	url := fmt.Sprintf("%s/%s/%s", p.BaseURL, config.Resource, config.ID)

	return NewSuccessResult("PokeAPI dry run completed", map[string]interface{}{
		"resource": config.Resource,
		"id":       config.ID,
		"url":      url,
		"api_info": "PokeAPI - https://pokeapi.co/",
		"note":     "This is a dry run - no actual PokeAPI call was made",
		"example_pokemon": map[string]interface{}{
			"name":   "bulbasaur",
			"id":     1,
			"height": 7,
			"weight": 69,
			"types":  []string{"grass", "poison"},
		},
	}, start)
}

