package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/alexmacdonald/simple-ipass/internal/models"
	"github.com/alexmacdonald/simple-ipass/internal/crypto"
	"github.com/google/uuid"
)

type Database struct {
	conn *sql.DB
}

// New creates a new database connection and initializes schema
func New(dbPath string) (*Database, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys
	if _, err := conn.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	db := &Database{conn: conn}

	// Initialize schema
	if err := db.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return db, nil
}

// initSchema creates tables from schema.sql
func (db *Database) initSchema() error {
	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema.sql: %w", err)
	}

	if _, err := db.conn.Exec(string(schema)); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}

// Close closes the database connection
func (db *Database) Close() error {
	return db.conn.Close()
}

// Ping verifies the database connection is alive and working
// Returns an error if the database is not accessible
func (db *Database) Ping() error {
	return db.conn.Ping()
}

// --- User Repository ---

// CreateUser creates a new user
func (db *Database) CreateUser(email, passwordHash string) (*models.User, error) {
	user := &models.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}

	query := `INSERT INTO users (id, email, password_hash, created_at) VALUES (?, ?, ?, ?)`
	_, err := db.conn.Exec(query, user.ID, user.Email, user.PasswordHash, user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (db *Database) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = ?`
	err := db.conn.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID retrieves a user by ID
func (db *Database) GetUserByID(id string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password_hash, created_at FROM users WHERE id = ?`
	err := db.conn.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// --- Credentials Repository ---
// TODO: MULTI-TENANT - Change user_id filter to tenant_id

// CreateCredential creates a new credential
func (db *Database) CreateCredential(userID, serviceName, apiKey string) (*models.Credential, error) {
	encryptedKey, err := crypto.Encrypt(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt key: %w", err)
	}

	cred := &models.Credential{
		ID:           uuid.New().String(),
		UserID:       userID,
		ServiceName:  serviceName,
		EncryptedKey: encryptedKey,
		CreatedAt:    time.Now(),
	}

	query := `INSERT INTO credentials (id, user_id, service_name, encrypted_key, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err = db.conn.Exec(query, cred.ID, cred.UserID, cred.ServiceName, cred.EncryptedKey, cred.CreatedAt)
	if err != nil {
		return nil, err
	}

	return cred, nil
}

// GetCredentialsByUserID retrieves all credentials for a user
func (db *Database) GetCredentialsByUserID(userID string) ([]models.Credential, error) {
	query := `SELECT id, user_id, service_name, encrypted_key, created_at FROM credentials WHERE user_id = ?`
	rows, err := db.conn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credentials []models.Credential
	for rows.Next() {
		var cred models.Credential
		err := rows.Scan(&cred.ID, &cred.UserID, &cred.ServiceName, &cred.EncryptedKey, &cred.CreatedAt)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, cred)
	}

	return credentials, nil
}

// GetCredentialByUserAndService retrieves a specific credential
func (db *Database) GetCredentialByUserAndService(userID, serviceName string) (*models.Credential, error) {
	cred := &models.Credential{}
	query := `SELECT id, user_id, service_name, encrypted_key, created_at FROM credentials WHERE user_id = ? AND service_name = ?`
	err := db.conn.QueryRow(query, userID, serviceName).Scan(&cred.ID, &cred.UserID, &cred.ServiceName, &cred.EncryptedKey, &cred.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Decrypt the key
	decryptedKey, err := crypto.Decrypt(cred.EncryptedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt key: %w", err)
	}
	cred.DecryptedKey = decryptedKey

	return cred, nil
}

// --- Workflows Repository ---
// TODO: MULTI-TENANT - Change user_id filter to tenant_id

// CreateWorkflow creates a new workflow with optional action chain
func (db *Database) CreateWorkflow(userID, name, triggerType, actionType, configJSON string) (*models.Workflow, error) {
	workflow := &models.Workflow{
		ID:          uuid.New().String(),
		UserID:      userID,
		Name:        name,
		TriggerType: triggerType,
		ActionType:  actionType,
		ConfigJSON:  configJSON,
		ActionChain: "", // Will be set separately if needed
		IsActive:    true,
		CreatedAt:   time.Now(),
	}

	query := `INSERT INTO workflows (id, user_id, name, trigger_type, action_type, config_json, action_chain, is_active, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.conn.Exec(query, workflow.ID, workflow.UserID, workflow.Name, workflow.TriggerType, workflow.ActionType, workflow.ConfigJSON, workflow.ActionChain, workflow.IsActive, workflow.CreatedAt)
	if err != nil {
		return nil, err
	}

	return workflow, nil
}

// CreateWorkflowWithChain creates a new workflow with an action chain
func (db *Database) CreateWorkflowWithChain(userID, name, triggerType, actionType, configJSON, actionChain string) (*models.Workflow, error) {
	workflow := &models.Workflow{
		ID:          uuid.New().String(),
		UserID:      userID,
		Name:        name,
		TriggerType: triggerType,
		ActionType:  actionType,
		ConfigJSON:  configJSON,
		ActionChain: actionChain,
		IsActive:    true,
		CreatedAt:   time.Now(),
	}

	query := `INSERT INTO workflows (id, user_id, name, trigger_type, action_type, config_json, action_chain, is_active, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.conn.Exec(query, workflow.ID, workflow.UserID, workflow.Name, workflow.TriggerType, workflow.ActionType, workflow.ConfigJSON, workflow.ActionChain, workflow.IsActive, workflow.CreatedAt)
	if err != nil {
		return nil, err
	}

	return workflow, nil
}

// GetWorkflowsByUserID retrieves all workflows for a user
func (db *Database) GetWorkflowsByUserID(userID string) ([]models.Workflow, error) {
	query := `SELECT id, user_id, name, trigger_type, action_type, config_json, action_chain, is_active, last_executed_at, created_at FROM workflows WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := db.conn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []models.Workflow
	for rows.Next() {
		var w models.Workflow
		var lastExecutedAt sql.NullTime
		var actionChain sql.NullString
		err := rows.Scan(&w.ID, &w.UserID, &w.Name, &w.TriggerType, &w.ActionType, &w.ConfigJSON, &actionChain, &w.IsActive, &lastExecutedAt, &w.CreatedAt)
		if err != nil {
			return nil, err
		}
		if lastExecutedAt.Valid {
			w.LastExecutedAt = &lastExecutedAt.Time
		}
		if actionChain.Valid {
			w.ActionChain = actionChain.String
		}
		workflows = append(workflows, w)
	}

	return workflows, nil
}

// GetWorkflowByID retrieves a workflow by ID
func (db *Database) GetWorkflowByID(workflowID string) (*models.Workflow, error) {
	w := &models.Workflow{}
	var lastExecutedAt sql.NullTime
	var actionChain sql.NullString
	query := `SELECT id, user_id, name, trigger_type, action_type, config_json, action_chain, is_active, last_executed_at, created_at FROM workflows WHERE id = ?`
	err := db.conn.QueryRow(query, workflowID).Scan(&w.ID, &w.UserID, &w.Name, &w.TriggerType, &w.ActionType, &w.ConfigJSON, &actionChain, &w.IsActive, &lastExecutedAt, &w.CreatedAt)
	if err != nil {
		return nil, err
	}
	if lastExecutedAt.Valid {
		w.LastExecutedAt = &lastExecutedAt.Time
	}
	if actionChain.Valid {
		w.ActionChain = actionChain.String
	}
	return w, nil
}

// UpdateWorkflowActive toggles workflow active status
func (db *Database) UpdateWorkflowActive(workflowID string, isActive bool) error {
	query := `UPDATE workflows SET is_active = ? WHERE id = ?`
	_, err := db.conn.Exec(query, isActive, workflowID)
	return err
}

// UpdateWorkflowLastExecuted updates the last execution time
func (db *Database) UpdateWorkflowLastExecuted(workflowID string, executedAt time.Time) error {
	query := `UPDATE workflows SET last_executed_at = ? WHERE id = ?`
	_, err := db.conn.Exec(query, executedAt, workflowID)
	return err
}

// DeleteWorkflow deletes a workflow
func (db *Database) DeleteWorkflow(workflowID string) error {
	query := `DELETE FROM workflows WHERE id = ?`
	_, err := db.conn.Exec(query, workflowID)
	return err
}

// GetActiveScheduledWorkflows retrieves all active scheduled workflows
func (db *Database) GetActiveScheduledWorkflows() ([]models.Workflow, error) {
	query := `SELECT id, user_id, name, trigger_type, action_type, config_json, action_chain, is_active, last_executed_at, created_at 
	          FROM workflows WHERE trigger_type = 'schedule' AND is_active = 1`
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []models.Workflow
	for rows.Next() {
		var w models.Workflow
		var lastExecutedAt sql.NullTime
		var actionChain sql.NullString
		err := rows.Scan(&w.ID, &w.UserID, &w.Name, &w.TriggerType, &w.ActionType, &w.ConfigJSON, &actionChain, &w.IsActive, &lastExecutedAt, &w.CreatedAt)
		if err != nil {
			return nil, err
		}
		if lastExecutedAt.Valid {
			w.LastExecutedAt = &lastExecutedAt.Time
		}
		if actionChain.Valid {
			w.ActionChain = actionChain.String
		}
		workflows = append(workflows, w)
	}

	return workflows, nil
}

// --- Logs Repository ---
// TODO: MULTI-TENANT - Join with workflows to filter by tenant_id

// CreateLog creates a new execution log
func (db *Database) CreateLog(workflowID, status, message string) error {
	log := &models.Log{
		ID:         uuid.New().String(),
		WorkflowID: workflowID,
		Status:     status,
		Message:    message,
		ExecutedAt: time.Now(),
	}

	query := `INSERT INTO logs (id, workflow_id, status, message, executed_at) VALUES (?, ?, ?, ?, ?)`
	_, err := db.conn.Exec(query, log.ID, log.WorkflowID, log.Status, log.Message, log.ExecutedAt)
	return err
}

// GetLogsByUserID retrieves all logs for a user's workflows
func (db *Database) GetLogsByUserID(userID string) ([]models.WorkflowLog, error) {
	query := `SELECT l.id, l.workflow_id, l.status, l.message, l.executed_at, w.name 
	          FROM logs l 
	          JOIN workflows w ON l.workflow_id = w.id 
	          WHERE w.user_id = ? 
	          ORDER BY l.executed_at DESC 
	          LIMIT 100`
	rows, err := db.conn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.WorkflowLog
	for rows.Next() {
		var log models.WorkflowLog
		err := rows.Scan(&log.ID, &log.WorkflowID, &log.Status, &log.Message, &log.ExecutedAt, &log.WorkflowName)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// GetLogsByWorkflowID retrieves logs for a specific workflow
func (db *Database) GetLogsByWorkflowID(workflowID string) ([]models.Log, error) {
	query := `SELECT id, workflow_id, status, message, executed_at FROM logs WHERE workflow_id = ? ORDER BY executed_at DESC LIMIT 50`
	rows, err := db.conn.Query(query, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.Log
	for rows.Next() {
		var log models.Log
		err := rows.Scan(&log.ID, &log.WorkflowID, &log.Status, &log.Message, &log.ExecutedAt)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

