package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexmacdonald/simple-ipass/internal/db"
	"github.com/alexmacdonald/simple-ipass/internal/engine"
	"github.com/alexmacdonald/simple-ipass/internal/handlers"
	"github.com/alexmacdonald/simple-ipass/internal/logger"
	"github.com/alexmacdonald/simple-ipass/internal/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Initialize structured logger (ELK-ready!)
	appLogger := logger.NewLogger("ipaas-api")
	appLogger.Info("Starting GoFlow API Server...", map[string]interface{}{
		"version": "0.4.0",
		"env":     getEnv("ENVIRONMENT", "development"),
	})

	// Initialize database with retry logic for Docker/production environments
	database, err := initializeDatabaseWithRetry(appLogger, 10, 2*time.Second)
	if err != nil {
		appLogger.Error("Failed to initialize database after retries", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	appLogger.Info("Database initialized successfully", nil)

	// Initialize executor with logger
	executor := engine.NewExecutor(database, appLogger)

	// Initialize scheduler with logger (tenant-aware ready!)
	scheduler := engine.NewScheduler(database, executor, appLogger)
	scheduler.Start(60 * time.Second) // Check every 60 seconds
	defer scheduler.Stop()

	// Setup router
	router := mux.NewRouter()

	// Add request logging middleware (tracks all HTTP requests with status codes & timing)
	router.Use(middleware.RequestLogger(appLogger))

	// Public routes
	authHandler := handlers.NewAuthHandler(database)
	router.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")
	
	// Dev mode endpoint (only enable in development)
	if getEnv("ENVIRONMENT", "development") == "development" {
		router.HandleFunc("/api/auth/dev-login", authHandler.DevLogin).Methods("POST")
		appLogger.Info("Dev mode enabled - /api/auth/dev-login endpoint available", nil)
	}

	// Webhook handler (public but workflow-specific)
	webhookHandler := handlers.NewWebhookHandler(database, executor)
	router.HandleFunc("/api/webhooks/{id}", webhookHandler.TriggerWebhook).Methods("POST")

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","version":"0.2.0"}`))
	}).Methods("GET")

	// Protected routes with tenant-aware middleware
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware(appLogger)) // Now logs user_id AND tenant_id!

	// Credentials routes
	credentialsHandler := handlers.NewCredentialsHandler(database)
	api.HandleFunc("/credentials", credentialsHandler.CreateCredential).Methods("POST")
	api.HandleFunc("/credentials", credentialsHandler.GetCredentials).Methods("GET")

	// Workflows routes
	workflowsHandler := handlers.NewWorkflowsHandler(database, executor)
	api.HandleFunc("/workflows", workflowsHandler.CreateWorkflow).Methods("POST")
	api.HandleFunc("/workflows", workflowsHandler.GetWorkflows).Methods("GET")
	api.HandleFunc("/workflows/dry-run", workflowsHandler.DryRunWorkflow).Methods("POST") // NEW: Dry run endpoint
	api.HandleFunc("/workflows/{id}/toggle", workflowsHandler.ToggleWorkflow).Methods("PUT")
	api.HandleFunc("/workflows/{id}", workflowsHandler.DeleteWorkflow).Methods("DELETE")

	// Logs routes
	logsHandler := handlers.NewLogsHandler(database)
	api.HandleFunc("/logs", logsHandler.GetLogs).Methods("GET")

	// Kong Gateway integration routes
	kongHandler := handlers.NewKongHandler(database, getEnv("KONG_ADMIN_URL", "http://kong:8001"))
	api.HandleFunc("/kong/services", kongHandler.CreateKongService).Methods("POST")
	api.HandleFunc("/kong/services", kongHandler.ListKongServices).Methods("GET")
	api.HandleFunc("/kong/services/{id}", kongHandler.DeleteKongService).Methods("DELETE")
	api.HandleFunc("/kong/routes", kongHandler.CreateKongRoute).Methods("POST")
	api.HandleFunc("/kong/plugins", kongHandler.AddKongPlugin).Methods("POST")
	api.HandleFunc("/kong/templates", kongHandler.CreateUseCaseTemplate).Methods("POST")

	// PRODUCTION FIX: Use battle-tested CORS library instead of manual headers
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: getAllowedOrigins(),
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
		},
		ExposedHeaders: []string{
			"Content-Length",
			"Content-Type",
		},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
		Debug:            getEnv("ENVIRONMENT", "development") == "development",
	}).Handler(router)

	// PRODUCTION FIX: Create HTTP server with proper timeouts
	port := getEnv("PORT", "8080")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsHandler,
		
		// Timeout configurations to prevent resource exhaustion
		ReadTimeout:       15 * time.Second, // Time to read request body
		ReadHeaderTimeout: 10 * time.Second, // Time to read request headers
		WriteTimeout:      30 * time.Second, // Time to write response (increased for long-running workflows)
		IdleTimeout:       120 * time.Second, // Time to keep connection open for next request
		
		// Maximum header size
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	appLogger.Info("Server configured with production timeouts", map[string]interface{}{
		"port":               port,
		"read_timeout":       srv.ReadTimeout.String(),
		"write_timeout":      srv.WriteTimeout.String(),
		"idle_timeout":       srv.IdleTimeout.String(),
		"max_header_bytes":   srv.MaxHeaderBytes,
	})

	// Setup graceful shutdown with context
	shutdownCtx, shutdownCancel := context.WithCancel(context.Background())
	defer shutdownCancel()

	// Channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		appLogger.Info("Server listening", map[string]interface{}{
			"port": port,
			"endpoints": map[string]interface{}{
				"health":   "/health",
				"auth":     "/api/auth/*",
				"webhooks": "/api/webhooks/:id",
				"api":      "/api/*",
			},
		})

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("Server failed to start", map[string]interface{}{
				"error": err.Error(),
			})
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	sig := <-sigChan
	appLogger.Info("Received shutdown signal", map[string]interface{}{
		"signal": sig.String(),
	})

	// Graceful shutdown with timeout
	shutdownTimeout := 30 * time.Second
	ctx, cancel := context.WithTimeout(shutdownCtx, shutdownTimeout)
	defer cancel()

	appLogger.Info("Initiating graceful shutdown...", map[string]interface{}{
		"timeout": shutdownTimeout.String(),
	})

	// Stop scheduler first
	scheduler.Stop()
	appLogger.Info("Scheduler stopped", nil)

	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error("Server shutdown error", map[string]interface{}{
			"error": err.Error(),
		})
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Close database
	database.Close()
	appLogger.Info("Database closed", nil)

	appLogger.Info("Graceful shutdown complete", nil)
}

