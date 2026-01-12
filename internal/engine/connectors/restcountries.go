package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// RESTCountriesConnector fetches country data from REST Countries API
// Reference: https://restcountries.com/
type RESTCountriesConnector struct {
	BaseURL string // Default: https://restcountries.com/v3.1
}

// RESTCountriesConfig represents REST Countries API connector configuration
type RESTCountriesConfig struct {
	SearchType string `json:"search_type"` // name, capital, currency, language, region, subregion
	Query      string `json:"query"`       // Search query (e.g., "united", "euro", "asia")
}

// ExecuteWithContext fetches country data from REST Countries API
func (r *RESTCountriesConnector) ExecuteWithContext(ctx context.Context, config RESTCountriesConfig) Result {
	start := time.Now()

	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled before REST Countries request: " + ctx.Err().Error())
	default:
	}

	// Set default base URL if not provided
	if r.BaseURL == "" {
		r.BaseURL = "https://restcountries.com/v3.1"
	}

	// Default values
	if config.SearchType == "" {
		config.SearchType = "all"
	}

	// Build URL
	var url string
	if config.SearchType == "all" {
		url = fmt.Sprintf("%s/all", r.BaseURL)
	} else if config.Query != "" {
		url = fmt.Sprintf("%s/%s/%s", r.BaseURL, config.SearchType, config.Query)
	} else {
		return NewFailureResult("Query is required for search type: "+config.SearchType, start)
	}

	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to create REST Countries request: %v", err), start)
	}

	// Execute request with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)

	// Check if context was cancelled during request
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled during REST Countries request: " + ctx.Err().Error())
	default:
	}

	if err != nil {
		return NewFailureResult(fmt.Sprintf("REST Countries request failed: %v", err), start)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to read REST Countries response: %v", err), start)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		return NewFailureResult(fmt.Sprintf("REST Countries returned HTTP error: %d - %s", resp.StatusCode, string(body)), start)
	}

	// Parse JSON response
	var countriesData interface{}
	if err := json.Unmarshal(body, &countriesData); err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to parse REST Countries response: %v", err), start)
	}

	// Count results
	resultCount := 0
	if countries, ok := countriesData.([]interface{}); ok {
		resultCount = len(countries)
	} else if countryMap, ok := countriesData.(map[string]interface{}); ok {
		// Single country result
		resultCount = 1
		countriesData = []interface{}{countryMap}
	}

	message := fmt.Sprintf("REST Countries data fetched: %d countries", resultCount)
	if config.Query != "" {
		message = fmt.Sprintf("REST Countries search '%s': %d countries", config.Query, resultCount)
	}

	return NewSuccessResult(message, map[string]interface{}{
		"search_type":    config.SearchType,
		"query":          config.Query,
		"country_count":  resultCount,
		"countries":      countriesData,
		"url":            url,
		"api_info":       "REST Countries API - https://restcountries.com/",
	}, start)
}

// SearchByName searches countries by name
func (r *RESTCountriesConnector) SearchByName(ctx context.Context, name string) Result {
	return r.ExecuteWithContext(ctx, RESTCountriesConfig{
		SearchType: "name",
		Query:      name,
	})
}

// SearchByCapital searches countries by capital city
func (r *RESTCountriesConnector) SearchByCapital(ctx context.Context, capital string) Result {
	return r.ExecuteWithContext(ctx, RESTCountriesConfig{
		SearchType: "capital",
		Query:      capital,
	})
}

// SearchByRegion searches countries by region
func (r *RESTCountriesConnector) SearchByRegion(ctx context.Context, region string) Result {
	return r.ExecuteWithContext(ctx, RESTCountriesConfig{
		SearchType: "region",
		Query:      region,
	})
}

// GetAllCountries fetches all countries
func (r *RESTCountriesConnector) GetAllCountries(ctx context.Context) Result {
	return r.ExecuteWithContext(ctx, RESTCountriesConfig{
		SearchType: "all",
	})
}

// DryRunRESTCountries simulates a REST Countries call without actually making the request
func (r *RESTCountriesConnector) DryRunRESTCountries(config RESTCountriesConfig) Result {
	start := time.Now()

	if r.BaseURL == "" {
		r.BaseURL = "https://restcountries.com/v3.1"
	}

	var url string
	if config.SearchType == "all" {
		url = fmt.Sprintf("%s/all", r.BaseURL)
	} else {
		url = fmt.Sprintf("%s/%s/%s", r.BaseURL, config.SearchType, config.Query)
	}

	return NewSuccessResult("REST Countries dry run completed", map[string]interface{}{
		"search_type": config.SearchType,
		"query":       config.Query,
		"url":         url,
		"api_info":    "REST Countries - https://restcountries.com/",
		"note":        "This is a dry run - no actual REST Countries call was made",
		"example_country": map[string]interface{}{
			"name":       map[string]string{"common": "United States", "official": "United States of America"},
			"capital":    []string{"Washington, D.C."},
			"region":     "Americas",
			"subregion":  "North America",
			"population": 329484123,
			"currencies": map[string]interface{}{"USD": map[string]string{"name": "United States dollar", "symbol": "$"}},
		},
	}, start)
}

