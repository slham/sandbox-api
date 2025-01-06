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

type getWorkoutsQuery struct {
	APIQuery
}

type getWorkoutsRequest struct {
	userID string
	query  getWorkoutsQuery
}

func getWorkoutsQueryParams(ctx context.Context, q url.Values) (getWorkoutsQuery, error) {
	gwq := getWorkoutsQuery{}
	apiQuery, err := getStandardQueryParams(ctx, q)
	if err != nil {
		return gwq, fmt.Errorf("failed to gather query params. %w", err)
	}
	gwq.APIQuery = apiQuery
	return gwq, nil
}

func handleGetWorkoutsError(ctx context.Context, w http.ResponseWriter, err error) {
	if errors.Is(err, ApiErrBadRequest) {
		slog.WarnContext(ctx, "error getting workouts", "err", err)
		request.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	slog.ErrorContext(ctx, "error getting workouts", "err", err)
	request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
	return
}

func (c *WorkoutController) GetWorkouts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.DebugContext(ctx, "get workouts request")
	vars := mux.Vars(r)
	userID := vars["user_id"]
	req := getWorkoutsRequest{userID: userID}
	query := r.URL.Query()
	q, err := getWorkoutsQueryParams(ctx, query)
	if err != nil {
		handleGetWorkoutsError(ctx, w, err)
		return
	}

	req.query = q

	workouts, err := c.getWorkouts(ctx, req)
	if err != nil {
		handleGetWorkoutsError(ctx, w, err)
		return
	}

	request.RespondWithJSON(w, http.StatusOK, workouts)
	return
}

func (c *WorkoutController) getWorkouts(ctx context.Context, req getWorkoutsRequest) ([]model.Workout, error) {
	q := dao.WorkoutQuery{
		UserID: req.userID,
		Query: dao.Query{
			SortCol: req.query.APIQuery.SortCol,
			Sort:    req.query.APIQuery.Sort,
			Limit:   req.query.APIQuery.Limit,
			Offset:  req.query.APIQuery.Offset,
		},
	}
	workouts, err := dao.GetWorkouts(ctx, q)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return workouts, nil
		}
		return workouts, fmt.Errorf("failed to get workouts. %w", err)
	}
	return workouts, nil
}
