package envelope_test

import (
	"strings"
	"testing"

	"github.com/vaultpipe/vaultpipe/internal/envelope"
)

var testKey16 = []byte("0123456789abcdef")          // 16 bytes
var testKey32 = []byte("0123456789abcdef0123456789abcdef") // 32 bytes

func TestNew_ValidKeyLengths(t *testing.T) {
	for _, key := range [][]byte{
		[]byte("0123456789abcdef"),           // 16
		[]byte("0123456789abcdef01234567"),   // 24
		testKey32,                            // 32
	} {
		_, err := envelope.New(key)
		if err != nil {
			t.Errorf("expected no error for key len %d, got %v", len(key), err)
		}
	}
}

func TestNew_InvalidKeyLength(t *testing.T) {
	_, err := envelope.New([]byte("short"))
	if err == nil {
		t.Fatal("expected error for invalid key length")
	}
	if err != envelope.ErrInvalidKey {
		t.Errorf("expected ErrInvalidKey, got %v", err)
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	e, err := envelope.New(testKey16)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	plaintext := "super-secret-value"
	enc, err := e.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if enc == plaintext {
		t.Error("encrypted value should differ from plaintext")
	}
	dec, err := e.Decrypt(enc)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if dec != plaintext {
		t.Errorf("expected %q, got %q", plaintext, dec)
	}
}

func TestEncrypt_ProducesUniqueOutputs(t *testing.T) {
	e, _ := envelope.New(testKey16)
	a, _ := e.Encrypt("same")
	b, _ := e.Encrypt("same")
	if a == b {
		t.Error("expected unique ciphertexts due to random nonce")
	}
}

func TestDecrypt_InvalidBase64_ReturnsError(t *testing.T) {
	e, _ := envelope.New(testKey16)
	_, err := e.Decrypt("!!!not-base64!!!")
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

func TestDecrypt_TamperedCiphertext_ReturnsError(t *testing.T) {
	e, _ := envelope.New(testKey16)
	enc, _ := e.Encrypt("value")
	tampered := strings.ToUpper(enc)
	_, err := e.Decrypt(tampered)
	if err == nil {
		t.Fatal("expected error for tampered ciphertext")
	}
}

func TestEncryptMap_DecryptMap_RoundTrip(t *testing.T) {
	e, _ := envelope.New(testKey32)
	secrets := map[string]string{
		"DB_PASS":  "hunter2",
		"API_KEY":  "abc123",
		"EMPTY_VAL": "",
	}
	encrypted, err := e.EncryptMap(secrets)
	if err != nil {
		t.Fatalf("EncryptMap: %v", err)
	}
	for k, v := range secrets {
		if encrypted[k] == v && v != "" {
			t.Errorf("key %q: expected encrypted value to differ", k)
		}
	}
	decrypted, err := e.DecryptMap(encrypted)
	if err != nil {
		t.Fatalf("DecryptMap: %v", err)
	}
	for k, want := range secrets {
		if got := decrypted[k]; got != want {
			t.Errorf("key %q: expected %q, got %q", k, want, got)
		}
	}
}
