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

type getSnapshotsQuery struct {
	APIQuery
}

type getSnapshotsRequest struct {
	userID     string
	calendarID string
	query      getSnapshotsQuery
}

func getSnapshotsQueryParams(ctx context.Context, q url.Values) (getSnapshotsQuery, error) {
	gwq := getSnapshotsQuery{}
	apiQuery, err := getStandardQueryParams(ctx, q)
	if err != nil {
		return gwq, fmt.Errorf("failed to gather query params. %w", err)
	}
	gwq.APIQuery = apiQuery
	return gwq, nil
}

func handleGetSnapshotsError(ctx context.Context, w http.ResponseWriter, err error) {
	if errors.Is(err, ApiErrBadRequest) {
		slog.WarnContext(ctx, "error getting snapshots", "err", err)
		request.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	slog.ErrorContext(ctx, "error getting snapshots", "err", err)
	request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
}

func (c *SnapshotController) GetSnapshots(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	userID := vars["user_id"]
	calendarID := vars["calendar_id"]
	req := getSnapshotsRequest{
		userID:     userID,
		calendarID: calendarID,
	}
	query := r.URL.Query()
	q, err := getSnapshotsQueryParams(ctx, query)
	if err != nil {
		handleGetSnapshotsError(ctx, w, err)
		return
	}

	req.query = q

	snapshots, err := c.getSnapshots(ctx, req)
	if err != nil {
		handleGetSnapshotsError(ctx, w, err)
		return
	}

	request.RespondWithJSON(w, http.StatusOK, snapshots)
}

func (c *SnapshotController) getSnapshots(ctx context.Context, req getSnapshotsRequest) ([]model.Snapshot, error) {
	q := dao.SnapshotQuery{
		Query: dao.Query{
			SortCol: req.query.APIQuery.SortCol,
			Sort:    req.query.APIQuery.Sort,
			Limit:   req.query.APIQuery.Limit,
			Offset:  req.query.APIQuery.Offset,
		},
	}
	snapshots, err := dao.GetSnapshots(ctx, q)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return snapshots, nil
		}
		return snapshots, fmt.Errorf("failed to get snapshots. %w", err)
	}
	return snapshots, nil
}
