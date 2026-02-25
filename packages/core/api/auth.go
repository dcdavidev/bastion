package api

import (
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dcdavidev/bastion/packages/core/auth"
	"github.com/dcdavidev/bastion/packages/core/crypto"
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// LoginHandler handles both admin (env-based) and collaborator (db-based) logins.
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var role string
	var userID string

	// 1. Check if it's a Collaborator Login (Database)
	if req.Username != "" {
		user, storedHashHex, saltHex, err := h.DB.GetUserByUsername(r.Context(), req.Username)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		salt, _ := hex.DecodeString(saltHex)
		storedHash, _ := hex.DecodeString(storedHashHex)
		computedHash := crypto.DeriveKey([]byte(req.Password), salt)

		if subtle.ConstantTimeCompare(computedHash, storedHash) != 1 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		role = user.Role
		userID = user.ID.String()
	} else {
		// 2. Fallback to Admin Login (Environment Variables)
		if !auth.VerifyAdmin(req.Password) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		role = "ADMIN"
		userID = "00000000-0000-0000-0000-000000000000" // Reserved Admin ID
	}

	// Generate JWT
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"admin":   role == "ADMIN",
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: tokenString})
}

// GetVaultConfigHandler returns the public vault configuration needed for client-side decryption.
func (h *Handler) GetVaultConfigHandler(w http.ResponseWriter, r *http.Request) {
	config, err := h.DB.GetVaultConfig(r.Context())
	if err != nil {
		http.Error(w, "Vault not initialized", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}
