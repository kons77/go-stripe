package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type Encryption struct {
	Key []byte // AES encryption key (must be 16, 24, or 32 bytes for AES-128/192/256)
}

// Encrypt takes a plaintext string and returns a base64-encoded AES-GCM encrypted string
func (e *Encryption) Encrypt(plaintext string) (string, error) {
	// Create a new AES cipher block
	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return "", err
	}

	// Create a GCM (Galois/Counter Mode) instance for authenticated encryption
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate a random nonce (number used once) required by GCM
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Seal appends the encrypted data to the nonce and authenticates it
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Encode to base64 for safe storage/transmission
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt takes a base64-encoded encrypted string and returns the original plaintext
func (e *Encryption) Decrypt(cryptoText string) (string, error) {
	// Decode base64-encoded ciphertext
	data, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	// Create AES cipher block
	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return "", err
	}

	// Create GCM instance
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// Extract nonce and actual ciphertext
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Decrypt and authenticate the data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
