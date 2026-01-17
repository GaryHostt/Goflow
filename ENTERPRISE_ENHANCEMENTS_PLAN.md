# 🚀 GoFlow v0.8.0 - Enterprise Enhancements

## Overview

Based on architectural review, implementing 4 key improvements to elevate GoFlow to enterprise-scale:

1. **Context Propagation** - Complete context lifecycle management
2. **Plugin System** - Extensible connector architecture
3. **Enhanced Logging** - Production-grade observability
4. **Runtime Parameters** - Dynamic workflow execution

---

## 1. Context Propagation Enhancement

### **Current State** ✅
- Context passed to executor
- Timeout management in place
- Graceful cancellation supported

### **Improvements Needed**
- [ ] Propagate context to ALL connector calls
- [ ] Add timeout configuration per workflow
- [ ] Implement context-aware database operations
- [ ] Add context to scheduler operations

### **Implementation**

#### **A. Workflow-Level Timeouts**

Add timeout configuration to workflows:

```go
// internal/models/models.go
type Workflow struct {
    ID              string    `json:"id"`
    UserID          string    `json:"user_id"`
    Name            string    `json:"name"`
    TriggerType     string    `json:"trigger_type"`
    ActionType      string    `json:"action_type"`
    ConfigJSON      string    `json:"config_json"`
    IsActive        bool      `json:"is_active"`
    LastExecutedAt  time.Time `json:"last_executed_at"`
    ActionChain     []ChainedAction `json:"action_chain,omitempty"`
    ExecutionTimeout int      `json:"execution_timeout"` // NEW: seconds, default 300
    CreatedAt       time.Time `json:"created_at"`
}
```

#### **B. Context-Aware Database Operations**

```go
// internal/db/store.go - Update interface
type Store interface {
    // Context-aware operations
    CreateWorkflowWithContext(ctx context.Context, wf models.Workflow) error
    GetWorkflowByIDWithContext(ctx context.Context, id string) (models.Workflow, error)
    // ... all operations get context parameter
}
```

#### **C. Propagate Through Execution Chain**

```go
// internal/engine/executor.go
func (e *Executor) ExecuteWorkflow(ctx context.Context, workflow models.Workflow) error {
    // Create workflow-specific context with timeout
    timeout := time.Duration(workflow.ExecutionTimeout) * time.Second
    if timeout == 0 {
        timeout = 300 * time.Second // default 5 minutes
    }
    
    workflowCtx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()
    
    // Pass context to all operations
    return e.executeWithContext(workflowCtx, workflow)
}
```

---

## 2. Plugin System

### **Architecture: Go Plugin Interface**

Allow users to add custom connectors without rebuilding GoFlow.

#### **A. Plugin Interface Definition**

```go
// internal/plugins/interface.go
package plugins

import (
    "context"
    "github.com/alexmacdonald/simple-ipass/internal/models"
)

// ConnectorPlugin defines the interface for custom connectors
type ConnectorPlugin interface {
    // Name returns the unique identifier for this connector
    Name() string
    
    // Execute runs the connector action
    Execute(ctx context.Context, config map[string]interface{}, data map[string]interface{}) (map[string]interface{}, error)
    
    // Validate checks if the configuration is valid
    Validate(config map[string]interface{}) error
    
    // Metadata returns information about the connector
    Metadata() PluginMetadata
}

type PluginMetadata struct {
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Version     string   `json:"version"`
    Author      string   `json:"author"`
    ConfigFields []Field `json:"config_fields"`
}

type Field struct {
    Name        string `json:"name"`
    Type        string `json:"type"` // string, int, bool, etc.
    Required    bool   `json:"required"`
    Description string `json:"description"`
}
```

#### **B. Plugin Loader**

