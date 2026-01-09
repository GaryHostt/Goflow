package logger

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// LogLevel represents the severity of a log entry
type LogLevel string

const (
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
	LevelDebug LogLevel = "debug"
)

// LogEntry represents a structured log entry ready for ELK
type LogEntry struct {
	Timestamp  time.Time              `json:"timestamp"`
	Level      LogLevel               `json:"level"`
	Message    string                 `json:"message"`
	UserID     string                 `json:"user_id,omitempty"`
	TenantID   string                 `json:"tenant_id,omitempty"`   // Multi-tenant ready!
	WorkflowID string                 `json:"workflow_id,omitempty"`
	Service    string                 `json:"service"`
	Meta       map[string]interface{} `json:"meta,omitempty"`
}

// Logger provides structured logging for ELK integration
type Logger struct {
	service string
}

// NewLogger creates a new structured logger
func NewLogger(service string) *Logger {
	return &Logger{service: service}
}

// Info logs an info-level message
func (l *Logger) Info(message string, meta map[string]interface{}) {
	l.log(LevelInfo, message, meta)
}

// Warn logs a warning-level message
func (l *Logger) Warn(message string, meta map[string]interface{}) {
	l.log(LevelWarn, message, meta)
}

// Error logs an error-level message
func (l *Logger) Error(message string, meta map[string]interface{}) {
	l.log(LevelError, message, meta)
}

// Debug logs a debug-level message
func (l *Logger) Debug(message string, meta map[string]interface{}) {
	l.log(LevelDebug, message, meta)
}

// InfoWithContext logs with user/tenant context for filtering in Kibana
func (l *Logger) InfoWithContext(message, userID, tenantID string, meta map[string]interface{}) {
	entry := l.buildEntry(LevelInfo, message, meta)
	entry.UserID = userID
	entry.TenantID = tenantID
	l.output(entry)
}

// WorkflowLog logs workflow execution events (highly queryable in ELK)
func (l *Logger) WorkflowLog(level LogLevel, message, workflowID, userID, tenantID string, meta map[string]interface{}) {
	entry := l.buildEntry(level, message, meta)
	entry.WorkflowID = workflowID
	entry.UserID = userID
	entry.TenantID = tenantID
	l.output(entry)
}

// log is the internal logging method
func (l *Logger) log(level LogLevel, message string, meta map[string]interface{}) {
	entry := l.buildEntry(level, message, meta)
	l.output(entry)
}

// buildEntry constructs a log entry
func (l *Logger) buildEntry(level LogLevel, message string, meta map[string]interface{}) LogEntry {
	return LogEntry{
		Timestamp: time.Now().UTC(),
		Level:     level,
		Message:   message,
		Service:   l.service,
		Meta:      meta,
	}
}

// output writes the log entry as JSON (ELK-ready format)
func (l *Logger) output(entry LogEntry) {
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		// Fallback to standard logging if JSON marshal fails
		log.Printf("[ERROR] Failed to marshal log entry: %v", err)
		return
	}

	// Output to stdout (captured by Docker/ELK)
	os.Stdout.Write(jsonBytes)
	os.Stdout.Write([]byte("\n"))

	// TODO: In production, also send directly to Elasticsearch
	// if elasticClient != nil {
	//     elasticClient.Index("ipaas-logs", entry)
	// }
}

// GetElasticSearchQuery generates a sample ES query for Kibana
func GetElasticSearchQuery() string {
	return `
// Example Kibana queries:

// 1. All logs for a specific tenant
GET /ipaas-logs/_search
{
  "query": {
    "match": { "tenant_id": "tenant_123" }
  }
}

// 2. Failed workflow executions
GET /ipaas-logs/_search
{
  "query": {
    "bool": {
      "must": [
        { "match": { "workflow_id": "*" } },
        { "match": { "level": "error" } }
      ]
    }
  }
}

// 3. Activity for a specific user
GET /ipaas-logs/_search
{
  "query": {
    "match": { "user_id": "user_123" }
  }
}
`
}

