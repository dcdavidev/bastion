package api

import "github.com/dcdavidev/bastion/internal/db"

// Handler holds the dependencies for the API endpoints.
type Handler struct {
	DB *db.DB
}

// NewHandler creates a new API handler with the provided database.
func NewHandler(database *db.DB) *Handler {
	return &Handler{DB: database}
}
