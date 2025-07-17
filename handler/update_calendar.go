package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/slham/sandbox-api/dao"
	"github.com/slham/sandbox-api/model"
	"github.com/slham/sandbox-api/request"
)

type updateCalendarRequest struct {
	UserID     string
	CalendarID string
	Name       string `json:"name"`
}

func handleUpdateCalendarError(ctx context.Context, w http.ResponseWriter, err error) {
	if errors.Is(err, ApiErrBadRequest) {
		slog.WarnContext(ctx, "error creating calendar", "err", err)
		request.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if errors.Is(err, ApiErrConflict) {
		slog.WarnContext(ctx, "error creating calendar", "err", err)
		request.RespondWithError(w, http.StatusConflict, err.Error())
		return
	}

	slog.ErrorContext(ctx, "error creating calendar", "err", err)
	request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
}

func (c *CalendarController) UpdateCalendar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.DebugContext(ctx, "update calendar request")
	req := updateCalendarRequest{}
	vars := mux.Vars(r)
	userID := vars["user_id"]
	calendarID := vars["calendar_id"]
	req.UserID = userID
	req.CalendarID = calendarID

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(ctx, "error decoding create calendar request", "err", err)
		request.RespondWithError(w, http.StatusBadRequest, "malformed request body")
		return
	}

	if req.UserID != userID {
		slog.WarnContext(ctx, "user not allowed to modify this calendar")
		request.RespondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED")
		return
	}

	calendar, err := c.updateCalendar(ctx, req)
	if err != nil {
		handleUpdateCalendarError(ctx, w, err)
		return
	}

	request.RespondWithJSON(w, http.StatusOK, calendar)
}

func (c *CalendarController) updateCalendar(ctx context.Context, req updateCalendarRequest) (model.Calendar, error) {
	calendar := model.Calendar{}
	if _, err := dao.GetUserByID(ctx, req.UserID); err != nil {
		slog.Warn("failed to find user", "err", err)
		return calendar, NewApiError(404, ApiErrNotFound).Append("user does not exist")
	}

	calendar, err := dao.GetCalendarByID(ctx, req.UserID, req.CalendarID)
	if err != nil {
		slog.Warn("failed to find calendar", "err", err)
		return calendar, NewApiError(404, ApiErrNotFound).Append("calendar does not exist")
	}

	if err := validateUpdateCalendarRequest(ctx, req); err != nil {
		return calendar, fmt.Errorf("failed to validate update calendar request. %w", err)
	}

	calendar.Name = req.Name

	if err := dao.UpdateCalendar(ctx, calendar); err != nil {
		if errors.Is(err, dao.ErrConflictCalendarName) {
			return calendar, NewApiError(409, ApiErrConflict).Append("calendar name already exists")
		}
		return calendar, fmt.Errorf("failed to update calendar. %w", err)
	}

	return calendar, nil
}

func validateUpdateCalendarRequest(ctx context.Context, req updateCalendarRequest) error {
	apiErr := NewApiError(400, ApiErrBadRequest)

	if req.Name == "" {
		apiErr = apiErr.Append("calendar must have a name")
	}

	if apiErr.HasError() {
		return apiErr
	}

	return nil
}
