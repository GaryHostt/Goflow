package handlers

import (
	"encoding/json"
	"net/http"
)

// JSONResponse is a standardized API response envelope
// Provides consistent structure for all API responses
type JSONResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *MetaData   `json:"meta,omitempty"`
}

// MetaData provides additional response metadata
type MetaData struct {
	RequestID string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Version   string `json:"version,omitempty"`
}

// SendJSON sends a standardized JSON response
func SendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	response := JSONResponse{
		Success: status >= 200 && status < 300,
		Data:    data,
	}
	
	json.NewEncoder(w).Encode(response)
}

// SendError sends a standardized error response
func SendError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	response := JSONResponse{
		Success: false,
		Error:   message,
	}
	
	json.NewEncoder(w).Encode(response)
}

// SendSuccess sends a standardized success response with data
func SendSuccess(w http.ResponseWriter, data interface{}) {
	SendJSON(w, http.StatusOK, data)
}

// SendCreated sends a 201 Created response
func SendCreated(w http.ResponseWriter, data interface{}) {
	SendJSON(w, http.StatusCreated, data)
}

// SendNoContent sends a 204 No Content response
func SendNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// SendBadRequest sends a 400 Bad Request error
func SendBadRequest(w http.ResponseWriter, message string) {
	SendError(w, http.StatusBadRequest, message)
}

// SendUnauthorized sends a 401 Unauthorized error
func SendUnauthorized(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	SendError(w, http.StatusUnauthorized, message)
}

// SendForbidden sends a 403 Forbidden error
func SendForbidden(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Forbidden"
	}
	SendError(w, http.StatusForbidden, message)
}

// SendNotFound sends a 404 Not Found error
func SendNotFound(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Resource not found"
	}
	SendError(w, http.StatusNotFound, message)
}

// SendInternalError sends a 500 Internal Server Error
func SendInternalError(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Internal server error"
	}
	SendError(w, http.StatusInternalServerError, message)
}

// SendValidationError sends a 422 Unprocessable Entity error
func SendValidationError(w http.ResponseWriter, message string) {
	SendError(w, http.StatusUnprocessableEntity, message)
}