```go
// internal/plugins/loader.go
package plugins

import (
    "fmt"
    "plugin"
    "sync"
)

type PluginManager struct {
    plugins map[string]ConnectorPlugin
    mu      sync.RWMutex
}

func NewPluginManager() *PluginManager {
    return &PluginManager{
        plugins: make(map[string]ConnectorPlugin),
    }
}

func (pm *PluginManager) LoadPlugin(path string) error {
    pm.mu.Lock()
    defer pm.mu.Unlock()
    
    // Load the plugin
    p, err := plugin.Open(path)
    if err != nil {
        return fmt.Errorf("failed to open plugin: %w", err)
    }
    
    // Look for the NewConnector symbol
    symbol, err := p.Lookup("NewConnector")
    if err != nil {
        return fmt.Errorf("plugin must export NewConnector function: %w", err)
    }
    
    // Assert it implements ConnectorPlugin
    newConnector, ok := symbol.(func() ConnectorPlugin)
    if !ok {
        return fmt.Errorf("NewConnector has wrong signature")
    }
    
    connector := newConnector()
    pm.plugins[connector.Name()] = connector
    
    return nil
}

func (pm *PluginManager) GetPlugin(name string) (ConnectorPlugin, error) {
    pm.mu.RLock()
    defer pm.mu.RUnlock()
    
    plugin, exists := pm.plugins[name]
    if !exists {
        return nil, fmt.Errorf("plugin %s not found", name)
    }
    
    return plugin, nil
}

func (pm *PluginManager) ListPlugins() []PluginMetadata {
    pm.mu.RLock()
    defer pm.mu.RUnlock()
    
    metadata := make([]PluginMetadata, 0, len(pm.plugins))
    for _, p := range pm.plugins {
        metadata = append(metadata, p.Metadata())
    }
    
    return metadata
}
```

#### **C. Example Plugin**

```go
// examples/custom-connector/main.go
package main

import (
    "context"
    "fmt"
    "github.com/alexmacdonald/simple-ipass/internal/plugins"
)

type CustomConnector struct{}

func (c *CustomConnector) Name() string {
    return "custom_api"
}

func (c *CustomConnector) Execute(ctx context.Context, config map[string]interface{}, data map[string]interface{}) (map[string]interface{}, error) {
    // Custom logic here
    apiKey := config["api_key"].(string)
    endpoint := config["endpoint"].(string)
    
    // Make API call, process data, etc.
    
    return map[string]interface{}{
        "status": "success",
        "data": "custom result",
    }, nil
}

func (c *CustomConnector) Validate(config map[string]interface{}) error {
    if _, ok := config["api_key"]; !ok {
        return fmt.Errorf("api_key is required")
    }
    return nil
}

func (c *CustomConnector) Metadata() plugins.PluginMetadata {
    return plugins.PluginMetadata{
        Name:        "Custom API Connector",
        Description: "Connects to a custom third-party API",
        Version:     "1.0.0",
        Author:      "Your Name",
        ConfigFields: []plugins.Field{
            {Name: "api_key", Type: "string", Required: true, Description: "API Key for authentication"},
            {Name: "endpoint", Type: "string", Required: true, Description: "API endpoint URL"},
        },
    }
}

// NewConnector is the entry point for the plugin
func NewConnector() plugins.ConnectorPlugin {
    return &CustomConnector{}
}
```

**Build plugin:**
```bash
go build -buildmode=plugin -o custom_api.so examples/custom-connector/main.go
```

#### **D. Integration with Executor**

```go
// internal/engine/executor.go - Add plugin support
func (e *Executor) ExecutePluginAction(ctx context.Context, pluginName string, config, data map[string]interface{}) (map[string]interface{}, error) {
    plugin, err := e.pluginManager.GetPlugin(pluginName)
    if err != nil {
        return nil, err
    }
    
    // Validate config
    if err := plugin.Validate(config); err != nil {
        return nil, fmt.Errorf("invalid config: %w", err)
    }
    
    // Execute with context
    return plugin.Execute(ctx, config, data)
}
```

---

## 3. Enhanced Logging

### **Current State** ✅
- Using `slog` for structured logging
- JSON output for ELK
- Context logging (user_id, tenant_id)

