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

func handleCreateSnapshotError(ctx context.Context, w http.ResponseWriter, err error) {
	if errors.Is(err, ApiErrBadRequest) {
		slog.WarnContext(ctx, "error creating snapshot", "err", err)
		request.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if errors.Is(err, ApiErrConflict) {
		slog.WarnContext(ctx, "error creating snapshot", "err", err)
		request.RespondWithError(w, http.StatusConflict, err.Error())
		return
	}

	slog.ErrorContext(ctx, "error creating snapshot", "err", err)
	request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
}

func (c *SnapshotController) CreateSnapshot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.DebugContext(ctx, "create snapshot request")
	snapshot := model.Snapshot{}
	vars := mux.Vars(r)
	userID := vars["user_id"]
	calendarID := vars["calendar_id"]

	if err := json.NewDecoder(r.Body).Decode(&snapshot); err != nil {
		slog.WarnContext(ctx, "error decoding create snapshot request", "err", err)
		request.RespondWithError(w, http.StatusBadRequest, "malformed request body")
		return
	}

	snapshot.UserID = userID
	snapshot.CalendarID = calendarID
	snapshot, err := c.createSnapshot(ctx, snapshot)
	if err != nil {
		handleCreateSnapshotError(ctx, w, err)
		return
	}

	request.RespondWithJSON(w, http.StatusCreated, snapshot)
}

func (c *SnapshotController) createSnapshot(ctx context.Context, snapshot model.Snapshot) (model.Snapshot, error) {
	slog.DebugContext(ctx, "createSnapshot", "calendarID", snapshot.CalendarID)
	if _, err := dao.GetCalendarByID(ctx, snapshot.UserID, snapshot.CalendarID); err != nil {
		return snapshot, NewApiError(404, ApiErrNotFound).Append("calendar does not exist")
	}

	if err := validateCreateSnapshotRequest(ctx, snapshot); err != nil {
		return snapshot, fmt.Errorf("failed to validate create snapshot request. %w", err)
	}

	snapshot.ID = newSnapshotID()

	snapshot, err := dao.InsertSnapshot(ctx, snapshot)
	if err != nil {
		return snapshot, fmt.Errorf("failed to insert snapshot. %w", err)
	}

	return snapshot, nil
}

func validateCreateSnapshotRequest(ctx context.Context, snapshot model.Snapshot) error {
	apiErr := NewApiError(http.StatusBadRequest, ApiErrBadRequest)

	if snapshot.Done.IsZero() {
		apiErr = apiErr.Append("snapshot must have a done")
	}

	if err := validateCreateWorkoutRequest(ctx, snapshot.Workout); err != nil {
		message := fmt.Sprintf("invalid snapshot workout. %s", err.Error())
		apiErr = apiErr.Append(message)
	}

	if apiErr.HasError() {
		return apiErr
	}

	return nil
}

func newSnapshotID() string {
	return fmt.Sprintf("snap_%s", ksuid.New().String())
}
