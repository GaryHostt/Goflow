package engine

import (
	"encoding/json"
	"log"
	"time"

	"github.com/alexmacdonald/simple-ipass/internal/db"
	"github.com/alexmacdonald/simple-ipass/internal/models"
)

// Scheduler handles scheduled workflow execution
type Scheduler struct {
	db       *db.Database
	executor *Executor
	ticker   *time.Ticker
	done     chan bool
}

// NewScheduler creates a new scheduler
func NewScheduler(database *db.Database, executor *Executor) *Scheduler {
	return &Scheduler{
		db:       database,
		executor: executor,
		done:     make(chan bool),
	}
}

// Start begins the scheduler loop
func (s *Scheduler) Start(interval time.Duration) {
	s.ticker = time.NewTicker(interval)
	log.Printf("Scheduler started with interval: %v", interval)

	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.checkAndExecute()
			case <-s.done:
				log.Println("Scheduler stopped")
				return
			}
		}
	}()
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.done <- true
}

// checkAndExecute checks for scheduled workflows that need to run
func (s *Scheduler) checkAndExecute() {
	workflows, err := s.db.GetActiveScheduledWorkflows()
	if err != nil {
		log.Printf("Error fetching scheduled workflows: %v", err)
		return
	}

	now := time.Now()

	for _, workflow := range workflows {
		// Parse config to get interval
		var config models.WorkflowConfig
		if err := json.Unmarshal([]byte(workflow.ConfigJSON), &config); err != nil {
			log.Printf("Failed to parse config for workflow %s: %v", workflow.ID, err)
			continue
		}

		// Default interval is 10 minutes if not specified
		interval := config.Interval
		if interval <= 0 {
			interval = 10
		}

		// Check if enough time has passed since last execution
		shouldExecute := false
		if workflow.LastExecutedAt == nil {
			shouldExecute = true
		} else {
			timeSinceLastExecution := now.Sub(*workflow.LastExecutedAt)
			if timeSinceLastExecution >= time.Duration(interval)*time.Minute {
				shouldExecute = true
			}
		}

		if shouldExecute {
			log.Printf("Triggering scheduled workflow: %s (ID: %s)", workflow.Name, workflow.ID)
			s.executor.ExecuteWorkflow(workflow)
		}
	}
}

