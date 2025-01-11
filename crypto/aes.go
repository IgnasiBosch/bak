package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

const EncryptionMagicBytes = "ENCRYPTED:" // or any unique identifier

func Encrypt(data []byte, key []byte) ([]byte, error) {
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

	encrypted := gcm.Seal(nonce, nonce, data, nil)
	// Prepend magic bytes
	return append([]byte(EncryptionMagicBytes), encrypted...), nil
}

func IsEncrypted(data []byte) bool {
	return len(data) >= len(EncryptionMagicBytes) &&
		string(data[:len(EncryptionMagicBytes)]) == EncryptionMagicBytes
}

func Decrypt(encrypted []byte, key []byte) ([]byte, error) {
	// Check if the data is encrypted
	if !IsEncrypted(encrypted) {
		return nil, fmt.Errorf("data is not encrypted")
	}

	// Remove magic bytes
	encrypted = encrypted[len(EncryptionMagicBytes):]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(encrypted) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
