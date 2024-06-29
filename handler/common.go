package handler

import (
	"errors"
	"strings"
)

var (
	ApiErrBadRequest = errors.New("bad request")
)

type ApiError struct {
	BaseError  error
	StatusCode int
	Errs       []string
}

func NewApiError(statusCode int, baseError error) *ApiError {
	return &ApiError{
		BaseError:  baseError,
		StatusCode: statusCode,
		Errs:       []string{},
	}
}

func (a *ApiError) Error() string {
	return strings.Join(a.Errs, ". ")
}

func (a *ApiError) Unwrap() error {
	return a.BaseError
}

func (a *ApiError) Append(s string) *ApiError {
	a.Errs = append(a.Errs, s)
	return a
}

func (a *ApiError) HasError() bool {
	return len(a.Errs) > 0
}
