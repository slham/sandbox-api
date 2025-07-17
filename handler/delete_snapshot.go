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

type deleteSnapshotRequest struct {
	UserID     string
	CalendarID string
	SnapshotID string
}

func handleDeleteSnapshotError(ctx context.Context, w http.ResponseWriter, err error) {
	slog.ErrorContext(ctx, "error deleting snapshot", "err", err)
	request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
	return
}

func (c *SnapshotController) DeleteSnapshot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.DebugContext(ctx, "delete snapshot request")
	vars := mux.Vars(r)
	userID := vars["user_id"]
	calendarID := vars["calendar_id"]
	snapshotID := vars["snapshot_id"]
	req := deleteSnapshotRequest{
		UserID:     userID,
		CalendarID: calendarID,
		SnapshotID: snapshotID,
	}

	err := c.deleteSnapshot(ctx, req)
	if err != nil {
		handleDeleteSnapshotError(ctx, w, err)
		return
	}

	request.RespondWithJSON(w, http.StatusNoContent, nil)
	return
}

func (c *SnapshotController) deleteSnapshot(ctx context.Context, req deleteSnapshotRequest) error {
	_, err := c.getSnapshotByID(ctx, getSnapshotRequest{UserID: req.UserID, SnapshotID: req.SnapshotID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewApiError(404, ApiErrNotFound)
		}
		return fmt.Errorf("failed to delete snapshot. %w", err)
	}

	err = dao.DeleteSnapshot(ctx, req.SnapshotID)
	if err != nil {
		return fmt.Errorf("failed to delete snapshot. %w", err)
	}

	return nil
}
