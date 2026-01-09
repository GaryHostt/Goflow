package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/alexmacdonald/simple-ipass/internal/db"
	"github.com/alexmacdonald/simple-ipass/internal/engine"
	"github.com/alexmacdonald/simple-ipass/internal/middleware"
	"github.com/alexmacdonald/simple-ipass/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// WorkflowsHandler handles workflow-related HTTP requests
// PRODUCTION: Uses Store interface for testability
type WorkflowsHandler struct {
	store    db.Store // Interface, not concrete type!
	executor *engine.Executor
}

// NewWorkflowsHandler creates a new workflows handler
func NewWorkflowsHandler(store db.Store, executor *engine.Executor) *WorkflowsHandler {
	return &WorkflowsHandler{store: store, executor: executor}
}

type CreateWorkflowRequest struct {
	Name        string `json:"name"`
	TriggerType string `json:"trigger_type"` // 'webhook', 'schedule'
	ActionType  string `json:"action_type"`  // 'slack_message', 'discord_post', 'weather_check'
	ConfigJSON  string `json:"config_json"`
}

// DryRunRequest represents a test execution request without saving
type DryRunRequest struct {
	ActionType string `json:"action_type"` // 'slack_message', 'discord_post', 'weather_check'
	ConfigJSON string `json:"config_json"`
}

// DryRunResponse represents the result of a dry run
type DryRunResponse struct {
	Success   bool                   `json:"success"`
	Message   string                 `json:"message"`
	Duration  string                 `json:"duration"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Timestamp string                 `json:"timestamp"`
}

// CreateWorkflow creates a new workflow
func (h *WorkflowsHandler) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	// TODO: MULTI-TENANT - Filter by tenant_id
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateWorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Name == "" || req.TriggerType == "" || req.ActionType == "" {
		http.Error(w, "name, trigger_type, and action_type are required", http.StatusBadRequest)
		return
	}

	// Validate trigger and action types
	validTriggers := map[string]bool{"webhook": true, "schedule": true}
	validActions := map[string]bool{"slack_message": true, "discord_post": true, "weather_check": true}

	if !validTriggers[req.TriggerType] {
		http.Error(w, "Invalid trigger_type. Must be 'webhook' or 'schedule'", http.StatusBadRequest)
		return
	}

	if !validActions[req.ActionType] {
		http.Error(w, "Invalid action_type. Must be 'slack_message', 'discord_post', or 'weather_check'", http.StatusBadRequest)
		return
	}

	if req.ConfigJSON == "" {
		req.ConfigJSON = "{}"
	}

	workflow, err := h.store.CreateWorkflow(userID, req.Name, req.TriggerType, req.ActionType, req.ConfigJSON)
	if err != nil {
		http.Error(w, "Failed to create workflow", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(workflow)
}

// DryRunWorkflow tests a workflow configuration without saving it
// PRODUCT FEATURE: Allows users to verify their integration works before committing
func (h *WorkflowsHandler) DryRunWorkflow(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tenantID, _ := middleware.GetTenantIDFromContext(r.Context())

	var req DryRunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate action type
	validActions := map[string]bool{"slack_message": true, "discord_post": true, "weather_check": true}
	if !validActions[req.ActionType] {
		http.Error(w, "Invalid action_type", http.StatusBadRequest)
		return
	}

	if req.ConfigJSON == "" {
		req.ConfigJSON = "{}"
	}

	// Create a temporary workflow for dry run (not saved to database)
	tempWorkflow := models.Workflow{
		ID:          "dryrun_" + uuid.New().String(),
		UserID:      userID,
		Name:        "Dry Run Test",
		TriggerType: "webhook",
		ActionType:  req.ActionType,
		ConfigJSON:  req.ConfigJSON,
		IsActive:    true,
	}

	// Execute the workflow synchronously (blocking) for dry run
	result := h.executor.DryRun(tempWorkflow, userID, tenantID)

	// Build response
	response := DryRunResponse{
		Success:   result.Status == "success",
		Message:   result.Message,
		Duration:  result.Duration,
		Data:      result.Data,
		Timestamp: result.Timestamp,
	}

	if result.Status != "success" {
		response.Error = result.Message
	}

	w.Header().Set("Content-Type", "application/json")
	if response.Success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(response)
}

// GetWorkflows retrieves all workflows for the user
func (h *WorkflowsHandler) GetWorkflows(w http.ResponseWriter, r *http.Request) {
	// TODO: MULTI-TENANT - Filter by tenant_id
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	workflows, err := h.store.GetWorkflowsByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to fetch workflows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workflows)
}

// ToggleWorkflow enables or disables a workflow
func (h *WorkflowsHandler) ToggleWorkflow(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	workflowID := vars["id"]

	// Verify ownership
	workflow, err := h.store.GetWorkflowByID(workflowID)
	if err != nil {
		http.Error(w, "Workflow not found", http.StatusNotFound)
		return
	}

	if workflow.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Toggle active status
	newStatus := !workflow.IsActive
	if err := h.store.UpdateWorkflowActive(workflowID, newStatus); err != nil {
		http.Error(w, "Failed to update workflow", http.StatusInternalServerError)
		return
	}

	workflow.IsActive = newStatus
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workflow)
}

// DeleteWorkflow deletes a workflow
func (h *WorkflowsHandler) DeleteWorkflow(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	workflowID := vars["id"]

	// Verify ownership
	workflow, err := h.store.GetWorkflowByID(workflowID)
	if err != nil {
		http.Error(w, "Workflow not found", http.StatusNotFound)
		return
	}

	if workflow.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := h.store.DeleteWorkflow(workflowID); err != nil {
		http.Error(w, "Failed to delete workflow", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
