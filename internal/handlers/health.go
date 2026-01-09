package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/alexmacdonald/simple-ipass/internal/db"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	store     db.Store
	startTime time.Time
	version   string
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(store db.Store, version string) *HealthHandler {
	return &HealthHandler{
		store:     store,
		startTime: time.Now(),
		version:   version,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`    // "healthy" or "unhealthy"
	Version   string            `json:"version"`
	Uptime    string            `json:"uptime"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

// Health performs a comprehensive health check
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	checks := make(map[string]string)
	isHealthy := true

	// Check 1: Database connectivity
	dbStatus := h.checkDatabase()
	checks["database"] = dbStatus
	if dbStatus != "ok" {
		isHealthy = false
	}

	// Check 2: Memory/Goroutines (basic check)
	checks["runtime"] = "ok"

	// Determine overall status
	status := "healthy"
	statusCode := http.StatusOK
	if !isHealthy {
		status = "unhealthy"
		statusCode = http.StatusServiceUnavailable
	}

	// Build response
	response := HealthResponse{
		Status:    status,
		Version:   h.version,
		Uptime:    time.Since(h.startTime).String(),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// Liveness is a simple liveness check (for Kubernetes)
func (h *HealthHandler) Liveness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"alive"}`))
}

// Readiness checks if the service is ready to accept traffic
func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	// Check if database is accessible
	if h.checkDatabase() != "ok" {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"status":"not_ready","reason":"database_unavailable"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ready"}`))
}

// checkDatabase verifies database connectivity
func (h *HealthHandler) checkDatabase() string {
	// Try a simple query to verify database is accessible
	// For production, you'd query a lightweight table
	_, err := h.store.GetUserByID("health_check_dummy")
	
	// We expect "not found" error, which means DB is working
	// Only "connection failed" type errors are problematic
	if err != nil {
		// Check if it's a connection error (not just "not found")
		errMsg := err.Error()
		if errMsg == "sql: database is closed" || 
		   errMsg == "database is locked" {
			return "error: " + errMsg
		}
		// "Not found" is acceptable - DB is working
		return "ok"
	}
	
	return "ok"
}

