package engine

import (
	"context"
	"sync"
	"time"

	"github.com/alexmacdonald/simple-ipass/internal/logger"
	"github.com/alexmacdonald/simple-ipass/internal/models"
)

// WorkflowJob represents a workflow execution job
type WorkflowJob struct {
	Workflow models.Workflow
	Executor *Executor
}

// WorkerPool manages a fixed number of workers to prevent resource exhaustion
// PRODUCTION: Bounded concurrency instead of unlimited goroutines
type WorkerPool struct {
	jobQueue   chan WorkflowJob
	workerCount int
	log        *logger.Logger
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workerCount int, log *logger.Logger) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		jobQueue:    make(chan WorkflowJob, workerCount*10), // Buffer 10x workers
		workerCount: workerCount,
		log:         log,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start spawns the worker goroutines
func (wp *WorkerPool) Start() {
	wp.log.Info("Starting worker pool", map[string]interface{}{
		"workers":     wp.workerCount,
		"queue_size": cap(wp.jobQueue),
	})

	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker is the individual worker goroutine
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	wp.log.Debug("Worker started", map[string]interface{}{
		"worker_id": id,
	})

	for {
		select {
		case <-wp.ctx.Done():
			wp.log.Debug("Worker stopping", map[string]interface{}{
				"worker_id": id,
			})
			return

		case job, ok := <-wp.jobQueue:
			if !ok {
				wp.log.Debug("Worker queue closed", map[string]interface{}{
					"worker_id": id,
				})
				return
			}

			// Execute the job with timeout
			wp.executeJob(job, id)
		}
	}
}

// executeJob runs a single workflow job with context awareness
func (wp *WorkerPool) executeJob(job WorkflowJob, workerID int) {
	// Create context with timeout for this job
	ctx, cancel := context.WithTimeout(wp.ctx, 5*time.Minute)
	defer cancel()

	wp.log.Debug("Worker processing job", map[string]interface{}{
		"worker_id":   workerID,
		"workflow_id": job.Workflow.ID,
	})

	start := time.Now()

	// Execute with context awareness
	job.Executor.ExecuteWorkflowWithContext(ctx, job.Workflow)

	duration := time.Since(start)

	wp.log.Debug("Worker completed job", map[string]interface{}{
		"worker_id":   workerID,
		"workflow_id": job.Workflow.ID,
		"duration":    duration.String(),
	})
}

// Submit adds a job to the queue
// PRODUCTION: Non-blocking with queue full handling
func (wp *WorkerPool) Submit(job WorkflowJob) {
	select {
	case wp.jobQueue <- job:
		// Job submitted successfully
	case <-time.After(5 * time.Second):
		// Queue is full, log warning
		wp.log.Warn("Worker queue full, job dropped", map[string]interface{}{
			"workflow_id":  job.Workflow.ID,
			"queue_length": len(wp.jobQueue),
			"queue_cap":    cap(wp.jobQueue),
		})
	}
}

// Shutdown gracefully stops the worker pool
func (wp *WorkerPool) Shutdown(ctx context.Context) error {
	wp.log.Info("Shutting down worker pool", map[string]interface{}{
		"pending_jobs": len(wp.jobQueue),
	})

	// Stop accepting new jobs
	close(wp.jobQueue)

	// Signal workers to stop
	wp.cancel()

	// Wait for workers to finish with timeout
	done := make(chan struct{})
	go func() {
		wp.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		wp.log.Info("Worker pool shutdown complete", nil)
		return nil
	case <-ctx.Done():
		wp.log.Warn("Worker pool shutdown timeout", map[string]interface{}{
			"timeout": ctx.Err().Error(),
		})
		return ctx.Err()
	}
}

// QueueLength returns the current number of pending jobs
func (wp *WorkerPool) QueueLength() int {
	return len(wp.jobQueue)
}

// QueueCapacity returns the maximum queue size
func (wp *WorkerPool) QueueCapacity() int {
	return cap(wp.jobQueue)
}

