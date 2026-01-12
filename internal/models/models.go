package models

import "time"

// User represents a user in the system
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never serialize password
	CreatedAt    time.Time `json:"created_at"`
}

// Credential represents encrypted API keys/tokens for third-party services
type Credential struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	ServiceName  string    `json:"service_name"` // e.g., 'slack', 'discord', 'openweather'
	EncryptedKey string    `json:"-"`            // Never expose in API
	DecryptedKey string    `json:"api_key,omitempty"` // Only populated when needed
	CreatedAt    time.Time `json:"created_at"`
}

// Workflow represents an integration workflow
type Workflow struct {
	ID              string         `json:"id"`
	UserID          string         `json:"user_id"`
	Name            string         `json:"name"`
	TriggerType     string         `json:"trigger_type"`     // 'webhook', 'schedule'
	ActionType      string         `json:"action_type"`      // Primary action: 'slack_message', 'discord_post', 'weather_check', etc.
	ConfigJSON      string         `json:"config_json"`      // Primary action configuration
	ActionChain     string         `json:"action_chain"`     // JSON array of additional actions to execute sequentially
	ParsedChain     []ChainedAction `json:"parsed_chain,omitempty"` // Parsed action chain (not stored in DB)
	TriggerPayload  string         `json:"trigger_payload,omitempty"` // JSON payload from webhook trigger for template mapping
	IsActive        bool           `json:"is_active"`
	LastExecutedAt  *time.Time     `json:"last_executed_at,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
}

// ChainedAction represents an additional action in a workflow chain
type ChainedAction struct {
	ActionType string                 `json:"action_type"` // 'slack_message', 'discord_post', 'twilio_sms', etc.
	Config     map[string]interface{} `json:"config"`      // Action-specific configuration
	UseDataFrom string                 `json:"use_data_from,omitempty"` // 'previous' to use data from previous action
}

// Log represents an execution log entry
type Log struct {
	ID         string    `json:"id"`
	WorkflowID string    `json:"workflow_id"`
	Status     string    `json:"status"` // 'success', 'failed'
	Message    string    `json:"message"`
	ExecutedAt time.Time `json:"executed_at"`
}

// WorkflowWithDetails includes workflow name for log display
type WorkflowLog struct {
	Log
	WorkflowName string `json:"workflow_name"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// RegisterRequest represents registration data
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=128"`
}

// AuthResponse represents the JWT token response
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// WorkflowConfig represents the configuration for different workflow types
type WorkflowConfig struct {
	// For webhook triggers
	WebhookURL string `json:"webhook_url,omitempty"`
	
	// For schedule triggers
	Interval int `json:"interval,omitempty"` // in minutes
	
	// For Slack action (supports templates like "Hello {{user.name}}")
	SlackMessage string `json:"slack_message,omitempty"`
	
	// For Discord action (supports templates like "Order {{order.id}} placed!")
	DiscordMessage string `json:"discord_message,omitempty"`
	
	// For Twilio SMS action
	TwilioTo      string `json:"twilio_to,omitempty"`      // Recipient phone number (supports templates like "{{user.phone}}")
	TwilioMessage string `json:"twilio_message,omitempty"` // SMS message (supports templates)
	
	// For News API action
	NewsQuery    string `json:"news_query,omitempty"`     // Search query (e.g., "bitcoin")
	NewsCountry  string `json:"news_country,omitempty"`   // Country code (e.g., "us")
	NewsCategory string `json:"news_category,omitempty"`  // Category (e.g., "technology")
	NewsPageSize int    `json:"news_page_size,omitempty"` // Number of articles (default: 10)
	
	// For Cat API action
	CatLimit     int    `json:"cat_limit,omitempty"`      // Number of cat images (default: 1)
	CatHasBreeds bool   `json:"cat_has_breeds,omitempty"` // Filter to cats with breed info
	CatBreedID   string `json:"cat_breed_id,omitempty"`   // Specific breed (e.g., "beng")
	CatCategory  string `json:"cat_category,omitempty"`   // Category (e.g., "boxes", "hats")
	
	// For Fake Store API action
	FakeStoreEndpoint string `json:"fakestore_endpoint,omitempty"` // "products", "users", "carts"
	FakeStoreLimit    int    `json:"fakestore_limit,omitempty"`    // Number of items
	FakeStoreCategory string `json:"fakestore_category,omitempty"` // Product category
	
	// For Weather check
	City string `json:"city,omitempty"`
	
	// For SOAP connector (Legacy protocol bridge)
	SOAPEndpoint   string                 `json:"soap_endpoint,omitempty"`   // SOAP service URL
	SOAPAction     string                 `json:"soap_action,omitempty"`     // SOAPAction header (optional)
	SOAPMethod     string                 `json:"soap_method,omitempty"`     // SOAP method name
	SOAPNamespace  string                 `json:"soap_namespace,omitempty"`  // XML namespace
	SOAPParameters map[string]interface{} `json:"soap_parameters,omitempty"` // Method parameters
	SOAPHeaders    map[string]string      `json:"soap_headers,omitempty"`    // Custom HTTP headers
	
	// For SWAPI connector (Star Wars API)
	SWAPIResource string `json:"swapi_resource,omitempty"` // films, people, planets, species, vehicles, starships
	SWAPIID       string `json:"swapi_id,omitempty"`       // Resource ID (e.g., "1" for first film)
	SWAPISearch   string `json:"swapi_search,omitempty"`   // Search query
	
	// For Salesforce connector
	SalesforceOperation  string                 `json:"salesforce_operation,omitempty"`   // query, create, get, update, delete
	SalesforceObject     string                 `json:"salesforce_object,omitempty"`      // Account, Contact, Lead, etc.
	SalesforceRecordID   string                 `json:"salesforce_record_id,omitempty"`   // Record ID for get/update/delete
	SalesforceQuery      string                 `json:"salesforce_query,omitempty"`       // SOQL query
	SalesforceData       map[string]interface{} `json:"salesforce_data,omitempty"`        // Data for create/update
	SalesforceInstanceURL string                 `json:"salesforce_instance_url,omitempty"` // Override instance URL
	
	// General purpose field for custom data
	CustomData map[string]interface{} `json:"custom_data,omitempty"`
}

