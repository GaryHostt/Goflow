package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alexmacdonald/simple-ipass/internal/db"
	"github.com/alexmacdonald/simple-ipass/internal/engine/connectors"
	"github.com/alexmacdonald/simple-ipass/internal/logger"
	"github.com/alexmacdonald/simple-ipass/internal/models"
	"github.com/alexmacdonald/simple-ipass/internal/utils"
)

// Executor handles workflow execution with structured logging
// Uses dependency injection (Store interface) for testability
type Executor struct {
	store          db.Store // Interface, not concrete type!
	log            *logger.Logger
	pool           *WorkerPool       // Bounded concurrency
	templateEngine *utils.TemplateEngine // Dynamic field mapping
}

// NewExecutor creates a new executor
func NewExecutor(store db.Store, log *logger.Logger) *Executor {
	// Initialize worker pool with 10 workers
	pool := NewWorkerPool(10, log)
	pool.Start()

	return &Executor{
		store:          store,
		log:            log,
		pool:           pool,
		templateEngine: utils.NewTemplateEngine(),
	}
}

// ExecuteWorkflow runs a workflow asynchronously via worker pool
// PRODUCTION: Uses bounded concurrency instead of unbounded goroutines
func (e *Executor) ExecuteWorkflow(workflow models.Workflow) {
	// Submit to worker pool instead of spawning goroutine directly
	e.pool.Submit(WorkflowJob{
		Workflow: workflow,
		Executor: e,
	})
}

// ExecuteWorkflowWithContext runs a workflow with context awareness
// PRODUCTION: Respects cancellation and timeouts
func (e *Executor) ExecuteWorkflowWithContext(ctx context.Context, workflow models.Workflow) {
	tenantID := "tenant_" + workflow.UserID

	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		e.log.WorkflowLog(
			logger.LevelWarn,
			"Workflow cancelled before execution",
			workflow.ID,
			workflow.UserID,
			tenantID,
			map[string]interface{}{
				"reason": ctx.Err().Error(),
			},
		)
		return
	default:
	}

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
	e.store.UpdateWorkflowLastExecuted(workflow.ID, time.Now())

	// Execute with context awareness
	result := e.executeWorkflowInternal(ctx, workflow, workflow.UserID, tenantID)

	// Only log if context wasn't cancelled
	select {
	case <-ctx.Done():
		e.log.WorkflowLog(
			logger.LevelWarn,
			"Workflow execution cancelled",
			workflow.ID,
			workflow.UserID,
			tenantID,
			map[string]interface{}{
				"reason":          ctx.Err().Error(),
				"partial_result": result.Status,
			},
		)
		return
	default:
		// Log to database
		e.store.CreateLog(workflow.ID, result.Status, result.Message)
	}
}

// DryRun executes a workflow synchronously without saving to database
// PRODUCT FEATURE: Test integration before committing
func (e *Executor) DryRun(workflow models.Workflow, userID, tenantID string) connectors.Result {
	// Use background context with timeout for dry runs
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	e.log.WorkflowLog(
		logger.LevelInfo,
		"Dry run execution (test mode)",
		workflow.ID,
		userID,
		tenantID,
		map[string]interface{}{
			"action_type": workflow.ActionType,
			"mode":        "dry_run",
		},
	)

	// Execute synchronously (blocking for immediate response)
	result := e.executeWorkflowInternal(ctx, workflow, userID, tenantID)

	// Log result (but NOT to database - it's a test!)
	logLevel := logger.LevelInfo
	if result.Status == "failed" {
		logLevel = logger.LevelError
	}

	e.log.WorkflowLog(
		logLevel,
		fmt.Sprintf("Dry run complete: %s", result.Message),
		workflow.ID,
		userID,
		tenantID,
		map[string]interface{}{
			"status":   result.Status,
			"duration": result.Duration,
			"mode":     "dry_run",
		},
	)

	return result
}

