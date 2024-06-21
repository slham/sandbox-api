package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
)

var key []byte

func Initialize(k string, mode string) bool {
	key = []byte(k)
	log.Println("encryption key set")
	return true
}

func Encrypt(s string) (string, error) {
	plaintext := []byte(s)
	var e error
	ct, e := encrypt(plaintext, key)
	output := base64.StdEncoding.EncodeToString(ct)
	return output, e
}

func Decrypt(s string) (string, error) {
	d, e := base64.StdEncoding.DecodeString(s)
	if e != nil {
		return "", e
	}
	ciphertext := []byte(d)
	pt, e := decrypt(ciphertext, key)
	return string(pt), e
}

func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		msg := fmt.Sprintf("ciphertext too short|%v < %v", len(ciphertext), gcm.NonceSize())
		return nil, errors.New(msg)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
