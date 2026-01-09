package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// FakeStoreAPI handles Fake Store API integrations
// API Documentation: https://fakestoreapi.com/docs
type FakeStoreAPI struct{}

// FakeStoreConfig represents Fake Store API query configuration
type FakeStoreConfig struct {
	Endpoint string `json:"endpoint"` // "products", "users", "carts", "categories"
	Limit    int    `json:"limit"`    // Number of items (default: 10)
	Category string `json:"category"` // For products: "electronics", "jewelery", "men's clothing", "women's clothing"
}

// Product represents a product from Fake Store API
type Product struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Image       string  `json:"image"`
	Rating      struct {
		Rate  float64 `json:"rate"`
		Count int     `json:"count"`
	} `json:"rating"`
}

// ExecuteWithContext fetches data from Fake Store API
func (f *FakeStoreAPI) ExecuteWithContext(ctx context.Context, config FakeStoreConfig) Result {
	start := time.Now()

	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled before Fake Store API request: " + ctx.Err().Error())
	default:
	}

	// Default values
	if config.Endpoint == "" {
		config.Endpoint = "products"
	}
	if config.Limit == 0 {
		config.Limit = 10
	}

	// Build API URL
	apiURL := fmt.Sprintf("https://fakestoreapi.com/%s", config.Endpoint)
	
	// Add category filter for products
	if config.Endpoint == "products" && config.Category != "" {
		apiURL = fmt.Sprintf("https://fakestoreapi.com/products/category/%s", config.Category)
	}

	// Add limit parameter
	if config.Limit > 0 && config.Limit < 20 {
		apiURL += fmt.Sprintf("?limit=%d", config.Limit)
	}

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to create Fake Store API request: %v", err), start)
	}

	// Execute request with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)

	// Check if context was cancelled during request
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled during Fake Store API request: " + ctx.Err().Error())
	default:
	}

	if err != nil {
		return NewFailureResult(fmt.Sprintf("Fake Store API request failed: %v", err), start)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode >= 400 {
		return NewFailureResult(fmt.Sprintf("Fake Store API returned error status: %d", resp.StatusCode), start)
	}

	// Parse response based on endpoint
	var data interface{}
	if config.Endpoint == "products" || config.Category != "" {
		var products []Product
		if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
			return NewFailureResult(fmt.Sprintf("Failed to parse Fake Store API response: %v", err), start)
		}
		data = products
	} else {
		// Generic JSON parsing for other endpoints
		var genericData interface{}
		if err := json.NewDecoder(resp.Body).Decode(&genericData); err != nil {
			return NewFailureResult(fmt.Sprintf("Failed to parse Fake Store API response: %v", err), start)
		}
		data = genericData
	}

	return NewSuccessResult("Fake Store data fetched successfully", map[string]interface{}{
		"endpoint": config.Endpoint,
		"data":     data,
	}, start)
}

// GetCategories is a helper to fetch available categories
func (f *FakeStoreAPI) GetCategories(ctx context.Context) Result {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://fakestoreapi.com/products/categories", nil)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to create request: %v", err), start)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Request failed: %v", err), start)
	}
	defer resp.Body.Close()

	var categories []string
	if err := json.NewDecoder(resp.Body).Decode(&categories); err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to parse response: %v", err), start)
	}

	return NewSuccessResult("Categories fetched successfully", map[string]interface{}{
		"categories": categories,
	}, start)
}

