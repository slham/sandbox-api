package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
	"github.com/slham/sandbox-api/dao"
	"github.com/slham/sandbox-api/model"
	"github.com/slham/sandbox-api/request"
)

func handleCreateCalendarError(ctx context.Context, w http.ResponseWriter, err error) {
	errMsg := "error creating calendar"
	if errors.Is(err, ApiErrBadRequest) {
		slog.WarnContext(ctx, errMsg, "err", err)
		request.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if errors.Is(err, ApiErrConflict) {
		slog.WarnContext(ctx, errMsg, "err", err)
		request.RespondWithError(w, http.StatusConflict, err.Error())
		return
	}

	slog.ErrorContext(ctx, errMsg, "err", err)
	request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
}

func (c *CalendarController) CreateCalendar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.DebugContext(ctx, "create calendar request")
	calendar := model.Calendar{}
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if err := json.NewDecoder(r.Body).Decode(&calendar); err != nil {
		slog.WarnContext(ctx, "error decoding create calendar request", "err", err)
		request.RespondWithError(w, http.StatusBadRequest, "malformed request body")
		return
	}

	calendar.UserID = userID
	calendar, err := c.createCalendar(ctx, calendar)
	if err != nil {
		handleCreateCalendarError(ctx, w, err)
		return
	}

	request.RespondWithJSON(w, http.StatusCreated, calendar)
}

func (c *CalendarController) createCalendar(ctx context.Context, calendar model.Calendar) (model.Calendar, error) {
	slog.DebugContext(ctx, "createCalendar", "userID", calendar.UserID)
	if _, err := dao.GetUserByID(ctx, calendar.UserID); err != nil {
		return calendar, NewApiError(404, ApiErrNotFound).Append("user does not exist")
	}

	if err := validateCreateCalendarRequest(ctx, calendar); err != nil {
		return calendar, fmt.Errorf("failed to validate create calendar request. %w", err)
	}

	calendar.ID = newCalendarID()

	calendar, err := dao.InsertCalendar(ctx, calendar)
	if err != nil {
		if errors.Is(err, dao.ErrConflictCalendarName) {
			return calendar, NewApiError(http.StatusConflict, ApiErrConflict).Append("calendar name already exists")
		}
		return calendar, fmt.Errorf("failed to insert calendar. %w", err)
	}

	return calendar, nil
}

func validateCreateCalendarRequest(ctx context.Context, calendar model.Calendar) error {
	apiErr := NewApiError(http.StatusBadRequest, ApiErrBadRequest)

	if calendar.Name == "" {
		apiErr = apiErr.Append("calendar must have a name")
	}

	if apiErr.HasError() {
		return apiErr
	}

	return nil
}

func newCalendarID() string {
	return fmt.Sprintf("cal_%s", ksuid.New().String())
}
