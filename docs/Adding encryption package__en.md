### Adding encryption package: fixing the Security Issue with Password Reset

(section 8 video 109)

Last time, we discussed a serious security issue in our application. Let's analyze it in detail and find a solution.

#### The Problem
On the password reset page, a user enters their email, and the system sends a reset link. However, when viewing the page source code, the email used for password reset is visible. This means that an attacker could modify it to another email and reset someone else's password.

#### The Solution
We will create an encryption package accessible from both the frontend and backend. This will allow us to encrypt the email before sending it to the server and decrypt it on the backend.

##### Step 1: Creating a Secret Key
Encryption algorithms require a fixed-length key of 32 characters. We define this key in the API file (`api.go`) and in the main file (`main.go`).

##### Step 2: Creating a New Package
In the `internal` folder, we create a new package named `encryption` and a file `encryption.go`.

##### Step 3: Defining the Encryption Structure
We create a struct `Encryption`, which contains the encryption key as a byte slice.

##### Step 4: Encryption Function
- Convert the input text into a byte slice.
- Create an encryption block using the key.
- Generate a random initialization vector (IV).
- Use AES encryption (CFB mode) to process the text.
- Encode the result in Base64 for safe use in URLs and HTML.

##### Step 5: Decryption Function
- Decode the Base64 string back into a byte slice.
- Create a decryption block.
- Validate the length of the encrypted text.
- Extract the IV and decrypt the data.
- Convert the result back to a string and return it.

### Detailed Code Analysis with Comments

```go
package encryption

import (
	"crypto/aes"       // AES encryption library
	"crypto/cipher"    // Methods for handling ciphers
	"crypto/rand"      // Random value generation
	"encoding/base64"  // Base64 encoding for safe transmission
	"io"               // Input/output operations
)

// Struct for handling encryption
type Encryption struct {
	Key []byte // Encryption key (must be 32 bytes)
}

// Encryption function
func (e *Encryption) Encrypt(text string) (string, error) {
	plaintext := []byte(text) // Convert text to byte slice

	block, err := aes.NewCipher(e.Key) // Create a new AES encryption block
	if err != nil {
		return "", err // Return an error if block creation fails
	}

	// Create a byte slice for storing encrypted text + IV
	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize] // First part of the array is IV

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

// Decryption function
func (e *Encryption) Decrypt(cryptoText string) (string, error) {
	// Decode from Base64
	cipherText, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(e.Key) // Create decryption block
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

	return string(cipherText), nil // Convert result back to string
}
```

### Conclusion
Now, the user's email is transmitted in an encrypted format, preventing attackers from modifying it. We used AES encryption with CFB mode and Base64 encoding for secure transmission. This mechanism helps protect the password reset system from attacks.