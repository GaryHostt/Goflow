package connectors

import (
	"time"
)

// Result represents the outcome of a connector execution
type Result struct {
	Status    string                 `json:"status"`    // "success", "failed", or "cancelled"
	Message   string                 `json:"message"`   // Human-readable message
	Data      map[string]interface{} `json:"data,omitempty"`
	Duration  string                 `json:"duration,omitempty"`
	Timestamp string                 `json:"timestamp"` // ISO8601 format
}

// NewSuccessResult creates a success result
func NewSuccessResult(message string, data map[string]interface{}, start time.Time) Result {
	return Result{
		Status:    "success",
		Message:   message,
		Data:      data,
		Duration:  time.Since(start).String(),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// NewFailureResult creates a failure result
func NewFailureResult(message string, start time.Time) Result {
	return Result{
		Status:    "failed",
		Message:   message,
		Duration:  time.Since(start).String(),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// NewCancelledResult creates a cancellation result
func NewCancelledResult(message string) Result {
	return Result{
		Status:    "cancelled",
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
