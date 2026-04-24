// Package envelope provides secret envelope encryption and decryption
// for protecting sensitive values before writing them to disk.
package envelope

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// ErrInvalidKey is returned when the key is not a valid AES key length.
var ErrInvalidKey = errors.New("envelope: key must be 16, 24, or 32 bytes")

// ErrInvalidCiphertext is returned when decryption fails due to malformed input.
var ErrInvalidCiphertext = errors.New("envelope: invalid ciphertext")

// Envelope encrypts and decrypts secret values using AES-GCM.
type Envelope struct {
	key []byte
}

// New creates a new Envelope with the provided AES key.
// The key must be 16, 24, or 32 bytes for AES-128, AES-192, or AES-256.
func New(key []byte) (*Envelope, error) {
	switch len(key) {
	case 16, 24, 32:
		// valid
	default:
		return nil, ErrInvalidKey
	}
	clone := make([]byte, len(key))
	copy(clone, key)
	return &Envelope{key: clone}, nil
}

// Encrypt encrypts plaintext and returns a base64-encoded ciphertext string.
func (e *Envelope) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("envelope: create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("envelope: create gcm: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("envelope: generate nonce: %w", err)
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a base64-encoded ciphertext string and returns the plaintext.
func (e *Envelope) Decrypt(encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", ErrInvalidCiphertext
	}
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("envelope: create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("envelope: create gcm: %w", err)
	}
	if len(data) < gcm.NonceSize() {
		return "", ErrInvalidCiphertext
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", ErrInvalidCiphertext
	}
	return string(plaintext), nil
}

// EncryptMap encrypts all values in a map, returning a new map of encrypted values.
func (e *Envelope) EncryptMap(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		enc, err := e.Encrypt(v)
		if err != nil {
			return nil, fmt.Errorf("envelope: encrypt key %q: %w", k, err)
		}
		out[k] = enc
	}
	return out, nil
}

// DecryptMap decrypts all values in a map, returning a new map of plaintext values.
func (e *Envelope) DecryptMap(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		dec, err := e.Decrypt(v)
		if err != nil {
			return nil, fmt.Errorf("envelope: decrypt key %q: %w", k, err)
		}
		out[k] = dec
	}
	return out, nil
}