// executeWorkflowInternal contains the core execution logic with context awareness
// PRODUCTION: Respects context cancellation throughout execution
func (e *Executor) executeWorkflowInternal(ctx context.Context, workflow models.Workflow, userID, tenantID string) connectors.Result {
	start := time.Now()

	// Check context before parsing
	select {
	case <-ctx.Done():
		return connectors.Result{
			Status:    "cancelled",
			Message:   "Execution cancelled: " + ctx.Err().Error(),
			Duration:  time.Since(start).String(),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	default:
	}

	// Parse config
	var config models.WorkflowConfig
	if err := json.Unmarshal([]byte(workflow.ConfigJSON), &config); err != nil {
		return connectors.Result{
			Status:    "failed",
			Message:   fmt.Sprintf("Failed to parse config: %v", err),
			Duration:  time.Since(start).String(),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// Check context before executing action
	select {
	case <-ctx.Done():
		return connectors.Result{
			Status:    "cancelled",
			Message:   "Execution cancelled before action: " + ctx.Err().Error(),
			Duration:  time.Since(start).String(),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	default:
	}

	// Execute the action based on action type
	var result connectors.Result

	switch workflow.ActionType {
	case "slack_message":
		result = e.executeSlackAction(ctx, userID, tenantID, config, workflow.TriggerPayload)
	case "discord_post":
		result = e.executeDiscordAction(ctx, userID, tenantID, config, workflow.TriggerPayload)
	case "twilio_sms":
		result = e.executeTwilioAction(ctx, userID, tenantID, config, workflow.TriggerPayload)
	case "news_fetch":
		result = e.executeNewsAPIAction(ctx, userID, tenantID, config)
	case "cat_fetch":
		result = e.executeCatAPIAction(ctx, userID, tenantID, config)
	case "fakestore_fetch":
		result = e.executeFakeStoreAction(ctx, userID, tenantID, config)
	case "weather_check":
		result = e.executeWeatherAction(ctx, userID, tenantID, config)
	default:
		result = connectors.Result{
			Status:    "failed",
			Message:   fmt.Sprintf("Unknown action type: %s", workflow.ActionType),
			Duration:  time.Since(start).String(),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// Add total duration if not already set
	if result.Duration == "" {
		result.Duration = time.Since(start).String()
	}

	return result
}

// executeSlackAction sends a message to Slack with context awareness and dynamic templates
func (e *Executor) executeSlackAction(ctx context.Context, userID, tenantID string, config models.WorkflowConfig, triggerPayload string) connectors.Result {
	// Check context before fetching credentials
	select {
	case <-ctx.Done():
		return connectors.Result{
			Status:    "cancelled",
			Message:   ctx.Err().Error(),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	default:
	}

	// Get Slack credentials
	cred, err := e.store.GetCredentialByUserAndService(userID, "slack")
	if err != nil {
		e.log.Error("Slack credentials not found", map[string]interface{}{
			"user_id":   userID,
			"tenant_id": tenantID,
			"error":     err.Error(),
		})
		return connectors.Result{
			Status:    "failed",
			Message:   fmt.Sprintf("Slack not connected: %v", err),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	slack := &connectors.SlackWebhook{
		WebhookURL: cred.DecryptedKey,
	}

	message := config.SlackMessage
	if message == "" {
		message = "Hello from GoFlow! 🚀"
	}

	// Apply dynamic template mapping if trigger payload exists
	if triggerPayload != "" {
		message = e.templateEngine.Render(message, triggerPayload)
	}

	// Execute with context (connector should respect cancellation)
	return slack.ExecuteWithContext(ctx, message)
}

// executeDiscordAction sends a message to Discord
func (e *Executor) executeDiscordAction(ctx context.Context, userID, tenantID string, config models.WorkflowConfig) connectors.Result {
	select {
	case <-ctx.Done():
		return connectors.Result{
			Status:    "cancelled",
			Message:   ctx.Err().Error(),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	default:
	}

	cred, err := e.store.GetCredentialByUserAndService(userID, "discord")
	if err != nil {
		e.log.Error("Discord credentials not found", map[string]interface{}{
			"user_id":   userID,
			"tenant_id": tenantID,
			"error":     err.Error(),
		})
		return connectors.Result{
			Status:    "failed",
			Message:   fmt.Sprintf("Discord not connected: %v", err),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
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
func (e *Executor) executeWeatherAction(ctx context.Context, userID, tenantID string, config models.WorkflowConfig) connectors.Result {
	select {
	case <-ctx.Done():
		return connectors.Result{
			Status:    "cancelled",
			Message:   ctx.Err().Error(),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	default:
	}

	cred, err := e.store.GetCredentialByUserAndService(userID, "openweather")
	if err != nil {
		e.log.Error("OpenWeather credentials not found", map[string]interface{}{
			"user_id":   userID,
			"tenant_id": tenantID,
			"error":     err.Error(),
		})
		return connectors.Result{
			Status:    "failed",
			Message:   fmt.Sprintf("OpenWeather not connected: %v", err),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	weather := &connectors.OpenWeatherAPI{
		APIKey: cred.DecryptedKey,
	}

	city := config.City
	if city == "" {
		city = "London"
	}

	return weather.FetchWeather(city)
}

// executeTwilioAction sends an SMS via Twilio with dynamic templates
func (e *Executor) executeTwilioAction(ctx context.Context, userID, tenantID string, config models.WorkflowConfig, triggerPayload string) connectors.Result {
	select {
	case <-ctx.Done():
		return connectors.Result{
			Status:    "cancelled",
			Message:   ctx.Err().Error(),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	default:
	}

	// Get Twilio credentials
	cred, err := e.store.GetCredentialByUserAndService(userID, "twilio")
	if err != nil {
		e.log.Error("Twilio credentials not found", map[string]interface{}{
			"user_id":   userID,
			"tenant_id": tenantID,
			"error":     err.Error(),
		})
		return connectors.Result{
			Status:    "failed",
			Message:   fmt.Sprintf("Twilio not connected: %v", err),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// Parse Twilio credentials from JSON
	var twilioConfig struct {
		AccountSID string `json:"account_sid"`
		AuthToken  string `json:"auth_token"`
		FromNumber string `json:"from_number"`
	}
	if err := json.Unmarshal([]byte(cred.DecryptedKey), &twilioConfig); err != nil {
		return connectors.Result{
			Status:    "failed",
			Message:   fmt.Sprintf("Invalid Twilio credentials format: %v", err),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	twilio := &connectors.TwilioSMS{
		AccountSID: twilioConfig.AccountSID,
		AuthToken:  twilioConfig.AuthToken,
		FromNumber: twilioConfig.FromNumber,
	}

	// Prepare SMS config
	smsConfig := connectors.TwilioConfig{
		To:      config.TwilioTo,
		Message: config.TwilioMessage,
	}

	// Apply dynamic template mapping
	if triggerPayload != "" {
		smsConfig.Message = e.templateEngine.Render(smsConfig.Message, triggerPayload)
		smsConfig.To = e.templateEngine.Render(smsConfig.To, triggerPayload)
	}

	return twilio.ExecuteWithContext(ctx, smsConfig)
}

// executeNewsAPIAction fetches news articles
func (e *Executor) executeNewsAPIAction(ctx context.Context, userID, tenantID string, config models.WorkflowConfig) connectors.Result {
	select {
	case <-ctx.Done():
		return connectors.Result{
			Status:    "cancelled",
			Message:   ctx.Err().Error(),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	default:
	}

	// Get News API credentials
	cred, err := e.store.GetCredentialByUserAndService(userID, "newsapi")
	if err != nil {
		e.log.Error("News API credentials not found", map[string]interface{}{
			"user_id":   userID,
			"tenant_id": tenantID,
			"error":     err.Error(),
		})
		return connectors.Result{
			Status:    "failed",
			Message:   fmt.Sprintf("News API not connected: %v", err),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	newsAPI := &connectors.NewsAPI{
		APIKey: cred.DecryptedKey,
	}

	newsConfig := connectors.NewsConfig{
		Query:    config.NewsQuery,
		Country:  config.NewsCountry,
		Category: config.NewsCategory,
		PageSize: config.NewsPageSize,
	}

	return newsAPI.ExecuteWithContext(ctx, newsConfig)
}

// executeCatAPIAction fetches cat images
func (e *Executor) executeCatAPIAction(ctx context.Context, userID, tenantID string, config models.WorkflowConfig) connectors.Result {
	select {
	case <-ctx.Done():
		return connectors.Result{
			Status:    "cancelled",
			Message:   ctx.Err().Error(),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	default:
	}

	// Cat API key is optional, but we'll check for it
	var apiKey string
	cred, err := e.store.GetCredentialByUserAndService(userID, "catapi")
	if err == nil {
		apiKey = cred.DecryptedKey
	}

	catAPI := &connectors.CatAPI{
		APIKey: apiKey,
	}

	catConfig := connectors.CatConfig{
		Limit:     config.CatLimit,
		HasBreeds: config.CatHasBreeds,
		BreedID:   config.CatBreedID,
		Category:  config.CatCategory,
	}

	return catAPI.ExecuteWithContext(ctx, catConfig)
}

// executeFakeStoreAction fetches data from Fake Store API
func (e *Executor) executeFakeStoreAction(ctx context.Context, userID, tenantID string, config models.WorkflowConfig) connectors.Result {
	select {
	case <-ctx.Done():
		return connectors.Result{
			Status:    "cancelled",
			Message:   ctx.Err().Error(),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	default:
	}

	// Fake Store API doesn't require authentication
	fakeStore := &connectors.FakeStoreAPI{}

	storeConfig := connectors.FakeStoreConfig{
		Endpoint: config.FakeStoreEndpoint,
		Limit:    config.FakeStoreLimit,
		Category: config.FakeStoreCategory,
	}

	return fakeStore.ExecuteWithContext(ctx, storeConfig)
}

// Shutdown gracefully stops the executor
func (e *Executor) Shutdown(ctx context.Context) error {
	return e.pool.Shutdown(ctx)
}
