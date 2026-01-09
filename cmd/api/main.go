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
	"github.com/alexmacdonald/simple-ipass/internal/middleware"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting iPaaS API Server...")

	// Initialize database
	database, err := db.New("ipaas.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	log.Println("Database initialized successfully")

	// Initialize executor
	executor := engine.NewExecutor(database)

	// Initialize scheduler for polling triggers
	scheduler := engine.NewScheduler(database, executor)
	scheduler.Start(60 * time.Second) // Check every 60 seconds
	defer scheduler.Stop()

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down gracefully...")
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

	// Protected routes
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

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

	// Add CORS middleware
	corsRouter := enableCORS(router)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, corsRouter))
}

// enableCORS adds CORS headers to all responses
func enableCORS(router *mux.Router) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		router.ServeHTTP(w, r)
	})
}

