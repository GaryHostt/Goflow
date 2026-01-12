package connectors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// SalesforceConnector interacts with Salesforce REST API
type SalesforceConnector struct {
	InstanceURL string // e.g., https://yourcompany.my.salesforce.com
	AccessToken string // OAuth2 access token
	APIVersion  string // Default: v59.0
}

// SalesforceConfig represents Salesforce connector configuration
type SalesforceConfig struct {
	Operation    string                 `json:"operation"`     // query, create, update, delete, get
	Object       string                 `json:"object"`        // Account, Contact, Lead, Opportunity, etc.
	RecordID     string                 `json:"record_id"`     // For get/update/delete operations
	Query        string                 `json:"query"`         // SOQL query for query operation
	Data         map[string]interface{} `json:"data"`          // Data for create/update operations
	InstanceURL  string                 `json:"instance_url"`  // Override instance URL
	AccessToken  string                 `json:"access_token"`  // Override access token
}

// SalesforceAuthConfig represents OAuth2 authentication config
type SalesforceAuthConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	SecurityToken string `json:"security_token"` // Appended to password
	LoginURL     string `json:"login_url"`      // Default: https://login.salesforce.com
}

// SalesforceTokenResponse represents OAuth2 token response
type SalesforceTokenResponse struct {
	AccessToken string `json:"access_token"`
	InstanceURL string `json:"instance_url"`
	ID          string `json:"id"`
	TokenType   string `json:"token_type"`
	IssuedAt    string `json:"issued_at"`
	Signature   string `json:"signature"`
}

// ExecuteWithContext performs Salesforce operations
func (s *SalesforceConnector) ExecuteWithContext(ctx context.Context, config SalesforceConfig) Result {
	start := time.Now()

	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled before Salesforce request: " + ctx.Err().Error())
	default:
	}

	// Override connector settings with config if provided
	instanceURL := s.InstanceURL
	if config.InstanceURL != "" {
		instanceURL = config.InstanceURL
	}

	accessToken := s.AccessToken
	if config.AccessToken != "" {
		accessToken = config.AccessToken
	}

	// Validate required fields
	if instanceURL == "" {
		return NewFailureResult("Salesforce instance URL is required", start)
	}
	if accessToken == "" {
		return NewFailureResult("Salesforce access token is required", start)
	}

	// Set default API version
	apiVersion := s.APIVersion
	if apiVersion == "" {
		apiVersion = "v59.0"
	}

	// Execute operation based on type
	switch config.Operation {
	case "query":
		return s.executeQuery(ctx, instanceURL, accessToken, apiVersion, config.Query, start)
	case "create":
		return s.executeCreate(ctx, instanceURL, accessToken, apiVersion, config.Object, config.Data, start)
	case "get":
		return s.executeGet(ctx, instanceURL, accessToken, apiVersion, config.Object, config.RecordID, start)
	case "update":
		return s.executeUpdate(ctx, instanceURL, accessToken, apiVersion, config.Object, config.RecordID, config.Data, start)
	case "delete":
		return s.executeDelete(ctx, instanceURL, accessToken, apiVersion, config.Object, config.RecordID, start)
	default:
		return NewFailureResult(fmt.Sprintf("Invalid Salesforce operation: %s. Valid: query, create, get, update, delete", config.Operation), start)
	}
}

