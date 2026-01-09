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
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Name            string    `json:"name"`
	TriggerType     string    `json:"trigger_type"` // 'webhook', 'schedule'
	ActionType      string    `json:"action_type"`  // 'slack_message', 'discord_post', 'weather_check'
	ConfigJSON      string    `json:"config_json"`
	IsActive        bool      `json:"is_active"`
	LastExecutedAt  *time.Time `json:"last_executed_at,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
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
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest represents registration data
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
	
	// For Slack action
	SlackMessage string `json:"slack_message,omitempty"`
	
	// For Discord action
	DiscordMessage string `json:"discord_message,omitempty"`
	
	// For Weather check
	City string `json:"city,omitempty"`
	
	// General purpose field for custom data
	CustomData map[string]interface{} `json:"custom_data,omitempty"`
}

