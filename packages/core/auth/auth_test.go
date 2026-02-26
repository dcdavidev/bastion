package auth

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/dcdavidev/bastion/packages/core/crypto"
	"github.com/stretchr/testify/assert"
)

func TestVerifyAdmin(t *testing.T) {
	// Setup test credentials
	password := "secure-admin-pass"
	salt, _ := crypto.GenerateSalt()
	hash := crypto.DeriveKey([]byte(password), salt)

	os.Setenv(EnvAdminHash, hex.EncodeToString(hash))
	os.Setenv(EnvAdminSalt, hex.EncodeToString(salt))
	defer os.Unsetenv(EnvAdminHash)
	defer os.Unsetenv(EnvAdminSalt)

	// Test correct password
	assert.True(t, VerifyAdmin(password))

	// Test incorrect password
	assert.False(t, VerifyAdmin("wrong-pass"))

	// Test missing environment variables
	os.Unsetenv(EnvAdminHash)
	assert.False(t, VerifyAdmin(password))
}

func TestVerifyAdmin_InvalidHex(t *testing.T) {
	os.Setenv(EnvAdminHash, "invalid-hex-string")
	os.Setenv(EnvAdminSalt, "not-hex")
	defer os.Unsetenv(EnvAdminHash)
	defer os.Unsetenv(EnvAdminSalt)

	assert.False(t, VerifyAdmin("any-pass"))
}
