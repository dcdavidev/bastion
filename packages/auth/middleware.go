package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const (
	AdminContextKey contextKey = "admin_claims"
	UserKey         contextKey = "user_id"
)

// GenerateToken creates a new JWT for a user.
func GenerateToken(userID uuid.UUID, username string, isAdmin bool) (string, error) {
	secret := os.Getenv("BASTION_JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("BASTION_JWT_SECRET not set")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID.String(),
		"username": username,
		"admin":    isAdmin,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(secret))
}

// JWTMiddleware validates the JWT token and ensures the user is authenticated.
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		secret := os.Getenv("BASTION_JWT_SECRET")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		uidStr, _ := claims["user_id"].(string)
		uid, _ := uuid.Parse(uidStr)

		// Add claims and user ID to context
		ctx := context.WithValue(r.Context(), AdminContextKey, claims)
		ctx = context.WithValue(ctx, UserKey, uid)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminMiddleware ensures that the authenticated user has admin privileges.
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(AdminContextKey).(jwt.MapClaims)
		if !ok || claims["admin"] != true {
			http.Error(w, "Forbidden: Admin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
