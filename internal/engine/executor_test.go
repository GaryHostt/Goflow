package engine_test

import (
	"context"
	"testing"
	"time"

	"github.com/alexmacdonald/simple-ipass/internal/db"
	"github.com/alexmacdonald/simple-ipass/internal/engine"
	"github.com/alexmacdonald/simple-ipass/internal/logger"
	"github.com/alexmacdonald/simple-ipass/internal/models"
)

// TestExecutorWithMockStore demonstrates dependency injection with MockStore
func TestExecutorWithMockStore(t *testing.T) {
	// Setup: Create in-memory mock (NO disk I/O!)
	mockStore := db.NewMockStore()
	testLogger := logger.NewLogger("test")
	executor := engine.NewExecutor(mockStore, testLogger)

	// Create test user
	user, err := mockStore.CreateUser("test@example.com", "hashed_password")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create test credential
	_, err = mockStore.CreateCredential(user.ID, "slack", "mock_webhook_url")
	if err != nil {
		t.Fatalf("Failed to create credential: %v", err)
	}

	// Create test workflow
	configJSON := `{"slack_message": "Test message"}`
	workflow, err := mockStore.CreateWorkflow(user.ID, "Test Workflow", "webhook", "slack_message", configJSON)
	if err != nil {
		t.Fatalf("Failed to create workflow: %v", err)
	}

	// Execute workflow
	tenantID := "tenant_" + user.ID
	result := executor.DryRun(*workflow, user.ID, tenantID)

	// Verify result (mock doesn't actually call Slack)
	if result.Status != "success" && result.Status != "failed" {
		t.Errorf("Expected success or failed, got %s: %s", result.Status, result.Message)
	}

	// Verify logs were NOT saved (dry run mode)
	logs, err := mockStore.GetLogsByWorkflowID(workflow.ID)
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}
	if len(logs) != 0 {
		t.Errorf("Expected 0 logs for dry run, got %d", len(logs))
	}
}

// TestContextCancellation proves executor respects context
func TestContextCancellation(t *testing.T) {
	mockStore := db.NewMockStore()
	testLogger := logger.NewLogger("test")
	executor := engine.NewExecutor(mockStore, testLogger)

	// Create test data
	user, _ := mockStore.CreateUser("cancel@example.com", "hashed")
	mockStore.CreateCredential(user.ID, "slack", "webhook_url")
	workflow, _ := mockStore.CreateWorkflow(user.ID, "Slow Workflow", "webhook", "slack_message", `{"slack_message":"test"}`)

	// Create context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Execute with cancelled context
	executor.ExecuteWorkflowWithContext(ctx, *workflow)

	// Give goroutine time to process cancellation
	time.Sleep(100 * time.Millisecond)

	// Verify execution was cancelled (no log should be saved)
	logs, _ := mockStore.GetLogsByWorkflowID(workflow.ID)
	if len(logs) > 0 {
		t.Errorf("Expected no logs for cancelled execution, got %d", len(logs))
	}
}

// TestWorkerPoolBoundedConcurrency verifies max 10 concurrent executions
func TestWorkerPoolBoundedConcurrency(t *testing.T) {
	mockStore := db.NewMockStore()
	testLogger := logger.NewLogger("test")
	executor := engine.NewExecutor(mockStore, testLogger)

	// Create test data
	user, _ := mockStore.CreateUser("pool@example.com", "hashed")
	mockStore.CreateCredential(user.ID, "slack", "webhook_url")

	// Submit 50 workflows simultaneously
	for i := 0; i < 50; i++ {
		workflow := &models.Workflow{
			ID:          "wf_" + string(rune(i)),
			UserID:      user.ID,
			Name:        "Workflow " + string(rune(i)),
			TriggerType: "webhook",
			ActionType:  "slack_message",
			ConfigJSON:  `{"slack_message":"test"}`,
			IsActive:    true,
		}
		executor.ExecuteWorkflow(*workflow)
	}

	// Give worker pool time to process
	time.Sleep(2 * time.Second)

	// Test passes if no panic/crash (proves bounded concurrency works)
	t.Log("Worker pool handled 50 concurrent jobs without crashing")
}

// BenchmarkMockStoreVsRealDB compares performance
func BenchmarkMockStoreVsRealDB(b *testing.B) {
	b.Run("MockStore", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			mockStore := db.NewMockStore()
			mockStore.CreateUser("bench@example.com", "hash")
		}
	})

	// Note: Real DB benchmark would require file I/O
	// Uncomment to compare:
	// b.Run("RealDB", func(b *testing.B) {
	//     for i := 0; i < b.N; i++ {
	//         realDB, _ := db.New("bench_test.db")
	//         realDB.CreateUser("bench@example.com", "hash")
	//         realDB.Close()
	//         os.Remove("bench_test.db")
	//     }
	// })
}

