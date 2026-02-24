package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dcdavidev/bastion/internal/api"
	"github.com/dcdavidev/bastion/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize Database
	database, err := db.NewConnection()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer database.Close()

	// Run Migrations
	if err := database.RunMigrations(); err != nil {
		log.Fatalf("Could not run migrations: %v", err)
	}

	// Initialize API Handler
	h := api.NewHandler(database)

	r := chi.NewRouter()

	// Standard middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// API Routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Post("/auth/login", h.LoginHandler)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(auth.JWTMiddleware)
			
			r.Get("/vault/config", h.GetVaultConfigHandler)
			r.Get("/audit", h.ListAuditLogs)
			r.Get("/clients", h.ListClients)
			r.Post("/clients", h.CreateClient)
			r.Delete("/clients/{id}", h.DeleteClient)

			r.Get("/projects", h.ListProjectsByClient)
			r.Get("/projects/{id}", h.GetProject)
			r.Get("/projects/{id}/key", h.GetProjectKey)
			r.Post("/projects", h.CreateProject)
			r.Delete("/projects/{id}", h.DeleteProject)

			r.Post("/collaborators", h.CreateCollaborator)

			r.Get("/secrets", h.ListSecretsByProject)
			r.Post("/secrets", h.CreateSecret)
			r.Get("/secrets/history", h.GetSecretHistory)
		})
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"up"}`))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Bastion server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
