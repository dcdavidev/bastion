package auth

import (
	"crypto/subtle"
	"encoding/hex"
	"os"

	"github.com/dcdavidev/bastion/internal/crypto"
)

// AdminCredentials defines the expected environment variables for the admin user.
const (
	EnvAdminHash = "ADMIN_PASSWORD_HASH"
	EnvAdminSalt = "ADMIN_PASSWORD_SALT"
)

// VerifyAdmin checks if the provided password matches the hash stored in environment variables.
func VerifyAdmin(password string) bool {
	storedHashHex := os.Getenv(EnvAdminHash)
	storedSaltHex := os.Getenv(EnvAdminSalt)

	if storedHashHex == "" || storedSaltHex == "" {
		return false
	}

	storedHash, err := hex.DecodeString(storedHashHex)
	if err != nil {
		return false
	}

	salt, err := hex.DecodeString(storedSaltHex)
	if err != nil {
		return false
	}

	// Derive hash from the provided password using the stored salt
	computedHash := crypto.DeriveKey([]byte(password), salt)

	// Use ConstantTimeCompare to prevent timing attacks
	return subtle.ConstantTimeCompare(computedHash, storedHash) == 1
}
