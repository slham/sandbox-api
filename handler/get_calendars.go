package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/slham/sandbox-api/dao"
	"github.com/slham/sandbox-api/model"
	"github.com/slham/sandbox-api/request"
)

type getCalendarsQuery struct {
	APIQuery
}

type getCalendarsRequest struct {
	userID string
	query  getCalendarsQuery
}

func getCalendarsQueryParams(ctx context.Context, q url.Values) (getCalendarsQuery, error) {
	gcq := getCalendarsQuery{}
	apiQuery, err := getStandardQueryParams(ctx, q)
	if err != nil {
		return gcq, fmt.Errorf("failed to gather query params. %w", err)
	}
	gcq.APIQuery = apiQuery
	return gcq, nil
}

func handleGetCalendarsError(ctx context.Context, w http.ResponseWriter, err error) {
	if errors.Is(err, ApiErrBadRequest) {
		slog.WarnContext(ctx, "error getting calendars", "err", err)
		request.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	slog.ErrorContext(ctx, "error getting calendars", "err", err)
	request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
}

func (c *CalendarController) GetCalendars(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	userID := vars["user_id"]
	req := getCalendarsRequest{userID: userID}
	query := r.URL.Query()
	q, err := getCalendarsQueryParams(ctx, query)
	if err != nil {
		handleGetCalendarsError(ctx, w, err)
		return
	}

	req.query = q

	calendars, err := c.getCalendars(ctx, req)
	if err != nil {
		handleGetCalendarsError(ctx, w, err)
		return
	}

	request.RespondWithJSON(w, http.StatusOK, calendars)
}

func (c *CalendarController) getCalendars(ctx context.Context, req getCalendarsRequest) ([]model.Calendar, error) {
	q := dao.CalendarQuery{
		UserID: req.userID,
		Query: dao.Query{
			SortCol: req.query.APIQuery.SortCol,
			Sort:    req.query.APIQuery.Sort,
			Limit:   req.query.APIQuery.Limit,
			Offset:  req.query.APIQuery.Offset,
		},
	}
	calendars, err := dao.GetCalendars(ctx, q)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return calendars, nil
		}
		return calendars, fmt.Errorf("failed to get calendars. %w", err)
	}
	return calendars, nil
}
