package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/alexmacdonald/simple-ipass/internal/db"
	"github.com/alexmacdonald/simple-ipass/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	log.Println("Generating test data for iPaaS...")

	// Initialize database
	database, err := db.New("ipaas.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Create test user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	user, err := database.CreateUser("demo@ipaas.com", string(hashedPassword))
	if err != nil {
		log.Printf("User might already exist: %v", err)
		// Try to get existing user
		user, err = database.GetUserByEmail("demo@ipaas.com")
		if err != nil {
			log.Fatalf("Failed to get user: %v", err)
		}
	}
	log.Printf("Created/found user: %s", user.Email)

	// Create test credentials (mock)
	credentials := []struct {
		service string
		key     string
	}{
		{"slack", "https://hooks.slack.com/services/MOCK/WEBHOOK/URL"},
		{"discord", "https://discord.com/api/webhooks/MOCK/WEBHOOK"},
		{"openweather", "mock_api_key_12345"},
	}

	for _, cred := range credentials {
		_, err := database.CreateCredential(user.ID, cred.service, cred.key)
		if err != nil {
			log.Printf("Credential for %s might already exist: %v", cred.service, err)
		} else {
			log.Printf("Created credential for: %s", cred.service)
		}
	}

	// Create test workflows
	workflows := []struct {
		name        string
		triggerType string
		actionType  string
		config      models.WorkflowConfig
	}{
		{
			name:        "Daily Weather Alert",
			triggerType: "schedule",
			actionType:  "weather_check",
			config: models.WorkflowConfig{
				City:     "New York",
				Interval: 60,
			},
		},
		{
			name:        "Webhook to Slack",
			triggerType: "webhook",
			actionType:  "slack_message",
			config: models.WorkflowConfig{
				SlackMessage: "Webhook triggered! 🎉",
			},
		},
		{
			name:        "Weather to Discord",
			triggerType: "schedule",
			actionType:  "discord_post",
			config: models.WorkflowConfig{
				DiscordMessage: "Daily weather update posted!",
				Interval:       30,
			},
		},
	}

	var workflowIDs []string
	for _, wf := range workflows {
		configJSON, _ := json.Marshal(wf.config)
		workflow, err := database.CreateWorkflow(user.ID, wf.name, wf.triggerType, wf.actionType, string(configJSON))
		if err != nil {
			log.Printf("Failed to create workflow %s: %v", wf.name, err)
			continue
		}
		workflowIDs = append(workflowIDs, workflow.ID)
		log.Printf("Created workflow: %s", wf.name)
	}

	// Generate historical logs
	statuses := []string{"success", "success", "success", "failed"} // 75% success rate
	messages := map[string][]string{
		"success": {
			"Message sent to Slack successfully",
			"Weather in New York: Clear (clear sky), Temperature: 18.5°C, Humidity: 65%",
			"Message sent to Discord successfully",
			"Workflow executed successfully",
			"Weather data fetched and processed",
		},
		"failed": {
			"Slack not connected: sql: no rows in result set",
			"Discord returned error status: 404",
			"OpenWeather API returned error status: 401",
			"Failed to execute action: network timeout",
		},
	}

	rand.Seed(time.Now().UnixNano())
	now := time.Now()

	for i := 0; i < 50; i++ {
		if len(workflowIDs) == 0 {
			break
		}

		workflowID := workflowIDs[rand.Intn(len(workflowIDs))]
		status := statuses[rand.Intn(len(statuses))]
		messageList := messages[status]
		message := messageList[rand.Intn(len(messageList))]

		// Create logs spread over the last 7 days
		hoursAgo := rand.Intn(24 * 7)
		executedAt := now.Add(-time.Duration(hoursAgo) * time.Hour)

		// Manually insert log with custom timestamp
		logID := fmt.Sprintf("log_%d_%d", i, time.Now().UnixNano())
		query := `INSERT INTO logs (id, workflow_id, status, message, executed_at) VALUES (?, ?, ?, ?, ?)`
		_, err := database.CreateLog(workflowID, status, message)
		if err != nil {
			log.Printf("Failed to create log: %v", err)
		}

		// Update workflow last executed time
		database.UpdateWorkflowLastExecuted(workflowID, executedAt)
	}

	log.Printf("Generated 50 historical log entries")
	log.Println("✅ Test data generation complete!")
	log.Println("---")
	log.Println("Test credentials:")
	log.Println("Email: demo@ipaas.com")
	log.Println("Password: password123")
}

