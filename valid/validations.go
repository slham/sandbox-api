package valid

import (
	"fmt"
	"regexp"

	"github.com/slham/sandbox-api/constant"
)

func IsAlphaNumeric(s string) (bool, error) {
	return regexp.Match(constant.AlphaNumericRegex.String(), []byte(s))
}

func IsNumber(s string) (bool, error) {
	return regexp.Match(constant.NumberRegex.String(), []byte(s))
}

func IsMediumPassword(s string) (bool, error) {
	return regexp.Match(constant.MediumPasswordRegex.String(), []byte(s))
}

func IsStrongPassword(s string) (bool, error) {
	return regexp.Match(constant.StrongPasswordRegex.String(), []byte(s))
}

func IsEmail(s string) (bool, error) {
	return regexp.Match(constant.EmailRegex.String(), []byte(s))
}

func validateWithRegex(input string, f func(string) (bool, error)) (bool, error) {
	b, err := f(input)

	if !b && err == nil {
		return false, fmt.Errorf("invalid format")
	}

	return true, nil
}
