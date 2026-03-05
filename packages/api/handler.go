package api

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"

	"github.com/dcdavidev/bastion/packages/db"
)

// Handler holds the dependencies for the API endpoints.
type Handler struct {
	DB db.Database
}

// NewHandler creates a new API handler with the provided database.
func NewHandler(database db.Database) *Handler {
	return &Handler{DB: database}
}

type StatusResponse struct {
	ConnectedToDB  bool     `json:"connected_to_db"`
	MissingEnvVars []string `json:"missing_env_vars"`
	JwtSecretStatus string  `json:"jwt_secret_status"` // "strong", "weak", "missing"
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
		JwtSecretStatus: "missing",
	}

	jwtSecret := os.Getenv("BASTION_JWT_SECRET")
	if jwtSecret == "" {
		resp.MissingEnvVars = append(resp.MissingEnvVars, "BASTION_JWT_SECRET")
		resp.JwtSecretStatus = "missing"
	} else {
		// A "real" key in Bastion should be a hex-encoded string of at least 32 bytes (64 hex chars)
		// and not one of the default placeholders.
		isHex := true
		decoded, err := hex.DecodeString(jwtSecret)
		if err != nil {
			isHex = false
		}

		isDefault := jwtSecret == "bastion_very_secret_key_change_me" || 
					 jwtSecret == "your_super_secret_jwt_key_change_me"

		if isDefault || !isHex || len(decoded) < 32 {
			resp.JwtSecretStatus = "weak"
		} else {
			resp.JwtSecretStatus = "strong"
		}
	}

	if os.Getenv("BASTION_DATABASE_URL") == "" && os.Getenv("DATABASE_URL") == "" {
		resp.MissingEnvVars = append(resp.MissingEnvVars, "BASTION_DATABASE_URL")
	}

	if h.DB != nil {
		err := h.DB.Ping(r.Context())
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
