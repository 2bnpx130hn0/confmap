package encryptor_test

import (
	"strings"
	"testing"

	"github.com/example/confmap/encryptor"
)

var testKey = []byte("12345678901234567890123456789012") // 32 bytes

func TestNew_InvalidKey(t *testing.T) {
	_, err := encryptor.New([]byte("short"))
	if err == nil {
		t.Fatal("expected error for short key")
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	enc, err := encryptor.New(testKey)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	plaintext := "super-secret-value"
	cipher, err := enc.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if !strings.HasPrefix(cipher, "enc:") {
		t.Errorf("expected enc: prefix, got %q", cipher)
	}
	got, err := enc.Decrypt(cipher)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if got != plaintext {
		t.Errorf("expected %q, got %q", plaintext, got)
	}
}

func TestDecrypt_MissingPrefix(t *testing.T) {
	enc, _ := encryptor.New(testKey)
	_, err := enc.Decrypt("not-encrypted")
	if err == nil {
		t.Fatal("expected error for missing prefix")
	}
}

func TestDecryptConfig_DecryptsStringValues(t *testing.T) {
	enc, _ := encryptor.New(testKey)
	cipher, _ := enc.Encrypt("my-password")
	cfg := map[string]any{
		"host":     "localhost",
		"password": cipher,
		"port":     5432,
	}
	out, err := enc.DecryptConfig(cfg)
	if err != nil {
		t.Fatalf("DecryptConfig: %v", err)
	}
	if out["password"] != "my-password" {
		t.Errorf("expected 'my-password', got %v", out["password"])
	}
	if out["host"] != "localhost" {
		t.Errorf("plain string should be preserved")
	}
	if out["port"] != 5432 {
		t.Errorf("non-string value should be preserved")
	}
}

func TestDecryptConfig_Nested(t *testing.T) {
	enc, _ := encryptor.New(testKey)
	cipher, _ := enc.Encrypt("nested-secret")
	cfg := map[string]any{
		"database": map[string]any{
			"user":     "admin",
			"password": cipher,
		},
	}
	out, err := enc.DecryptConfig(cfg)
	if err != nil {
		t.Fatalf("DecryptConfig: %v", err)
	}
	db, ok := out["database"].(map[string]any)
	if !ok {
		t.Fatal("expected nested map")
	}
	if db["password"] != "nested-secret" {
		t.Errorf("expected 'nested-secret', got %v", db["password"])
	}
}

func TestEncrypt_UniqueEachCall(t *testing.T) {
	enc, _ := encryptor.New(testKey)
	a, _ := enc.Encrypt("same")
	b, _ := enc.Encrypt("same")
	if a == b {
		t.Error("expected different ciphertexts due to random nonce")
	}
}
