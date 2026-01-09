package db

import (
	"time"

	"github.com/alexmacdonald/simple-ipass/internal/models"
)

// MockStore is a mock implementation of Store for testing
// This allows E2E tests to run without touching the filesystem
type MockStore struct {
	Users       map[string]*models.User
	Credentials map[string]*models.Credential
	Workflows   map[string]*models.Workflow
	Logs        []models.Log
}

// NewMockStore creates a new in-memory mock store
func NewMockStore() *MockStore {
	return &MockStore{
		Users:       make(map[string]*models.User),
		Credentials: make(map[string]*models.Credential),
		Workflows:   make(map[string]*models.Workflow),
		Logs:        make([]models.Log, 0),
	}
}

// User operations
func (m *MockStore) CreateUser(email, passwordHash string) (*models.User, error) {
	user := &models.User{
		ID:           "mock_user_" + email,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}
	m.Users[user.ID] = user
	return user, nil
}

func (m *MockStore) GetUserByEmail(email string) (*models.User, error) {
	for _, user := range m.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, ErrNotFound
}

func (m *MockStore) GetUserByID(id string) (*models.User, error) {
	if user, ok := m.Users[id]; ok {
		return user, nil
	}
	return nil, ErrNotFound
}

// Credential operations
func (m *MockStore) CreateCredential(userID, serviceName, apiKey string) (*models.Credential, error) {
	cred := &models.Credential{
		ID:           "mock_cred_" + serviceName,
		UserID:       userID,
		ServiceName:  serviceName,
		EncryptedKey: "encrypted_" + apiKey, // Mock encryption
		CreatedAt:    time.Now(),
	}
	m.Credentials[cred.ID] = cred
	return cred, nil
}

func (m *MockStore) GetCredentialsByUserID(userID string) ([]models.Credential, error) {
	var creds []models.Credential
	for _, cred := range m.Credentials {
		if cred.UserID == userID {
			creds = append(creds, *cred)
		}
	}
	return creds, nil
}

func (m *MockStore) GetCredentialByUserAndService(userID, serviceName string) (*models.Credential, error) {
	for _, cred := range m.Credentials {
		if cred.UserID == userID && cred.ServiceName == serviceName {
			// Mock decryption
			cred.DecryptedKey = "mock_webhook_url"
			return cred, nil
		}
	}
	return nil, ErrNotFound
}

// Workflow operations
func (m *MockStore) CreateWorkflow(userID, name, triggerType, actionType, configJSON string) (*models.Workflow, error) {
	workflow := &models.Workflow{
		ID:          "mock_wf_" + name,
		UserID:      userID,
		Name:        name,
		TriggerType: triggerType,
		ActionType:  actionType,
		ConfigJSON:  configJSON,
		IsActive:    true,
		CreatedAt:   time.Now(),
	}
	m.Workflows[workflow.ID] = workflow
	return workflow, nil
}

func (m *MockStore) GetWorkflowsByUserID(userID string) ([]models.Workflow, error) {
	var workflows []models.Workflow
	for _, wf := range m.Workflows {
		if wf.UserID == userID {
			workflows = append(workflows, *wf)
		}
	}
	return workflows, nil
}

func (m *MockStore) GetWorkflowByID(workflowID string) (*models.Workflow, error) {
	if wf, ok := m.Workflows[workflowID]; ok {
		return wf, nil
	}
	return nil, ErrNotFound
}

func (m *MockStore) UpdateWorkflowActive(workflowID string, isActive bool) error {
	if wf, ok := m.Workflows[workflowID]; ok {
		wf.IsActive = isActive
		return nil
	}
	return ErrNotFound
}

func (m *MockStore) UpdateWorkflowLastExecuted(workflowID string, executedAt time.Time) error {
	if wf, ok := m.Workflows[workflowID]; ok {
		wf.LastExecutedAt = &executedAt
		return nil
	}
	return ErrNotFound
}

func (m *MockStore) DeleteWorkflow(workflowID string) error {
	delete(m.Workflows, workflowID)
	return nil
}

func (m *MockStore) GetActiveScheduledWorkflows() ([]models.Workflow, error) {
	var workflows []models.Workflow
	for _, wf := range m.Workflows {
		if wf.TriggerType == "schedule" && wf.IsActive {
			workflows = append(workflows, *wf)
		}
	}
	return workflows, nil
}

// Log operations
func (m *MockStore) CreateLog(workflowID, status, message string) error {
	log := models.Log{
		ID:         "mock_log_" + workflowID,
		WorkflowID: workflowID,
		Status:     status,
		Message:    message,
		ExecutedAt: time.Now(),
	}
	m.Logs = append(m.Logs, log)
	return nil
}

func (m *MockStore) GetLogsByUserID(userID string) ([]models.WorkflowLog, error) {
	var logs []models.WorkflowLog
	for _, log := range m.Logs {
		// Find workflow to get user_id
		if wf, ok := m.Workflows[log.WorkflowID]; ok {
			if wf.UserID == userID {
				logs = append(logs, models.WorkflowLog{
					Log:          log,
					WorkflowName: wf.Name,
				})
			}
		}
	}
	return logs, nil
}

func (m *MockStore) GetLogsByWorkflowID(workflowID string) ([]models.Log, error) {
	var logs []models.Log
	for _, log := range m.Logs {
		if log.WorkflowID == workflowID {
			logs = append(logs, log)
		}
	}
	return logs, nil
}

// Lifecycle
func (m *MockStore) Close() error {
	// No-op for in-memory mock
	return nil
}

// Common errors
var (
	ErrNotFound = &StoreError{Code: "not_found", Message: "Resource not found"}
)

// StoreError represents a database error
type StoreError struct {
	Code    string
	Message string
}

func (e *StoreError) Error() string {
	return e.Message
}

