package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/alexmacdonald/simple-ipass/internal/db"
	"github.com/alexmacdonald/simple-ipass/internal/engine"
	"github.com/gorilla/mux"
)

// WebhookHandler handles webhook-related HTTP requests  
// PRODUCTION: Uses Store interface for testability
type WebhookHandler struct {
	store    db.Store // Interface, not concrete type!
	executor *engine.Executor
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(store db.Store, executor *engine.Executor) *WebhookHandler {
	return &WebhookHandler{store: store, executor: executor}
}

// TriggerWebhook handles incoming webhook requests
func (h *WebhookHandler) TriggerWebhook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	workflowID := vars["id"]

	// Lookup the workflow
	workflow, err := h.store.GetWorkflowByID(workflowID)
	if err != nil {
		http.Error(w, "Workflow not found", http.StatusNotFound)
		return
	}

	// Check if workflow is active
	if !workflow.IsActive {
		http.Error(w, "Workflow is not active", http.StatusBadRequest)
		return
	}

	// Check if trigger type is webhook
	if workflow.TriggerType != "webhook" {
		http.Error(w, "This workflow does not support webhook triggers", http.StatusBadRequest)
		return
	}

	// Execute the workflow asynchronously
	h.executor.ExecuteWorkflow(*workflow)

	// Return immediate response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "triggered",
		"message": "Workflow execution started",
	})
}

