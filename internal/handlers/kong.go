package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/alexmacdonald/simple-ipass/internal/db"
	"github.com/alexmacdonald/simple-ipass/internal/middleware"
	"github.com/alexmacdonald/simple-ipass/internal/models"
	"github.com/alexmacdonald/simple-ipass/internal/utils"
	"github.com/gorilla/mux"
)

// KongHandler handles Kong Gateway integration
type KongHandler struct {
	store       db.Store
	kongAdminURL string // Kong Admin API URL (default: http://kong:8001)
}

// NewKongHandler creates a new Kong handler
func NewKongHandler(store db.Store, kongAdminURL string) *KongHandler {
	if kongAdminURL == "" {
		kongAdminURL = "http://kong:8001" // Default in Docker
	}
	return &KongHandler{
		store:       store,
		kongAdminURL: kongAdminURL,
	}
}

// KongService represents a Kong service
type KongService struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Protocol string `json:"protocol,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Path     string `json:"path,omitempty"`
}

// KongRoute represents a Kong route
type KongRoute struct {
	ID      string   `json:"id,omitempty"`
	Name    string   `json:"name"`
	Paths   []string `json:"paths"`
	Methods []string `json:"methods,omitempty"`
	Service struct {
		ID string `json:"id"`
	} `json:"service"`
}

// KongPlugin represents a Kong plugin
type KongPlugin struct {
	ID      string                 `json:"id,omitempty"`
	Name    string                 `json:"name"`
	Service struct {
		ID string `json:"id"`
	} `json:"service,omitempty"`
	Config map[string]interface{} `json:"config"`
}

// CreateKongService creates a Kong service that proxies to GoFlow
func (h *KongHandler) CreateKongService(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name        string `json:"name"`
		WorkflowID  string `json:"workflow_id"`
		UseCaseName string `json:"use_case"` // protocol_bridge, webhook_handler, aggregator, etc.
	}

	if err := utils.DecodeJSONStrict(w, r, &req); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify workflow ownership
	workflow, err := h.store.GetWorkflowByID(req.WorkflowID)
	if err != nil {
		utils.WriteJSONError(w, "Workflow not found", http.StatusNotFound)
		return
	}

	if workflow.UserID != userID {
		utils.WriteJSONError(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Create Kong service pointing to GoFlow webhook
	kongService := KongService{
		Name: req.Name,
		URL:  fmt.Sprintf("http://backend:8080/api/webhooks/%s", req.WorkflowID),
	}

	// Call Kong Admin API
	serviceResp, err := h.callKongAdmin("POST", "/services", kongService)
	if err != nil {
		utils.WriteJSONError(w, fmt.Sprintf("Failed to create Kong service: %v", err), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, serviceResp, http.StatusCreated)
}

// CreateKongRoute creates a route for a Kong service
func (h *KongHandler) CreateKongRoute(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ServiceID string   `json:"service_id"`
		Name      string   `json:"name"`
		Paths     []string `json:"paths"`
		Methods   []string `json:"methods"`
	}

	if err := utils.DecodeJSONStrict(w, r, &req); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	kongRoute := KongRoute{
		Name:    req.Name,
		Paths:   req.Paths,
		Methods: req.Methods,
	}
	kongRoute.Service.ID = req.ServiceID

	// Call Kong Admin API
	routeResp, err := h.callKongAdmin("POST", "/routes", kongRoute)
	if err != nil {
		utils.WriteJSONError(w, fmt.Sprintf("Failed to create Kong route: %v", err), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, routeResp, http.StatusCreated)
}

// AddKongPlugin adds a plugin to a Kong service
func (h *KongHandler) AddKongPlugin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ServiceID  string                 `json:"service_id"`
		PluginName string                 `json:"plugin_name"` // rate-limiting, key-auth, oauth2, etc.
		Config     map[string]interface{} `json:"config"`
	}

	if err := utils.DecodeJSONStrict(w, r, &req); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	kongPlugin := KongPlugin{
		Name:   req.PluginName,
		Config: req.Config,
	}
	kongPlugin.Service.ID = req.ServiceID

	// Call Kong Admin API
	pluginResp, err := h.callKongAdmin("POST", "/plugins", kongPlugin)
	if err != nil {
		utils.WriteJSONError(w, fmt.Sprintf("Failed to add Kong plugin: %v", err), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, pluginResp, http.StatusCreated)
}

// ListKongServices lists all Kong services
func (h *KongHandler) ListKongServices(w http.ResponseWriter, r *http.Request) {
	services, err := h.callKongAdmin("GET", "/services", nil)
	if err != nil {
		utils.WriteJSONError(w, fmt.Sprintf("Failed to list Kong services: %v", err), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, services, http.StatusOK)
}

// DeleteKongService deletes a Kong service
func (h *KongHandler) DeleteKongService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceID := vars["id"]

	_, err := h.callKongAdmin("DELETE", fmt.Sprintf("/services/%s", serviceID), nil)
	if err != nil {
		utils.WriteJSONError(w, fmt.Sprintf("Failed to delete Kong service: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// callKongAdmin makes a request to Kong Admin API
func (h *KongHandler) callKongAdmin(method, path string, body interface{}) (map[string]interface{}, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, h.kongAdminURL+path, reqBody)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Kong API error: %d - %s", resp.StatusCode, string(responseBody))
	}

	// For DELETE requests, return empty response
	if method == "DELETE" {
		return map[string]interface{}{"success": true}, nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// CreateUseCaseTemplate creates a Kong setup for common use cases
func (h *KongHandler) CreateUseCaseTemplate(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		WorkflowID string `json:"workflow_id"`
		UseCase    string `json:"use_case"` // protocol_bridge, webhook_handler, aggregator, auth_overlay, monetization
	}

	if err := utils.DecodeJSONStrict(w, r, &req); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify workflow ownership
	workflow, err := h.store.GetWorkflowByID(req.WorkflowID)
	if err != nil {
		utils.WriteJSONError(w, "Workflow not found", http.StatusNotFound)
		return
	}

	if workflow.UserID != userID {
		utils.WriteJSONError(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Create service + route + plugins based on use case
	result, err := h.setupUseCase(req.UseCase, workflow)
	if err != nil {
		utils.WriteJSONError(w, fmt.Sprintf("Failed to setup use case: %v", err), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, result, http.StatusCreated)
}

// setupUseCase configures Kong for specific use cases
func (h *KongHandler) setupUseCase(useCase string, workflow *models.Workflow) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	switch useCase {
	case "protocol_bridge":
		// SOAP to REST bridge
		service := KongService{
			Name: fmt.Sprintf("bridge-%s", workflow.ID),
			URL:  fmt.Sprintf("http://backend:8080/api/webhooks/%s", workflow.ID),
		}
		serviceResp, err := h.callKongAdmin("POST", "/services", service)
		if err != nil {
			return nil, err
		}
		result["service"] = serviceResp

		// Add request transformer to convert REST to workflow format
		// Add response transformer to format response
		
	case "webhook_handler":
		// High-throughput webhook processing with rate limiting
		service := KongService{
			Name: fmt.Sprintf("webhook-%s", workflow.ID),
			URL:  fmt.Sprintf("http://backend:8080/api/webhooks/%s", workflow.ID),
		}
		serviceResp, err := h.callKongAdmin("POST", "/services", service)
		if err != nil {
			return nil, err
		}
		result["service"] = serviceResp

		// Add rate limiting plugin
		serviceID := serviceResp["id"].(string)
		plugin := KongPlugin{
			Name: "rate-limiting",
			Config: map[string]interface{}{
				"second": 100,
				"hour":   10000,
			},
		}
		plugin.Service.ID = serviceID
		pluginResp, err := h.callKongAdmin("POST", "/plugins", plugin)
		if err != nil {
			return nil, err
		}
		result["rate_limiting"] = pluginResp

	case "aggregator":
		// API aggregation with caching
		service := KongService{
			Name: fmt.Sprintf("aggregator-%s", workflow.ID),
			URL:  fmt.Sprintf("http://backend:8080/api/webhooks/%s", workflow.ID),
		}
		serviceResp, err := h.callKongAdmin("POST", "/services", service)
		if err != nil {
			return nil, err
		}
		result["service"] = serviceResp

		// Add proxy cache plugin
		serviceID := serviceResp["id"].(string)
		plugin := KongPlugin{
			Name: "proxy-cache",
			Config: map[string]interface{}{
				"response_code": []int{200, 301, 404},
				"request_method": []string{"GET", "HEAD"},
				"content_type":   []string{"application/json"},
				"cache_ttl":      300,
			},
		}
		plugin.Service.ID = serviceID
		pluginResp, err := h.callKongAdmin("POST", "/plugins", plugin)
		if err != nil {
			return nil, err
		}
		result["cache"] = pluginResp

	case "auth_overlay":
		// OAuth2/Key auth overlay
		service := KongService{
			Name: fmt.Sprintf("auth-%s", workflow.ID),
			URL:  fmt.Sprintf("http://backend:8080/api/webhooks/%s", workflow.ID),
		}
		serviceResp, err := h.callKongAdmin("POST", "/services", service)
		if err != nil {
			return nil, err
		}
		result["service"] = serviceResp

		// Add key-auth plugin
		serviceID := serviceResp["id"].(string)
		plugin := KongPlugin{
			Name: "key-auth",
			Config: map[string]interface{}{
				"key_names": []string{"apikey", "X-API-Key"},
			},
		}
		plugin.Service.ID = serviceID
		pluginResp, err := h.callKongAdmin("POST", "/plugins", plugin)
		if err != nil {
			return nil, err
		}
		result["auth"] = pluginResp

	case "monetization":
		// Usage tracking for billing
		service := KongService{
			Name: fmt.Sprintf("usage-%s", workflow.ID),
			URL:  fmt.Sprintf("http://backend:8080/api/webhooks/%s", workflow.ID),
		}
		serviceResp, err := h.callKongAdmin("POST", "/services", service)
		if err != nil {
			return nil, err
		}
		result["service"] = serviceResp

		// Add request size limiting and rate limiting for billing
		serviceID := serviceResp["id"].(string)
		
		// Rate limiting for usage tracking
		rateLimitPlugin := KongPlugin{
			Name: "rate-limiting",
			Config: map[string]interface{}{
				"minute": 60,
				"hour":   1000,
				"policy": "local",
			},
		}
		rateLimitPlugin.Service.ID = serviceID
		rateLimitResp, err := h.callKongAdmin("POST", "/plugins", rateLimitPlugin)
		if err != nil {
			return nil, err
		}
		result["rate_limiting"] = rateLimitResp

		// Request size limiting
		sizeLimitPlugin := KongPlugin{
			Name: "request-size-limiting",
			Config: map[string]interface{}{
				"allowed_payload_size": 1,
			},
		}
		sizeLimitPlugin.Service.ID = serviceID
		sizeLimitResp, err := h.callKongAdmin("POST", "/plugins", sizeLimitPlugin)
		if err != nil {
			return nil, err
		}
		result["size_limiting"] = sizeLimitResp
	}

	return result, nil
}

