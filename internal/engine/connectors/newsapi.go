package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// NewsAPI handles News API integrations
// API Documentation: https://newsapi.org/docs
type NewsAPI struct {
	APIKey string
}

// NewsConfig represents News API query configuration
type NewsConfig struct {
	Query    string `json:"query"`     // Search query (e.g., "bitcoin")
	Country  string `json:"country"`   // Country code (e.g., "us")
	Category string `json:"category"`  // Category (e.g., "technology")
	PageSize int    `json:"page_size"` // Number of articles (default: 10)
}

// NewsArticle represents a single news article
type NewsArticle struct {
	Source struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"source"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PublishedAt string `json:"publishedAt"`
}

// NewsAPIResponse represents the News API response
type NewsAPIResponse struct {
	Status       string        `json:"status"`
	TotalResults int           `json:"totalResults"`
	Articles     []NewsArticle `json:"articles"`
}

// ExecuteWithContext fetches news articles from News API
func (n *NewsAPI) ExecuteWithContext(ctx context.Context, config NewsConfig) Result {
	start := time.Now()

	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled before News API request: " + ctx.Err().Error())
	default:
	}

	// Default values
	if config.PageSize == 0 {
		config.PageSize = 10
	}
	if config.PageSize > 100 {
		config.PageSize = 100 // News API limit
	}

	// Build API URL
	var apiURL string
	if config.Query != "" {
		// Search everything
		apiURL = fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&pageSize=%d&apiKey=%s",
			config.Query, config.PageSize, n.APIKey)
	} else {
		// Top headlines
		apiURL = fmt.Sprintf("https://newsapi.org/v2/top-headlines?pageSize=%d&apiKey=%s",
			config.PageSize, n.APIKey)
		if config.Country != "" {
			apiURL += "&country=" + config.Country
		}
		if config.Category != "" {
			apiURL += "&category=" + config.Category
		}
	}

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to create News API request: %v", err), start)
	}

	// Execute request with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)

	// Check if context was cancelled during request
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled during News API request: " + ctx.Err().Error())
	default:
	}

	if err != nil {
		return NewFailureResult(fmt.Sprintf("News API request failed: %v", err), start)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode >= 400 {
		return NewFailureResult(fmt.Sprintf("News API returned error status: %d", resp.StatusCode), start)
	}

	// Parse response
	var newsResp NewsAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&newsResp); err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to parse News API response: %v", err), start)
	}

	if newsResp.Status != "ok" {
		return NewFailureResult(fmt.Sprintf("News API returned error status: %s", newsResp.Status), start)
	}

	return NewSuccessResult("News articles fetched successfully", map[string]interface{}{
		"total_results": newsResp.TotalResults,
		"articles":      newsResp.Articles,
		"count":         len(newsResp.Articles),
	}, start)
}

