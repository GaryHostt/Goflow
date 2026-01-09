package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/alexmacdonald/simple-ipass/internal/db"
	"github.com/alexmacdonald/simple-ipass/internal/middleware"
)

type CredentialsHandler struct {
	db *db.Database
}

func NewCredentialsHandler(database *db.Database) *CredentialsHandler {
	return &CredentialsHandler{db: database}
}

type CreateCredentialRequest struct {
	ServiceName string `json:"service_name"`
	APIKey      string `json:"api_key"`
}

// CreateCredential saves encrypted API keys/webhooks
func (h *CredentialsHandler) CreateCredential(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ServiceName == "" || req.APIKey == "" {
		http.Error(w, "service_name and api_key are required", http.StatusBadRequest)
		return
	}

	// Create credential with encryption
	cred, err := h.db.CreateCredential(userID, req.ServiceName, req.APIKey)
	if err != nil {
		http.Error(w, "Failed to save credential", http.StatusInternalServerError)
		return
	}

	// Don't return the encrypted key
	cred.EncryptedKey = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cred)
}

// GetCredentials lists user's connections (without exposing keys)
func (h *CredentialsHandler) GetCredentials(w http.ResponseWriter, r *http.Request) {
	// TODO: MULTI-TENANT - Filter by tenant_id instead of user_id
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	creds, err := h.db.GetCredentialsByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to fetch credentials", http.StatusInternalServerError)
		return
	}

	// Remove encrypted keys from response
	for i := range creds {
		creds[i].EncryptedKey = ""
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(creds)
}

