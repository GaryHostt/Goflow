package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/alexmacdonald/simple-ipass/internal/db"
	"github.com/alexmacdonald/simple-ipass/internal/engine"
	"github.com/gorilla/mux"
)

type WebhookHandler struct {
	db       *db.Database
	executor *engine.Executor
}

func NewWebhookHandler(database *db.Database, executor *engine.Executor) *WebhookHandler {
	return &WebhookHandler{db: database, executor: executor}
}

// TriggerWebhook handles incoming webhook requests
func (h *WebhookHandler) TriggerWebhook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	workflowID := vars["id"]

	// Lookup the workflow
	workflow, err := h.db.GetWorkflowByID(workflowID)
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

