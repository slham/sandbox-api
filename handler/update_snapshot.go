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

type updateSnapshotRequest struct {
	UserID     string
	SnapshotID string
	CalendarID string         `json:"calendar_id"`
	Done       model.UnixTime `json:"done"`
	Workout    model.Workout  `json:"workout"`
}

func handleUpdateSnapshotError(ctx context.Context, w http.ResponseWriter, err error) {
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

func (c *SnapshotController) UpdateSnapshot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.DebugContext(ctx, "update snapshot request")
	req := updateSnapshotRequest{}
	vars := mux.Vars(r)
	userID := vars["user_id"]
	snapshotID := vars["snapshot_id"]
	calendarID := vars["calendar_id"]
	req.UserID = userID
	req.CalendarID = calendarID
	req.SnapshotID = snapshotID

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(ctx, "error decoding create snapshot request", "err", err)
		request.RespondWithError(w, http.StatusBadRequest, "malformed request body")
		return
	}

	if req.UserID != userID {
		slog.WarnContext(ctx, "user not allowed to modify this snapshot")
		request.RespondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED")
		return
	}

	snapshot, err := c.updateSnapshot(ctx, req)
	if err != nil {
		handleUpdateSnapshotError(ctx, w, err)
		return
	}

	request.RespondWithJSON(w, http.StatusOK, snapshot)
}

func (c *SnapshotController) updateSnapshot(ctx context.Context, req updateSnapshotRequest) (model.Snapshot, error) {
	snapshot := model.Snapshot{}
	if _, err := dao.GetUserByID(ctx, req.UserID); err != nil {
		slog.Warn("failed to find user", "err", err)
		return snapshot, NewApiError(404, ApiErrNotFound).Append("user does not exist")
	}
	if _, err := dao.GetCalendarByID(ctx, req.UserID, req.CalendarID); err != nil {
		slog.Warn("failed to find calendar", "err", err)
		return snapshot, NewApiError(404, ApiErrNotFound).Append("calendar does not exist")
	}

	snapshot, err := dao.GetSnapshotByID(ctx, req.SnapshotID)
	if err != nil {
		slog.Warn("failed to find snapshot", "err", err)
		return snapshot, NewApiError(404, ApiErrNotFound).Append("snapshot does not exist")
	}

	if err := validateUpdateSnapshotRequest(ctx, req); err != nil {
		return snapshot, fmt.Errorf("failed to validate update snapshot request. %w", err)
	}

	snapshot.Done = req.Done
	snapshot.Workout = req.Workout

	if err := dao.UpdateSnapshot(ctx, snapshot); err != nil {
		return snapshot, fmt.Errorf("failed to update snapshot. %w", err)
	}

	return snapshot, nil
}

func validateUpdateSnapshotRequest(ctx context.Context, req updateSnapshotRequest) error {
	apiErr := NewApiError(400, ApiErrBadRequest)

	if req.Done.IsZero() {
		apiErr = apiErr.Append("snapshot must have a done")
	}

	if err := validateCreateWorkoutRequest(ctx, req.Workout); err != nil {
		message := fmt.Sprintf("invalid snapshot workout. %s", err.Error())
		apiErr = apiErr.Append(message)
	}

	if apiErr.HasError() {
		return apiErr
	}

	return nil
}
