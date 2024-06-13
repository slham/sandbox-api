package valid

import (
	"fmt"
	"regexp"

	"github.com/slham/sandbox-api/constant"
)

func isAlphaNumeric(s string) (bool, error) {
	return regexp.Match(constant.AlphaNumericRegex.String(), []byte(s))
}

func isNumber(s string) (bool, error) {
	return regexp.Match(constant.NumberRegex.String(), []byte(s))
}

func isMediumPassword(s string) (bool, error) {
	return regexp.Match(constant.MediumPasswordRegex.String(), []byte(s))
}

func isStrongPassword(s string) (bool, error) {
	return regexp.Match(constant.StrongPasswordRegex.String(), []byte(s))
}

func isEmail(s string) (bool, error) {
	return regexp.Match(constant.EmailRegex.String(), []byte(s))
}

func validateWithRegex(input string, f func(string) (bool, error)) (bool, error) {
	b, err := f(input)

	if !b && err == nil {
		return false, fmt.Errorf("invalid format")
	}

	return true, nil
}
