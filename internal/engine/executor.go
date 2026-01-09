package engine

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/alexmacdonald/simple-ipass/internal/db"
	"github.com/alexmacdonald/simple-ipass/internal/engine/connectors"
	"github.com/alexmacdonald/simple-ipass/internal/models"
)

// Executor handles workflow execution
type Executor struct {
	db *db.Database
}

// NewExecutor creates a new executor
func NewExecutor(database *db.Database) *Executor {
	return &Executor{db: database}
}

// ExecuteWorkflow runs a workflow asynchronously
func (e *Executor) ExecuteWorkflow(workflow models.Workflow) {
	go func() {
		log.Printf("Executing workflow: %s (ID: %s)", workflow.Name, workflow.ID)

		// Update last executed time
		e.db.UpdateWorkflowLastExecuted(workflow.ID, time.Now())

		// Parse config
		var config models.WorkflowConfig
		if err := json.Unmarshal([]byte(workflow.ConfigJSON), &config); err != nil {
			e.db.CreateLog(workflow.ID, "failed", fmt.Sprintf("Failed to parse config: %v", err))
			return
		}

		// Execute the action based on action type
		var result connectors.Result

		switch workflow.ActionType {
		case "slack_message":
			result = e.executeSlackAction(workflow.UserID, config)
		case "discord_post":
			result = e.executeDiscordAction(workflow.UserID, config)
		case "weather_check":
			result = e.executeWeatherAction(workflow.UserID, config)
		default:
			result = connectors.Result{
				Status:  "failed",
				Message: fmt.Sprintf("Unknown action type: %s", workflow.ActionType),
			}
		}

		// Log the result
		e.db.CreateLog(workflow.ID, result.Status, result.Message)
		log.Printf("Workflow %s completed with status: %s", workflow.ID, result.Status)
	}()
}

// executeSlackAction sends a message to Slack
func (e *Executor) executeSlackAction(userID string, config models.WorkflowConfig) connectors.Result {
	// Get Slack credentials
	cred, err := e.db.GetCredentialByUserAndService(userID, "slack")
	if err != nil {
		return connectors.Result{
			Status:  "failed",
			Message: fmt.Sprintf("Slack not connected: %v", err),
		}
	}

	slack := &connectors.SlackWebhook{
		WebhookURL: cred.DecryptedKey,
	}

	message := config.SlackMessage
	if message == "" {
		message = "Hello from iPaaS! 🚀"
	}

	return slack.Execute(message)
}

// executeDiscordAction sends a message to Discord
func (e *Executor) executeDiscordAction(userID string, config models.WorkflowConfig) connectors.Result {
	// Get Discord credentials
	cred, err := e.db.GetCredentialByUserAndService(userID, "discord")
	if err != nil {
		return connectors.Result{
			Status:  "failed",
			Message: fmt.Sprintf("Discord not connected: %v", err),
		}
	}

	discord := &connectors.DiscordWebhook{
		WebhookURL: cred.DecryptedKey,
	}

	message := config.DiscordMessage
	if message == "" {
		message = "Hello from iPaaS! 🎮"
	}

	return discord.Execute(message)
}

// executeWeatherAction fetches weather data
func (e *Executor) executeWeatherAction(userID string, config models.WorkflowConfig) connectors.Result {
	// Get OpenWeather credentials
	cred, err := e.db.GetCredentialByUserAndService(userID, "openweather")
	if err != nil {
		return connectors.Result{
			Status:  "failed",
			Message: fmt.Sprintf("OpenWeather not connected: %v", err),
		}
	}

	weather := &connectors.OpenWeatherAPI{
		APIKey: cred.DecryptedKey,
	}

	city := config.City
	if city == "" {
		city = "London" // Default city
	}

	return weather.FetchWeather(city)
}