### **Improvements**

#### **A. Log Levels Per Component**

```go
// internal/logger/logger.go - Enhanced
type Logger struct {
    logger    *slog.Logger
    component string
    level     slog.Level
}

func (l *Logger) WithComponent(component string) *Logger {
    return &Logger{
        logger:    l.logger,
        component: component,
        level:     l.level,
    }
}

func (l *Logger) WithLevel(level slog.Level) *Logger {
    return &Logger{
        logger:    l.logger,
        component: l.component,
        level:     level,
    }
}
```

#### **B. Distributed Tracing Support**

```go
// internal/logger/tracing.go
package logger

import (
    "context"
    "github.com/google/uuid"
)

type traceIDKey struct{}

func WithTraceID(ctx context.Context) context.Context {
    return context.WithValue(ctx, traceIDKey{}, uuid.New().String())
}

func GetTraceID(ctx context.Context) string {
    if traceID, ok := ctx.Value(traceIDKey{}).(string); ok {
        return traceID
    }
    return ""
}

func (l *Logger) LogWithTrace(ctx context.Context, level slog.Level, msg string, args map[string]interface{}) {
    if args == nil {
        args = make(map[string]interface{})
    }
    
    // Add trace ID
    if traceID := GetTraceID(ctx); traceID != "" {
        args["trace_id"] = traceID
    }
    
    // Add component
    if l.component != "" {
        args["component"] = l.component
    }
    
    l.logger.Log(ctx, level, msg, slogArgs(args)...)
}
```

#### **C. Performance Metrics**

```go
// internal/logger/metrics.go
package logger

import (
    "context"
    "time"
)

func (l *Logger) LogDuration(ctx context.Context, operation string, start time.Time, err error) {
    duration := time.Since(start)
    
    args := map[string]interface{}{
        "operation":    operation,
        "duration_ms":  duration.Milliseconds(),
        "duration_sec": duration.Seconds(),
    }
    
    if err != nil {
        args["error"] = err.Error()
        l.LogWithTrace(ctx, slog.LevelError, "Operation failed", args)
    } else {
        l.LogWithTrace(ctx, slog.LevelInfo, "Operation completed", args)
    }
}

// Usage
func (e *Executor) ExecuteWorkflow(ctx context.Context, workflow models.Workflow) error {
    start := time.Now()
    defer e.logger.LogDuration(ctx, "execute_workflow", start, nil)
    
    // ... execution logic
}
```

---

## 4. Runtime Parameters

### **Feature: Parameterized Workflow Execution**

Allow users to trigger workflows with custom JSON inputs.

#### **A. Add Parameters to Workflow Model**

```go
// internal/models/models.go
type Workflow struct {
    // ... existing fields
    Parameters []WorkflowParameter `json:"parameters,omitempty"` // NEW
}

type WorkflowParameter struct {
    Name         string `json:"name"`
    Type         string `json:"type"` // string, number, boolean, object
    Required     bool   `json:"required"`
    DefaultValue string `json:"default_value,omitempty"`
    Description  string `json:"description,omitempty"`
}
```

#### **B. Trigger with Parameters API**

```go
// internal/handlers/workflows.go - New endpoint
func (h *WorkflowsHandler) TriggerWorkflowWithParams(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    workflowID := vars["id"]
    
    // Parse parameters
    var req struct {
        Parameters map[string]interface{} `json:"parameters"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // Get workflow
    workflow, err := h.store.GetWorkflowByID(workflowID)
    if err != nil {
        http.Error(w, "Workflow not found", http.StatusNotFound)
        return
    }
    
    // Validate parameters
    if err := validateParameters(workflow.Parameters, req.Parameters); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Execute with parameters
    ctx := r.Context()
    go h.executor.ExecuteWorkflowWithParams(ctx, workflow, req.Parameters)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status": "triggered",
        "workflow_id": workflowID,
        "parameters": req.Parameters,
    })
}

