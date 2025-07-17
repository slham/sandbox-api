package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/slham/sandbox-api/dao"
	"github.com/slham/sandbox-api/model"
	"github.com/slham/sandbox-api/request"
)

type getCalendarRequest struct {
	UserID     string
	CalendarID string
}

func handleGetCalendarError(ctx context.Context, w http.ResponseWriter, err error) {
	if errors.Is(err, ApiErrNotFound) {
		slog.WarnContext(ctx, "error getting calendar by id", "err", err)
		request.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	slog.ErrorContext(ctx, "error getting calendar by id", "err", err)
	request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
}

func (c *CalendarController) GetCalendar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.DebugContext(ctx, "get calendar by id request")
	vars := mux.Vars(r)
	userID := vars["user_id"]
	calendarID := vars["calendar_id"]

	req := getCalendarRequest{UserID: userID, CalendarID: calendarID}
	calendar, err := c.getCalendarByID(ctx, req)
	if err != nil {
		handleGetCalendarError(ctx, w, err)
		return
	}

	request.RespondWithJSON(w, http.StatusOK, calendar)
	return
}

func (c *CalendarController) getCalendarByID(ctx context.Context, req getCalendarRequest) (model.Calendar, error) {
	calendar, err := dao.GetCalendarByID(ctx, req.UserID, req.CalendarID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return calendar, NewApiError(404, ApiErrNotFound)
		}
		return calendar, fmt.Errorf("failed to get calendar by id. %w", err)
	}
	return calendar, nil
}
