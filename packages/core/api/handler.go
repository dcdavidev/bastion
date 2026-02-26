package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/dcdavidev/bastion/packages/core/db"
)

// Handler holds the dependencies for the API endpoints.
type Handler struct {
	DB *db.DB
}

// NewHandler creates a new API handler with the provided database.
func NewHandler(database *db.DB) *Handler {
	return &Handler{DB: database}
}

type StatusResponse struct {
	ConnectedToDB  bool     `json:"connected_to_db"`
	MissingEnvVars []string `json:"missing_env_vars"`
	Migrations     struct {
		CurrentVersion uint `json:"current_version"`
		HasPending     bool `json:"has_pending"`
		IsDirty        bool `json:"is_dirty"`
	} `json:"migrations"`
	HasAdmin bool   `json:"has_admin"`
	Version  string `json:"version"`
}

func (h *Handler) StatusHandler(w http.ResponseWriter, r *http.Request) {
	resp := StatusResponse{
		Version: "1.0.0", // Replace with version constant if available
	}

	requiredVars := []string{"BASTION_DATABASE_URL", "BASTION_JWT_SECRET"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			resp.MissingEnvVars = append(resp.MissingEnvVars, v)
		}
	}

	if h.DB != nil && h.DB.Pool != nil {
		err := h.DB.Pool.Ping(r.Context())
		if err == nil {
			resp.ConnectedToDB = true
			
			// Check migrations
			version, pending, err := h.DB.GetMigrationStatus()
			resp.Migrations.CurrentVersion = version
			resp.Migrations.HasPending = pending
			if err != nil {
				resp.Migrations.IsDirty = true
			}

			// Check admin
			hasAdmin, _ := h.DB.HasAdmin(r.Context())
			resp.HasAdmin = hasAdmin
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