// executeQuery runs a SOQL query
func (s *SalesforceConnector) executeQuery(ctx context.Context, instanceURL, accessToken, apiVersion, query string, start time.Time) Result {
	if query == "" {
		return NewFailureResult("SOQL query is required", start)
	}

	// Build URL with encoded query
	queryURL := fmt.Sprintf("%s/services/data/%s/query?q=%s", instanceURL, apiVersion, url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, "GET", queryURL, nil)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to create Salesforce request: %v", err), start)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)

	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled during Salesforce query: " + ctx.Err().Error())
	default:
	}

	if err != nil {
		return NewFailureResult(fmt.Sprintf("Salesforce query failed: %v", err), start)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to read Salesforce response: %v", err), start)
	}

	if resp.StatusCode >= 400 {
		return NewFailureResult(fmt.Sprintf("Salesforce returned HTTP error: %d - %s", resp.StatusCode, string(body)), start)
	}

	var queryResult map[string]interface{}
	if err := json.Unmarshal(body, &queryResult); err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to parse Salesforce response: %v", err), start)
	}

	recordCount := 0
	if records, ok := queryResult["records"].([]interface{}); ok {
		recordCount = len(records)
	}

	return NewSuccessResult(fmt.Sprintf("Salesforce query returned %d records", recordCount), map[string]interface{}{
		"operation":    "query",
		"query":        query,
		"record_count": recordCount,
		"data":         queryResult,
	}, start)
}

// executeCreate creates a new record
func (s *SalesforceConnector) executeCreate(ctx context.Context, instanceURL, accessToken, apiVersion, object string, data map[string]interface{}, start time.Time) Result {
	if object == "" {
		return NewFailureResult("Salesforce object type is required", start)
	}
	if len(data) == 0 {
		return NewFailureResult("Data is required to create Salesforce record", start)
	}

	createURL := fmt.Sprintf("%s/services/data/%s/sobjects/%s", instanceURL, apiVersion, object)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to marshal data: %v", err), start)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", createURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to create Salesforce request: %v", err), start)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)

	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled during Salesforce create: " + ctx.Err().Error())
	default:
	}

	if err != nil {
		return NewFailureResult(fmt.Sprintf("Salesforce create failed: %v", err), start)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to read Salesforce response: %v", err), start)
	}

	if resp.StatusCode >= 400 {
		return NewFailureResult(fmt.Sprintf("Salesforce returned HTTP error: %d - %s", resp.StatusCode, string(body)), start)
	}

	var createResult map[string]interface{}
	if err := json.Unmarshal(body, &createResult); err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to parse Salesforce response: %v", err), start)
	}

	recordID := ""
	if id, ok := createResult["id"].(string); ok {
		recordID = id
	}

	return NewSuccessResult(fmt.Sprintf("Salesforce %s created successfully: %s", object, recordID), map[string]interface{}{
		"operation": "create",
		"object":    object,
		"record_id": recordID,
		"data":      createResult,
	}, start)
}

// executeGet retrieves a record by ID
func (s *SalesforceConnector) executeGet(ctx context.Context, instanceURL, accessToken, apiVersion, object, recordID string, start time.Time) Result {
	if object == "" || recordID == "" {
		return NewFailureResult("Salesforce object type and record ID are required", start)
	}

	getURL := fmt.Sprintf("%s/services/data/%s/sobjects/%s/%s", instanceURL, apiVersion, object, recordID)

	req, err := http.NewRequestWithContext(ctx, "GET", getURL, nil)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to create Salesforce request: %v", err), start)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)

	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled during Salesforce get: " + ctx.Err().Error())
	default:
	}

	if err != nil {
		return NewFailureResult(fmt.Sprintf("Salesforce get failed: %v", err), start)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to read Salesforce response: %v", err), start)
	}

	if resp.StatusCode >= 400 {
		return NewFailureResult(fmt.Sprintf("Salesforce returned HTTP error: %d - %s", resp.StatusCode, string(body)), start)
	}

	var record map[string]interface{}
	if err := json.Unmarshal(body, &record); err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to parse Salesforce response: %v", err), start)
	}

	return NewSuccessResult(fmt.Sprintf("Salesforce %s retrieved: %s", object, recordID), map[string]interface{}{
		"operation": "get",
		"object":    object,
		"record_id": recordID,
		"data":      record,
	}, start)
}

