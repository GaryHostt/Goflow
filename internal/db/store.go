package db

import (
	"time"

	"github.com/alexmacdonald/simple-ipass/internal/models"
)

// Store defines the interface for data persistence
// This allows for easy testing with mocks and potential database swaps
type Store interface {
	// User operations
	CreateUser(email, passwordHash string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)

	// Credential operations
	CreateCredential(userID, serviceName, apiKey string) (*models.Credential, error)
	GetCredentialsByUserID(userID string) ([]models.Credential, error)
	GetCredentialByUserAndService(userID, serviceName string) (*models.Credential, error)

	// Workflow operations
	CreateWorkflow(userID, name, triggerType, actionType, configJSON string) (*models.Workflow, error)
	GetWorkflowsByUserID(userID string) ([]models.Workflow, error)
	GetWorkflowByID(workflowID string) (*models.Workflow, error)
	UpdateWorkflowActive(workflowID string, isActive bool) error
	UpdateWorkflowLastExecuted(workflowID string, executedAt time.Time) error
	DeleteWorkflow(workflowID string) error
	GetActiveScheduledWorkflows() ([]models.Workflow, error)

	// Log operations
	CreateLog(workflowID, status, message string) error
	GetLogsByUserID(userID string) ([]models.WorkflowLog, error)
	GetLogsByWorkflowID(workflowID string) ([]models.Log, error)

	// Lifecycle
	Close() error
}

// Ensure Database implements Store interface
var _ Store = (*Database)(nil)

