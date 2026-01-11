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
	"github.com/slham/sandbox-api/request"
)

type deleteCalendarRequest struct {
	UserID     string
	CalendarID string
}

func handleDeleteCalendarError(ctx context.Context, w http.ResponseWriter, err error) {
	slog.ErrorContext(ctx, "error deleting calendar", "err", err)
	request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
}

func (c *CalendarController) DeleteCalendar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.DebugContext(ctx, "delete calendar request")
	vars := mux.Vars(r)
	userID := vars["user_id"]
	calendarID := vars["calendar_id"]
	req := deleteCalendarRequest{
		UserID:     userID,
		CalendarID: calendarID,
	}

	err := c.deleteCalendar(ctx, req)
	if err != nil {
		handleDeleteCalendarError(ctx, w, err)
		return
	}

	request.RespondWithJSON(w, http.StatusNoContent, nil)
}

func (c *CalendarController) deleteCalendar(ctx context.Context, req deleteCalendarRequest) error {
	_, err := c.getCalendarByID(ctx, getCalendarRequest{UserID: req.UserID, CalendarID: req.CalendarID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewApiError(404, ApiErrNotFound)
		}
		return fmt.Errorf("failed to delete calendar. %w", err)
	}

	err = dao.DeleteCalendar(ctx, req.UserID, req.CalendarID)
	if err != nil {
		return fmt.Errorf("failed to delete calendar. %w", err)
	}

	return nil
}
