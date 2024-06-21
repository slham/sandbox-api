package handler

import "fmt"

type ApiError struct {
	StatusCode int
	Errs       []string
}

func (a ApiError) Error() string {
	return fmt.Sprintf("%+v", a.Errs)
}

func (a ApiError) Append(s string) ApiError {
	a.Errs = append(a.Errs, s)
	return a
}

func (a ApiError) HasError() bool {
	return len(a.Errs) > 0
}
