//go:build unit
// +build unit

package valid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func TestIsAlphaNumeric(t *testing.T) {
	tables := []struct {
		input     string
		expected1 bool
		expected2 error
	}{
		{"456", true, nil},
		{"abc", true, nil},
		{"abc159", true, nil},
		{"753abc159", true, nil},
		{"753abc", true, nil},
		{"753-abc", false, nil},
		{"753 abc", false, nil},
		{"", false, nil},
	}

	for _, table := range tables {
		status, err := IsAlphaNumeric(table.input)
		assert.Equal(t, table.expected1, status, "status does not match")
		assert.Equal(t, table.expected2, err, "errors do not match")
	}
}

func TestIsNumber(t *testing.T) {
	tables := []struct {
		input     string
		expected1 bool
		expected2 error
	}{
		{"456", true, nil},
		{"abc", false, nil},
		{"abc159", false, nil},
		{"753abc159", false, nil},
		{"753abc", false, nil},
		{"753-abc", false, nil},
		{"753 abc", false, nil},
		{"", false, nil},
	}

	for _, table := range tables {
		status, err := IsNumber(table.input)
		assert.Equal(t, table.expected1, status, "status does not match")
		assert.Equal(t, table.expected2, err, "errors do not match")
	}
}

func TestIsMediumPassword(t *testing.T) {
	tables := []struct {
		input     string
		expected1 bool
	}{
		{"456", false},
		{"abc", false},
		{"abc159", false},
		{"753abc159!", false},
		{"7#53abc!", false},
		{"753-Abcb!", true},
		{"753 ABC", false},
		{"", false},
	}

	for _, table := range tables {
		status := IsMediumPassword(table.input)
		assert.Equal(t, table.expected1, status, "status does not match")
	}
}

func TestIsStrongPassword(t *testing.T) {
	tables := []struct {
		input     string
		expected1 bool
	}{
		{"456", false},
		{"abc", false},
		{"abc159", false},
		{"753abc159!753abc159!", false},
		{"7#53abc!7#53abc!", false},
		{"753-Abcb!753-Abcb!", true},
		{"753 ABC753 ABC753 ABC", false},
		{"", false},
	}

	for _, table := range tables {
		status := IsStrongPassword(table.input)
		assert.Equal(t, table.expected1, status, "status does not match")
	}
}
