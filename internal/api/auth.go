package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dcdavidev/bastion/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// LoginHandler handles the admin login request.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !auth.VerifyAdmin(req.Password) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Generate JWT
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: tokenString})
}