func validateParameters(defined []models.WorkflowParameter, provided map[string]interface{}) error {
    for _, param := range defined {
        if param.Required {
            if _, exists := provided[param.Name]; !exists {
                return fmt.Errorf("required parameter %s is missing", param.Name)
            }
        }
    }
    return nil
}
```

#### **C. Use Parameters in Execution**

```go
// internal/engine/executor.go
func (e *Executor) ExecuteWorkflowWithParams(ctx context.Context, workflow models.Workflow, params map[string]interface{}) error {
    // Merge parameters into trigger payload
    triggerPayload := map[string]interface{}{
        "runtime_params": params,
        "workflow_id":    workflow.ID,
        "triggered_at":   time.Now(),
    }
    
    // Use template engine to substitute parameters
    config := workflow.Config
    for key, value := range params {
        // Replace {{params.key}} in config
        placeholder := fmt.Sprintf("{{params.%s}}", key)
        config = strings.ReplaceAll(config, placeholder, fmt.Sprint(value))
    }
    
    // Execute with modified config
    return e.executeWorkflowInternal(ctx, workflow, config, triggerPayload)
}
```

#### **D. Frontend Support**

```typescript
// frontend/lib/api.ts - Add parameter support
export const workflows = {
  // ... existing methods
  
  triggerWithParams: async (id: string, parameters: Record<string, any>) => {
    const response = await apiClient(`/workflows/${id}/trigger`, {
      method: 'POST',
      body: JSON.stringify({ parameters }),
    })
    return response.json()
  },
}
```

---

## Implementation Priority

### **Phase 1: Quick Wins** (1-2 days)
1. ✅ Context Propagation - Extend existing implementation
2. ✅ Enhanced Logging - Add trace IDs and metrics
3. 🆕 Runtime Parameters - Basic implementation

### **Phase 2: Advanced Features** (3-5 days)
4. 🆕 Plugin System - Full implementation with examples

---

## Benefits

### **Context Propagation**
- ✅ Better timeout management
- ✅ Graceful cancellation of long tasks
- ✅ Resource cleanup guarantees
- ✅ Prevents orphaned goroutines

### **Plugin System**
- ✅ Add connectors without rebuilding
- ✅ Community contributions
- ✅ Customer-specific integrations
- ✅ Versioned plugin updates

### **Enhanced Logging**
- ✅ Distributed tracing (trace_id)
- ✅ Performance metrics per operation
- ✅ Better ELK/Loki integration
- ✅ Component-level log filtering

### **Runtime Parameters**
- ✅ Dynamic workflow execution
- ✅ Reusable workflow templates
- ✅ API-driven automation
- ✅ A/B testing workflows

---

## Migration Guide

### **Existing Workflows**
- ✅ **Backward compatible** - no changes needed
- ✅ New features are opt-in
- ✅ Default values preserve current behavior

### **Custom Code**
- ⚠️  Database interface updated (adds context)
- ⚠️  Executor signatures changed (adds context)
- ✅ Easy migration: just add `context.Background()` initially

---

## Testing Strategy

### **Unit Tests**
- Test parameter validation
- Test plugin loading
- Test context cancellation
- Test log output format

### **Integration Tests**
- Test plugin execution
- Test parameterized workflows
- Test context timeout behavior
- Test trace ID propagation

### **Load Tests**
- Test plugin performance
- Test context overhead
- Test concurrent parameter execution

---

## Documentation

Will create:
- `PLUGIN_SYSTEM_GUIDE.md` - How to build plugins
- `RUNTIME_PARAMETERS_GUIDE.md` - How to use parameters
- `CONTEXT_BEST_PRACTICES.md` - Context usage patterns
- `ENHANCED_LOGGING_GUIDE.md` - Logging strategies

---

## Next Steps

1. **Review this proposal**
2. **Prioritize features**
3. **Start implementation**
4. **Test thoroughly**
5. **Document for users**

Would you like me to proceed with implementing any of these features?

