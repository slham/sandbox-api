package valid

import (
	"fmt"
	"net/mail"
	"regexp"
	"unicode"

	"github.com/slham/sandbox-api/constant"
)

func IsAlphaNumeric(s string) (bool, error) {
	return regexp.MatchString(constant.AlphaNumericRegex.String(), s)
}

func IsNumber(s string) (bool, error) {
	return regexp.MatchString(constant.NumberRegex.String(), s)
}

func IsMediumPassword(s string) bool {
	return isValidPassword(s, 8)
}

func IsStrongPassword(s string) bool {
	return isValidPassword(s, 16)
}

func IsEmail(s string) error {
	_, err := mail.ParseAddress(s)
	return err
}

func validateWithRegex(input string, f func(string) (bool, error)) (bool, error) {
	b, err := f(input)

	if !b && err == nil {
		return false, fmt.Errorf("invalid format")
	}

	return true, nil
}

func isValidPassword(s string, minLen int) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= minLen {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
