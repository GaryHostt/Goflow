package main

import (
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
)

func main() {
	// Initialize structured logger (ELK-ready!)
	appLogger := logger.NewLogger("ipaas-api")
	appLogger.Info("Starting iPaaS API Server...", map[string]interface{}{
		"version": "0.1.0",
		"env":     os.Getenv("ENVIRONMENT"),
	})

	// Initialize database
	database, err := db.New("ipaas.db")
	if err != nil {
		appLogger.Error("Failed to initialize database", map[string]interface{}{
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

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		appLogger.Info("Shutting down gracefully...", map[string]interface{}{
			"signal": sig.String(),
		})
		scheduler.Stop()
		database.Close()
		os.Exit(0)
	}()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(database)
	credentialsHandler := handlers.NewCredentialsHandler(database)
	workflowsHandler := handlers.NewWorkflowsHandler(database, executor)
	webhookHandler := handlers.NewWebhookHandler(database, executor)
	logsHandler := handlers.NewLogsHandler(database)

	// Setup router
	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/api/webhooks/{id}", webhookHandler.TriggerWebhook).Methods("POST")

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}).Methods("GET")

	// Protected routes with tenant-aware middleware
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware(appLogger)) // Now logs user_id AND tenant_id!

	// Credentials routes
	api.HandleFunc("/credentials", credentialsHandler.CreateCredential).Methods("POST")
	api.HandleFunc("/credentials", credentialsHandler.GetCredentials).Methods("GET")

	// Workflows routes
	api.HandleFunc("/workflows", workflowsHandler.CreateWorkflow).Methods("POST")
	api.HandleFunc("/workflows", workflowsHandler.GetWorkflows).Methods("GET")
	api.HandleFunc("/workflows/{id}/toggle", workflowsHandler.ToggleWorkflow).Methods("PUT")
	api.HandleFunc("/workflows/{id}", workflowsHandler.DeleteWorkflow).Methods("DELETE")

	// Logs routes
	api.HandleFunc("/logs", logsHandler.GetLogs).Methods("GET")

	// Add CORS middleware with logging
	corsRouter := enableCORS(router, appLogger)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	appLogger.Info("Server listening", map[string]interface{}{
		"port": port,
		"endpoints": map[string]interface{}{
			"health":   "/health",
			"auth":     "/api/auth/*",
			"webhooks": "/api/webhooks/:id",
			"api":      "/api/*",
		},
	})

	log.Fatal(http.ListenAndServe(":"+port, corsRouter))
}

// enableCORS adds CORS headers with logging
func enableCORS(router *mux.Router, appLogger *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Log all requests (ELK will capture these)
		appLogger.Debug("HTTP Request", map[string]interface{}{
			"method": r.Method,
			"path":   r.URL.Path,
			"ip":     r.RemoteAddr,
		})

		router.ServeHTTP(w, r)
	})
}