// getAllowedOrigins returns CORS allowed origins based on environment
func getAllowedOrigins() []string {
	env := getEnv("ENVIRONMENT", "development")
	
	if env == "production" {
		// Production: Only allow specific domains
		allowedOrigins := getEnv("CORS_ALLOWED_ORIGINS", "")
		if allowedOrigins != "" {
			// Parse comma-separated list
			return parseCSV(allowedOrigins)
		}
		// Default production origins
		return []string{
			"https://app.ipaas.com",
			"https://dashboard.ipaas.com",
		}
	}
	
	// Development: Allow localhost and common dev ports
	return []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"http://localhost:8080",
		"http://127.0.0.1:3000",
	}
}

// parseCSV splits a comma-separated string into a slice
func parseCSV(s string) []string {
	if s == "" {
		return []string{}
	}
	var result []string
	for _, item := range splitString(s, ",") {
		trimmed := trimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// splitString splits a string by delimiter
func splitString(s, delimiter string) []string {
	if s == "" {
		return []string{}
	}
	// Simple split implementation
	var result []string
	current := ""
	for _, char := range s {
		if string(char) == delimiter {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// trimSpace removes leading and trailing whitespace
func trimSpace(s string) string {
	start := 0
	end := len(s)
	
	// Trim leading spaces
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	
	// Trim trailing spaces
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	
	return s[start:end]
}

// getEnv gets an environment variable with a default fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// initializeDatabaseWithRetry attempts to initialize the database with exponential backoff
// This is critical for Docker environments where the DB container might not be ready immediately
func initializeDatabaseWithRetry(logger *logger.Logger, maxRetries int, initialDelay time.Duration) (*db.Database, error) {
	dbPath := getEnv("DB_PATH", "ipaas.db")
	delay := initialDelay

	logger.Info("Initializing database with retry logic", map[string]interface{}{
		"db_path":     dbPath,
		"max_retries": maxRetries,
	})

	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		database, err := db.New(dbPath)
		if err == nil {
			// Success! Test the connection with a simple query
			pingErr := database.Ping()
			if pingErr == nil {
				logger.Info("Database connection successful", map[string]interface{}{
					"attempt": attempt,
				})
				return database, nil
			}
			// Close the database if ping failed
			database.Close()
			err = pingErr
		}

		// Log the failure
		logger.Warn("Database initialization failed, retrying...", map[string]interface{}{
			"attempt":      attempt,
			"max_retries":  maxRetries,
			"error":        err.Error(),
			"retry_in":     delay.String(),
		})

		// If this was the last attempt, return the error
		if attempt == maxRetries {
			return nil, err
		}

		// Wait before retrying with exponential backoff
		time.Sleep(delay)
		
		// Exponential backoff: double the delay each time (max 30 seconds)
		delay *= 2
		if delay > 30*time.Second {
			delay = 30 * time.Second
		}
	}

	return nil, err
}
