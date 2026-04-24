// Package envelope provides AES-GCM envelope encryption for secret values
// managed by vaultpipe.
//
// It allows secrets fetched from HashiCorp Vault to be encrypted at rest
// before being written to .env files, providing an additional layer of
// protection for sensitive data on disk.
//
// # Usage
//
//	enc, err := envelope.New([]byte("your-32-byte-aes-key-here!!!!!!!"))
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	ciphertext, err := enc.Encrypt("my-secret-value")
//	plaintext, err := enc.Decrypt(ciphertext)
//
// Keys must be 16, 24, or 32 bytes for AES-128, AES-192, or AES-256 respectively.
// Each encryption call produces a unique ciphertext via a random nonce.
package envelope
