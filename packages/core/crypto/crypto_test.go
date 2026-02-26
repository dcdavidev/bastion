package crypto

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeriveKey(t *testing.T) {
	password := []byte("secret-password")
	salt := []byte("1234567890123456") // 16 bytes

	key1 := DeriveKey(password, salt)
	key2 := DeriveKey(password, salt)

	assert.Len(t, key1, keyLen)
	assert.Equal(t, key1, key2)

	// Different password
	key3 := DeriveKey([]byte("other-password"), salt)
	assert.NotEqual(t, key1, key3)

	// Different salt
	key4 := DeriveKey(password, []byte("6543210987654321"))
	assert.NotEqual(t, key1, key4)
}

func TestEncryptDecrypt(t *testing.T) {
	key, err := GenerateRandomKey()
	require.NoError(t, err)

	plaintext := []byte("hello-bastion-world")
	ciphertext, err := Encrypt(key, plaintext)
	require.NoError(t, err)
	assert.NotEqual(t, plaintext, ciphertext)

	decrypted, err := Decrypt(key, ciphertext)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)

	// Test with wrong key
	wrongKey, _ := GenerateRandomKey()
	_, err = Decrypt(wrongKey, ciphertext)
	assert.Error(t, err)
}

func TestGenerateSalt(t *testing.T) {
	salt1, err := GenerateSalt()
	require.NoError(t, err)
	assert.Len(t, salt1, saltLen)

	salt2, err := GenerateSalt()
	require.NoError(t, err)
	assert.NotEqual(t, salt1, salt2)
}

func TestGenerateRandomKey(t *testing.T) {
	key1, err := GenerateRandomKey()
	require.NoError(t, err)
	assert.Len(t, key1, keyLen)

	key2, err := GenerateRandomKey()
	require.NoError(t, err)
	assert.NotEqual(t, key1, key2)
}

func TestGenerateRandomKeyInto(t *testing.T) {
	key := make([]byte, 32)
	err := GenerateRandomKeyInto(key)
	require.NoError(t, err)
	
	// Check it's not all zeros
	isZero := true
	for _, b := range key {
		if b != 0 {
			isZero = false
			break
		}
	}
	assert.False(t, isZero)
}

func TestWrapUnwrapKey(t *testing.T) {
	wrapperKey, _ := GenerateRandomKey()
	targetKey, _ := GenerateRandomKey()

	wrapped, err := WrapKey(wrapperKey, targetKey)
	require.NoError(t, err)

	unwrapped, err := UnwrapKey(wrapperKey, wrapped)
	require.NoError(t, err)
	assert.Equal(t, targetKey, unwrapped)
}

func TestHexEncoding(t *testing.T) {
	// Simple test to verify our assumptions about hex strings
	data := []byte("hello")
	hexStr := hex.EncodeToString(data)
	assert.Equal(t, "68656c6c6f", hexStr)

	decoded, err := hex.DecodeString(hexStr)
	require.NoError(t, err)
	assert.Equal(t, data, decoded)
}
