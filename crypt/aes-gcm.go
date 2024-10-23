package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

var key []byte

func Initialize(k string) {
	key = []byte(k)
}

func Encrypt(s string) (string, error) {
	plaintext := []byte(s)

	ct, err := encrypt(plaintext, key)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt. %w", err)
	}

	output := base64.StdEncoding.EncodeToString(ct)

	return output, nil
}

func Decrypt(s string) (string, error) {
	d, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", fmt.Errorf("failed to decode. %w", err)
	}

	pt, err := decrypt(d, key)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt. %w", err)
	}

	return string(pt), nil
}

func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher. %w", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("failed to create gcm. %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to read. %w", err)
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher. %w", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("failed to create gcm. %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		msg := fmt.Sprintf("ciphertext too short|%v < %v", len(ciphertext), gcm.NonceSize())
		return nil, fmt.Errorf(msg)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	output, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("fail to open gcm. %w", err)
	}

	return output, nil
}
