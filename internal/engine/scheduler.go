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
// PRODUCTION: Uses Store interface for testability
type Scheduler struct {
	store    db.Store // Interface, not concrete type!
	executor *Executor
	ticker   *time.Ticker
	done     chan bool
	log      *logger.Logger
	// MULTI-TENANT: Future fields for rate limiting
	// tenantRateLimits map[string]time.Duration
}

// NewScheduler creates a new scheduler
func NewScheduler(store db.Store, executor *Executor, log *logger.Logger) *Scheduler {
	return &Scheduler{
		store:    store,
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
// PRODUCTION: Uses panic recovery to prevent one bad workflow from crashing scheduler
func (s *Scheduler) checkAndExecute() {
	// PRODUCTION FIX: Recover from panics to keep scheduler running
	defer func() {
		if r := recover(); r != nil {
			s.log.Error("Scheduler panic recovered", map[string]interface{}{
				"panic": r,
			})
		}
	}()

	workflows, err := s.store.GetActiveScheduledWorkflows()
	if err != nil {
		s.log.Error("Failed to fetch scheduled workflows", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	now := time.Now()
	executedCount := 0

	for _, workflow := range workflows {
		// PRODUCTION FIX: Wrap each workflow execution in its own recovery
		func() {
			defer func() {
				if r := recover(); r != nil {
					s.log.Error("Workflow execution panic", map[string]interface{}{
						"workflow_id": workflow.ID,
						"user_id":     workflow.UserID,
						"panic":       r,
					})
				}
			}()

			// PRODUCTION FIX: Re-check is_active before execution
			// (user might have disabled it milliseconds ago)
			currentWorkflow, err := s.store.GetWorkflowByID(workflow.ID)
			if err != nil {
				s.log.Warn("Workflow not found during execution check", map[string]interface{}{
					"workflow_id": workflow.ID,
					"error":       err.Error(),
				})
				return
			}

			if !currentWorkflow.IsActive {
				s.log.Debug("Workflow disabled before execution", map[string]interface{}{
					"workflow_id": workflow.ID,
				})
				return
			}

			// Parse config to get interval
			var config models.WorkflowConfig
			if err := json.Unmarshal([]byte(workflow.ConfigJSON), &config); err != nil {
				s.log.Error("Failed to parse workflow config", map[string]interface{}{
					"workflow_id": workflow.ID,
					"user_id":     workflow.UserID,
					"error":       err.Error(),
				})
				return
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
				s.executor.ExecuteWorkflow(*currentWorkflow)
				executedCount++
			}
		}() // End of panic-recovery wrapper
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
