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
		status, err := isAlphaNumeric(table.input)
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
		status, err := isNumber(table.input)
		assert.Equal(t, table.expected1, status, "status does not match")
		assert.Equal(t, table.expected2, err, "errors do not match")
	}
}

func TestIsMediumPassword(t *testing.T) {
	tables := []struct {
		input     string
		expected1 bool
		expected2 error
	}{
		{"456", false, nil},
		{"abc", false, nil},
		{"abc159", false, nil},
		{"753abc159!", true, nil},
		{"7#53abc!", true, nil},
		{"753-Abcb!", true, nil},
		{"753 ABC", false, nil},
		{"", false, nil},
	}

	for _, table := range tables {
		status, err := isMediumPassword(table.input)
		assert.Equal(t, table.expected1, status, "status does not match")
		assert.Equal(t, table.expected2, err, "errors do not match")
	}
}

func TestIsStrongPassword(t *testing.T) {
	tables := []struct {
		input     string
		expected1 bool
		expected2 error
	}{
		{"456", false, nil},
		{"abc", false, nil},
		{"abc159", false, nil},
		{"753abc159!753abc159!", true, nil},
		{"7#53abc!7#53abc!", true, nil},
		{"753-Abcb!753-Abcb!", true, nil},
		{"753 ABC753 ABC753 ABC", false, nil},
		{"", false, nil},
	}

	for _, table := range tables {
		status, err := isStrongPassword(table.input)
		assert.Equal(t, table.expected1, status, "status does not match")
		assert.Equal(t, table.expected2, err, "errors do not match")
	}
}

func TestIsEmail(t *testing.T) {
	tables := []struct {
		input     string
		expected1 bool
		expected2 error
	}{
		{"this@that.com", true, nil},
		{"this@that.org", true, nil},
		{"this@that.net", true, nil},
		{"this@that.gov", true, nil},
		{"this@that.info", true, nil},
		{"this@that.blah", true, nil},
		{"abc", false, nil},
		{"abc159", false, nil},
		{"753abc159!753abc159!", false, nil},
		{"7#53@bc!7#53abc!", false, nil},
		{"753-Abcb!753-Abcb!", false, nil},
		{"753 ABC7.53 ABC753 ABC", false, nil},
		{"", false, nil},
	}

	for _, table := range tables {
		status, err := isEmail(table.input)
		assert.Equal(t, table.expected1, status, "status does not match")
		assert.Equal(t, table.expected2, err, "errors do not match")
	}
}
