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

type getSnapshotRequest struct {
	UserID     string
	SnapshotID string
}

func handleGetSnapshotError(ctx context.Context, w http.ResponseWriter, err error) {
	if errors.Is(err, ApiErrNotFound) {
		slog.WarnContext(ctx, "error getting snapshot by id", "err", err)
		request.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	slog.ErrorContext(ctx, "error getting snapshot by id", "err", err)
	request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
}

func (c *SnapshotController) GetSnapshot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.DebugContext(ctx, "get snapshot by id request")
	vars := mux.Vars(r)
	userID := vars["user_id"]
	snapshotID := vars["snapshot_id"]

	req := getSnapshotRequest{UserID: userID, SnapshotID: snapshotID}
	snapshot, err := c.getSnapshotByID(ctx, req)
	if err != nil {
		handleGetSnapshotError(ctx, w, err)
		return
	}

	request.RespondWithJSON(w, http.StatusOK, snapshot)
	return
}

func (c *SnapshotController) getSnapshotByID(ctx context.Context, req getSnapshotRequest) (model.Snapshot, error) {
	snapshot, err := dao.GetSnapshotByID(ctx, req.SnapshotID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return snapshot, NewApiError(404, ApiErrNotFound)
		}
		return snapshot, fmt.Errorf("failed to get snapshot by id. %w", err)
	}
	return snapshot, nil
}
