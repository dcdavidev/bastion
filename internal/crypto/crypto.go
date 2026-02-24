package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	saltLen    = 16
	keyLen     = 32 // AES-256
	timeParams = 1
	memParams  = 64 * 1024
	threads    = 4
)

// DeriveKey generates a 32-byte key from a password and salt using Argon2id.
func DeriveKey(password []byte, salt []byte) []byte {
	return argon2.IDKey(password, salt, timeParams, memParams, threads, keyLen)
}

// Encrypt encrypts plain text using a key with AES-GCM.
func Encrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Returns nonce + ciphertext
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt decrypts ciphertext using a key with AES-GCM.
func Decrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, encryptedData := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, encryptedData, nil)
}

// GenerateSalt creates a random 16-byte salt.
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	return salt, nil
}

// GenerateRandomKey creates a random 32-byte key for AES-256.
func GenerateRandomKey() ([]byte, error) {
	key := make([]byte, keyLen)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

// WrapKey encrypts a target key using a wrapper key.
func WrapKey(wrapperKey, targetKey []byte) ([]byte, error) {
	return Encrypt(wrapperKey, targetKey)
}

// UnwrapKey decrypts a target key using a wrapper key.
func UnwrapKey(wrapperKey, wrappedKey []byte) ([]byte, error) {
	return Decrypt(wrapperKey, wrappedKey)
}
