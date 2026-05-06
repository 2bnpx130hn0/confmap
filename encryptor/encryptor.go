// Package encryptor provides utilities for encrypting and decrypting
// sensitive values within a config map using AES-GCM encryption.
package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

const encPrefix = "enc:"

// Encryptor encrypts and decrypts config values using a symmetric key.
type Encryptor struct {
	gcm cipher.AEAD
}

// New creates a new Encryptor. key must be 16, 24, or 32 bytes.
func New(key []byte) (*Encryptor, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("encryptor: invalid key: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("encryptor: failed to create GCM: %w", err)
	}
	return &Encryptor{gcm: gcm}, nil
}

// Encrypt encrypts plaintext and returns a prefixed base64 string.
func (e *Encryptor) Encrypt(plaintext string) (string, error) {
	nonce := make([]byte, e.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("encryptor: nonce generation failed: %w", err)
	}
	ciphertext := e.gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return encPrefix + base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a prefixed base64 string produced by Encrypt.
func (e *Encryptor) Decrypt(ciphertext string) (string, error) {
	if !strings.HasPrefix(ciphertext, encPrefix) {
		return "", errors.New("encryptor: value is not encrypted")
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(ciphertext, encPrefix))
	if err != nil {
		return "", fmt.Errorf("encryptor: base64 decode failed: %w", err)
	}
	ns := e.gcm.NonceSize()
	if len(decoded) < ns {
		return "", errors.New("encryptor: ciphertext too short")
	}
	plain, err := e.gcm.Open(nil, decoded[:ns], decoded[ns:], nil)
	if err != nil {
		return "", fmt.Errorf("encryptor: decryption failed: %w", err)
	}
	return string(plain), nil
}

// DecryptConfig walks a config map and decrypts any string values with the enc: prefix.
func (e *Encryptor) DecryptConfig(cfg map[string]any) (map[string]any, error) {
	out := make(map[string]any, len(cfg))
	for k, v := range cfg {
		switch val := v.(type) {
		case string:
			if strings.HasPrefix(val, encPrefix) {
				decrypted, err := e.Decrypt(val)
				if err != nil {
					return nil, fmt.Errorf("encryptor: key %q: %w", k, err)
				}
				out[k] = decrypted
			} else {
				out[k] = val
			}
		case map[string]any:
			nested, err := e.DecryptConfig(val)
			if err != nil {
				return nil, err
			}
			out[k] = nested
		default:
			out[k] = v
		}
	}
	return out, nil
}
