package engine

import (
	"encoding/json"
	"fmt"

	"github.com/alexmacdonald/simple-ipass/internal/db"
	"github.com/alexmacdonald/simple-ipass/internal/engine/connectors"
	"github.com/alexmacdonald/simple-ipass/internal/logger"
	"github.com/alexmacdonald/simple-ipass/internal/models"
	"time"
)

// Executor handles workflow execution with structured logging
type Executor struct {
	db  *db.Database
	log *logger.Logger
}

// NewExecutor creates a new executor
func NewExecutor(database *db.Database, log *logger.Logger) *Executor {
	return &Executor{
		db:  database,
		log: log,
	}
}

// ExecuteWorkflow runs a workflow asynchronously with full logging
func (e *Executor) ExecuteWorkflow(workflow models.Workflow) {
	go func() {
		tenantID := "tenant_" + workflow.UserID // Phase 1: derive tenant from user

		e.log.WorkflowLog(
			logger.LevelInfo,
			"Executing workflow",
			workflow.ID,
			workflow.UserID,
			tenantID,
			map[string]interface{}{
				"workflow_name": workflow.Name,
				"trigger_type":  workflow.TriggerType,
				"action_type":   workflow.ActionType,
			},
		)

		// Update last executed time
		e.db.UpdateWorkflowLastExecuted(workflow.ID, time.Now())

		// Parse config
		var config models.WorkflowConfig
		if err := json.Unmarshal([]byte(workflow.ConfigJSON), &config); err != nil {
			e.log.WorkflowLog(
				logger.LevelError,
				fmt.Sprintf("Failed to parse config: %v", err),
				workflow.ID,
				workflow.UserID,
				tenantID,
				map[string]interface{}{
					"error": err.Error(),
				},
			)
			e.db.CreateLog(workflow.ID, "failed", fmt.Sprintf("Failed to parse config: %v", err))
			return
		}

		// Execute the action based on action type
		var result connectors.Result

		switch workflow.ActionType {
		case "slack_message":
			result = e.executeSlackAction(workflow.UserID, tenantID, config)
		case "discord_post":
			result = e.executeDiscordAction(workflow.UserID, tenantID, config)
		case "weather_check":
			result = e.executeWeatherAction(workflow.UserID, tenantID, config)
		default:
			result = connectors.Result{
				Status:  "failed",
				Message: fmt.Sprintf("Unknown action type: %s", workflow.ActionType),
			}
		}

		// Log the result with full context
		logLevel := logger.LevelInfo
		if result.Status == "failed" {
			logLevel = logger.LevelError
		}

		e.log.WorkflowLog(
			logLevel,
			result.Message,
			workflow.ID,
			workflow.UserID,
			tenantID,
			map[string]interface{}{
				"status":      result.Status,
				"action_type": workflow.ActionType,
				"data":        result.Data,
			},
		)

		// Log to database
		e.db.CreateLog(workflow.ID, result.Status, result.Message)
	}()
}

// executeSlackAction sends a message to Slack
func (e *Executor) executeSlackAction(userID, tenantID string, config models.WorkflowConfig) connectors.Result {
	// Get Slack credentials
	cred, err := e.db.GetCredentialByUserAndService(userID, "slack")
	if err != nil {
		e.log.Error("Slack credentials not found", map[string]interface{}{
			"user_id":   userID,
			"tenant_id": tenantID,
			"error":     err.Error(),
		})
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
func (e *Executor) executeDiscordAction(userID, tenantID string, config models.WorkflowConfig) connectors.Result {
	// Get Discord credentials
	cred, err := e.db.GetCredentialByUserAndService(userID, "discord")
	if err != nil {
		e.log.Error("Discord credentials not found", map[string]interface{}{
			"user_id":   userID,
			"tenant_id": tenantID,
			"error":     err.Error(),
		})
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
func (e *Executor) executeWeatherAction(userID, tenantID string, config models.WorkflowConfig) connectors.Result {
	// Get OpenWeather credentials
	cred, err := e.db.GetCredentialByUserAndService(userID, "openweather")
	if err != nil {
		e.log.Error("OpenWeather credentials not found", map[string]interface{}{
			"user_id":   userID,
			"tenant_id": tenantID,
			"error":     err.Error(),
		})
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
