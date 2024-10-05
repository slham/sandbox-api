//go:build unit
// +build unit

package crypt

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestCrypt(t *testing.T) {
	key = []byte("qwertyuiopasdfghjklzxcvbnm098765")
	blah := "encryptMe"
	e, err := Encrypt(blah)
	if err != nil {
		t.Fatal(err.Error())
	}
	d, err := Decrypt(e)
	if err != nil {
		t.Fatal(err.Error())
	}
	assert.Equal(t, blah, d)
}
