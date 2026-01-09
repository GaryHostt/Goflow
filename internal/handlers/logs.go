package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/alexmacdonald/simple-ipass/internal/db"
	"github.com/alexmacdonald/simple-ipass/internal/middleware"
)

// LogsHandler handles log retrieval HTTP requests
// PRODUCTION: Uses Store interface for testability
type LogsHandler struct {
	store db.Store // Interface, not concrete type!
}

// NewLogsHandler creates a new logs handler
func NewLogsHandler(store db.Store) *LogsHandler {
	return &LogsHandler{store: store}
}

// GetLogs retrieves logs for the user's workflows
func (h *LogsHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	// TODO: MULTI-TENANT - Filter by tenant_id
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if filtering by specific workflow
	workflowID := r.URL.Query().Get("workflow_id")

	if workflowID != "" {
		// Verify ownership of workflow
		workflow, err := h.store.GetWorkflowByID(workflowID)
		if err != nil {
			http.Error(w, "Workflow not found", http.StatusNotFound)
			return
		}

		if workflow.UserID != userID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Get logs for this workflow
		logs, err := h.store.GetLogsByWorkflowID(workflowID)
		if err != nil {
			http.Error(w, "Failed to fetch logs", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logs)
		return
	}

	// Get all logs for user's workflows
	logs, err := h.store.GetLogsByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to fetch logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

