package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

type Encryption struct {
	Key []byte // Encryption key (must be 32 bytes)
}

func (e *Encryption) Encrypt(text string) (string, error) {
	plaintext := []byte(text)

	block, err := aes.NewCipher(e.Key) // Create a new AES encryption block
	if err != nil {
		return "", err
	}

	// Create a byte slice for storing encrypted text + IV
	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize] // iv = Initialization Vector, First part of the array is IV

	// Generate a random IV
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Create CFB (Cipher Feedback Mode) encrypter
	stream := cipher.NewCFBEncrypter(block, iv)

	// Encrypt text
	stream.XORKeyStream(cipherText[aes.BlockSize:], plaintext)

	// Encode encrypted text in Base64 for safe transmission
	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func (e *Encryption) Decrypt(cryptoText string) (string, error) {
	// Decode from Base64
	cipherText, _ := base64.URLEncoding.DecodeString(cryptoText)

	//Create decryption block
	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return "", err
	}

	// Validate the length of the encrypted text
	if len(cipherText) < aes.BlockSize {
		return "", err
	}

	// Extract IV from the first part of the array
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	// Create CFB decrypter
	stream := cipher.NewCFBDecrypter(block, iv)

	// Decrypt text
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
