package engine

import (
	"encoding/json"
	"log"
	"time"

	"github.com/alexmacdonald/simple-ipass/internal/db"
	"github.com/alexmacdonald/simple-ipass/internal/logger"
	"github.com/alexmacdonald/simple-ipass/internal/models"
)

// Scheduler handles scheduled workflow execution with tenant-aware rate limiting
type Scheduler struct {
	db       *db.Database
	executor *Executor
	ticker   *time.Ticker
	done     chan bool
	log      *logger.Logger
	// MULTI-TENANT: Future fields for rate limiting
	// tenantRateLimits map[string]time.Duration
}

// NewScheduler creates a new scheduler
func NewScheduler(database *db.Database, executor *Executor, log *logger.Logger) *Scheduler {
	return &Scheduler{
		db:       database,
		executor: executor,
		done:     make(chan bool),
		log:      log,
	}
}

// Start begins the scheduler loop
func (s *Scheduler) Start(interval time.Duration) {
	s.ticker = time.NewTicker(interval)
	s.log.Info("Scheduler started", map[string]interface{}{
		"interval": interval.String(),
	})

	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.checkAndExecute()
			case <-s.done:
				s.log.Info("Scheduler stopped", nil)
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
		s.log.Error("Failed to fetch scheduled workflows", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	now := time.Now()
	executedCount := 0

	for _, workflow := range workflows {
		// Parse config to get interval
		var config models.WorkflowConfig
		if err := json.Unmarshal([]byte(workflow.ConfigJSON), &config); err != nil {
			s.log.Error("Failed to parse workflow config", map[string]interface{}{
				"workflow_id": workflow.ID,
				"user_id":     workflow.UserID,
				"error":       err.Error(),
			})
			continue
		}

		// Default interval is 10 minutes if not specified
		interval := config.Interval
		if interval <= 0 {
			interval = 10
		}

		// MULTI-TENANT: Check tenant-specific rate limits
		// tenantID := "tenant_" + workflow.UserID // Phase 1
		// if customInterval := s.getTenantRateLimit(tenantID); customInterval > 0 {
		//     interval = customInterval
		// }

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
			s.log.InfoWithContext(
				"Triggering scheduled workflow",
				workflow.UserID,
				"tenant_"+workflow.UserID, // Phase 1: user is tenant
				map[string]interface{}{
					"workflow_id":   workflow.ID,
					"workflow_name": workflow.Name,
					"interval":      interval,
				},
			)
			s.executor.ExecuteWorkflow(workflow)
			executedCount++
		}
	}

	if executedCount > 0 {
		s.log.Info("Scheduler tick completed", map[string]interface{}{
			"total_workflows": len(workflows),
			"executed":        executedCount,
		})
	}
}

// MULTI-TENANT: Future method for tenant-specific rate limits
// func (s *Scheduler) getTenantRateLimit(tenantID string) int {
//     // Query tenant settings from database
//     // SELECT polling_interval_minutes FROM tenant_settings WHERE tenant_id = ?
//     // 
//     // Example:
//     // - Free tier: 60 minutes
//     // - Pro tier: 10 minutes
//     // - Enterprise: 1 minute
//     return 0 // 0 means use default
// }
