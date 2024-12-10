package handler

import (
	"errors"
	"log/slog"
	"net/url"
	"strconv"
	"strings"
)

var (
	ApiErrBadRequest = errors.New("bad request")
	ApiErrForbidden  = errors.New("forbidden")
	ApiErrNotFound   = errors.New("not found")
	ApiErrConflict   = errors.New("conflict")
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

type APIQuery struct {
	SortCol string
	Sort    string
	Limit   int
	Offset  int
}

func getStandardQueryParams(query url.Values) (APIQuery, error) {
	apiQuery := APIQuery{}
	if qSortCol := query.Get("sort_column"); qSortCol != "" {
		apiQuery.SortCol = qSortCol
	}
	if qSort := query.Get("sort"); qSort != "" {
		apiQuery.Sort = qSort
	}
	if qLimit := query.Get("limit"); qLimit != "" {
		limit, err := strconv.Atoi(qLimit)
		if err != nil {
			slog.Warn("invalid limit", "limit", qLimit)
			return apiQuery, NewApiError(400, ApiErrBadRequest).Append("invalid limit")
		}
		apiQuery.Limit = limit
	}
	if qOffset := query.Get("offset"); qOffset != "" {
		offset, err := strconv.Atoi(qOffset)
		if err != nil {
			slog.Warn("invalid offset", "offset", qOffset)
			return apiQuery, NewApiError(400, ApiErrBadRequest).Append("invalid offset")
		}
		apiQuery.Offset = offset
	}
	return apiQuery, nil
}