// executeUpdate updates an existing record
func (s *SalesforceConnector) executeUpdate(ctx context.Context, instanceURL, accessToken, apiVersion, object, recordID string, data map[string]interface{}, start time.Time) Result {
	if object == "" || recordID == "" {
		return NewFailureResult("Salesforce object type and record ID are required", start)
	}
	if len(data) == 0 {
		return NewFailureResult("Data is required to update Salesforce record", start)
	}

	updateURL := fmt.Sprintf("%s/services/data/%s/sobjects/%s/%s", instanceURL, apiVersion, object, recordID)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to marshal data: %v", err), start)
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", updateURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to create Salesforce request: %v", err), start)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)

	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled during Salesforce update: " + ctx.Err().Error())
	default:
	}

	if err != nil {
		return NewFailureResult(fmt.Sprintf("Salesforce update failed: %v", err), start)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to read Salesforce response: %v", err), start)
	}

	if resp.StatusCode >= 400 {
		return NewFailureResult(fmt.Sprintf("Salesforce returned HTTP error: %d - %s", resp.StatusCode, string(body)), start)
	}

	return NewSuccessResult(fmt.Sprintf("Salesforce %s updated: %s", object, recordID), map[string]interface{}{
		"operation": "update",
		"object":    object,
		"record_id": recordID,
	}, start)
}

// executeDelete deletes a record
func (s *SalesforceConnector) executeDelete(ctx context.Context, instanceURL, accessToken, apiVersion, object, recordID string, start time.Time) Result {
	if object == "" || recordID == "" {
		return NewFailureResult("Salesforce object type and record ID are required", start)
	}

	deleteURL := fmt.Sprintf("%s/services/data/%s/sobjects/%s/%s", instanceURL, apiVersion, object, recordID)

	req, err := http.NewRequestWithContext(ctx, "DELETE", deleteURL, nil)
	if err != nil {
		return NewFailureResult(fmt.Sprintf("Failed to create Salesforce request: %v", err), start)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)

	select {
	case <-ctx.Done():
		return NewCancelledResult("Context cancelled during Salesforce delete: " + ctx.Err().Error())
	default:
	}

	if err != nil {
		return NewFailureResult(fmt.Sprintf("Salesforce delete failed: %v", err), start)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return NewFailureResult(fmt.Sprintf("Salesforce returned HTTP error: %d - %s", resp.StatusCode, string(body)), start)
	}

	return NewSuccessResult(fmt.Sprintf("Salesforce %s deleted: %s", object, recordID), map[string]interface{}{
		"operation": "delete",
		"object":    object,
		"record_id": recordID,
	}, start)
}

// Authenticate obtains an OAuth2 access token using password grant
func (s *SalesforceConnector) Authenticate(ctx context.Context, config SalesforceAuthConfig) (*SalesforceTokenResponse, error) {
	if config.LoginURL == "" {
		config.LoginURL = "https://login.salesforce.com"
	}

	tokenURL := config.LoginURL + "/services/oauth2/token"

	// Build form data
	formData := url.Values{}
	formData.Set("grant_type", "password")
	formData.Set("client_id", config.ClientID)
	formData.Set("client_secret", config.ClientSecret)
	formData.Set("username", config.Username)
	// Append security token to password
	formData.Set("password", config.Password+config.SecurityToken)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Salesforce authentication failed: %d - %s", resp.StatusCode, string(body))
	}

	var tokenResp SalesforceTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// DryRunSalesforce simulates a Salesforce call without actually making the request
func (s *SalesforceConnector) DryRunSalesforce(config SalesforceConfig) Result {
	start := time.Now()

	apiVersion := s.APIVersion
	if apiVersion == "" {
		apiVersion = "v59.0"
	}

	return NewSuccessResult("Salesforce dry run completed", map[string]interface{}{
		"operation":   config.Operation,
		"object":      config.Object,
		"record_id":   config.RecordID,
		"query":       config.Query,
		"api_version": apiVersion,
		"note":        "This is a dry run - no actual Salesforce call was made",
		"example_operations": map[string]string{
			"query":  "SELECT Id, Name FROM Account LIMIT 10",
			"create": "Create new Account: {Name: 'Acme Corp', Industry: 'Technology'}",
			"get":    "Retrieve Account record by ID",
			"update": "Update Account: {Phone: '+1-555-1234'}",
			"delete": "Delete Account record by ID",
		},
	}, start)
}

